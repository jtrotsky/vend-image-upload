// Package vendapi handles interactions with the Vend API.
package vendapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/jtrotsky/govend/vend"
	log "github.com/sirupsen/logrus"
)

// UploadImage uploads a single product image to Vend.
func UploadImage(authToken, domainPrefix, imagePath string, product ProductUpload) error {
	var err error

	// This checks we actually have an image to post.
	if len(product.ImageURL) > 0 {

		// First grab and save the image from the URL.
		imageURL := fmt.Sprintf("%s", product.ImageURL)

		var body bytes.Buffer
		// Start multipart writer.
		writer := multipart.NewWriter(&body)

		// Key "image" value is the image binary.
		var part io.Writer
		part, err = writer.CreateFormFile("image", imageURL)
		if err != nil {
			log.WithError(err).Warning("Error creating multipart form file")
			return err
		}

		// Open image file.
		var file *os.File
		file, err = os.Open(imagePath)
		if err != nil {
			log.WithError(err).Warning("Error opening image file")
			return err
		}

		// Make sure file is closed and then removed at end.
		defer file.Close()
		defer os.Remove(imageURL)

		// Copying image binary to form file.
		// Ignoring number of bytes copied.
		_, err = io.Copy(part, file)
		if err != nil {
			log.Fatalf("Error copying file for requst body: %s", err)
			return err
		}

		err = writer.Close()
		if err != nil {
			log.WithError(err).Warning("Error closing writer")
			return err
		}

		// Create the Vend URL to send our image to.
		url := vend.ImageUploadURLFactory(domainPrefix, product.ID)

		log.WithField("url", url).Info("Uploading to Vend")

		req, _ := http.NewRequest("POST", url, &body)

		req.Header.Set("User-agent", "vend-image-upload")
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))

		client := &http.Client{}

		// Make the request.
		var attempt int
		var res *http.Response
		for {
			time.Sleep(time.Second)
			res, err = client.Do(req)
			// Catch error.
			if err != nil || !vend.ResponseCheck(res.StatusCode) {
				log.WithError(err).WithField("status", res.StatusCode).Info("Error performing request")
				// Delays between attempts will be exponentially longer each time.
				attempt++
				delay := vend.BackoffDuration(attempt)
				time.Sleep(delay)
			} else {
				// Ensure that image file is removed after it's uploaded.
				os.Remove(imagePath)
				break
			}
		}

		// Make sure response body is closed at end.
		defer res.Body.Close()

		var resBody []byte
		resBody, err = ioutil.ReadAll(res.Body)
		if err != nil {
			log.WithError(err).Warning("Error reading response body")
			return err
		}

		// fmt.Println(string(resBody))

		// Unmarshal JSON response into our respone struct.
		// from this we can find info about the image status.
		response := ImageUpload{}
		err = json.Unmarshal(resBody, &response)
		if err != nil {
			log.WithError(err).Warning("Error unmarhsalling response body")
			return err
		}

		payload := response.Data

		log.WithField("position", *payload.Position).Info("Image created")
	}
	return err
}
