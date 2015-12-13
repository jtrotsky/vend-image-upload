// Package vendapi handles interactions with the Vend API.
package vendapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/jtrotsky/govend/vend"
)

// TODO: condensed version of below.
// func Upload(imagePath string) error {
// 	var err error
// 	return err
// }

// ImageUpload uploads a single product image to Vend.
func ImageUpload(authToken, domainPrefix, imagePath string, product UploadProduct) error {

	var err error

	// Check an image exists.
	if product.ImageURL == nil {
		return err
	}

	// This checks we actually have an image to post.
	if len(*product.ImageURL) > 0 {

		// First grab and save the image from the URL.
		imageURL := fmt.Sprintf("%s", *product.ImageURL)

		var body bytes.Buffer
		// Start multipart writer.
		writer := multipart.NewWriter(&body)

		// Key "image" value is the image binary.
		part, err := writer.CreateFormFile("image", imageURL)
		if err != nil {
			log.Fatalf("Error creating multipart form file: %s", err)
			return err
		}

		// Open image file.
		file, err := os.Open(imagePath)
		if err != nil {
			log.Fatalf("Error opening image file: %s", err)
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
			log.Fatalf("Error closing writer: %s", err)
			return err
		}

		// Create the Vend URL to send our image to.
		// TODO: domainprefix/productID
		url := vend.ImageUploadURLFactory(domainPrefix, product.ID)

		fmt.Printf("\nUploading image to: %s\n", url)

		req, err := http.NewRequest("POST", url, &body)

		req.Header.Set("User-agent", "vend-image-upload")
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))

		client := &http.Client{}

		// Make the request.
		// TODO:
		var attempt int
		var res *http.Response
		for {
			time.Sleep(time.Second)

			res, err = client.Do(req)
			// Catch error.
			// TODO: !vend.ResponseCheck(res.StatusCode)
			if err != nil {
				fmt.Println(err)
				// fmt.Printf("\nError performing request: %s Status code: %d", err, res.StatusCode)
				// Delays between attempts will be exponentially longer each time.
				attempt++
				delay := vend.BackoffDuration(attempt)
				time.Sleep(delay)
			} else {
				break
			}
		}

		// Make sure response body is closed at end.
		defer res.Body.Close()

		resBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Printf("Error while reading response body: %s\n", err)
			return err
		}

		fmt.Println(string(resBody))

		// Unmarshal JSON response into our respone struct.
		// from this we can find info about the image status.
		response := ImageUploadResponse{}
		err = json.Unmarshal(resBody, &response)
		if err != nil {
			log.Fatalf("Error unmarhsalling response body: %s", err)
			return err
		}

		payload := response.Data

		fmt.Printf("\nImage created at position: %d\n\n", *payload.Position)
	}
	return err
}
