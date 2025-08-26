package models

// Ruling represents a ruling on a Magic card from the Scryfall API.
type Ruling struct {
	Object      string `json:"object"`
	OracleID    string `json:"oracle_id"`
	Source      string `json:"source"`
	PublishedAt string `json:"published_at"`
	Comment     string `json:"comment"`
}
