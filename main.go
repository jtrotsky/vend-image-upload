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
	authToken       string
	domainPrefix    string
	logFilePath     string
	productFilePath string
)

func main() {
	// Calculate time for use in log filename.
	currentTimeUTC := time.Now()
	// Set log filepath.
	logFilePath = fmt.Sprintf("./%d_vendimageupload_errors.csv", currentTimeUTC.Unix())

	// Start CSV logfile in current directory with unix timestamp in file name.
	logFile := logger.NewLogFile(logFilePath)
	logFile.CreateLog()

	// Invoke new Vend client.
	// Timezone argument left blank as unused.
	vendClient := vend.NewClient(authToken, domainPrefix, "")
	manager := manager.NewManager(vendClient)

	manager.Run(productFilePath, logFile)
}

func init() {
	// Set default log output to terminal.
	log.SetOutput(os.Stdout)

	// Get store info from command line flags.
	flag.StringVar(&authToken, "t", "", "API Authentication Token.")
	flag.StringVar(&domainPrefix, "d", "", "Vend store name.")
	flag.StringVar(&productFilePath, "f", "", "Path to product CSV file.")
	flag.Parse()

	// Check all arguments are given.
	if authToken == "" {
		log.Println(
			"Authentication token not given. Expected like: oe1R9xoQeJRUdyVkz6trbcf9GnUTBovJWKRSBCEf")
		os.Exit(0)
	}
	if domainPrefix == "" {
		log.Println(
			"Domain prefix not given. Expected like: store-name.vendhq.com")
		os.Exit(0)
	}
	if productFilePath == "" {
		log.Println(
			"Path to file not given. Expected like: /Users/<user>/Documents/product-file.csv")
		os.Exit(0)
	}

	// To save people who write DomainPrefix.vendhq.com.
	// Split DomainPrefix on the "." period character then grab the first part.
	parts := strings.Split(domainPrefix, ".")
	domainPrefix = parts[0]
}
