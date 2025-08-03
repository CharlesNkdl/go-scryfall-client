package card

type ImageStatus string

const (
	ImageStatusHighRes     ImageStatus = "highres_scan"
	ImageStatusLowRes      ImageStatus = "lowres"
	ImageStatusMissing     ImageStatus = "missing"
	ImageStatusPlaceholder ImageStatus = "placeholder"
)
