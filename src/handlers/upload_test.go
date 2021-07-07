package handlers

import (
	"bytes"
	"fmt"
	"github.com/ACLzz/ImageCropper/src/config"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"testing"
)

var serverUrl = fmt.Sprintf("http://%s:%d/", config.ConfigObj.Host, config.ConfigObj.Port)

func TestUploadImage(_t *testing.T) {
	_t.Run("file was not sent", func(t *testing.T) {								// request without file
		resp := Post(fmt.Sprint(serverUrl, "upload_image"), "", t)
		expectStatus := 400
		if resp.StatusCode != expectStatus {
			t.Errorf("expected status was %d but it is %d", expectStatus, resp.StatusCode)
		}
	})

	_t.Run("successful file sending", func(t *testing.T) {
		testFilePath := fmt.Sprint(config.ConfigObj.ExtraFolder, "test.png")
		resp := Post(fmt.Sprint(serverUrl, "upload_image"), testFilePath, t)

		t.Run("status code", func(t *testing.T) {	// check status code
			expectStatus := 200
			if resp.StatusCode != expectStatus {
				t.Errorf("expected status was %d but it is %d", expectStatus, resp.StatusCode)
			}
		})

		t.Run("file uploaded", func(t *testing.T) {	// check if orig was saved
			dir, _ := os.ReadDir(config.ConfigObj.OrigPicsDest)
			if len(dir) < 1 {
				t.Error("file wasn't saved")
			} else {
				filepath := fmt.Sprint(config.ConfigObj.OrigPicsDest, dir[0].Name())
				if err := os.Remove(filepath); err != nil {
					t.Error("cannot remove file with path ", filepath)
				}
			}
		})
	})
}

func Post(url string, filepath string, t *testing.T) *http.Response {
	// Post request to upload image
	var Client = &http.Client{}
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		t.Error(err)
	}

	if len(filepath) > 0 {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		fw, err := writer.CreateFormFile("file", filepath)
		if err != nil {
			t.Error(err)
		}
		file, err := os.Open(filepath)
		if err != nil {
			t.Error(err)
		}
		if _, err := io.Copy(fw, file); err != nil {
			t.Error(err)
		}
		writer.Close()

		req, err = http.NewRequest("POST", url, bytes.NewReader(body.Bytes()))
		if err != nil {
			t.Error(err)
		}
		req.Header.Set("Content-Type", writer.FormDataContentType())
	}

	resp, err := Client.Do(req)
	if err != nil {
		t.Error(err)
	}

	return resp
}