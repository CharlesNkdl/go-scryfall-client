package services

import (
	"context"
	"fmt"
	"net/url"

	"github.com/CharlesNkdl/go-scryfall-client/scryfall/models"
)

type SymbolService struct {
	Client HTTPClient
}

func (s *SymbolService) GetSymbol(ctx context.Context, symbol string) (*models.Symbol, error) {
	path := fmt.Sprintf("/symbology/%s", url.QueryEscape(symbol))
	req, err := s.Client.NewRequest(ctx, "GET", path)
	if err != nil {
		return nil, err
	}
	var symbolResult models.Symbol
	if err := s.Client.Do(req, &symbolResult); err != nil {
		return nil, err
	}
	return &symbolResult, nil
}
