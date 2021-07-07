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

var cropperSizes = [3]int{64, 128, 256}

func cropImages() {
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
		_imageExtension := strings.Split(imageMsg.Filename, ".")
		imageExtension := strings.ToLower(_imageExtension[len(_imageExtension)-1])

		for _, res := range cropperSizes {
			filepath := fmt.Sprintf("%s%dx%d/%s", config.ConfigObj.CroppedPicsDest, res, res, imageMsg.Filename)

			data, err := base64.StdEncoding.DecodeString(imageMsg.Image)
			if err != nil {
				logrus.Error(err)
				continue
			}

			reader := bytes.NewReader(data)
			_, err = reader.Seek(0, 0)
			if err != nil {
				logrus.Error(err)
				continue
			}

			src, _, err := image.Decode(reader)
			if err != nil {
				logrus.Error(err)
				continue
			}

			dst := image.NewRGBA(image.Rect(0, 0, res, res))
			draw.CatmullRom.Scale(dst, dst.Bounds(),
				src, src.Bounds(),
				draw.Over, nil)

			file, err := os.Create(filepath)
			if err != nil {
				logrus.Error(err)
				continue
			}

			switch imageExtension {
			case "png":
				if err := png.Encode(file, dst); err != nil {
					logrus.Error(err)
				}
			case "jpg":
				if err := jpeg.Encode(file, dst, nil); err != nil {
					logrus.Error(err)
				}
			}

			file.Close()
		}
	}
}

func StartCropService() {
	go func() {
		logrus.Info("starting cropper service")
		cropImages()
	}()
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