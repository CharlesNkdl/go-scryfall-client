package errors

import (
	"fmt"
	"github.com/CharlesNkdl/go-scryfall-client/scryfall/models"
)

type ApiError struct {
	ErrInfo models.ScryfallError
}

func (e *ApiError) Error() string {
	return fmt.Sprintf("Scryfall API Error: %s (status : %d , code : %s)",
		e.ErrInfo.Detail, e.ErrInfo.Status, e.ErrInfo.Code)
}
