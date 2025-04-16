// Package crawler provides functionality to crawl websites starting from a base URL up to a specified depth.
// It tracks visited URLs, collects metadata such as status codes and response times, and builds a map of discovered links.
package crawler

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ckshitij/web_crawler/pkg/queue"
	"golang.org/x/net/html"
)

// NewCrawler initializes and returns a new Crawler instance.
func NewCrawler(baseURL string, maxDepth int) (*Crawler, error) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}
	crawler := &Crawler{
		baseURL:     parsedURL,
		visitedURLs: sync.Map{},
		maxDepth:    maxDepth,
		treeMap:     make(map[int][]EndpointResponse),
	}
	return crawler, nil
}

// Crawls performs breadth-first crawling using a custom queue.
func (c *Crawler) Crawls(newURL EndpointResponse) {
	queue := queue.NewQueue[*EndpointResponse]()
	queue.Enqueue(&newURL)
	for !queue.IsEmpty() {
		currentResponse, err := queue.Front()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		if currentResponse.Depth > c.maxDepth {
			break
		}

		c.Lock()
		c.treeMap[currentResponse.Depth] = append(c.treeMap[currentResponse.Depth], *currentResponse)
		c.Unlock()
		for _, link := range currentResponse.Links {
			if _, ok := c.visitedURLs.Load(link); !ok {
				c.visitedURLs.Store(link, true)
				newResponse, err := c.getURLInfo(link)
				if err != nil {
					fmt.Println("Error:", err)
					continue
				}
				newResponse.Depth = currentResponse.Depth + 1
				queue.Enqueue(newResponse)
			}
		}

		queue.Dequeue()
	}
}

// CrawlWorkers processes URL responses concurrently from a channel.
func (c *Crawler) CrawlWorkers(wg *sync.WaitGroup, ch chan *EndpointResponse) {
	for response := range ch {
		c.Lock()
		c.treeMap[response.Depth] = append(c.treeMap[response.Depth], *response)
		c.Unlock()

		if response.Depth >= c.maxDepth {
			wg.Done()
			continue
		}

		for _, link := range response.Links {
			if _, ok := c.visitedURLs.Load(link); !ok {
				c.visitedURLs.Store(link, true)
				newResp, err := c.getURLInfo(link)
				if err != nil {
					continue
				}
				newResp.Parent = response.URL
				newResp.Depth = response.Depth + 1
				wg.Add(1)
				ch <- newResp
			}
		}
		wg.Done()
	}
}

// CrawlSite begins the site crawl using multiple workers.
func (c *Crawler) CrawlSite() {
	rootSite, err := c.getURLInfo(c.baseURL.String())
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	rootSite.Depth = 0
	c.visitedURLs.Store(rootSite.URL, true)

	var wg sync.WaitGroup
	resCh := make(chan *EndpointResponse, 1000)

	for range 10 {
		go c.CrawlWorkers(&wg, resCh)
	}

	wg.Add(1)
	resCh <- rootSite

	wg.Wait()
	close(resCh)
	fmt.Println("Crawling completed.")
}

// PrintSiteMap prints the sitemap in a hierarchical format based on parent-child relationships.
func (c *Crawler) PrintSiteMap() {
	// Build parent -> children map
	parentMap := make(map[string][]EndpointResponse)
	for i, responses := range c.treeMap {
		if i == 0 {
			continue // Skip the root level
		}
		for _, resp := range responses {
			parentMap[resp.Parent] = append(parentMap[resp.Parent], resp)
		}
	}

	root := c.treeMap[0][0] // Assume the root is the only entry at depth 0
	fmt.Printf("Main Domain: %s\n", root.URL)
	printChildren(parentMap, root.URL, 1)
}

// printChildren recursively prints children of a parent URL with indentation.
func printChildren(parentMap map[string][]EndpointResponse, parent string, level int) {
	children, exists := parentMap[parent]
	if !exists {
		return
	}
	for _, child := range children {
		val, _ := stripHostname(child.URL)
		fmt.Printf("%s├── %s\n", strings.Repeat("  ", level), val)
		printChildren(parentMap, child.URL, level+1)
	}
}

// stripHostname removes the scheme and hostname from a URL, leaving the path, query, and fragment.
func stripHostname(rawURL string) (string, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	stripped := parsed.Path
	if parsed.RawQuery != "" {
		stripped += "?" + parsed.RawQuery
	}
	if parsed.Fragment != "" {
		stripped += "#" + parsed.Fragment
	}
	return stripped, nil
}

// processToken extracts anchor links from HTML tokens.
func (r *Crawler) processToken(token *html.Tokenizer) []string {
	links := make([]string, 0)
	for {
		tt := token.Next()
		switch tt {
		case html.ErrorToken:
			return links
		case html.StartTagToken, html.SelfClosingTagToken:
			t := token.Token()
			if t.Data == "a" {
				for _, attr := range t.Attr {
					if attr.Key == "href" {
						link, err := r.baseURL.Parse(attr.Val)
						if err == nil && link.Hostname() == r.baseURL.Hostname() && len(links) < 4 {
							links = append(links, link.String())
						}
					}
				}
			}
		}
	}
}

// getURLInfo performs an HTTP GET request and parses basic metadata and links.
func (r *Crawler) getURLInfo(url string) (*EndpointResponse, error) {
	startTime := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching URL:", err)
		return nil, err
	}
	defer resp.Body.Close()
	responseTime := time.Since(startTime)

	response := EndpointResponse{
		URL:          url,
		StatusCode:   resp.StatusCode,
		ResponseTime: responseTime,
	}
	if resp.StatusCode != http.StatusOK {
		return &response, nil
	}
	tokens := html.NewTokenizer(resp.Body)
	links := r.processToken(tokens)
	response.Links = links

	return &response, err
}

func (c *Crawler) ExportSiteMapJSON(filename string) error {
	root := c.buildSiteTree()
	data, err := json.MarshalIndent(root, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func (c *Crawler) ExportSiteMapXML(filename string) error {
	root := c.buildSiteTree()
	data, err := xml.MarshalIndent(root, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func (c *Crawler) buildSiteTree() SiteMapNode {
	var parentMap = map[string][]EndpointResponse{}
	for depth, nodes := range c.treeMap {
		if depth == 0 {
			continue
		}
		for _, node := range nodes {
			parentMap[node.Parent] = append(parentMap[node.Parent], node)
		}
	}

	root := c.treeMap[0][0]
	return buildSiteNode(root, parentMap)
}

func buildSiteNode(node EndpointResponse, parentMap map[string][]EndpointResponse) SiteMapNode {
	SiteMapNode := SiteMapNode{
		ResponseTime: node.ResponseTime,
		StatusCode:   node.StatusCode,
		URL:          node.URL,
		Childerns:    []*SiteMapNode{},
	}
	childers := parentMap[node.URL]
	for _, child := range childers {
		childNode := buildSiteNode(child, parentMap)
		SiteMapNode.Childerns = append(SiteMapNode.Childerns, &childNode)
	}
	return SiteMapNode
}
