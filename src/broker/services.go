package broker

import (
	"github.com/ACLzz/ImageCropper/src/config"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

const CropperQueueName = "cropper"

func StartServices()  {
	// prepares queues for services and starts them
	ch, cls := GetChannel()
	defer cls()
	if _, err := ch.
		QueueDeclare(CropperQueueName, true, false, false, false, nil); err != nil {
		logrus.Fatal("error declaring queue: " + err.Error())
	}

	if err := ch.Qos(1, 0, false); err != nil {
		logrus.Error("cannot enable qos")
	}

	StartCropService()
}

func GetChannel() (*amqp.Channel, func()) {
	// returns channel interface and close function
	conn, err := amqp.Dial(config.ConfigObj.BrokerUrl)
	if err != nil {
		logrus.Fatal("rabbitMQ not available:", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		logrus.Fatal("couldn't open rabbitMQ channel: ",err)
	}

	cls := func() {
		ch.Close()
		conn.Close()
	}
	return ch, cls
}