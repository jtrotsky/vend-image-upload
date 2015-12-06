// Package reader handles all input from CSV files.
package reader

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/jtrotsky/vend-image-upload/vendapi"
)

// ReadCSV a line of provided CSV file.
func ReadCSV(path string) (*[]vendapi.UploadProduct, error) {

	// SKU and Handle should be unique identifiers.
	exampleHeader := []string{"sku", "handle", "image_url"}

	// Open our provided CSV file.
	f, err := os.Open(path)
	if err != nil {
		log.Printf("Could not read from CSV file: %s", err)
		os.Exit(1)
	}
	// Make sure to close at end.
	defer f.Close()

	// Create CSV reader on our file.
	r := csv.NewReader(f)

	// Read and store our header line.
	headerRow, err := r.Read()

	// Check each string in the header row is same as our template.
	for i := range headerRow {
		if headerRow[i] != exampleHeader[i] {
			fmt.Println("Found error in header row. Check log.")
			log.Fatalf("No header match for: %s Instead got: %s.",
				string(exampleHeader[i]), string(headerRow[i]))
		}
	}

	// Read the rest of the data from the CSV.
	rawData, err := r.ReadAll()

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
