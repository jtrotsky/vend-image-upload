// Package reader handles all input from CSV files.
package reader

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jtrotsky/vend-image-upload/logger"
	"github.com/jtrotsky/vend-image-upload/vendapi"
)

// ReadCSV reads the provided CSV file and stores the input as product objects.
func ReadCSV(productFilePath string, logFile *logger.LogFile) (*[]vendapi.ProductUpload, error) {
	// SKU and Handle combo should be a unique identifier.
	header := []string{"sku", "handle", "image_url"}

	// Open our provided CSV file.
	file, err := os.Open(productFilePath)
	if err != nil {
		log.Fatalf("Could not read from CSV file: %s", err)
		return &[]vendapi.ProductUpload{}, err
	}
	// Make sure to close at end.
	defer file.Close()

	// Create CSV reader on our file.
	reader := csv.NewReader(file)

	// Read and store our header line.
	headerRow, err := reader.Read()
	if err != nil {
		return &[]vendapi.ProductUpload{}, err
	}

	// Check each string in the header row is same as our template.
	for i, row := range headerRow {
		if strings.ToLower(row) != header[i] {
			log.Fatalf(
				"No header match for: %q. Instead got: %q. Expected: %s, %s, %s",
				header[i], strings.ToLower(row), header[0], header[1], header[2])
			return &[]vendapi.ProductUpload{}, err
		}
	}

	// Read the rest of the data from the CSV.
	rawData, err := reader.ReadAll()
	if err != nil {
		return &[]vendapi.ProductUpload{}, err
	}

	var product vendapi.ProductUpload
	var productList []vendapi.ProductUpload
	var rowNumber int

	// Loop through rows and assign them to product.
	for _, row := range rawData {
		rowNumber++
		product, err = readRow(row)
		if err != nil {
			// TODO: Can this be shortened?
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
			logFile.WriteEntry(
				logger.RowError{
					"read", rowNumber, productID, productSKU, productHandle, imageURL, err})
			log.Printf("Error reading row %d from CSV for product: %s. Error: %s",
				rowNumber, row, err)
			continue
		}

		// Append each product to our list.
		productList = append(productList, product)
	}

	// Check how many rows we successfully read and stored.
	if len(productList) > 0 {
		if len(productList) != len(rawData) {
			fmt.Printf("%d of %d rows successful, see error file for details.\n",
				len(productList), len(rawData))
		} else {
			fmt.Printf("%d rows successful.\n", len(productList))
		}
	} else {
		fmt.Printf("No valid products.\n")
		os.Exit(0)
	}

	return &productList, err
}

// Read a single row of a CSV file and check for errors.
func readRow(row []string) (vendapi.ProductUpload, error) {
	var product vendapi.ProductUpload

	product.SKU = &row[0]
	product.Handle = &row[1]
	product.ImageURL = &row[2]

	for i := range row {
		if len(row[i]) < 1 {
			err := errors.New("Missing field")
			return product, err
		}
	}
	return product, nil
}
