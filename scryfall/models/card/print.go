package card

type Print struct {
	Artist           *string       `json:"artist,omitempty"`
	ArtistIds        []string      `json:"artist_ids,omitempty"`
	AttractionLights []string      `json:"attraction_lights,omitempty"`
	Booster          bool          `json:"booster"`
	BorderColor      BorderColor   `json:"border_color"`
	CardBackId       *string       `json:"card_back_id,omitempty"`
	CollectorNumber  string        `json:"collector_number"`
	ContentWarning   *bool         `json:"content_warning,omitempty"`
	Digital          bool          `json:"digital"`
	Finishes         []Finish      `json:"finishes,omitempty"`
	FlavorName       *string       `json:"flavor_name,omitempty"`
	FlavorText       *string       `json:"flavor_text,omitempty"`
	FrameEffects     []FrameEffect `json:"frame_effects,omitempty"`
	Frame            Frame         `json:"frame"`
	FullArt          bool          `json:"full_art"`
	Games            []Game        `json:"games,omitempty"`
}
