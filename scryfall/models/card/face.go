package card

type Face struct {
	Artist          *string    `json:"artist,omitempty"`
	ArtistId        *string    `json:"artist_id,omitempty"`
	Cmc             *float32   `json:"cmc,omitempty"`
	ColorIndicator  []Color    `json:"color_indicator,omitempty"`
	Colors          []Color    `json:"colors,omitempty"`
	Defense         *string    `json:"defense,omitempty"`
	FlavorText      *string    `json:"flavor_text,omitempty"`
	IllustrationId  *string    `json:"illustration_id,omitempty"`
	ImageUris       *ImageUris `json:"image_uris,omitempty"`
	Layout          *Layout    `json:"layout,omitempty"`
	Loyalty         *string    `json:"loyalty,omitempty"`
	ManaCost        string     `json:"mana_cost"`
	Name            string     `json:"name"`
	Object          string     `json:"object"`
	OracleId        *string    `json:"oracle_id,omitempty"`
	OracleText      *string    `json:"oracle_text,omitempty"`
	Power           *string    `json:"power,omitempty"`
	PrintedName     *string    `json:"printed_name,omitempty"`
	PrintedText     *string    `json:"printed_text,omitempty"`
	PrintedTypeLine *string    `json:"printed_type_line,omitempty"`
	Toughness       *string    `json:"toughness,omitempty"`
	TypeLine        *string    `json:"type_line,omitempty"`
	Watermark       *string    `json:"watermark,omitempty"`
}
