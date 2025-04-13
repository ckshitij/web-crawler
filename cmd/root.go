package cmd

import (
	"fmt"
	"os"

	"github.com/ckshitij/web_crawler/pkg/crawler"
	"github.com/ckshitij/web_crawler/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	jsonOutput string
	xmlOutput  string
	maxDepth   int
)

var rootCmd = &cobra.Command{
	Use:   "web_crawler [url]",
	Short: "A simple website crawler that builds a site map",
	Long:  `Crawls a website and generates a site map, with optional export to JSON or XML.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]

		crawl, err := crawler.NewCrawler(url, maxDepth)
		if err != nil {
			fmt.Println("Error initializing crawler:", err)
			return
		}
		crawl.CrawlSite()
		crawl.PrintTreeMap()

		if jsonOutput != "" {
			utils.JSON.Export(jsonOutput)
		}

		if xmlOutput != "" {
			utils.XML.Export(xmlOutput)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVar(&jsonOutput, "json", "", "Export site map as JSON file")
	rootCmd.Flags().StringVar(&xmlOutput, "xml", "", "Export site map as XML file")
	rootCmd.Flags().IntVar(&maxDepth, "max_depth", 2, "Maximum depth for crawling")
}
