package services

import (
	"context"
	"fmt"
	"net/url"

	"github.com/cnkdl/go-scryfall-client/scryfall/models"
	"github.com/cnkdl/go-scryfall-client/scryfall/models/card"
	cardreq "github.com/cnkdl/go-scryfall-client/scryfall/models/request/cards"
)

type CardService struct {
	Client HTTPClient
}

func (s *CardService) GetById(ctx context.Context, params *cardreq.IdCardParams) (*card.CardCore, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("invalid named parameters: %w", err)
	}
	urlValues, err := params.ToURLValues()
	if err != nil {
		return nil, fmt.Errorf("failed to get URL values: %w", err)
	}
	path := fmt.Sprintf("/cards/%s?%s", url.QueryEscape(params.Id), urlValues.Encode())
	fmt.Println("Requesting card by ID:", path)
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

func (s *CardService) GetByName(ctx context.Context, params *cardreq.NamedCardParams) (*card.CardCore, error) {
	urlValues, err := params.ToURLValues()
	if err != nil {
		return nil, fmt.Errorf("failed to get URL values: %w", err)
	}
	path := fmt.Sprintf("/cards/named?%s", urlValues.Encode())
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

func (s *CardService) Search(ctx context.Context, params *cardreq.SearchParams) (*models.List[card.CardCore], error) {
	urlValues, err := params.ToURLValues()
	if err != nil {
		return nil, fmt.Errorf("failed to get URL values: %w", err)
	}
	path := fmt.Sprintf("/cards/search?%s", urlValues.Encode())
	req, err := s.Client.NewRequest(ctx, "GET", path)
	if err != nil {
		return nil, err
	}
	var searchResult models.List[card.CardCore]
	if err := s.Client.Do(req, &searchResult); err != nil {
		return nil, err
	}
	return &searchResult, nil
}

func (s *CardService) Autocomplete(ctx context.Context, params *cardreq.AutocompleteCardParams) (*models.Catalog, error) {
	urlValues, err := params.ToURLValues()
	if err != nil {
		return nil, fmt.Errorf("failed to get URL values: %w", err)
	}
	path := fmt.Sprintf("/cards/autocomplete?%s", urlValues.Encode())
	req, err := s.Client.NewRequest(ctx, "GET", path)
	if err != nil {
		return nil, err
	}
	var autocompleteResult models.Catalog
	if err := s.Client.Do(req, &autocompleteResult); err != nil {
		return nil, err
	}
	return &autocompleteResult, nil
}

func (s *CardService) GetRandom(ctx context.Context, params *cardreq.RandomCardParams) (*card.CardCore, error) {
	urlValues, err := params.ToURLValues()
	if err != nil {
		return nil, fmt.Errorf("failed to get URL values: %w", err)
	}
	path := fmt.Sprintf("/cards/random?%s", urlValues.Encode())
	req, err := s.Client.NewRequest(ctx, "GET", path)
	if err != nil {
		return nil, err
	}
	var randomCardResult card.CardCore
	if err := s.Client.Do(req, &randomCardResult); err != nil {
		return nil, err
	}
	return &randomCardResult, nil
}

func (s *CardService) GetByCodeNumberLang(ctx context.Context, params *cardreq.CardsByCodeNumberLangParams) (*models.List[card.CardCore], error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("invalid named parameters: %w", err)
	}
	path := params.BuildPath()
	req, err := s.Client.NewRequest(ctx, "GET", path)
	if err != nil {
		return nil, err
	}
	var setCardsResult models.List[card.CardCore]
	if err := s.Client.Do(req, &setCardsResult); err != nil {
		return nil, err
	}
	return &setCardsResult, nil
}
