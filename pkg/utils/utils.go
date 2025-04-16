package utils

import (
	"fmt"

	"github.com/ckshitij/web_crawler/pkg/crawler"
)

type OutputFormatType string

const (
	JSON OutputFormatType = "json"
	XML  OutputFormatType = "xml"
)

func (j OutputFormatType) Export(filename string, c *crawler.Crawler) {
	switch j {
	case JSON:
		fmt.Println("Exporting to JSON file:", filename)
		c.ExportSiteMapJSON(filename)
		// Implement JSON export logic here
	case XML:
		fmt.Println("Exporting to XML file:", filename)
		c.ExportSiteMapXML(filename)
		// Implement XML export logic here
	default:
		fmt.Println("Unknown output type")
	}
}
