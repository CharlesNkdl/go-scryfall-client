package services

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/cnkdl/go-scryfall-client/scryfall/models"
)

func TestSymbolService_GetSymbol(t *testing.T) {
	tests := []struct {
		name       string
		symbol     string
		mockClient *MockHTTPClient
		want       *models.Symbol
		wantErr    bool
		errMsg     string
	}{
		{
			name:   "successful get mana symbol",
			symbol: "{R}",
			mockClient: &MockHTTPClient{
				DoFunc: func(req *http.Request, v interface{}) error {
					if sym, ok := v.(*models.Symbol); ok {
						sym.Object = "card_symbol"
						sym.Symbol = "{R}"
						sym.LooseVariant = stringPtr("{r}")
						sym.EnglishName = "one red mana"
						sym.Transposable = true
						sym.RepresentsMana = true
						sym.AppearsInManaCosts = true
						sym.ManaValue = floatPtr(1)
						sym.Funny = false
						sym.Colors = []models.Color{"R"}
					}
					return nil
				},
			},
			want: &models.Symbol{
				Object:             "card_symbol",
				Symbol:             "{R}",
				LooseVariant:       stringPtr("{r}"),
				EnglishName:        "one red mana",
				Transposable:       true,
				RepresentsMana:     true,
				AppearsInManaCosts: true,
				ManaValue:          floatPtr(1),
				Funny:              false,
				Colors:             []models.Color{"R"},
			},
			wantErr: false,
		},
		{
			name:   "successful get generic mana symbol",
			symbol: "{1}",
			mockClient: &MockHTTPClient{
				DoFunc: func(req *http.Request, v interface{}) error {
					if sym, ok := v.(*models.Symbol); ok {
						sym.Object = "card_symbol"
						sym.Symbol = "{1}"
						sym.EnglishName = "one generic mana"
						sym.Transposable = false
						sym.RepresentsMana = true
						sym.AppearsInManaCosts = true
						sym.ManaValue = floatPtr(1)
						sym.Funny = false
						sym.Colors = []models.Color{}
					}
					return nil
				},
			},
			want: &models.Symbol{
				Object:             "card_symbol",
				Symbol:             "{1}",
				EnglishName:        "one generic mana",
				Transposable:       false,
				RepresentsMana:     true,
				AppearsInManaCosts: true,
				ManaValue:          floatPtr(1),
				Funny:              false,
				Colors:             []models.Color{},
			},
			wantErr: false,
		},
		{
			name:   "successful get tap symbol",
			symbol: "{T}",
			mockClient: &MockHTTPClient{
				DoFunc: func(req *http.Request, v interface{}) error {
					if sym, ok := v.(*models.Symbol); ok {
						sym.Object = "card_symbol"
						sym.Symbol = "{T}"
						sym.EnglishName = "tap this permanent"
						sym.Transposable = false
						sym.RepresentsMana = false
						sym.AppearsInManaCosts = false
						sym.Funny = false
						sym.Colors = []models.Color{}
					}
					return nil
				},
			},
			want: &models.Symbol{
				Object:             "card_symbol",
				Symbol:             "{T}",
				EnglishName:        "tap this permanent",
				Transposable:       false,
				RepresentsMana:     false,
				AppearsInManaCosts: false,
				Funny:              false,
				Colors:             []models.Color{},
			},
			wantErr: false,
		},
		{
			name:   "empty symbol",
			symbol: "",
			mockClient: &MockHTTPClient{
				DoFunc: func(req *http.Request, v interface{}) error {
					if sym, ok := v.(*models.Symbol); ok {
						sym.Object = "card_symbol"
					}
					return nil
				},
			},
			want: &models.Symbol{
				Object: "card_symbol",
			},
			wantErr: false,
		},
		{
			name:   "invalid symbol",
			symbol: "{INVALID}",
			mockClient: &MockHTTPClient{
				DoFunc: func(req *http.Request, v interface{}) error {
					return errors.New("symbol not found")
				},
			},
			want:    nil,
			wantErr: true,
			errMsg:  "symbol not found",
		},
		{
			name:   "new request error",
			symbol: "{R}",
			mockClient: &MockHTTPClient{
				NewRequestFunc: func(ctx context.Context, method, path string) (*http.Request, error) {
					return nil, errors.New("failed to create request")
				},
			},
			want:    nil,
			wantErr: true,
			errMsg:  "failed to create request",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SymbolService{
				Client: tt.mockClient,
			}

			got, err := s.GetSymbol(context.Background(), tt.symbol)

			if tt.wantErr {
				if err == nil {
					t.Errorf("SymbolService.GetSymbol() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errMsg != "" && !contains(err.Error(), tt.errMsg) {
					t.Errorf("SymbolService.GetSymbol() error = %v, want error containing %v", err, tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("SymbolService.GetSymbol() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !equalSymbols(got, tt.want) {
				t.Errorf("SymbolService.GetSymbol() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test that the correct path is being constructed and URL encoding works
func TestSymbolService_GetSymbol_PathConstruction(t *testing.T) {
	tests := []struct {
		symbol       string
		expectedPath string
	}{
		{"{R}", "/symbology/%7BR%7D"},     // URL encoded braces
		{"{U}", "/symbology/%7BU%7D"},     // URL encoded braces
		{"{B}", "/symbology/%7BB%7D"},     // URL encoded braces
		{"{G}", "/symbology/%7BG%7D"},     // URL encoded braces
		{"{W}", "/symbology/%7BW%7D"},     // URL encoded braces
		{"{T}", "/symbology/%7BT%7D"},     // URL encoded braces
		{"{1}", "/symbology/%7B1%7D"},     // URL encoded braces
		{"{X}", "/symbology/%7BX%7D"},     // URL encoded braces
		{"{W/U}", "/symbology/%7BW%2FU%7D"}, // URL encoded braces and slash
		{"{2/W}", "/symbology/%7B2%2FW%7D"}, // URL encoded braces and slash
		{"{R/P}", "/symbology/%7BR%2FP%7D"}, // URL encoded braces and slash for Phyrexian mana
	}

	for _, tt := range tests {
		t.Run("path_for_"+tt.symbol, func(t *testing.T) {
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

			s := &SymbolService{
				Client: mockClient,
			}

			_, err := s.GetSymbol(context.Background(), tt.symbol)
			if err != nil {
				t.Errorf("SymbolService.GetSymbol() unexpected error: %v", err)
			}
		})
	}
}

// Test URL encoding specifically
func TestSymbolService_GetSymbol_URLEncoding(t *testing.T) {
	symbol := "{W/U}"
	expectedEncodedSymbol := url.QueryEscape(symbol)

	mockClient := &MockHTTPClient{
		NewRequestFunc: func(ctx context.Context, method, path string) (*http.Request, error) {
			expectedPath := "/symbology/" + expectedEncodedSymbol
			if path != expectedPath {
				t.Errorf("Expected path %s, got %s", expectedPath, path)
			}
			return &http.Request{}, nil
		},
		DoFunc: func(req *http.Request, v interface{}) error {
			return nil
		},
	}

	s := &SymbolService{
		Client: mockClient,
	}

	_, err := s.GetSymbol(context.Background(), symbol)
	if err != nil {
		t.Errorf("SymbolService.GetSymbol() unexpected error: %v", err)
	}
}

// Test various special symbols that might need encoding
func TestSymbolService_GetSymbol_SpecialSymbols(t *testing.T) {
	specialSymbols := []string{
		"{W/U}", // Hybrid mana
		"{2/W}", // Hybrid mana with generic
		"{R/P}", // Phyrexian mana
		"{∞}",   // Infinity symbol (if used)
		"{½}",   // Half symbol (if used)
		"{E}",   // Energy symbol
		"{Q}",   // Untap symbol
		"{CHAOS}", // Chaos symbol for planechase
	}

	for _, symbol := range specialSymbols {
		t.Run("special_symbol_"+symbol, func(t *testing.T) {
			mockClient := &MockHTTPClient{
				NewRequestFunc: func(ctx context.Context, method, path string) (*http.Request, error) {
					// Verify that the symbol was URL encoded properly
					expectedPath := "/symbology/" + url.QueryEscape(symbol)
					if path != expectedPath {
						t.Errorf("For symbol %s, expected path %s, got %s", symbol, expectedPath, path)
					}
					return &http.Request{}, nil
				},
				DoFunc: func(req *http.Request, v interface{}) error {
					return nil
				},
			}

			s := &SymbolService{
				Client: mockClient,
			}

			_, err := s.GetSymbol(context.Background(), symbol)
			if err != nil {
				t.Errorf("SymbolService.GetSymbol() unexpected error for symbol %s: %v", symbol, err)
			}
		})
	}
}

// Helper function for comparing symbols
func equalSymbols(a, b *models.Symbol) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	
	// Compare basic fields
	if a.Object != b.Object ||
		a.Symbol != b.Symbol ||
		a.EnglishName != b.EnglishName ||
		a.Transposable != b.Transposable ||
		a.RepresentsMana != b.RepresentsMana ||
		a.AppearsInManaCosts != b.AppearsInManaCosts ||
		a.Funny != b.Funny {
		return false
	}
	
	// Compare optional fields
	if (a.LooseVariant == nil) != (b.LooseVariant == nil) {
		return false
	}
	if a.LooseVariant != nil && b.LooseVariant != nil && *a.LooseVariant != *b.LooseVariant {
		return false
	}
	
	if (a.ManaValue == nil) != (b.ManaValue == nil) {
		return false
	}
	if a.ManaValue != nil && b.ManaValue != nil && *a.ManaValue != *b.ManaValue {
		return false
	}
	
	// Compare colors slice
	if len(a.Colors) != len(b.Colors) {
		return false
	}
	for i := range a.Colors {
		if a.Colors[i] != b.Colors[i] {
			return false
		}
	}
	
	return true
}

func floatPtr(f float32) *float32 {
	return &f
}