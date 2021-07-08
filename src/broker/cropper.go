package broker

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/ACLzz/ImageCropper/src/config"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"golang.org/x/image/draw"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type AddImage struct {
	Filename	string
	Image		string
}

func StartCropService() {
	go func() {
		logrus.Info("starting cropper service")
		cropService()
	}()
}

func cropService() {
	ch, cls := GetChannel()
	defer cls()
	qCh, err := ch.Consume(CropperQueueName, "",
		true, false, false, false, nil)
	if err != nil {
		logrus.Errorf("error in consuming %s queue channel", CropperQueueName)
		return
	}

	for m := range qCh {
		imageMsg := &AddImage{}
		err := json.Unmarshal(m.Body, imageMsg)
		if err != nil {
			log.Printf("Error decoding JSON: %s", err)
			continue
		}
		logrus.Infof("recived file with name %s", imageMsg.Filename)

		for _, res := range config.ConfigObj.CropperSizes {
			// get image from string
			src := decodeImage(imageMsg.Image)
			if src == nil { continue }

			// crop and save it
			cropImage(res, res, imageMsg.Filename, *src)
		}
	}
}

func ConvertFileToMessage(filename string, file io.Reader) *amqp.Publishing {
	data, err := ioutil.ReadAll(file)
	if err != nil {
		logrus.Error(err)
		return nil
	}

	imageMessage := AddImage{
		Filename: filename,
		Image: base64.StdEncoding.EncodeToString(data),
	}

	body, err := json.Marshal(imageMessage)
	if err != nil {
		logrus.Error("error in encoding message", err)
		return nil
	}

	message := amqp.Publishing{
		Body: body,
		ContentType:  "text/plain",
		DeliveryMode: amqp.Persistent,
	}
	return &message
}

func decodeImage(imgStr string) *image.Image {
	data, err := base64.StdEncoding.DecodeString(imgStr)
	if err != nil {
		logrus.Error(err)
		return nil
	}

	reader := bytes.NewReader(data)
	_, err = reader.Seek(0, 0)
	if err != nil {
		logrus.Error(err)
		return nil
	}

	src, _, err := image.Decode(reader)
	if err != nil {
		logrus.Error(err)
		return nil
	}

	return &src
}

func cropImage(x, y int, filename string, src  image.Image) {
	// crop image
	dst := image.NewRGBA(image.Rect(0, 0, x, y))
	draw.CatmullRom.Scale(dst, dst.Bounds(),
		src, src.Bounds(),
		draw.Over, nil)

	// create directory for that resolution
	resolutionPath := fmt.Sprintf("%s%dx%d/", config.ConfigObj.CroppedPicsDest, x, y)
	if err := os.Mkdir(resolutionPath, 0755); err != nil {
		logrus.Error(err)
	}

	// write image to file
	filepath := resolutionPath + filename
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		logrus.Error(err)
		return
	}

	_imageExtension := strings.Split(filename, ".")
	imageExtension := strings.ToLower(_imageExtension[len(_imageExtension)-1])

	switch imageExtension {
	case "png":
		if err := png.Encode(file, dst); err != nil {
			logrus.Error(err)
		}
	case "jpg":
		if err := jpeg.Encode(file, dst, nil); err != nil {
			logrus.Error(err)
		}
	default:
		logrus.Errorf("format '%s' is unsupported", imageExtension)
	}
	file.Close()
}