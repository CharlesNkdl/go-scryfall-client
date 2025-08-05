package services

import (
	"context"
	"fmt"

	"github.com/CharlesNkdl/go-scryfall-client/scryfall/models"
)

type RulingService struct {
	Client HTTPClient
}

func (s *RulingService) GetRulings(ctx context.Context, cardID string) ([]models.Ruling, error) {
	path := fmt.Sprintf("/cards/%s/rulings", cardID)
	req, err := s.Client.NewRequest(ctx, "GET", path)
	if err != nil {
		return nil, err
	}

	var rulingsResult []models.Ruling
	if err := s.Client.Do(req, &rulingsResult); err != nil {
		return nil, err
	}
	return rulingsResult, nil
}
