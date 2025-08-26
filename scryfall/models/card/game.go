package card

// Game represents the game format where a Magic card is available.
type Game string

// Game formats.
const (
	GamePaper Game = "paper"
	GameArena Game = "arena"
	GameMTGO  Game = "mtgo"
)
