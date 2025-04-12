package crawler

import (
	"fmt"
	"net/http"
	"time"
)

type EndpointResponse struct {
	URL          string
	StatusCode   int
	ResponseTime time.Duration
}

func GetURLInfo(url string) (*EndpointResponse, error) {
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
	return &EndpointResponse{
		URL:          url,
		StatusCode:   resp.StatusCode,
		ResponseTime: responseTime, // Example of calculating response time
	}, err
}
