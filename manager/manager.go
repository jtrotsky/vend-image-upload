package manager

import (
	"fmt"
	"log"

	"github.com/jtrotsky/govend/vend"
	"github.com/jtrotsky/vend-image-upload/image"
	"github.com/jtrotsky/vend-image-upload/reader"
	"github.com/jtrotsky/vend-image-upload/vendapi"
)

// Manager contains the Vend client.
type Manager struct {
	client vend.Client
}

// NewManager creates an instance of manager.
func NewManager(client vend.Client) *Manager {
	return &Manager{client}
}

// TODO: Comment syntax.

// Run reads the product CSV, gets all products from Vend, then posts their images.
func (manager *Manager) Run(filePath string) {
	// Log opening timestamp.
	log.Printf("BEGIN\n")

	imagePayload, err := reader.ReadCSV(filePath)
	if err != nil {
		log.Fatalf("Error reading CSV file: %s", err)
	}

	// TODO:
	// fmt.Printf("%s", imagePayload)

	fmt.Printf("\nGrabbing products.\n")
	// Get all products, ignore productMap.
	_, productMap, err := manager.client.Products()
	if err != nil {
		log.Fatalf("Failed to get products.: %s", err)
	}

	products, err := matchVendProduct(productMap, imagePayload)

	for _, i := range *imagePayload {
		// Loop through each CSV line, match it's respective product ID
		// scrape its image, then post its image.
		imagePath := image.Grab(i)
		fmt.Printf("\n\n%s\n", i)
		vendapi.ImageUpload(imagePath, *products)
	}
}

func matchVendProduct(productMap *map[string]vend.Product, imagePayload *[]vendapi.UploadProduct) (*[]vendapi.UploadProduct, error) {
	// TODO: Loop through handle/sku combos to get product ids from Product then
	// build them into the uploadProduct payload.
	var err error
	return err
}
