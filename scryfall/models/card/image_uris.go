package card

// ImageUris represents image URIs for a Magic card from the Scryfall API.
type ImageUris struct {
	Type     *ImageUrisType `json:"type,omitempty"`
	ImageURI *string        `json:"image_uri,omitempty"`
}
