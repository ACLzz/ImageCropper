package config

import (
	"fmt"
	"os"
)

type Config struct {
	OrigPicsDest    string
	CroppedPicsDest string
	ExtraFolder		string
	Host			string
	Port			int
	BrokerUrl		string
}

func getFolder(path string) string {
	configFolder := fmt.Sprint(os.Getenv("HOME"), "/.imcr/")
	return fmt.Sprint(configFolder, path)
}

var ConfigObj = Config{
	OrigPicsDest:		getFolder("orig/"),
	CroppedPicsDest:	getFolder("cropped/"),
	ExtraFolder:		getFolder("extra/"),
	Host:				"localhost",
	Port:				8888,
	BrokerUrl:			"amqp://guest:guest@localhost:5672",
}
