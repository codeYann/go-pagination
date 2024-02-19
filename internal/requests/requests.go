// Package requests is a package that contains the request handlers for the application
package requests

import (
	"fmt"
	"io"
	"net/http"
)

// HTTPRequest is a struct which deals with request operations.
type HTTPRequest struct {
	BaseURL string
	Method  string
}

// NewHTTPRequest is factory builder.
func NewHTTPRequest(baseURL string, method string) *HTTPRequest {
	return &HTTPRequest{BaseURL: baseURL, Method: method}
}

// GetBaseURL returns base url
func (r HTTPRequest) GetBaseURL() string {
	return r.BaseURL
}

// GetMethod returns a method used in HTTPRequest
func (r HTTPRequest) GetMethod() string {
	return r.Method
}

func (r HTTPRequest) setURL(url string) (string, error) {
	if r.GetBaseURL() == "" {
		return "", fmt.Errorf("base url not set")
	}

	if url == "" {
		return "", fmt.Errorf("url not set")
	}

	// if  baseURL does not end with a slash and url does not start with a slash
	if r.GetBaseURL()[len(r.GetBaseURL())-1:] != "/" && url[:1] != "/" {
		return fmt.Sprintf("%s/%s", r.GetBaseURL(), url), nil
	}

	return fmt.Sprintf("%s%s", r.GetBaseURL(), url), nil
}

// Request makes an HTTP request.
func (r HTTPRequest) Request(url string) ([]byte, error) {
	finalURL, err := r.setURL(url)
	if err != nil {
		return nil, err
	}

	client := http.Client{}
	response := make(chan []byte)

	go func() {
		request, err := client.Get(finalURL)
		if err != nil {
			panic(err)
		}

		data, err := io.ReadAll(request.Body)
		if err != nil {
			panic(err)
		}

		response <- data
		request.Body.Close()
	}()

	data := <-response
	close(response)

	return data, nil
}
