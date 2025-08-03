package card

type ImageUris struct {
	Type     *ImageUrisType `json:"type,omitempty"`
	ImageUri *string        `json:"image_uri,omitempty"`
}
