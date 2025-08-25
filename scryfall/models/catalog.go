package models

// Catalog represents a catalog of values from the Scryfall API.
type Catalog struct {
	Object      string   `json:"object"`
	URI         string   `json:"uri"`
	TotalValues int      `json:"total_values"`
	Data        []string `json:"data"`
}
