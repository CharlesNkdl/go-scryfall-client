package card

type Print struct {
	Artist           *string        `json:"artist,omitempty"`
	ArtistIds        []string       `json:"artist_ids,omitempty"`
	AttractionLights []string       `json:"attraction_lights,omitempty"`
	IsInBooster      bool           `json:"booster"`
	BorderColor      BorderColor    `json:"border_color"`
	CardBackId       *string        `json:"card_back_id,omitempty"`
	CollectorNumber  string         `json:"collector_number"`
	IsContentWarning *bool          `json:"content_warning,omitempty"`
	IsDigital        bool           `json:"digital"`
	Finishes         []Finish       `json:"finishes,omitempty"`
	FlavorName       *string        `json:"flavor_name,omitempty"`
	FlavorText       *string        `json:"flavor_text,omitempty"`
	FrameEffects     []FrameEffect  `json:"frame_effects,omitempty"`
	Frame            Frame          `json:"frame"`
	IsFullArt        bool           `json:"full_art"`
	Games            []Game         `json:"games"`
	IsHighresImage   bool           `json:"highres_image"`
	IllustrationId   *string        `json:"illustration_id,omitempty"`
	ImageStatus      ImageStatus    `json:"image_status"`
	ImageUris        []ImageUris    `json:"image_uris,omitempty"`
	IsOversized      bool           `json:"oversized"`
	Prices           []Price        `json:"prices,omitempty"`
	Rarity           Rarity         `json:"rarity"`
	RelatedUris      []string       `json:"related_uris,omitempty"`
	ReleasedAt       string         `json:"released_at"`
	IsReprint        bool           `json:"is_reprint"`
	ScryfallSetUri   string         `json:"scryfall_set_uri"`
	SetName          string         `json:"set_name"`
	SetSearchUri     string         `json:"set_search_uri"`
	SetType          string         `json:"set_type"`
	SetUri           string         `json:"set_uri"`
	Set              string         `json:"set"`
	SetId            string         `json:"set_id"`
	IsStorySpotlight bool           `json:"story_spotlight"`
	IsVariation      bool           `json:"variation"`
	VariationOf      *string        `json:"variation_of,omitempty"`
	SecurityStamp    *SecurityStamp `json:"security_stamp,omitempty"`
	Watermark        *string        `json:"watermark,omitempty"`
	Preview          []string       `json:"preview,omitempty"`
}
