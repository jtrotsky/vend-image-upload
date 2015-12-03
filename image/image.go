package image

// Grabs and saves image locally, then returns its name.
import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/jtrotsky/govend/vend"
)

// Grab ...
func Grab(productID, URL string) string {

	// Grab the image and write it to a file.
	image, err := urlGet(URL)
	if err != nil {
		fmt.Printf("Something when wrong grabbing image: %s", err)
	}

	// Split the Shopify URL up to make it easier to grab the file extension.
	// Shopify supports jpg, gif, and png.
	parts := strings.Split(URL, ".")
	extension := parts[3]

	fileName := fmt.Sprintf("%s.%s", productID, extension[:len(extension)-2])

	// Write product data to file
	err = ioutil.WriteFile(fileName, image, 0666)
	if err != nil {
		fmt.Printf("Something went wrong writing image to file.")
	}

	return fileName
}

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

	vend.ResponseCheck(res.StatusCode)

	// Read what we got back.
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("\nError while reading response body: %s\n", err)
		return nil, err
	}

	return body, err
}
