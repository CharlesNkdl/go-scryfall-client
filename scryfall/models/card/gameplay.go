package card

import (
	"github.com/cnkdl/go-scryfall-client/scryfall/models"
)

type Gameplay struct {
	AllParts       []Related      `json:"all_parts"`
	CardFaces      []Face         `json:"card_faces,omitempty"`
	Cmc            *float32       `json:"cmc,omitempty"`
	ColorIdentity  []models.Color `json:"color_identity,omitempty"`
	ColorIndicator []models.Color `json:"color_indicator,omitempty"`
	Colors         []models.Color `json:"colors,omitempty"`
	Defense        *string        `json:"defense,omitempty"`
	EdhrecRank     *int           `json:"edhrec_rank,omitempty"`
	GameChanger    *bool          `json:"game_changer,omitempty"`
	HandModifier   *string        `json:"hand_modifier,omitempty"`
	Keywords       []string       `json:"keywords"`
	Legalities     []Legality     `json:"legalities"`
	LifeModifier   *string        `json:"life_modifier,omitempty"`
	Loyalty        *string        `json:"loyalty,omitempty"`
	ManaCost       *string        `json:"mana_cost,omitempty"`
	Name           string         `json:"name"`
	OracleText     *string        `json:"oracle_text,omitempty"`
	PennyRank      *int           `json:"penny_rank,omitempty"`
	Power          *string        `json:"power,omitempty"`
	ProducedMana   []models.Color `json:"produced_mana,omitempty"`
	Reserved       bool           `json:"reserved"`
	Toughness      *string        `json:"toughness,omitempty"`
	TypeLine       string         `json:"type_line"`
}
