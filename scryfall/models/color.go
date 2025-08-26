package models

// Color represents a Magic card color from the Scryfall API.
type Color string

// Magic card colors.
const (
	ColorWhite Color = "W"
	ColorBlue  Color = "U"
	ColorBlack Color = "B"
	ColorRed   Color = "R"
	ColorGreen Color = "G"
)
