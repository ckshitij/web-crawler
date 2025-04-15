package crawler

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/ckshitij/web_crawler/pkg/queue"
	"golang.org/x/net/html"
)

type Crawler struct {
	baseURL     *url.URL
	visitedURLs sync.Map
	maxDepth    int
	treeMap     map[int][]EndpointResponse
	sync.Mutex
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
		visitedURLs: sync.Map{},
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

func (c *Crawler) CrawlWorkers(wg *sync.WaitGroup, ch chan *EndpointResponse) {
	for response := range ch {
		// No wg.Add here

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
				newResp.Depth = response.Depth + 1

				wg.Add(1) // ✅ Safe to add here
				ch <- newResp
			}
		}
		wg.Done()
	}
}

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

	wg.Add(1) // ✅ Add before sending to channel
	resCh <- rootSite

	wg.Wait()
	close(resCh) // ✅ After workers are done
	fmt.Println("Crawling completed.")
}

func (c *Crawler) PrintTreeMap() {
	for depth, responses := range c.treeMap {
		for _, response := range responses {
			fmt.Printf(" Depth: %d  URL: %s, Status Code: %d, Response Time: %s\n", depth, response.URL, response.StatusCode, response.ResponseTime)
		}
	}
}

func stripHostname(rawURL string) (string, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	// Combine Path + RawQuery + Fragment (if any)
	stripped := parsed.Path
	if parsed.RawQuery != "" {
		stripped += "?" + parsed.RawQuery
	}
	if parsed.Fragment != "" {
		stripped += "#" + parsed.Fragment
	}

	return stripped, nil
}

func (c *Crawler) PrintSiteMap() {
	// Print the treeMap in a structured way
	fmt.Println("Main Domain:", c.treeMap[0][0].URL)
	for i := 1; i < c.maxDepth; i++ {
		for _, response := range c.treeMap[i] {
			val, _ := stripHostname(response.URL)
			fmt.Printf("Level %d  URL: %s, Status Code: %d, Response Time: %s\n", i, val, response.StatusCode, response.ResponseTime)
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
