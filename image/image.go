// Package image grabs and saves image locally, then returns its name.
package image

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/jtrotsky/govend/vend"
	"github.com/jtrotsky/vend-image-upload/vendapi"
	log "github.com/sirupsen/logrus"
)

// Grab downloads a product image and writes it to a file.
func Grab(products vendapi.ProductUpload) (string, error) {

	// Grab the image based on provided URL.
	image, err := urlGet(products.ImageURL)
	if err != nil {
		return "", err
	}

	// Split the URL up to make it easier to grab the file extension.
	parts := strings.Split(products.ImageURL, ".")
	extension := parts[len(parts)-1]
	// If the extension looks about the right length then use it for the
	// filename, otherwise do not.
	var fileName string
	if len(extension) == 3 {
		fileName = fmt.Sprintf("%s.%s", products.ID, extension)
	} else {
		fileName = fmt.Sprintf("%s", products.ID)
	}

	// Write product data to file
	err = ioutil.WriteFile(fileName, image, 0666)
	if err != nil {
		log.WithError(err).Error("Something went wrong writing image to file")
		return "", err
	}

	return fileName, err
}

// TODO: url get body
// Or split to return either response or body
// Get body takes response and returns body.
func urlGet(url string) ([]byte, error) {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	log.WithField("external_url", url).Info("Getting image from external url")

	// Doing the request.
	res, err := client.Get(url)
	if err != nil {
		log.WithError(err).Error("Error performing request")
		return nil, err
	}
	// Make sure response body is closed at end.
	defer res.Body.Close()

	// Check HTTP response.
	if !vend.ResponseCheck(res.StatusCode) {
		log.WithField("status", res.StatusCode).Error("Bad status code")
		return nil, err
	}

	// Read what we got back.
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.WithError(err).Error("Error while reading response body")
		return nil, err
	}

	return body, err
}
