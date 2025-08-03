package card

type Price struct {
	Type  PriceType `json:"type"`
	Value *string   `json:"value,omitempty"`
}
