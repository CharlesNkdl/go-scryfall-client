package services

import (
	"context"
	"fmt"
	"net/http"
)

// MockHTTPClient is a mock implementation of HTTPClient for testing
type MockHTTPClient struct {
	NewRequestFunc func(ctx context.Context, method, path string) (*http.Request, error)
	DoFunc         func(req *http.Request, v interface{}) error
}

func (m *MockHTTPClient) NewRequest(ctx context.Context, method, path string) (*http.Request, error) {
	if m.NewRequestFunc != nil {
		return m.NewRequestFunc(ctx, method, path)
	}
	// Default implementation
	req, err := http.NewRequestWithContext(ctx, method, "https://api.scryfall.com"+path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	return req, nil
}

func (m *MockHTTPClient) Do(req *http.Request, v interface{}) error {
	if m.DoFunc != nil {
		return m.DoFunc(req, v)
	}
	// Default implementation - returns nil (success)
	return nil
}