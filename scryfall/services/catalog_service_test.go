package services

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/cnkdl/go-scryfall-client/scryfall/models"
)

func TestCatalogService_GetCatalog(t *testing.T) {
	tests := []struct {
		name        string
		catalogType string
		mockClient  *MockHTTPClient
		want        *models.Catalog
		wantErr     bool
		errMsg      string
	}{
		{
			name:        "successful get card names catalog",
			catalogType: "card-names",
			mockClient: &MockHTTPClient{
				DoFunc: func(req *http.Request, v interface{}) error {
					if catalog, ok := v.(*models.Catalog); ok {
						catalog.Object = "catalog"
						catalog.URI = "https://api.scryfall.com/catalog/card-names"
						catalog.TotalValues = 3
						catalog.Data = []string{
							"Lightning Bolt",
							"Lightning Strike", 
							"Lightning Helix",
						}
					}
					return nil
				},
			},
			want: &models.Catalog{
				Object:      "catalog",
				URI:         "https://api.scryfall.com/catalog/card-names",
				TotalValues: 3,
				Data: []string{
					"Lightning Bolt",
					"Lightning Strike",
					"Lightning Helix",
				},
			},
			wantErr: false,
		},
		{
			name:        "successful get artist names catalog",
			catalogType: "artist-names",
			mockClient: &MockHTTPClient{
				DoFunc: func(req *http.Request, v interface{}) error {
					if catalog, ok := v.(*models.Catalog); ok {
						catalog.Object = "catalog"
						catalog.URI = "https://api.scryfall.com/catalog/artist-names"
						catalog.TotalValues = 2
						catalog.Data = []string{
							"Christopher Rush",
							"Richard Garfield",
						}
					}
					return nil
				},
			},
			want: &models.Catalog{
				Object:      "catalog",
				URI:         "https://api.scryfall.com/catalog/artist-names",
				TotalValues: 2,
				Data: []string{
					"Christopher Rush",
					"Richard Garfield",
				},
			},
			wantErr: false,
		},
		{
			name:        "empty catalog type",
			catalogType: "",
			mockClient: &MockHTTPClient{
				DoFunc: func(req *http.Request, v interface{}) error {
					if catalog, ok := v.(*models.Catalog); ok {
						catalog.Object = "catalog"
						catalog.Data = []string{}
					}
					return nil
				},
			},
			want: &models.Catalog{
				Object: "catalog",
				Data:   []string{},
			},
			wantErr: false,
		},
		{
			name:        "invalid catalog type",
			catalogType: "invalid-catalog-type",
			mockClient: &MockHTTPClient{
				DoFunc: func(req *http.Request, v interface{}) error {
					return errors.New("catalog not found")
				},
			},
			want:    nil,
			wantErr: true,
			errMsg:  "catalog not found",
		},
		{
			name:        "new request error",
			catalogType: "card-names",
			mockClient: &MockHTTPClient{
				NewRequestFunc: func(ctx context.Context, method, path string) (*http.Request, error) {
					return nil, errors.New("failed to create request")
				},
			},
			want:    nil,
			wantErr: true,
			errMsg:  "failed to create request",
		},
		{
			name:        "empty catalog response",
			catalogType: "word-bank",
			mockClient: &MockHTTPClient{
				DoFunc: func(req *http.Request, v interface{}) error {
					if catalog, ok := v.(*models.Catalog); ok {
						catalog.Object = "catalog"
						catalog.URI = "https://api.scryfall.com/catalog/word-bank"
						catalog.TotalValues = 0
						catalog.Data = []string{}
					}
					return nil
				},
			},
			want: &models.Catalog{
				Object:      "catalog",
				URI:         "https://api.scryfall.com/catalog/word-bank",
				TotalValues: 0,
				Data:        []string{},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &CatalogService{
				Client: tt.mockClient,
			}

			got, err := s.GetCatalog(context.Background(), tt.catalogType)

			if tt.wantErr {
				if err == nil {
					t.Errorf("CatalogService.GetCatalog() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errMsg != "" && !contains(err.Error(), tt.errMsg) {
					t.Errorf("CatalogService.GetCatalog() error = %v, want error containing %v", err, tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("CatalogService.GetCatalog() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !equalCatalogsDetailed(got, tt.want) {
				t.Errorf("CatalogService.GetCatalog() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test that the correct path is being constructed
func TestCatalogService_GetCatalog_PathConstruction(t *testing.T) {
	tests := []struct {
		catalogType  string
		expectedPath string
	}{
		{"card-names", "/catalog/card-names"},
		{"artist-names", "/catalog/artist-names"},
		{"word-bank", "/catalog/word-bank"},
		{"creature-types", "/catalog/creature-types"},
		{"planeswalker-types", "/catalog/planeswalker-types"},
		{"land-types", "/catalog/land-types"},
		{"artifact-types", "/catalog/artifact-types"},
		{"enchantment-types", "/catalog/enchantment-types"},
		{"spell-types", "/catalog/spell-types"},
		{"powers", "/catalog/powers"},
		{"toughnesses", "/catalog/toughnesses"},
		{"loyalties", "/catalog/loyalties"},
		{"watermarks", "/catalog/watermarks"},
		{"keyword-abilities", "/catalog/keyword-abilities"},
		{"keyword-actions", "/catalog/keyword-actions"},
		{"ability-words", "/catalog/ability-words"},
	}

	for _, tt := range tests {
		t.Run("path_for_"+tt.catalogType, func(t *testing.T) {
			mockClient := &MockHTTPClient{
				NewRequestFunc: func(ctx context.Context, method, path string) (*http.Request, error) {
					if path != tt.expectedPath {
						t.Errorf("Expected path %s, got %s", tt.expectedPath, path)
					}
					if method != "GET" {
						t.Errorf("Expected method GET, got %s", method)
					}
					return &http.Request{}, nil
				},
				DoFunc: func(req *http.Request, v interface{}) error {
					return nil
				},
			}

			s := &CatalogService{
				Client: mockClient,
			}

			_, err := s.GetCatalog(context.Background(), tt.catalogType)
			if err != nil {
				t.Errorf("CatalogService.GetCatalog() unexpected error: %v", err)
			}
		})
	}
}

// Helper function for comparing catalogs with more detail
func equalCatalogsDetailed(a, b *models.Catalog) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if a.Object != b.Object || a.URI != b.URI || a.TotalValues != b.TotalValues {
		return false
	}
	if len(a.Data) != len(b.Data) {
		return false
	}
	for i := range a.Data {
		if a.Data[i] != b.Data[i] {
			return false
		}
	}
	return true
}