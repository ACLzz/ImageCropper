package broker

import (
	"encoding/base64"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"io"
	"io/ioutil"
	"log"
)

type AddImage struct {
	Filename	string
	Image		string
}

func cropImages() {
	ch, cls := GetChannel()
	defer cls()
	qCh, err := ch.Consume(CropperQueueName, "",
		false, false, false, false, nil)
	if err != nil {
		logrus.Errorf("error in consuming %s queue channel", CropperQueueName)
		return
	}

	for m := range qCh {
		image := &AddImage{}
		err := json.Unmarshal(m.Body, image)
		if err != nil {
			log.Printf("Error decoding JSON: %s", err)
		}

		logrus.Infof("recived file with name %s", image.Filename)
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