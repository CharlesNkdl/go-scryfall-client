package models

type Ruling struct {
	Object       string `json:"object"`
	Oracle_id    string `json:"oracle_id"`
	Source       string `json:"source"`
	Published_at string `json:"published_at"`
	Comment      string `json:"comment"`
}
