package card

type Legality string

const (
	LegalityLegal      Legality = "legal"
	LegalityNotLegal   Legality = "not_legal"
	LegalityRestricted Legality = "restricted"
	LegalityBanned     Legality = "banned"
)
