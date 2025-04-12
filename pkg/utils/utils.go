package utils

import "fmt"

type OutputFormatType string

const (
	JSON OutputFormatType = "json"
	XML  OutputFormatType = "xml"
)

func (j OutputFormatType) Export(filename string) {
	switch j {
	case JSON:
		fmt.Println("Exporting to JSON file:", filename)
		// Implement JSON export logic here
	case XML:
		fmt.Println("Exporting to XML file:", filename)
		// Implement XML export logic here
	default:
		fmt.Println("Unknown output type")
	}
}
