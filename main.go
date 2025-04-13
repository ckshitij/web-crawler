package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ckshitij/web_crawler/pkg/crawler"
	"github.com/ckshitij/web_crawler/pkg/utils"
)

func main() {
	// Define flags
	jsonOutput := flag.String("json", "", "Export site map as JSON file")
	xmlOutput := flag.String("xml", "", "Export site map as XML file")

	// Parse flags
	flag.Parse()

	// Get the positional argument (URL)
	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("Usage: go run main.go <url> [--json filename] [--xml filename]")
		os.Exit(1)
	}
	url := args[0]

	// Print the received values
	fmt.Println("Target URL:", url)

	if *jsonOutput != "" {
		utils.JSON.Export(*jsonOutput)
	}

	if *xmlOutput != "" {
		utils.XML.Export(*xmlOutput)
	}

	// Now you can call your crawl logic with `url`
	// and pass json/xml file paths if needed.
	crawl, err := crawler.NewCrawler(url, 2)
	if err != nil {
		fmt.Println("Error creating crawler:", err)
		return
	}
	// Call the crawl method
	crawl.CrawlSite()
	crawl.PrintTreeMap()
}
