package manager

import (
	"errors"

	"github.com/jtrotsky/govend/vend"
	"github.com/jtrotsky/vend-image-upload/image"
	"github.com/jtrotsky/vend-image-upload/reader"
	"github.com/jtrotsky/vend-image-upload/vendapi"
	log "github.com/sirupsen/logrus"
)

// Manager contains the Vend client.
type Manager struct {
	client vend.Client
}

// NewManager creates an instance of manager.
func NewManager(client vend.Client) *Manager {
	return &Manager{client}
}

// Run reads the product CSV, gets all products from Vend, then posts their images.
func (manager *Manager) Run(productFilePath string) {
	// Log opening timestamp.
	log.Info("BEGIN")

	log.Info("Reading products from CSV file")

	// Read provided CSV file and store product info.
	productsFromCSV, err := reader.ReadCSV(productFilePath)
	if err != nil {
		log.WithError(err).Fatal("Error reading CSV file")
	}

	log.Info("Grabbing products from Vend")

	// Get all products from Vend.
	_, productsFromVend, err := manager.client.Products()
	if err != nil {
		log.WithError(err).Fatal("Failed to get products")
	}

	log.Info("Looking for product matches")

	// Match products from Vend with those from the provided CSV file.
	matchedProducts := matchVendProduct(productsFromVend, productsFromCSV)
	if err != nil {
		log.WithError(err).Error("Error matching product from Vend to CSV input")
	}

	log.Info("Getting and posting images")

	for _, product := range *matchedProducts {
		// For each product match, first grab the image from the URL, then post that
		// image to the product on Vend.
		imagePath, err := image.Grab(product)
		if err != nil {
			log.WithError(err).WithFields(log.Fields{
				"type":           "network",
				"bla":            0,
				"product_id":     product.ID,
				"product_sku":    product.SKU,
				"product_handle": product.Handle,
				"image_url":      product.ImageURL,
			})
			log.WithFields(log.Fields{
				"handle": product.Handle,
				"sku":    product.SKU,
			}).Warning("Ignoring product")
			// Ignore product if image grabbing errored.
			continue
		}
		vendapi.UploadImage(manager.client.Token, manager.client.DomainPrefix,
			imagePath, product)
	}

	// Log closing timestamp.
	log.Info("FINISHED")
}

func matchVendProduct(productsFromVend *map[string]vend.Product, productsFromCSV *[]vendapi.ProductUpload) *[]vendapi.ProductUpload {

	var products []vendapi.ProductUpload

	// Loop through each product from the store, and add the ID field
	// to any product from the CSV file that matches.
Match:
	for _, csvProduct := range *productsFromCSV {
		for _, vendProduct := range *productsFromVend {
			// Ignore if any required values are empty.
			if vendProduct.SKU == nil || vendProduct.Handle == nil ||
				csvProduct.SKU == "" || csvProduct.Handle == "" {
				continue
			}
			// Ignore if product deleted.
			if vendProduct.DeletedAt != nil {
				continue
			}
			// Make sure we have a unique handle/sku match, then add product to list.
			if *vendProduct.SKU == csvProduct.SKU &&
				*vendProduct.Handle == csvProduct.Handle {
				products = append(products,
					vendapi.ProductUpload{
						ID:       *vendProduct.ID,
						Handle:   csvProduct.Handle,
						SKU:      csvProduct.SKU,
						ImageURL: csvProduct.ImageURL,
					})
				continue Match
			}
		}
		// Record product from CSV as error if no match to Vend products.
		err := errors.New("No handle/sku match")
		log.WithError(err).WithFields(log.Fields{
			"type":                  "match",
			"csv_product_sku":       csvProduct.SKU,
			"csv_product_handle":    csvProduct.Handle,
			"csv_product_image_url": csvProduct.ImageURL,
		})
	}

	// Check how many matches we got.
	if len(products) > 0 {
		log.WithFields(log.Fields{
			"matched": len(products),
			"total":   len(*productsFromCSV),
		}).Info("Products matched")
	} else {
		log.Error("No product matches")
		return nil
	}
	return &products
}
