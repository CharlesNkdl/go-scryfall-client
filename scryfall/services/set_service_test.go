package services

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/cnkdl/go-scryfall-client/scryfall/models/set"
)

func TestSetService_GetById(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		mockClient *MockHTTPClient
		want       *set.Set
		wantErr    bool
		errMsg     string
	}{
		{
			name: "successful get by id",
			id:   "khm",
			mockClient: &MockHTTPClient{
				DoFunc: func(req *http.Request, v interface{}) error {
					if setResult, ok := v.(*set.Set); ok {
						setResult.Id = "khm"
						setResult.Object = "set"
						setResult.Code = "khm"
						setResult.Name = "Kaldheim"
					}
					return nil
				},
			},
			want: &set.Set{
				Id:     "khm",
				Object: "set",
				Code:   "khm",
				Name:   "Kaldheim",
			},
			wantErr: false,
		},
		{
			name: "empty id",
			id:   "",
			mockClient: &MockHTTPClient{
				DoFunc: func(req *http.Request, v interface{}) error {
					return nil
				},
			},
			want: &set.Set{},
			wantErr: false,
		},
		{
			name: "http client error",
			id:   "invalid",
			mockClient: &MockHTTPClient{
				DoFunc: func(req *http.Request, v interface{}) error {
					return errors.New("set not found")
				},
			},
			want:    nil,
			wantErr: true,
			errMsg:  "set not found",
		},
		{
			name: "new request error",
			id:   "khm",
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
			s := &SetService{
				Client: tt.mockClient,
			}

			got, err := s.GetById(context.Background(), tt.id)

			if tt.wantErr {
				if err == nil {
					t.Errorf("SetService.GetById() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errMsg != "" && !contains(err.Error(), tt.errMsg) {
					t.Errorf("SetService.GetById() error = %v, want error containing %v", err, tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("SetService.GetById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !equalSets(got, tt.want) {
				t.Errorf("SetService.GetById() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test that the correct path is being constructed
func TestSetService_GetById_PathConstruction(t *testing.T) {
	mockClient := &MockHTTPClient{
		NewRequestFunc: func(ctx context.Context, method, path string) (*http.Request, error) {
			expectedPath := "/sets/khm"
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

	s := &SetService{
		Client: mockClient,
	}

	_, err := s.GetById(context.Background(), "khm")
	if err != nil {
		t.Errorf("SetService.GetById() unexpected error: %v", err)
	}
}

// Helper function for comparing sets
func equalSets(a, b *set.Set) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.Id == b.Id && a.Object == b.Object && a.Code == b.Code && a.Name == b.Name
}