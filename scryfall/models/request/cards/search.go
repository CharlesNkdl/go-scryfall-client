package cards

import (
	"errors"
	"net/url"
	"strconv"
	"strings"
)

type UniqueMode string

const (
	UniqueCards  UniqueMode = "cards"
	UniqueArt    UniqueMode = "art"
	UniquePrints UniqueMode = "prints"
)

type OrderMode string

const (
	OrderName      OrderMode = "name"
	OrderSet       OrderMode = "set"
	OrderReleased  OrderMode = "released"
	OrderRarity    OrderMode = "rarity"
	OrderColor     OrderMode = "color"
	OrderUSD       OrderMode = "usd"
	OrderTix       OrderMode = "tix"
	OrderEUR       OrderMode = "eur"
	OrderCMC       OrderMode = "cmc"
	OrderPower     OrderMode = "power"
	OrderToughness OrderMode = "toughness"
	OrderEDHRec    OrderMode = "edhrec"
	OrderPenny     OrderMode = "penny"
	OrderArtist    OrderMode = "artist"
	OrderReview    OrderMode = "review"
)

type Direction string

const (
	DirAuto Direction = "auto"
	DirAsc  Direction = "asc"
	DirDesc Direction = "desc"
)

type Format string

const (
	FormatJSON Format = "json"
	FormatCSV  Format = "csv"
)

type SearchParams struct {
	Query string `json:"q"`

	Unique              UniqueMode `json:"unique,omitempty"`
	IncludeExtras       bool       `json:"include_extras,omitempty"`
	IncludeMultilingual bool       `json:"include_multilingual,omitempty"`
	IncludeVariations   bool       `json:"include_variations,omitempty"`

	Order OrderMode `json:"order,omitempty"`
	Dir   Direction `json:"dir,omitempty"`

	Page int `json:"page,omitempty"`

	Format Format `json:"format,omitempty"`
	Pretty bool   `json:"pretty,omitempty"`
}

func (p *SearchParams) Validate() error {
	if strings.TrimSpace(p.Query) == "" {
		return errors.New("query parameter 'q' is required")
	}

	if len([]rune(p.Query)) > 1000 {
		return errors.New("query parameter 'q' must not exceed 1000 Unicode characters")
	}

	if p.Unique != "" {
		switch p.Unique {
		case UniqueCards, UniqueArt, UniquePrints:
		default:
			return errors.New("invalid unique mode: must be 'cards', 'art', or 'prints'")
		}
	}

	if p.Order != "" {
		switch p.Order {
		case OrderName, OrderSet, OrderReleased, OrderRarity, OrderColor,
			OrderUSD, OrderTix, OrderEUR, OrderCMC, OrderPower,
			OrderToughness, OrderEDHRec, OrderPenny, OrderArtist, OrderReview:
		default:
			return errors.New("invalid order mode")
		}
	}

	if p.Dir != "" {
		switch p.Dir {
		case DirAuto, DirAsc, DirDesc:
		default:
			return errors.New("invalid direction: must be 'auto', 'asc', or 'desc'")
		}
	}

	if p.Format != "" {
		switch p.Format {
		case FormatJSON, FormatCSV:
		default:
			return errors.New("invalid format: must be 'json' or 'csv'")
		}
	}

	if p.Page < 0 {
		return errors.New("page must be positive")
	}

	return nil
}

func (p *SearchParams) ToURLValues() url.Values {
	values := url.Values{}

	values.Set("q", p.Query)

	if p.Unique != "" {
		values.Set("unique", string(p.Unique))
	}

	if p.Order != "" {
		values.Set("order", string(p.Order))
	}

	if p.Dir != "" {
		values.Set("dir", string(p.Dir))
	}

	if p.IncludeExtras {
		values.Set("include_extras", "true")
	}

	if p.IncludeMultilingual {
		values.Set("include_multilingual", "true")
	}

	if p.IncludeVariations {
		values.Set("include_variations", "true")
	}

	if p.Page > 0 {
		values.Set("page", strconv.Itoa(p.Page))
	}

	if p.Format != "" {
		values.Set("format", string(p.Format))
	}

	if p.Pretty {
		values.Set("pretty", "true")
	}

	return values
}

func NewSearchParams(query string) *SearchParams {
	return &SearchParams{
		Query:  query,
		Unique: UniqueCards,
		Order:  OrderName,
		Dir:    DirAuto,
		Format: FormatJSON,
		Page:   1,
	}
}

func (p *SearchParams) WithUnique(unique UniqueMode) *SearchParams {
	p.Unique = unique
	return p
}

func (p *SearchParams) WithOrder(order OrderMode, dir Direction) *SearchParams {
	p.Order = order
	p.Dir = dir
	return p
}

func (p *SearchParams) WithIncludes(extras, multilingual, variations bool) *SearchParams {
	p.IncludeExtras = extras
	p.IncludeMultilingual = multilingual
	p.IncludeVariations = variations
	return p
}

func (p *SearchParams) WithPage(page int) *SearchParams {
	p.Page = page
	return p
}

func (p *SearchParams) WithFormat(format Format, pretty bool) *SearchParams {
	p.Format = format
	p.Pretty = pretty
	return p
}
