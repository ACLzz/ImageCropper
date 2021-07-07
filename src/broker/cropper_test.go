package broker

import (
	"fmt"
	"github.com/ACLzz/ImageCropper/src/config"
	"os"
	"testing"
)

func TestCropper(t *testing.T) {
	StartCropService()
	ch, cls := GetChannel()
	defer cls()

	fn := "test.png"
	filepath := config.ConfigObj.ExtraFolder + fn
	file, err := os.OpenFile(filepath, os.O_RDONLY, 0644)
	if err != nil {
		t.Error(err)
	}
	defer file.Close()

	image := ConvertFileToMessage(fn, file)
	err = ch.Publish("", CropperQueueName, false, false, *image)

	cropRes := [3]string{"256x256", "128x128", "64x64"}
	for _, res := range cropRes {
		croppedFP := fmt.Sprint(config.ConfigObj.CroppedPicsDest, res, "/", fn)
		if _, err := os.ReadFile(croppedFP); err != nil {
			t.Error(err)
		} else {
			if err := os.Remove(croppedFP); err != nil {
				t.Error("cannot delete ", croppedFP)
			}
		}
	}
}
