package services

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/cnkdl/go-scryfall-client/scryfall/models"
	"github.com/cnkdl/go-scryfall-client/scryfall/models/card"
	cardreq "github.com/cnkdl/go-scryfall-client/scryfall/models/request/cards"
)

func TestCardService_GetById(t *testing.T) {
	tests := []struct {
		name       string
		params     *cardreq.IdCardParams
		mockClient *MockHTTPClient
		want       *card.CardCore
		wantErr    bool
		errMsg     string
	}{
		{
			name: "successful get by id",
			params: &cardreq.IdCardParams{
				ID: "f2f0e0e6-8eee-4001-a6a8-2c8f1b9a3c3e",
			},
			mockClient: &MockHTTPClient{
				DoFunc: func(req *http.Request, v interface{}) error {
					if cardCore, ok := v.(*card.CardCore); ok {
						cardCore.ID = "f2f0e0e6-8eee-4001-a6a8-2c8f1b9a3c3e"
						cardCore.Object = "card"
						cardCore.GameplayFields.Name = "Lightning Bolt"
					}
					return nil
				},
			},
			want: &card.CardCore{
				ID:     "f2f0e0e6-8eee-4001-a6a8-2c8f1b9a3c3e",
				Object: "card",
				GameplayFields: card.Gameplay{
					Name: "Lightning Bolt",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid id parameter",
			params: &cardreq.IdCardParams{
				ID: "",
			},
			mockClient: &MockHTTPClient{},
			want:       nil,
			wantErr:    true,
			errMsg:     "invalid named parameters",
		},
		{
			name: "invalid format parameter",
			params: &cardreq.IdCardParams{
				ID:     "f2f0e0e6-8eee-4001-a6a8-2c8f1b9a3c3e",
				Format: stringPtr("invalid"),
			},
			mockClient: &MockHTTPClient{},
			want:       nil,
			wantErr:    true,
			errMsg:     "invalid named parameters",
		},
		{
			name: "http client error",
			params: &cardreq.IdCardParams{
				ID: "f2f0e0e6-8eee-4001-a6a8-2c8f1b9a3c3e",
			},
			mockClient: &MockHTTPClient{
				DoFunc: func(req *http.Request, v interface{}) error {
					return errors.New("network error")
				},
			},
			want:    nil,
			wantErr: true,
			errMsg:  "network error",
		},
		{
			name: "with format parameter",
			params: &cardreq.IdCardParams{
				ID:     "f2f0e0e6-8eee-4001-a6a8-2c8f1b9a3c3e",
				Format: stringPtr("json"),
			},
			mockClient: &MockHTTPClient{
				DoFunc: func(req *http.Request, v interface{}) error {
					if cardCore, ok := v.(*card.CardCore); ok {
						cardCore.ID = "f2f0e0e6-8eee-4001-a6a8-2c8f1b9a3c3e"
						cardCore.Object = "card"
					}
					return nil
				},
			},
			want: &card.CardCore{
				ID:     "f2f0e0e6-8eee-4001-a6a8-2c8f1b9a3c3e",
				Object: "card",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &CardService{
				Client: tt.mockClient,
			}

			got, err := s.GetById(context.Background(), tt.params)

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetById() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errMsg != "" && err.Error() != tt.errMsg {
					// For flexible error message matching
					if len(tt.errMsg) > 0 && len(err.Error()) > 0 {
						// Just check if the expected error message is contained in the actual error
						if !contains(err.Error(), tt.errMsg) {
							t.Errorf("GetById() error = %v, want error containing %v", err, tt.errMsg)
						}
					}
				}
				return
			}

			if err != nil {
				t.Errorf("GetById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !equalCards(got, tt.want) {
				t.Errorf("GetById() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCardService_GetByName(t *testing.T) {
	tests := []struct {
		name       string
		params     *cardreq.NamedCardParams
		mockClient *MockHTTPClient
		want       *card.CardCore
		wantErr    bool
		errMsg     string
	}{
		{
			name: "successful exact search",
			params: &cardreq.NamedCardParams{
				Exact: stringPtr("Lightning Bolt"),
			},
			mockClient: &MockHTTPClient{
				DoFunc: func(req *http.Request, v interface{}) error {
					if cardCore, ok := v.(*card.CardCore); ok {
						cardCore.ID = "test-id"
						cardCore.Object = "card"
						cardCore.GameplayFields.Name = "Lightning Bolt"
					}
					return nil
				},
			},
			want: &card.CardCore{
				ID:     "test-id",
				Object: "card",
				GameplayFields: card.Gameplay{
					Name: "Lightning Bolt",
				},
			},
			wantErr: false,
		},
		{
			name: "successful fuzzy search",
			params: &cardreq.NamedCardParams{
				Fuzzy: stringPtr("Lighning Bolt"),
			},
			mockClient: &MockHTTPClient{
				DoFunc: func(req *http.Request, v interface{}) error {
					if cardCore, ok := v.(*card.CardCore); ok {
						cardCore.ID = "test-id"
						cardCore.Object = "card"
						cardCore.GameplayFields.Name = "Lightning Bolt"
					}
					return nil
				},
			},
			want: &card.CardCore{
				ID:     "test-id",
				Object: "card",
				GameplayFields: card.Gameplay{
					Name: "Lightning Bolt",
				},
			},
			wantErr: false,
		},
		{
			name: "no search parameters provided",
			params: &cardreq.NamedCardParams{},
			mockClient: &MockHTTPClient{},
			want:       nil,
			wantErr:    true,
			errMsg:     "either 'exact' or 'fuzzy' parameter is required",
		},
		{
			name: "both exact and fuzzy provided",
			params: &cardreq.NamedCardParams{
				Exact: stringPtr("Lightning Bolt"),
				Fuzzy: stringPtr("Lighning Bolt"),
			},
			mockClient: &MockHTTPClient{},
			want:       nil,
			wantErr:    true,
			errMsg:     "'exact' and 'fuzzy' parameters are mutually exclusive",
		},
		{
			name: "empty exact parameter",
			params: &cardreq.NamedCardParams{
				Exact: stringPtr(""),
			},
			mockClient: &MockHTTPClient{},
			want:       nil,
			wantErr:    true,
			errMsg:     "exact name cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &CardService{
				Client: tt.mockClient,
			}

			got, err := s.GetByName(context.Background(), tt.params)

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetByName() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errMsg != "" && !contains(err.Error(), tt.errMsg) {
					t.Errorf("GetByName() error = %v, want error containing %v", err, tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("GetByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !equalCards(got, tt.want) {
				t.Errorf("GetByName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCardService_Search(t *testing.T) {
	tests := []struct {
		name       string
		params     *cardreq.SearchParams
		mockClient *MockHTTPClient
		want       *models.List[card.CardCore]
		wantErr    bool
		errMsg     string
	}{
		{
			name: "successful search",
			params: &cardreq.SearchParams{
				Query: "lightning",
			},
			mockClient: &MockHTTPClient{
				DoFunc: func(req *http.Request, v interface{}) error {
					if list, ok := v.(*models.List[card.CardCore]); ok {
						list.Object = "list"
						list.TotalCards = intPtr(1)
						list.Data = []card.CardCore{
							{
								ID:     "test-id",
								Object: "card",
								GameplayFields: card.Gameplay{
									Name: "Lightning Bolt",
								},
							},
						}
					}
					return nil
				},
			},
			want: &models.List[card.CardCore]{
				Object:     "list",
				TotalCards: intPtr(1),
				Data: []card.CardCore{
					{
						ID:     "test-id",
						Object: "card",
						GameplayFields: card.Gameplay{
							Name: "Lightning Bolt",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "empty query",
			params: &cardreq.SearchParams{
				Query: "",
			},
			mockClient: &MockHTTPClient{},
			want:       nil,
			wantErr:    true,
			errMsg:     "query parameter 'q' is required",
		},
		{
			name: "query too long",
			params: &cardreq.SearchParams{
				Query: string(make([]byte, 1001)), // Create a string longer than 1000 characters
			},
			mockClient: &MockHTTPClient{},
			want:       nil,
			wantErr:    true,
			errMsg:     "query parameter 'q' must not exceed 1000 Unicode characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &CardService{
				Client: tt.mockClient,
			}

			got, err := s.Search(context.Background(), tt.params)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Search() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errMsg != "" && !contains(err.Error(), tt.errMsg) {
					t.Errorf("Search() error = %v, want error containing %v", err, tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !equalLists(got, tt.want) {
				t.Errorf("Search() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCardService_Autocomplete(t *testing.T) {
	tests := []struct {
		name       string
		params     *cardreq.AutocompleteCardParams
		mockClient *MockHTTPClient
		want       *models.Catalog
		wantErr    bool
		errMsg     string
	}{
		{
			name: "successful autocomplete",
			params: &cardreq.AutocompleteCardParams{
				Query: "light",
			},
			mockClient: &MockHTTPClient{
				DoFunc: func(req *http.Request, v interface{}) error {
					if catalog, ok := v.(*models.Catalog); ok {
						catalog.Object = "catalog"
						catalog.Data = []string{"Lightning Bolt", "Lightning Strike"}
					}
					return nil
				},
			},
			want: &models.Catalog{
				Object: "catalog",
				Data:   []string{"Lightning Bolt", "Lightning Strike"},
			},
			wantErr: false,
		},
		{
			name: "empty query",
			params: &cardreq.AutocompleteCardParams{
				Query: "",
			},
			mockClient: &MockHTTPClient{},
			want:       nil,
			wantErr:    true,
			errMsg:     "query cannot be empty",
		},
		{
			name: "invalid format",
			params: &cardreq.AutocompleteCardParams{
				Query:  "light",
				Format: stringPtr("xml"),
			},
			mockClient: &MockHTTPClient{},
			want:       nil,
			wantErr:    true,
			errMsg:     "invalid format: must be 'json'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &CardService{
				Client: tt.mockClient,
			}

			got, err := s.Autocomplete(context.Background(), tt.params)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Autocomplete() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errMsg != "" && !contains(err.Error(), tt.errMsg) {
					t.Errorf("Autocomplete() error = %v, want error containing %v", err, tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("Autocomplete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !equalCatalogs(got, tt.want) {
				t.Errorf("Autocomplete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCardService_GetRandom(t *testing.T) {
	tests := []struct {
		name       string
		params     *cardreq.RandomCardParams
		mockClient *MockHTTPClient
		want       *card.CardCore
		wantErr    bool
		errMsg     string
	}{
		{
			name:   "successful random card",
			params: &cardreq.RandomCardParams{},
			mockClient: &MockHTTPClient{
				DoFunc: func(req *http.Request, v interface{}) error {
					if cardCore, ok := v.(*card.CardCore); ok {
						cardCore.ID = "random-id"
						cardCore.Object = "card"
						cardCore.GameplayFields.Name = "Random Card"
					}
					return nil
				},
			},
			want: &card.CardCore{
				ID:     "random-id",
				Object: "card",
				GameplayFields: card.Gameplay{
					Name: "Random Card",
				},
			},
			wantErr: false,
		},
		{
			name: "with query filter",
			params: &cardreq.RandomCardParams{
				Query: stringPtr("type:creature"),
			},
			mockClient: &MockHTTPClient{
				DoFunc: func(req *http.Request, v interface{}) error {
					if cardCore, ok := v.(*card.CardCore); ok {
						cardCore.ID = "creature-id"
						cardCore.Object = "card"
						cardCore.GameplayFields.Name = "Random Creature"
					}
					return nil
				},
			},
			want: &card.CardCore{
				ID:     "creature-id",
				Object: "card",
				GameplayFields: card.Gameplay{
					Name: "Random Creature",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid format",
			params: &cardreq.RandomCardParams{
				Format: stringPtr("invalid"),
			},
			mockClient: &MockHTTPClient{},
			want:       nil,
			wantErr:    true,
			errMsg:     "invalid format",
		},
		{
			name: "invalid face without image format",
			params: &cardreq.RandomCardParams{
				Face: stringPtr("back"),
			},
			mockClient: &MockHTTPClient{},
			want:       nil,
			wantErr:    true,
			errMsg:     "face=back can only be used with format=image",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &CardService{
				Client: tt.mockClient,
			}

			got, err := s.GetRandom(context.Background(), tt.params)

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetRandom() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errMsg != "" && !contains(err.Error(), tt.errMsg) {
					t.Errorf("GetRandom() error = %v, want error containing %v", err, tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("GetRandom() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !equalCards(got, tt.want) {
				t.Errorf("GetRandom() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCardService_GetByCodeNumberLang(t *testing.T) {
	tests := []struct {
		name       string
		params     *cardreq.CardsByCodeNumberLangParams
		mockClient *MockHTTPClient
		want       *models.List[card.CardCore]
		wantErr    bool
		errMsg     string
	}{
		{
			name: "successful get by code and number",
			params: &cardreq.CardsByCodeNumberLangParams{
				Code:   "khm",
				Number: "123",
			},
			mockClient: &MockHTTPClient{
				DoFunc: func(req *http.Request, v interface{}) error {
					if list, ok := v.(*models.List[card.CardCore]); ok {
						list.Object = "list"
						list.TotalCards = intPtr(1)
						list.Data = []card.CardCore{
							{
								ID:     "card-id",
								Object: "card",
								GameplayFields: card.Gameplay{
									Name: "Test Card",
								},
								PrintFields: card.Print{
									Set: "khm",
								},
							},
						}
					}
					return nil
				},
			},
			want: &models.List[card.CardCore]{
				Object:     "list",
				TotalCards: intPtr(1),
				Data: []card.CardCore{
					{
						ID:     "card-id",
						Object: "card",
						GameplayFields: card.Gameplay{
							Name: "Test Card",
						},
						PrintFields: card.Print{
							Set: "khm",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid code length",
			params: &cardreq.CardsByCodeNumberLangParams{
				Code:   "ab",
				Number: "123",
			},
			mockClient: &MockHTTPClient{},
			want:       nil,
			wantErr:    true,
			errMsg:     "code must be between 3 and 5 characters",
		},
		{
			name: "empty number",
			params: &cardreq.CardsByCodeNumberLangParams{
				Code:   "khm",
				Number: "",
			},
			mockClient: &MockHTTPClient{},
			want:       nil,
			wantErr:    true,
			errMsg:     "number cannot be empty",
		},
		{
			name: "unsupported language",
			params: &cardreq.CardsByCodeNumberLangParams{
				Code:   "khm",
				Number: "123",
				Lang:   stringPtr("invalid"),
			},
			mockClient: &MockHTTPClient{},
			want:       nil,
			wantErr:    true,
			errMsg:     "unsupported language",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &CardService{
				Client: tt.mockClient,
			}

			got, err := s.GetByCodeNumberLang(context.Background(), tt.params)

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetByCodeNumberLang() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errMsg != "" && !contains(err.Error(), tt.errMsg) {
					t.Errorf("GetByCodeNumberLang() error = %v, want error containing %v", err, tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("GetByCodeNumberLang() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !equalLists(got, tt.want) {
				t.Errorf("GetByCodeNumberLang() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helper functions for testing

func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func contains(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	if len(s) == 0 {
		return false
	}
	
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func equalCards(a, b *card.CardCore) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.ID == b.ID && a.Object == b.Object && 
		a.GameplayFields.Name == b.GameplayFields.Name && 
		a.PrintFields.Set == b.PrintFields.Set
}

func equalLists(a, b *models.List[card.CardCore]) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if a.Object != b.Object || len(a.Data) != len(b.Data) {
		return false
	}
	if (a.TotalCards == nil) != (b.TotalCards == nil) {
		return false
	}
	if a.TotalCards != nil && b.TotalCards != nil && *a.TotalCards != *b.TotalCards {
		return false
	}
	for i := range a.Data {
		if !equalCards(&a.Data[i], &b.Data[i]) {
			return false
		}
	}
	return true
}

func equalCatalogs(a, b *models.Catalog) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if a.Object != b.Object || len(a.Data) != len(b.Data) {
		return false
	}
	for i := range a.Data {
		if a.Data[i] != b.Data[i] {
			return false
		}
	}
	return true
}