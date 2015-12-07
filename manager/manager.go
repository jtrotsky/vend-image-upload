package manager

import (
	"fmt"
	"log"

	"github.com/jtrotsky/govend/vend"
	"github.com/jtrotsky/vend-image-upload/image"
	"github.com/jtrotsky/vend-image-upload/reader"
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

// Run reads the product CSV, gets all Vend products, then posts images.
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
	_, _, err = manager.client.Products()
	if err != nil {
		log.Fatalf("Failed to get products.: %s", err)
	}

	for _, i := range *imagePayload {
		image.Grab(i)
		fmt.Printf("\n\n%s\n", i)
	}
}
