package manager

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/jtrotsky/govend/vend"
	"github.com/jtrotsky/vend-image-upload/image"
	"github.com/jtrotsky/vend-image-upload/logger"
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
func (manager *Manager) Run(productFilePath string, logFile *logger.LogFile) {
	// Log opening timestamp.
	log.Printf("BEGIN\n")

	fmt.Printf("\nReading products from CSV file.\n")
	// Read provided CSV file and store product info.
	productsFromCSV, err := reader.ReadCSV(productFilePath, logFile)
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
	matchedProducts := matchVendProduct(productsFromVend, productsFromCSV, logFile.FilePath)
	if err != nil {
		fmt.Printf("Error matching product from Vend to CSV input: %s", err)
	}

	fmt.Printf("\nGetting and posting images.\n")
	for _, product := range *matchedProducts {
		// For each product match, first grab the image from the URL, then post that
		// image to the product on Vend.
		imagePath, err := image.Grab(product, logFile.FilePath)
		if err != nil {
			var productID, productSKU, productHandle, imageURL string
			if product.ID != nil {
				productID = *product.ID
			} else {
				productID = ""
			}
			if product.SKU != nil {
				productSKU = *product.SKU
			} else {
				productSKU = ""
			}
			if product.Handle != nil {
				productHandle = *product.Handle
			} else {
				productHandle = ""
			}
			if product.ImageURL != nil {
				imageURL = *product.ImageURL
			} else {
				imageURL = ""
			}
			logFile := logger.NewLogFile(logFile.FilePath)
			logFile.WriteEntry(logger.RowError{
				"network", 0, productID, productSKU, productHandle, imageURL, err})
			fmt.Printf("<<FAILURE>> Ignoring product %s %s.\n\n",
				*product.Handle, *product.SKU)
			// Ignore product if image grabbing errored.
			continue
		}
		vendapi.ImageUpload(manager.client.Token, manager.client.DomainPrefix,
			imagePath, product)
	}

	// If no errors were recorded then remove the error file.
	if logFile.ErrorCount == 0 {
		os.Remove(logFile.FilePath)
	}

	// Log closing timestamp.
	log.Printf("FIN\n")
}

func matchVendProduct(productsFromVend *map[string]vend.Product,
	productsFromCSV *[]vendapi.ProductUpload, logFilePath string) *[]vendapi.ProductUpload {

	var products []vendapi.ProductUpload

	// Loop through each product from the store, and add the ID field
	// to any product from the CSV file that matches.
Match:
	for _, csvProduct := range *productsFromCSV {
		for _, vendProduct := range *productsFromVend {
			// Ignore if any required values are empty.
			if vendProduct.SKU == nil || vendProduct.Handle == nil ||
				csvProduct.SKU == nil || csvProduct.Handle == nil {
				continue
			}
			// Ignore if product deleted.
			if vendProduct.DeletedAt != nil {
				continue
			}
			// Make sure we have a unique handle/sku match.
			if *vendProduct.SKU == *csvProduct.SKU && *vendProduct.Handle == *csvProduct.Handle {
				products = append(products, vendapi.ProductUpload{vendProduct.ID, csvProduct.Handle, csvProduct.SKU, csvProduct.ImageURL})
				continue Match
			}
		}
		// Record product from CSV if no match to Vend products.
		var productSKU, productHandle, imageURL string
		if csvProduct.SKU != nil {
			productSKU = *csvProduct.SKU
		} else {
			productSKU = ""
		}
		if csvProduct.Handle != nil {
			productHandle = *csvProduct.Handle
		} else {
			productHandle = ""
		}
		if csvProduct.ImageURL != nil {
			imageURL = *csvProduct.ImageURL
		} else {
			imageURL = ""
		}
		err := errors.New("No handle/sku match")
		logFile := logger.NewLogFile(logFilePath)
		logFile.WriteEntry(
			logger.RowError{
				"match", 0, "", productSKU, productHandle, imageURL, err})
	}

	// Check how many matches we got.
	if len(products) > 0 {
		fmt.Printf("%d of %d products matched.\n", len(products), len(*productsFromCSV))
	} else {
		fmt.Printf("No product matches.\n")
		os.Exit(0)
	}

	return &products
}
