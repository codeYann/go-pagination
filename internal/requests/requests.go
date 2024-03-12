// Package requests is a package that contains the request handlers for the application
package requests

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTPRequest is a struct which deals with request operations.
type HTTPRequest struct {
	BaseURL    string
	Method     string
	MaxTimeout time.Duration
}

// request makes an HTTP request.
func (r HTTPRequest) request(url string) ([]byte, error) {
	finalURL := fmt.Sprintf("%s%s", r.BaseURL, url)
	req, err := http.NewRequest(r.Method, finalURL, nil)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), r.MaxTimeout)
	defer cancel()

	req = req.WithContext(ctx)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// MakeRequest makes a request
func (r HTTPRequest) MakeRequest(ctx context.Context, url string) ([]byte, error) {
	data, err := r.request(url)
	if err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return data, nil
	}
}
