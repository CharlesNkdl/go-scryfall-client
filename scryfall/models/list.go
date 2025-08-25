package models

// List represents a paginated list response from the Scryfall API.
type List[T any] struct {
	Object     string   `json:"object"`
	Data       []T      `json:"data"`
	HasMore    bool     `json:"has_more"`
	NextPage   *string  `json:"next_page,omitempty"`
	TotalCards *int     `json:"total_cards,omitempty"`
	Warnings   []string `json:"warnings,omitempty"`
}
