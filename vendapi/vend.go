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
	"github.com/jtrotsky/vend-image-upload/image"
)

// imageUpload uploads a single product image to Vend.
func imageUpload(products []Product) error {

	for _, product := range products {

		// This checks we actually have an image to post.
		if len(product.Images) > 0 {

			// If there are more than one image we will continue looping
			// until there are no more to loop through.
			for _, i := range product.Images {

				// First grab and save the image from the URL.
				// TODO: Grab and save image.
				imageName := image.Grab(*product.ID, *i.URL)
				path := fmt.Sprintf("%s", imageName)

				// Make sure image is removed at end.
				defer os.Remove(path)

				var body bytes.Buffer
				// Start multipart writer.
				writer := multipart.NewWriter(&body)

				// Key "image" value is the image binary.
				part, err := writer.CreateFormFile("image", path)
				if err != nil {
					log.Fatalf("Error creating multipart form file: %s", err)
					return err
				}

				// Open image file.
				file, err := os.Open(path)
				if err != nil {
					log.Fatalf("Error opening image file: %s", err)
					return err
				}

				// Make sure file is closed and then removed at end.
				defer file.Close()
				defer os.Remove(path)

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
				var domainPrefix, productID string
				url := vend.ImageUploadURLFactory(domainPrefix, productID)

				req, err := http.NewRequest("POST", url, &body)

				req.Header.Set("User-agent", "Support-tool: choppily - one of JOEYM8's tools.")
				req.Header.Set("Content-Type", writer.FormDataContentType())
				// TODO: Token
				token := ""
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

				client := &http.Client{}

				// Make the request.
				// TODO:
				var attempt int
				var res *http.Response
				for {
					time.Sleep(time.Second)
					res, err = client.Do(req)
					if err != nil || !vend.ResponseCheck(res.StatusCode) {
						if res.StatusCode == 404 {
							break
						}
						fmt.Printf("\nError performing request: %s Status code: %d", err, res.StatusCode)
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

				// Unmarshal JSON response into our respone struct.
				// from this we can find info about the image status.
				response := ImageUploadResponse{}
				err = json.Unmarshal(resBody, &response)
				if err != nil {
					log.Fatalf("Error unmarhsalling response body: %s", err)
					return err
				}

				payload := response.Data

				fmt.Printf("\nImage created at position: %d\n", *payload.Position)
			}
		}
	}
	return nil
}
