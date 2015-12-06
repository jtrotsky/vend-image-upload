package main

import (
	"flag"
	"strings"

	"github.com/jtrotsky/govend/vend"
	"github.com/jtrotsky/vend-image-upload/manager"
)

var (
	domainPrefix string
	path         string
	token        string
)

func main() {

	// Invoke new Vend client. Timezone argument left blank.
	v := vend.NewClient(token, domainPrefix, "")
	manager := manager.NewManager(v)

	manager.Run(path)
}

func init() {

	// Get store info from command line flags.
	flag.StringVar(&domainPrefix, "d", "",
		"The Vend store name (prefix of xxxx.vendhq.com)")
	flag.StringVar(&path, "f", "",
		"Path to product CSV file that contains image URLs.")
	flag.StringVar(&token, "t", "",
		"Personal API Access Token for the store, generated from Setup -> API Access.")
	flag.Parse()

	// To save people who write DomainPrefix.vendhq.com.
	// Split DomainPrefix on the "." period character then grab the first part.
	parts := strings.Split(domainPrefix, ".")
	domainPrefix = parts[0]
}
