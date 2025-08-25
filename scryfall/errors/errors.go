package errors

import (
	"fmt"
	"github.com/cnkdl/go-scryfall-client/scryfall/models"
)

// APIError represents an error response from the Scryfall API.
type APIError struct {
	ErrInfo models.ScryfallError
}

func (e *APIError) Error() string {
	return fmt.Sprintf("Scryfall API Error: %s (status : %d , code : %s)",
		e.ErrInfo.Detail, e.ErrInfo.Status, e.ErrInfo.Code)
}
