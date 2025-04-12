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

	fmt.Println("Crawling the site...", url)
	data, err := crawler.GetURLInfo(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Crawled URL: %s, Status Code: %d, Response Time: %s\n", data.URL, data.StatusCode, data.ResponseTime)
}
