package card

// Component represents the component type of a Magic card.
type Component string

// Component types for Magic cards.
const (
	ComponentToken      Component = "token"
	ComponentMeldPart   Component = "meld_part"
	ComponentMeldResult Component = "meld_result"
	ComponentComboPiece Component = "combo_piece"
)
