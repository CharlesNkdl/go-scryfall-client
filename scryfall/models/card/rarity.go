package card

type Rarity string

const (
	RarityCommon   Rarity = "common"
	RarityUncommon Rarity = "uncommon"
	RarityRare     Rarity = "rare"
	RaritySpecial  Rarity = "special"
	RarityMythic   Rarity = "mythic"
	RarityBonus    Rarity = "bonus"
)
