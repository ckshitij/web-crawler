package crawler

import (
	"net/url"
	"sync"
	"time"
)

// Crawler is a structure that performs recursive website crawling.
type Crawler struct {
	baseURL     *url.URL                   // Root URL to start crawling from
	visitedURLs sync.Map                   // Thread-safe store for visited URLs
	maxDepth    int                        // Maximum depth to crawl
	treeMap     map[int][]EndpointResponse // Map of depth level to response data
	sync.Mutex                             // Embedded mutex for synchronizing access
}

// EndpointResponse represents metadata collected from a URL.
type EndpointResponse struct {
	URL          string
	StatusCode   int
	Depth        int
	Links        []string
	ResponseTime time.Duration
	Parent       string // Parent URL for hierarchical representation
}

type SiteMapNode struct {
	URL          string         `json:"url" xml:"url"`
	StatusCode   int            `json:"status_code" xml:"status_code"`
	Childerns    []*SiteMapNode `json:"children" xml:"children"`
	ResponseTime time.Duration  `json:"response_time" xml:"response_time"`
}
