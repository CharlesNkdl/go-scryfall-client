package scryfall

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/time/rate"
	"net/http"
	"time"

	"github.com/cnkdl/go-scryfall-client/scryfall/errors"
	"github.com/cnkdl/go-scryfall-client/scryfall/models"
	"github.com/cnkdl/go-scryfall-client/scryfall/services"
)

const (
	defaultBaseURL = "https://api.scryfall.com"
)

// Client represents a Scryfall API client with rate limiting and HTTP configuration.
type Client struct {
	httpClient *http.Client
	BaseURL    string

	RateLimiter *rate.Limiter // asked by scryfall API documentation

	Cards    *services.CardService
	Sets     *services.SetService
	Rulings  *services.RulingService
	Catalogs *services.CatalogService
	Symbols  *services.SymbolService
}

// NewClient creates a new Scryfall API client with default configuration.
func NewClient() *Client {
	c := &Client{
		httpClient: &http.Client{
			Timeout: 20 * time.Second,
		},
		BaseURL:     defaultBaseURL,
		RateLimiter: rate.NewLimiter(1, 5),
	}
	c.Cards = &services.CardService{Client: c}
	c.Sets = &services.SetService{Client: c}
	c.Rulings = &services.RulingService{Client: c}
	c.Catalogs = &services.CatalogService{Client: c}
	c.Symbols = &services.SymbolService{Client: c}

	return c
}

// NewRequest creates a new HTTP request with the given method and path.
func (c *Client) NewRequest(ctx context.Context, method, path string) (*http.Request, error) {
	url := c.BaseURL + path
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "go-scryfall-client/1.0")
	req.Header.Set("Accept", "application/json")

	return req, nil
}

// Do executes an HTTP request and decodes the response into the provided interface.
func (c *Client) Do(req *http.Request, v interface{}) error {
	if err := c.RateLimiter.Wait(req.Context()); err != nil {
		return err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusTooManyRequests {
		apiErr := &models.ScryfallError{}
		json.NewDecoder(resp.Body).Decode(&apiErr.Detail)
		return &errors.APIError{
			ErrInfo: models.ScryfallError{
				Status: resp.StatusCode,
				Code:   "too_many_requests",
				Detail: apiErr.Detail,
			},
		}
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		apiErr := &models.ScryfallError{}
		if err := json.NewDecoder(resp.Body).Decode(apiErr); err != nil {
			return fmt.Errorf("failed to decode error response: %w", err)
		}
		apiErr.Status = resp.StatusCode
		return &errors.APIError{ErrInfo: *apiErr}
	}
	if v != nil {
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}
