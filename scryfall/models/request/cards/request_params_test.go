package cards

import (
	"net/url"
	"strings"
	"testing"
)

func TestIdCardParams_Validate(t *testing.T) {
	tests := []struct {
		name    string
		params  IdCardParams
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid id",
			params: IdCardParams{
				ID: "f2f0e0e6-8eee-4001-a6a8-2c8f1b9a3c3e",
			},
			wantErr: false,
		},
		{
			name: "empty id",
			params: IdCardParams{
				ID: "",
			},
			wantErr: true,
			errMsg:  "id must be a non-empty string",
		},
		{
			name: "valid format json",
			params: IdCardParams{
				ID:     "test-id",
				Format: stringPtr("json"),
			},
			wantErr: false,
		},
		{
			name: "valid format text",
			params: IdCardParams{
				ID:     "test-id",
				Format: stringPtr("text"),
			},
			wantErr: false,
		},
		{
			name: "valid format image",
			params: IdCardParams{
				ID:     "test-id",
				Format: stringPtr("image"),
			},
			wantErr: false,
		},
		{
			name: "invalid format",
			params: IdCardParams{
				ID:     "test-id",
				Format: stringPtr("xml"),
			},
			wantErr: true,
			errMsg:  "invalid format: must be 'json', 'text', or 'image'",
		},
		{
			name: "valid version",
			params: IdCardParams{
				ID:      "test-id",
				Version: stringPtr("normal"),
			},
			wantErr: false,
		},
		{
			name: "invalid version",
			params: IdCardParams{
				ID:      "test-id",
				Version: stringPtr("invalid"),
			},
			wantErr: true,
			errMsg:  "invalid version",
		},
		{
			name: "valid face",
			params: IdCardParams{
				ID:   "test-id",
				Face: stringPtr("front"),
			},
			wantErr: false,
		},
		{
			name: "invalid face",
			params: IdCardParams{
				ID:   "test-id",
				Face: stringPtr("invalid"),
			},
			wantErr: true,
			errMsg:  "invalid face: must be 'front' or 'back'",
		},
		{
			name: "valid pretty true",
			params: IdCardParams{
				ID:     "test-id",
				Pretty: boolPtr(true),
			},
			wantErr: false,
		},
		{
			name: "valid pretty false",
			params: IdCardParams{
				ID:     "test-id",
				Pretty: boolPtr(false),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.params.Validate()
			if tt.wantErr {
				if err == nil {
					t.Errorf("IdCardParams.Validate() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("IdCardParams.Validate() error = %v, want error containing %v", err, tt.errMsg)
				}
				return
			}
			if err != nil {
				t.Errorf("IdCardParams.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIdCardParams_ToURLValues(t *testing.T) {
	tests := []struct {
		name    string
		params  IdCardParams
		want    url.Values
		wantErr bool
	}{
		{
			name: "minimal valid params",
			params: IdCardParams{
				ID: "test-id",
			},
			want:    url.Values{},
			wantErr: false,
		},
		{
			name: "all params",
			params: IdCardParams{
				ID:      "test-id",
				Format:  stringPtr("json"),
				Version: stringPtr("normal"),
				Face:    stringPtr("front"),
				Pretty:  boolPtr(true),
			},
			want: url.Values{
				"format":  []string{"json"},
				"version": []string{"normal"},
				"face":    []string{"front"},
				"pretty":  []string{"true"},
			},
			wantErr: false,
		},
		{
			name: "invalid params",
			params: IdCardParams{
				ID: "",
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.params.ToURLValues()
			if tt.wantErr {
				if err == nil {
					t.Errorf("IdCardParams.ToURLValues() error = nil, wantErr %v", tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Errorf("IdCardParams.ToURLValues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !equalURLValues(got, tt.want) {
				t.Errorf("IdCardParams.ToURLValues() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNamedCardParams_Validate(t *testing.T) {
	tests := []struct {
		name    string
		params  NamedCardParams
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid exact search",
			params: NamedCardParams{
				Exact: stringPtr("Lightning Bolt"),
			},
			wantErr: false,
		},
		{
			name: "valid fuzzy search",
			params: NamedCardParams{
				Fuzzy: stringPtr("Lightning Bolt"),
			},
			wantErr: false,
		},
		{
			name:    "no search parameters",
			params:  NamedCardParams{},
			wantErr: true,
			errMsg:  "either 'exact' or 'fuzzy' parameter is required",
		},
		{
			name: "both exact and fuzzy",
			params: NamedCardParams{
				Exact: stringPtr("Lightning Bolt"),
				Fuzzy: stringPtr("Lightning Bolt"),
			},
			wantErr: true,
			errMsg:  "'exact' and 'fuzzy' parameters are mutually exclusive",
		},
		{
			name: "empty exact",
			params: NamedCardParams{
				Exact: stringPtr(""),
			},
			wantErr: true,
			errMsg:  "exact name cannot be empty",
		},
		{
			name: "empty fuzzy",
			params: NamedCardParams{
				Fuzzy: stringPtr(""),
			},
			wantErr: true,
			errMsg:  "fuzzy name cannot be empty",
		},
		{
			name: "valid set code",
			params: NamedCardParams{
				Exact: stringPtr("Lightning Bolt"),
				Set:   stringPtr("khm"),
			},
			wantErr: false,
		},
		{
			name: "invalid set code too short",
			params: NamedCardParams{
				Exact: stringPtr("Lightning Bolt"),
				Set:   stringPtr("ab"),
			},
			wantErr: true,
			errMsg:  "set code must be between 3 and 5 characters",
		},
		{
			name: "invalid set code too long",
			params: NamedCardParams{
				Exact: stringPtr("Lightning Bolt"),
				Set:   stringPtr("abcdef"),
			},
			wantErr: true,
			errMsg:  "set code must be between 3 and 5 characters",
		},
		{
			name: "face back without image format",
			params: NamedCardParams{
				Exact: stringPtr("Lightning Bolt"),
				Face:  stringPtr("back"),
			},
			wantErr: true,
			errMsg:  "face=back can only be used with format=image",
		},
		{
			name: "face back with image format",
			params: NamedCardParams{
				Exact:  stringPtr("Lightning Bolt"),
				Face:   stringPtr("back"),
				Format: stringPtr("image"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.params.Validate()
			if tt.wantErr {
				if err == nil {
					t.Errorf("NamedCardParams.Validate() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("NamedCardParams.Validate() error = %v, want error containing %v", err, tt.errMsg)
				}
				return
			}
			if err != nil {
				t.Errorf("NamedCardParams.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSearchParams_Validate(t *testing.T) {
	tests := []struct {
		name    string
		params  SearchParams
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid query",
			params: SearchParams{
				Query: "lightning",
			},
			wantErr: false,
		},
		{
			name: "empty query",
			params: SearchParams{
				Query: "",
			},
			wantErr: true,
			errMsg:  "query parameter 'q' is required",
		},
		{
			name: "whitespace only query",
			params: SearchParams{
				Query: "   ",
			},
			wantErr: true,
			errMsg:  "query parameter 'q' is required",
		},
		{
			name: "query too long",
			params: SearchParams{
				Query: strings.Repeat("a", 1001),
			},
			wantErr: true,
			errMsg:  "query parameter 'q' must not exceed 1000 Unicode characters",
		},
		{
			name: "negative page",
			params: SearchParams{
				Query: "lightning",
				Page:  -1,
			},
			wantErr: true,
			errMsg:  "page must be positive",
		},
		{
			name: "zero page is valid",
			params: SearchParams{
				Query: "lightning",
				Page:  0,
			},
			wantErr: false,
		},
		{
			name: "valid unique mode",
			params: SearchParams{
				Query:  "lightning",
				Unique: UniqueCards,
			},
			wantErr: false,
		},
		{
			name: "invalid order mode",
			params: SearchParams{
				Query: "lightning",
				Order: "invalid",
			},
			wantErr: true,
			errMsg:  "invalid order mode",
		},
		{
			name: "invalid direction",
			params: SearchParams{
				Query: "lightning",
				Dir:   "invalid",
			},
			wantErr: true,
			errMsg:  "invalid direction: must be 'auto', 'asc', or 'desc'",
		},
		{
			name: "invalid format",
			params: SearchParams{
				Query:  "lightning",
				Format: "xml",
			},
			wantErr: true,
			errMsg:  "invalid format: must be 'json' or 'csv'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.params.Validate()
			if tt.wantErr {
				if err == nil {
					t.Errorf("SearchParams.Validate() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("SearchParams.Validate() error = %v, want error containing %v", err, tt.errMsg)
				}
				return
			}
			if err != nil {
				t.Errorf("SearchParams.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRandomCardParams_Validate(t *testing.T) {
	tests := []struct {
		name    string
		params  RandomCardParams
		wantErr bool
		errMsg  string
	}{
		{
			name:    "empty params",
			params:  RandomCardParams{},
			wantErr: false,
		},
		{
			name: "valid query",
			params: RandomCardParams{
				Query: stringPtr("type:creature"),
			},
			wantErr: false,
		},
		{
			name: "empty query",
			params: RandomCardParams{
				Query: stringPtr(""),
			},
			wantErr: true,
			errMsg:  "search query 'q' cannot be empty when provided",
		},
		{
			name: "valid format",
			params: RandomCardParams{
				Format: stringPtr("json"),
			},
			wantErr: false,
		},
		{
			name: "invalid format",
			params: RandomCardParams{
				Format: stringPtr("xml"),
			},
			wantErr: true,
			errMsg:  "invalid format",
		},
		{
			name: "valid version",
			params: RandomCardParams{
				Version: stringPtr("normal"),
			},
			wantErr: false,
		},
		{
			name: "invalid version",
			params: RandomCardParams{
				Version: stringPtr("invalid"),
			},
			wantErr: true,
			errMsg:  "invalid version",
		},
		{
			name: "valid face with image format",
			params: RandomCardParams{
				Format: stringPtr("image"),
				Face:   stringPtr("back"),
			},
			wantErr: false,
		},
		{
			name: "invalid face without image format",
			params: RandomCardParams{
				Face: stringPtr("back"),
			},
			wantErr: true,
			errMsg:  "face=back can only be used with format=image",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.params.Validate()
			if tt.wantErr {
				if err == nil {
					t.Errorf("RandomCardParams.Validate() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("RandomCardParams.Validate() error = %v, want error containing %v", err, tt.errMsg)
				}
				return
			}
			if err != nil {
				t.Errorf("RandomCardParams.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAutocompleteCardParams_Validate(t *testing.T) {
	tests := []struct {
		name    string
		params  AutocompleteCardParams
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid query",
			params: AutocompleteCardParams{
				Query: "light",
			},
			wantErr: false,
		},
		{
			name: "empty query",
			params: AutocompleteCardParams{
				Query: "",
			},
			wantErr: true,
			errMsg:  "query cannot be empty",
		},
		{
			name: "valid format",
			params: AutocompleteCardParams{
				Query:  "light",
				Format: stringPtr("json"),
			},
			wantErr: false,
		},
		{
			name: "invalid format",
			params: AutocompleteCardParams{
				Query:  "light",
				Format: stringPtr("xml"),
			},
			wantErr: true,
			errMsg:  "invalid format: must be 'json'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.params.Validate()
			if tt.wantErr {
				if err == nil {
					t.Errorf("AutocompleteCardParams.Validate() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("AutocompleteCardParams.Validate() error = %v, want error containing %v", err, tt.errMsg)
				}
				return
			}
			if err != nil {
				t.Errorf("AutocompleteCardParams.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}

func equalURLValues(a, b url.Values) bool {
	if len(a) != len(b) {
		return false
	}
	for key, valuesA := range a {
		valuesB, exists := b[key]
		if !exists {
			return false
		}
		if len(valuesA) != len(valuesB) {
			return false
		}
		for i, valueA := range valuesA {
			if valueA != valuesB[i] {
				return false
			}
		}
	}
	return true
}