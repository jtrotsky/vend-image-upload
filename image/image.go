package image

// Grabs and saves image locally, then returns its name.
import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/jtrotsky/govend/vend"
	"github.com/jtrotsky/vend-image-upload/vendapi"
)

// Grab ...
func Grab(products vendapi.UploadProduct) (string, error) {

	// Grab the image and write it to a file.
	image, err := urlGet(*products.ImageURL)
	if err != nil {
		fmt.Printf("Something when wrong grabbing image: %s", err)
	}

	// Split the URL up to make it easier to grab the file extension.
	parts := strings.Split(*products.ImageURL, ".")
	// TODO: Confirm URL scheme. Cannot get referral urls, like goo.gl.
	extension := parts[len(parts)-1]

	fileName := fmt.Sprintf("%s.%s", *products.SKU, extension)

	// Write product data to file
	// TODO: Confirm correct chmod
	err = ioutil.WriteFile(fileName, image, 0666)
	if err != nil {
		// TODO: error? log?
		fmt.Printf("Something went wrong writing image to file.\n")
	}

	return fileName, err
}

// TODO: url get body
// Or split to return either response or body
// Get body takes response and returns body.
func urlGet(url string) ([]byte, error) {

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("\nError creating http request: %s", err)
		return nil, err
	}

	fmt.Printf("Grabbing: %s", url)
	// Doing the request.
	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("\nError performing request: %s", err)
		return nil, err
	}
	// Make sure response body is closed at end.
	defer res.Body.Close()

	// Check HTTP response.
	// TODO: More descriptive errors.
	if !vend.ResponseCheck(res.StatusCode) {
		log.Fatalf("Error: %d", res.StatusCode)
	}

	// Read what we got back.
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("\nError while reading response body: %s\n", err)
		return nil, err
	}

	return body, err
}
