package services

import (
	"context"
	"fmt"
	"net/url"

	"github.com/cnkdl/go-scryfall-client/scryfall/models/card"
)

type CardService struct {
	Client HTTPClient
}

func (s *CardService) GetById(ctx context.Context, id string) (*card.CardCore, error) {
	path := fmt.Sprintf("/cards/%s", id)
	req, err := s.Client.NewRequest(ctx, "GET", path)
	if err != nil {
		return nil, err
	}

	var cardResult card.CardCore
	if err := s.Client.Do(req, &cardResult); err != nil {
		return nil, err
	}
	return &cardResult, nil
}

func (s *CardService) GetByName(ctx context.Context, name string, fuzzy bool) (*card.CardCore, error) {
	params := url.Values{}
	if fuzzy {
		params.Set("fuzzy", name)
	} else {
		params.Set("exact", name)
	}
	path := fmt.Sprintf("/cards/named?%s", params.Encode())

	req, err := s.Client.NewRequest(ctx, "GET", path)
	if err != nil {
		return nil, err
	}

	var cardResult card.CardCore
	if err := s.Client.Do(req, &cardResult); err != nil {
		return nil, err
	}
	return &cardResult, nil
}
