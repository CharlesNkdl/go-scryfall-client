package card

type PriceType string

const (
	PriceTypeUSD       PriceType = "usd"
	PriceTypeUSDFoil   PriceType = "usd_foil"
	PriceTypeUSDEtched PriceType = "usd_etched"
	PriceTypeEUR       PriceType = "eur"
	PriceTypeEURFoil   PriceType = "eur_foil"
	PriceTypeEUREtched PriceType = "eur_etched"
	PriceTypeTix       PriceType = "tix"
)
