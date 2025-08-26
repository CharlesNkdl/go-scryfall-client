package card

// ImageUrisType represents different image URI types for card images.
type ImageUrisType string

// Image URI types for different versions of card images.
const (
	ImageUrisPng        ImageUrisType = "png"
	ImageUrisBorderCrop ImageUrisType = "border_crop"
	ImageUrisArtCrop    ImageUrisType = "art_crop"
	ImageUrisLarge      ImageUrisType = "large"
	ImageUrisNormal     ImageUrisType = "normal"
	ImageUrisSmall      ImageUrisType = "small"
)
