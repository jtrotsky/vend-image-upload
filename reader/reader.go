// Package reader handles all input from CSV files.
package reader

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jtrotsky/vend-image-upload/vendapi"
)

// ReadCSV a line of provided CSV file.
func ReadCSV(path string) (*[]vendapi.UploadProduct, error) {

	// SKU and Handle should be unique identifiers.
	header := []string{"sku", "handle", "image_url"}

	// Open our provided CSV file.
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("Could not read from CSV file: %s", err)
		os.Exit(0)
	}
	// Make sure to close at end.
	defer file.Close()

	// Create CSV reader on our file.
	reader := csv.NewReader(file)

	// Read and store our header line.
	headerRow, err := reader.Read()

	// Check each string in the header row is same as our template.
	for i, row := range headerRow {
		if strings.ToLower(row) != header[i] {
			log.Fatalf(
				"No header match for: %q. Instead got: %q. Expected: %s, %s, %s",
				header[i], strings.ToLower(row), header[0], header[1], header[2])
		}
	}

	// Read the rest of the data from the CSV.
	rawData, err := reader.ReadAll()

	var product vendapi.UploadProduct
	var productList []vendapi.UploadProduct

	// Loop through rows and assign them to product.
	for _, row := range rawData {
		product.SKU = &row[0]
		product.Handle = &row[1]
		product.ImageURL = &row[2]

		// Append each product to our list.
		productList = append(productList, product)
	}

	return &productList, err
}
