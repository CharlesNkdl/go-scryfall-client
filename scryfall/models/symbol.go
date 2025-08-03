package models

type Symbol struct {
	Object             string   `json:"object"`
	Symbol             string   `json:"symbol"`
	LooseVariant       *string  `json:"loose_variant,omitempty"`
	EnglishName        string   `json:"english"`
	Transposable       bool     `json:"transposable"`
	RepresentsMana     bool     `json:"represents_mana"`
	ManaValue          *float32 `json:"mana_value,omitempty"`
	AppearsInManaCosts bool     `json:"appears_in_mana_costs"`
	Funny              bool     `json:"funny"`
	Colors             []Color  `json:"colors"`
	Hybrid             bool     `json:"hybrid"`
	Phyrexian          bool     `json:"phyrexian"`
	GathererAlternates []string `json:"gatherer_alternates,omitempty"`
	SvgUri             *string  `json:"svg_uri,omitempty"`
}
