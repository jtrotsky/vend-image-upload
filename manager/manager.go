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
	vend vend.Client
}

// NewManager creates an instance of manager.
func NewManager(vend vend.Client) *Manager {
	return &Manager{vend}
}

// Run execues the app's main functions.
func (manager *Manager) Run(f string) {
	// Using log gives us an opening timestamp.
	log.Printf("BEGIN\n")

	imagePayload, err := reader.ReadCSV(f)
	if err != nil {
		log.Fatalf("Error reading CSV file: %s", err)
	}

	// TODO:
	// fmt.Printf("%s", imagePayload)

	fmt.Printf("\nGrabbing products.\n")
	// Get all products, ignore productMap.
	_, _, err = manager.vend.Products()
	if err != nil {
		log.Fatalf("Failed to get products.: %s", err)
	}

	for _, i := range *imagePayload {
		image.Grab(i)
		fmt.Printf("\n\n%s\n", i)
	}
}
