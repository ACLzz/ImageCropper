package main

import (
	"fmt"
	"github.com/ACLzz/ImageCropper/src/broker"
	"github.com/ACLzz/ImageCropper/src/config"
	"github.com/ACLzz/ImageCropper/src/handlers"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	r := gin.Default()
	handlers.MainRouter(r)
	broker.StartServices()

	if err := r.Run(fmt.Sprintf("%s:%d", config.ConfigObj.Host, config.ConfigObj.Port));
	err != nil {
		logrus.Fatal(err)
	}
}