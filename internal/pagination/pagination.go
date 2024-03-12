// Package pagination deals with pagination of an API.
package pagination

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/codeYann/go-pagination/internal/requests"
)

// Options struct defines the pagination options.
type Options struct {
	MaxRetries        int
	RetryTimeout      int
	MaxRequestTimeout int
	Threshold         int
}

// DefaultOptions defines the default pagination options.
var DefaultOptions = Options{
	MaxRetries:        4,
	RetryTimeout:      1000,
	MaxRequestTimeout: 1000,
	Threshold:         200,
}

// Pager interface defines the pagination methods.
type Pager interface {
	GetPaginated(url string, page int) ([]byte, error)
	HandleRequest(url string, page, retries int) ([]byte, error)
}

// Pagination struct holds the pagination options.
type Pagination struct {
	Options
	request *requests.HTTPRequest
}

// NewPagination is a build function that creates a new Pagination instance with the provided options.
func NewPagination(request *requests.HTTPRequest, options ...Options) *Pagination {
	pagination := &Pagination{
		Options: DefaultOptions,
		request: request,
	}

	applyOption := func(target *int, source int) {
		if source != 0 {
			*target = source
		}
	}

	for _, opt := range options {
		applyOption(&pagination.MaxRetries, opt.MaxRetries)
		applyOption(&pagination.RetryTimeout, opt.RetryTimeout)
		applyOption(&pagination.MaxRequestTimeout, opt.MaxRequestTimeout)
		applyOption(&pagination.Threshold, opt.Threshold)
	}

	return pagination
}

// HandleRequest handles a request.
func (p Pagination) HandleRequest(url string, page, retries int) ([]byte, error) {
	finalURL := fmt.Sprintf("%s?tid=%d", url, page)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(p.MaxRequestTimeout)*time.Millisecond,
	)
	defer cancel()

	result, err := p.request.MakeRequest(ctx, finalURL)
	if err != nil {
		if retries == p.MaxRetries {
			fmt.Println("Max retries reached")
			return nil, err
		}
		return p.HandleRequest(url, page, retries+1)
	}

	return result, nil
}

// GetPaginated makes a pagination for a certain request
func (p Pagination) GetPaginated(url string, page int) ([]byte, error) {
	result, err := p.HandleRequest(url, page, 1)
	if err != nil {
		return nil, err
	}

	// JSON parsing result
	var data []struct {
		Amount float32 `json:"amount"`
		Date   int64   `json:"date"`
		Price  float64 `json:"price"`
		Tid    int64   `json:"tid"`
		Type   string  `json:"type"`
	}

	// Unmarshal the result
	decoder := json.NewDecoder(bytes.NewReader(result))
	for decoder.More() {
		if err := decoder.Decode(&data); err != nil {
			return nil, err
		}
	}

	lastID := data[len(data)-1].Tid
	if lastID == 0 {
		return nil, nil
	}

	time.Sleep(time.Duration(p.Threshold) * time.Millisecond)

	return p.GetPaginated(url, int(lastID))
}
