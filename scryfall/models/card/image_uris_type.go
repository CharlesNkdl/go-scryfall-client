package card

type ImageUrisType string

const (
	ImageUrisPng        ImageUrisType = "png"
	ImageUrisBorderCrop ImageUrisType = "border_crop"
	ImageUrisArtCrop    ImageUrisType = "art_crop"
	ImageUrisLarge      ImageUrisType = "large"
	ImageUrisNormal     ImageUrisType = "normal"
	ImageUrisSmall      ImageUrisType = "small"
)
