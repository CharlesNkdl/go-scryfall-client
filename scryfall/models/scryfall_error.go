package models

// ScryfallError represents an error response structure from the Scryfall API.
type ScryfallError struct {
	Status   int      `json:"status"`
	Code     string   `json:"code"`
	Detail   string   `json:"detail"`
	Type     *string  `json:"type,omitempty"`
	Warnings []string `json:"warnings,omitempty"`
}
