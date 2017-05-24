// Package reader handles all input from CSV files.
package reader

import (
	"encoding/csv"
	"errors"
	"os"
	"strings"

	"github.com/jtrotsky/vend-image-upload/vendapi"
	log "github.com/sirupsen/logrus"
	"github.com/wallclockbuilder/stringutil"
)

// ReadCSV reads the provided CSV file and stores the input as product objects.
func ReadCSV(productFilePath string) (*[]vendapi.ProductUpload, error) {
	// SKU and Handle combo should be a unique identifier.
	header := []string{"sku", "handle", "image_url"}

	// Open our provided CSV file.
	file, err := os.Open(productFilePath)
	if err != nil {
		log.WithError(err).WithField("filepath", productFilePath).Error("Could not read from CSV file")
		return &[]vendapi.ProductUpload{}, err
	}
	// Make sure to close at end.
	defer file.Close()

	// Create CSV reader on our file.
	reader := csv.NewReader(file)

	// Read and store our header line.
	headerRow, err := reader.Read()
	if err != nil {
		log.Error("Failed to read headerow.")
		return &[]vendapi.ProductUpload{}, err
	}

	if len(headerRow) > 3 {
		log.WithField("header_row_error", headerRow).Warning("Header row longer than expected")
	}

	// Check each string in the header row is same as our template.
	for i, row := range headerRow {
		if stringutil.Strip(strings.ToLower(row)) != header[i] {
			log.WithFields(log.Fields{
				"type":     "header_mismatch",
				"expected": header[i],
				"got":      strings.ToLower(row),
			}).Error("Header mismatch")
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
			log.WithError(err).WithFields(log.Fields{
				"type":          "row_read",
				"row":           rowNumber,
				"product_id":    product.ID,
				"product_sku":   product.SKU,
				"produt_handle": product.Handle,
				"image_url":     product.ImageURL,
				"reason":        err,
			}).Error("Error reading row from CSV")
			continue
		}

		// Append each product to our list.
		productList = append(productList, product)
	}

	// Check how many rows we successfully read and stored.
	if len(productList) > 0 {
		log.WithFields(log.Fields{
			"successful": len(productList),
			"total":      len(rawData),
		}).Info("Products read from file")
	} else {
		log.Fatal("No valid products found")
	}

	return &productList, err
}

// Read a single row of a CSV file and check for errors.
func readRow(row []string) (vendapi.ProductUpload, error) {
	var product vendapi.ProductUpload

	product.SKU = row[0]
	product.Handle = row[1]
	product.ImageURL = row[2]

	for i := range row {
		if len(row[i]) < 1 {
			err := errors.New("Missing field")
			return product, err
		}
	}
	return product, nil
}
