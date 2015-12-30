package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jtrotsky/govend/vend"
	"github.com/jtrotsky/vend-image-upload/logger"
	"github.com/jtrotsky/vend-image-upload/manager"
)

var (
	domainPrefix    string
	productFilePath string
	logFilePath     string
	authToken       string
)

func main() {

	// Invoke new Vend client.
	// Timezone argument left blank as unused.
	vendClient := vend.NewClient(authToken, domainPrefix, "")
	manager := manager.NewManager(vendClient)

	manager.Run(productFilePath, logFilePath)
}

func init() {
	// Set default log output to terminal.
	log.SetOutput(os.Stdout)

	// Calculate time for use in log filename.
	currentTime := time.Now()
	// Set log filepath.
	logFilePath = fmt.Sprintf("./%d_vendimageupload_error.csv", currentTime.Unix())

	// Start CSV logfile in current directory with unix timestamp in file name.
	logFile := logger.NewLogFile(logFilePath)
	logFile.CreateLog()

	// Get store info from command line flags.
	flag.StringVar(&domainPrefix, "d", "", "Vend store name.")
	flag.StringVar(&productFilePath, "f", "", "Path to product CSV file.")
	flag.StringVar(&authToken, "t", "", "API Authentication Token.")
	flag.Parse()

	// Check all arguments are given.
	if domainPrefix == "" {
		log.Println(
			"Domain prefix not given. Expected like: store-name.vendhq.com")
		os.Exit(0)
	}
	if productFilePath == "" {
		log.Println(
			"Path to file not given. Expected like: ~/Documents/product-csv-file.csv")
		os.Exit(0)
	}
	if authToken == "" {
		log.Println(
			"Authentication token not given. Expected like: oe1R9xoQeJRUdyVkz6trbcf9GnUTBovJWKRSBCEf")
		os.Exit(0)
	}

	// To save people who write DomainPrefix.vendhq.com.
	// Split DomainPrefix on the "." period character then grab the first part.
	parts := strings.Split(domainPrefix, ".")
	domainPrefix = parts[0]
}
