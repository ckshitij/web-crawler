package crawler

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/ckshitij/web_crawler/pkg/queue"
	"golang.org/x/net/html"
)

type Crawler struct {
	baseURL     *url.URL
	visitedURLs map[string]bool
	maxDepth    int
	treeMap     map[int][]EndpointResponse
}

type EndpointResponse struct {
	URL          string
	StatusCode   int
	Depth        int
	Links        []string
	ResponseTime time.Duration
}

func NewCrawler(baseURL string, maxDepth int) (*Crawler, error) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}
	crawler := &Crawler{
		baseURL:     parsedURL,
		visitedURLs: make(map[string]bool),
		maxDepth:    maxDepth,
		treeMap:     make(map[int][]EndpointResponse),
	}
	return crawler, nil
}

// CrawlSite starts the crawling process from the base URL
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

		c.treeMap[currentResponse.Depth] = append(c.treeMap[currentResponse.Depth], *currentResponse)
		for _, link := range currentResponse.Links {
			if _, ok := c.visitedURLs[link]; !ok {
				c.visitedURLs[link] = true
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

func (c *Crawler) CrawlSite() {

	rootSite, err := c.getURLInfo(c.baseURL.String())
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	rootSite.Depth = 0
	c.visitedURLs[rootSite.URL] = true
	fmt.Printf("Crawled URL: %s, Status Code: %d, Response Time: %s, Links %+v\n", rootSite.URL, rootSite.StatusCode, rootSite.ResponseTime, rootSite.Links)
	c.Crawls(*rootSite)
	fmt.Println("Crawling completed.")
}

func (c *Crawler) PrintTreeMap() {
	for depth, responses := range c.treeMap {
		for _, response := range responses {
			fmt.Printf(" Depth: %d  URL: %s, Status Code: %d, Response Time: %s\n", depth, response.URL, response.StatusCode, response.ResponseTime)
		}
	}
}

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

func (r *Crawler) getURLInfo(url string) (*EndpointResponse, error) {
	// Simulate fetching URL info
	// In a real scenario, you would make an HTTP request and get the response
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
		ResponseTime: responseTime, // Example of calculating response time
	}
	if resp.StatusCode != http.StatusOK {
		return &response, nil
	}
	tokens := html.NewTokenizer(resp.Body)
	links := r.processToken(tokens)
	response.Links = links

	return &response, err
}
