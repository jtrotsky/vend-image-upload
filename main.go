package main

import (
	"flag"
	"fmt"
	"os"
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
		fmt.Println(
			"Domain prefix not given. Expected like: store-name.vendhq.com")
		os.Exit(0)
	}
	if filePath == "" {
		fmt.Println(
			"Path to file not given. Expected like: ~/Documents/product-csv-file.csv")
		os.Exit(0)
	}
	if authToken == "" {
		fmt.Println(
			"Authentication token not given. Expected like: oe1R9xoQeJRUdyVkz6trbcf9GnUTBovJWKRSBCEf")
		os.Exit(0)
	}

	// To save people who write DomainPrefix.vendhq.com.
	// Split DomainPrefix on the "." period character then grab the first part.
	parts := strings.Split(domainPrefix, ".")
	domainPrefix = parts[0]
}
