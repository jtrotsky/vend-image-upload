package manager

import (
	"fmt"
	"log"
	"os"

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

	// Read provided CSV file and store product info.
	imagePayload, err := reader.ReadCSV(filePath)
	if err != nil {
		log.Fatalf("Error reading CSV file: %s", err)
	}

	fmt.Printf("\nGrabbing products.\n")
	// Get all products from Vend.
	_, productMap, err := manager.client.Products()
	if err != nil {
		log.Fatalf("Failed to get products.: %s", err)
	}

	// Match products from Vend with those from the provided CSV file.
	products := matchVendProduct(productMap, imagePayload)
	if err != nil {
		fmt.Printf("Error matching product from Vend to CSV input: %s", err)
	}

	fmt.Printf("\nGetting and posting images.\n")
	for _, product := range *products {
		// For each product match, first grab the image from the URL, then post that
		// image to the product on Vend.
		imagePath, err := image.Grab(product)
		if err != nil {
			log.Fatalf("Failed to get image for product.: %s", err)
		}
		vendapi.ImageUpload(manager.client.Token, manager.client.DomainPrefix, imagePath, product)
	}

	// Log closing timestamp.
	log.Printf("FIN\n")
}

func matchVendProduct(productMap *map[string]vend.Product,
	imagePayload *[]vendapi.UploadProduct) *[]vendapi.UploadProduct {

	var products []vendapi.UploadProduct

	// Loop through each product from the store, and add the ID field
	// to any product from the CSV file that matches.
	for _, product := range *productMap {
		for _, upload := range *imagePayload {

			// Ignore if any required values are empty.
			if product.SKU == nil || product.Handle == nil ||
				upload.SKU == nil || upload.Handle == nil {
				continue
			}

			// Ignore if product deleted.
			if product.DeletedAt != nil {
				continue
			}

			// Make sure we have a unique handle/sku match.
			if *product.SKU == *upload.SKU && *product.Handle == *upload.Handle {
				products = append(products, vendapi.UploadProduct{product.ID, upload.Handle, upload.SKU, upload.ImageURL})
				break
			}
		}
	}

	// Check how many matches we got.
	if len(products) > 0 {
		fmt.Printf("\nGot %d matches.\n", len(products))
	} else {
		fmt.Printf("\nNo product matches.\n")
		os.Exit(0)
	}

	return &products
}
