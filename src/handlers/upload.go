package handlers

import (
	"fmt"
	"github.com/ACLzz/ImageCropper/src/broker"
	"github.com/ACLzz/ImageCropper/src/config"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/thanhpk/randstr"
	"net/http"
	"strings"
)

func UploadImage(c *gin.Context) {
	fileField := "file"
	file, err := c.FormFile(fileField)
	if err != nil {
		logrus.Error(err)
		c.String(http.StatusBadRequest, fmt.Sprintf("file in field '%s' was not sent", fileField))
		return
	}

	_fileExtension := strings.Split(file.Filename, ".")
	fileExtension := _fileExtension[len(_fileExtension)-1]
	fn := fmt.Sprint(randstr.Hex(16), ".", fileExtension)
	filepath := fmt.Sprint(config.ConfigObj.OrigPicsDest, fn)

	if err := c.SaveUploadedFile(file, filepath); err != nil {
		return
	}

	go func() {
		ch, cls := broker.GetChannel()
		defer cls()

		f, err := file.Open()
		if err != nil {
			logrus.Error(err)
			return
		}
		defer f.Close()
		image := broker.ConvertFileToMessage(fn, f)
		if image == nil {
			return
		}

		// publish image to queue
		err = ch.Publish("", broker.CropperQueueName, false, false, *image)
	}()

	c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
}
