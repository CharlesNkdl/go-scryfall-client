package services

import (
	"context"
	"fmt"

	"github.com/cnkdl/go-scryfall-client/scryfall/models"
)

type CatalogService struct {
	Client HTTPClient
}

func (s *CatalogService) GetCatalog(ctx context.Context, catalogType string) (*models.Catalog, error) {
	path := fmt.Sprintf("/catalog/%s", catalogType)
	req, err := s.Client.NewRequest(ctx, "GET", path)
	if err != nil {
		return nil, err
	}

	var catalogResult models.Catalog
	if err := s.Client.Do(req, &catalogResult); err != nil {
		return nil, err
	}
	return &catalogResult, nil
}
