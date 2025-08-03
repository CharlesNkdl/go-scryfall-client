package card

type CardCore struct {
	ArenaId           *int     `json:"arena_id,omitempty"`
	Id                string   `json:"id"`
	Lang              string   `json:"lang"`
	MTGOId            *int     `json:"mtgo_id,omitempty"`
	MTGOFoilId        *int     `json:"mtgo_foil_id,omitempty"`
	MultiverseIds     []int    `json:"multiverse_ids,omitempty"`
	TcgPlayerId       *int     `json:"tcgplayer_id,omitempty"`
	TcgPlayerEtchedId *int     `json:"tcgplayer_etched_id,omitempty"`
	CardmarketId      *int     `json:"cardmarket_id,omitempty"`
	Object            string   `json:"object"`
	Layout            Layout   `json:"layout"`
	OracleId          *string  `json:"oracle_id,omitempty"`
	PrintsSearchUri   string   `json:"prints_search_uri"`
	RulingsUri        string   `json:"rulings_uri"`
	ScryfallUri       string   `json:"scryfall_uri"`
	Uri               string   `json:"uri"`
	GameplayFields    Gameplay `json:"gameplay_fields"`
	PrintFields       Print    `json:"print_fields"`
	CardFaces         []Face   `json:"card_faces,omitempty"`
}
