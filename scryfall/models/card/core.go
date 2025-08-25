package card

// CardCore represents the core card data structure from the Scryfall API.
type CardCore struct {
	ArenaID           *int     `json:"arena_id,omitempty"`
	ID                string   `json:"id"`
	Lang              string   `json:"lang"`
	MTGOID            *int     `json:"mtgo_id,omitempty"`
	MTGOFoilID        *int     `json:"mtgo_foil_id,omitempty"`
	MultiverseIDs     []int    `json:"multiverse_ids,omitempty"`
	TcgPlayerID       *int     `json:"tcgplayer_id,omitempty"`
	TcgPlayerEtchedID *int     `json:"tcgplayer_etched_id,omitempty"`
	CardmarketID      *int     `json:"cardmarket_id,omitempty"`
	Object            string   `json:"object"`
	Layout            Layout   `json:"layout"`
	OracleID          *string  `json:"oracle_id,omitempty"`
	PrintsSearchURI   string   `json:"prints_search_uri"`
	RulingsURI        string   `json:"rulings_uri"`
	ScryfallURI       string   `json:"scryfall_uri"`
	URI               string   `json:"uri"`
	GameplayFields    Gameplay `json:"gameplay_fields"`
	PrintFields       Print    `json:"print_fields"`
	CardFaces         []Face   `json:"card_faces,omitempty"`
}
