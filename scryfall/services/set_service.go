package services

import (
	"context"
	"fmt"

	"github.com/cnkdl/go-scryfall-client/scryfall/models/set"
)

type SetService struct {
	Client HTTPClient
}

func (s *SetService) GetById(ctx context.Context, id string) (*set.Set, error) {
	path := fmt.Sprintf("/sets/%s", id)
	req, err := s.Client.NewRequest(ctx, "GET", path)
	if err != nil {
		return nil, err
	}

	var setResult set.Set
	if err := s.Client.Do(req, &setResult); err != nil {
		return nil, err
	}
	return &setResult, nil
}
