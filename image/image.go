package image

// Grabs and saves image locally, then returns its name.
import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/jtrotsky/govend/vend"
	"github.com/jtrotsky/vend-image-upload/vendapi"
)

// Grab ...
func Grab(products vendapi.UploadProduct) string {

	// Grab the image and write it to a file.
	image, err := urlGet(*products.ImageURL)
	if err != nil {
		fmt.Printf("Something when wrong grabbing image: %s", err)
	}

	// Split the URL up to make it easier to grab the file extension.
	parts := strings.Split(*products.ImageURL, ".")
	// TODO: Confirm URL scheme.
	extension := parts[len(parts)-1]

	fileName := fmt.Sprintf("%s.%s", *products.SKU, extension)

	// Write product data to file
	err = ioutil.WriteFile(fileName, image, 0666)
	if err != nil {
		// TODO: error? log?
		fmt.Printf("Something went wrong writing image to file.\n")
	}

	return fileName
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

	fmt.Printf("\nGrabbing: %s", url)
	// Doing the request.
	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("\nError performing request: %s", err)
		return nil, err
	}
	// Make sure response body is closed at end.
	defer res.Body.Close()

	// TODO: Check if true/false
	vend.ResponseCheck(res.StatusCode)

	// Read what we got back.
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("\nError while reading response body: %s\n", err)
		return nil, err
	}

	return body, err
}
