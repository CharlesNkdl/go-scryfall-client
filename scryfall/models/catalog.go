package models

type Catalog struct {
	Object      string   `json:"object"`
	Uri         string   `json:"uri"`
	TotalValues int      `json:"total_values"`
	Data        []string `json:"data"`
}
