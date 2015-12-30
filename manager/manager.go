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
func (manager *Manager) Run(productFilePath, logFilePath string) {
	// Log opening timestamp.
	log.Printf("BEGIN\n")

	fmt.Printf("\nReading products from CSV file.\n")
	// Read provided CSV file and store product info.
	productsFromCSV, err := reader.ReadCSV(productFilePath, logFilePath)
	if err != nil {
		log.Fatalf("Error reading CSV file: %s", err)
		os.Exit(0)
	}

	fmt.Printf("\nGrabbing products from Vend.\n")
	// Get all products from Vend.
	_, productsFromVend, err := manager.client.Products()
	if err != nil {
		log.Fatalf("Failed to get products.: %s", err)
		os.Exit(0)
	}

	fmt.Printf("\nLooking for product matches.\n")
	// Match products from Vend with those from the provided CSV file.
	matchedProducts := matchVendProduct(productsFromVend, productsFromCSV)
	if err != nil {
		fmt.Printf("Error matching product from Vend to CSV input: %s", err)
	}

	fmt.Printf("\nGetting and posting images.\n")
	for _, product := range *matchedProducts {
		// For each product match, first grab the image from the URL, then post that
		// image to the product on Vend.
		imagePath, err := image.Grab(product)
		if err != nil {
			fmt.Printf("<<FAILURE>> Ignoring product %s %s.\nFailed to get image: %s\n\n",
				*product.Handle, *product.SKU, err)
			// Ignore product if image grabbing errored.
			continue
		}
		vendapi.ImageUpload(manager.client.Token, manager.client.DomainPrefix,
			imagePath, product)
	}

	// Log closing timestamp.
	log.Printf("FIN\n")
}

func matchVendProduct(productMap *map[string]vend.Product,
	imagePayload *[]vendapi.ProductUpload) *[]vendapi.ProductUpload {

	var products []vendapi.ProductUpload

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
				products = append(products, vendapi.ProductUpload{product.ID, upload.Handle, upload.SKU, upload.ImageURL})
				break
			}
		}
	}

	// Check how many matches we got.
	if len(products) > 0 {
		fmt.Printf("%d products matched.\n", len(products))
	} else {
		fmt.Printf("No product matches.\n")
		os.Exit(0)
	}

	return &products
}
