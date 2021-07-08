package broker

import (
	"bytes"
	"fmt"
	"github.com/ACLzz/ImageCropper/src/config"
	img "image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"testing"
	"time"
)

func TestCropper(t *testing.T) {
	StartCropService()
	ch, cls := GetChannel()
	defer cls()

	fn := "test.png"
	filepath := config.ConfigObj.ExtraFolder + fn
	file, err := os.OpenFile(filepath, os.O_RDONLY, 0755)
	if err != nil {
		t.Error(err)
		return
	}
	defer file.Close()

	image := ConvertFileToMessage(fn, file)
	err = ch.Publish("", CropperQueueName, false, false, *image)
	<-time.After(2 * time.Second)	// wait for crop

	for _, res := range config.ConfigObj.CropperSizes {
		croppedFP := fmt.Sprint(config.ConfigObj.CroppedPicsDest, res,"x", res, "/", fn) // ~/.imrc/cropped/16x16/filename.ext
		file, err := os.OpenFile(croppedFP, os.O_RDONLY, 0755)
		if err != nil {
			t.Error(err)
			continue
		}
		data, err := io.ReadAll(file)
		file.Close()
		if err != nil {
			t.Error(err)
			continue
		} else {
			if err := os.Remove(croppedFP); err != nil {
				t.Error("cannot delete ", croppedFP)
				continue
			}
		}

		imgObj, _, err := img.Decode(bytes.NewReader(data))
		if err != nil {
			t.Error(err)
			continue
		}
		imgSize := imgObj.Bounds()
		if imgSize != img.Rect(0, 0, res, res) {
			t.Errorf("image boundes must be %dx%d but it is %dx%d", res, res, imgSize.Dx(), imgSize.Dy())
		}
	}
}
