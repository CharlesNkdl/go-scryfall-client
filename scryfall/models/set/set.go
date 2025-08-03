package set

type Set struct {
	Object        string  `json:"object"`
	Id            string  `json:"id"`
	Code          string  `json:"code"`
	MtgoCode      *string `json:"mtgo_code,omitempty"`
	ArenaCode     *string `json:"arena_code,omitempty"`
	TcgPlayerId   *int    `json:"tcgplayer_id,omitempty"`
	Name          string  `json:"name"`
	SetType       SetType `json:"set_type"`
	ReleasedAt    *string `json:"released_at,omitempty"`
	BlockCode     *string `json:"block_code,omitempty"`
	Block         *string `json:"block,omitempty"`
	ParentSetCode *string `json:"parent_set_code,omitempty"`
	CardCount     int     `json:"card_count"`
	PrintedSize   *int    `json:"printed_size,omitempty"`
	Digital       bool    `json:"digital"`
	FoilOnly      bool    `json:"foil_only"`
	NonFoilOnly   bool    `json:"nonfoil_only"`
	ScryfallUri   string  `json:"scryfall_uri"`
	Uri           string  `json:"uri"`
	IconSvgUri    string  `json:"icon_svg_uri"`
	SearchUri     string  `json:"search_uri"`
}
