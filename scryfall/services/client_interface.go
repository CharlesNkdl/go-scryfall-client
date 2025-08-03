package services

import (
	"context"
	"net/http"
)

type HTTPClient interface {
	NewRequest(ctx context.Context, method, path string) (*http.Request, error)
	Do(req *http.Request, v interface{}) error
}
