package main

import (
	"flag"
	"log"
	"strings"

	"github.com/jtrotsky/govend/vend"
	"github.com/jtrotsky/vend-image-upload/manager"
)

var (
	domainPrefix string
	filePath     string
	authToken    string
)

func main() {

	// Invoke new Vend client.
	// Timezone argument left blank as unused.
	vendClient := vend.NewClient(authToken, domainPrefix, "")
	manager := manager.NewManager(vendClient)

	manager.Run(filePath)
}

func init() {

	// Get store info from command line flags.
	flag.StringVar(&domainPrefix, "d", "", "Vend store name.")
	flag.StringVar(&filePath, "f", "", "Path to product CSV file.")
	flag.StringVar(&authToken, "t", "", "API Authentication Token.")
	flag.Parse()

	// Check all arguments are given.
	if domainPrefix == "" {
		log.Fatalf("Domain prefix not given.")
	}
	if filePath == "" {
		log.Fatalf("Path to file not given.")
	}
	if authToken == "" {
		log.Fatalf("Authentication token not given.")
	}

	// To save people who write DomainPrefix.vendhq.com.
	// Split DomainPrefix on the "." period character then grab the first part.
	parts := strings.Split(domainPrefix, ".")
	domainPrefix = parts[0]
}
