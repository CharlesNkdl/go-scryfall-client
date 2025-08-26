package services

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/cnkdl/go-scryfall-client/scryfall/models"
)

func TestRulingService_GetRulings(t *testing.T) {
	tests := []struct {
		name       string
		cardID     string
		mockClient *MockHTTPClient
		want       []models.Ruling
		wantErr    bool
		errMsg     string
	}{
		{
			name:   "successful get rulings",
			cardID: "f2f0e0e6-8eee-4001-a6a8-2c8f1b9a3c3e",
			mockClient: &MockHTTPClient{
				DoFunc: func(req *http.Request, v interface{}) error {
					if rulings, ok := v.(*[]models.Ruling); ok {
						*rulings = []models.Ruling{
							{
								Object:      "ruling",
								OracleID:    "oracle-123",
								Source:      "scryfall",
								PublishedAt: "2004-10-04",
								Comment:     "Lightning Bolt deals 3 damage to any target.",
							},
							{
								Object:      "ruling",
								OracleID:    "oracle-123",
								Source:      "wotc",
								PublishedAt: "2004-10-04",
								Comment:     "This spell can target players, creatures, or planeswalkers.",
							},
						}
					}
					return nil
				},
			},
			want: []models.Ruling{
				{
					Object:      "ruling",
					OracleID:    "oracle-123",
					Source:      "scryfall",
					PublishedAt: "2004-10-04",
					Comment:     "Lightning Bolt deals 3 damage to any target.",
				},
				{
					Object:      "ruling",
					OracleID:    "oracle-123",
					Source:      "wotc",
					PublishedAt: "2004-10-04",
					Comment:     "This spell can target players, creatures, or planeswalkers.",
				},
			},
			wantErr: false,
		},
		{
			name:   "empty card id",
			cardID: "",
			mockClient: &MockHTTPClient{
				DoFunc: func(req *http.Request, v interface{}) error {
					if rulings, ok := v.(*[]models.Ruling); ok {
						*rulings = []models.Ruling{}
					}
					return nil
				},
			},
			want:    []models.Ruling{},
			wantErr: false,
		},
		{
			name:   "card not found",
			cardID: "invalid-card-id",
			mockClient: &MockHTTPClient{
				DoFunc: func(req *http.Request, v interface{}) error {
					return errors.New("card not found")
				},
			},
			want:    nil,
			wantErr: true,
			errMsg:  "card not found",
		},
		{
			name:   "new request error",
			cardID: "f2f0e0e6-8eee-4001-a6a8-2c8f1b9a3c3e",
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
			name:   "no rulings available",
			cardID: "f2f0e0e6-8eee-4001-a6a8-2c8f1b9a3c3e",
			mockClient: &MockHTTPClient{
				DoFunc: func(req *http.Request, v interface{}) error {
					if rulings, ok := v.(*[]models.Ruling); ok {
						*rulings = []models.Ruling{}
					}
					return nil
				},
			},
			want:    []models.Ruling{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &RulingService{
				Client: tt.mockClient,
			}

			got, err := s.GetRulings(context.Background(), tt.cardID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("RulingService.GetRulings() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errMsg != "" && !contains(err.Error(), tt.errMsg) {
					t.Errorf("RulingService.GetRulings() error = %v, want error containing %v", err, tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("RulingService.GetRulings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !equalRulings(got, tt.want) {
				t.Errorf("RulingService.GetRulings() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test that the correct path is being constructed
func TestRulingService_GetRulings_PathConstruction(t *testing.T) {
	cardID := "f2f0e0e6-8eee-4001-a6a8-2c8f1b9a3c3e"
	expectedPath := "/cards/" + cardID + "/rulings"

	mockClient := &MockHTTPClient{
		NewRequestFunc: func(ctx context.Context, method, path string) (*http.Request, error) {
			if path != expectedPath {
				t.Errorf("Expected path %s, got %s", expectedPath, path)
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

	s := &RulingService{
		Client: mockClient,
	}

	_, err := s.GetRulings(context.Background(), cardID)
	if err != nil {
		t.Errorf("RulingService.GetRulings() unexpected error: %v", err)
	}
}

// Helper function for comparing rulings
func equalRulings(a, b []models.Ruling) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].Object != b[i].Object ||
			a[i].OracleID != b[i].OracleID ||
			a[i].Source != b[i].Source ||
			a[i].PublishedAt != b[i].PublishedAt ||
			a[i].Comment != b[i].Comment {
			return false
		}
	}
	return true
}