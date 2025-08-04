package cards

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type AutocompleteCardParams struct {
	Query         string  `json:"q" validate:"required"`
	Format        *string `json:"format,omitempty" validate:"omitempty,oneof=json"`
	Pretty        *bool   `json:"pretty,omitempty"`
	IncludeExtras *bool   `json:"include_extras,omitempty"`
}

func (p *AutocompleteCardParams) Validate() error {
	if strings.TrimSpace(p.Query) == "" {
		return fmt.Errorf("query cannot be empty")
	}
	if p.Format != nil && *p.Format != "json" {
		return fmt.Errorf("invalid format: must be 'json'")
	}
	if p.Pretty != nil && !*p.Pretty && *p.Pretty {
		return fmt.Errorf("pretty must be a boolean value")
	}
	if p.IncludeExtras != nil && !*p.IncludeExtras && *p.IncludeExtras {
		return fmt.Errorf("include_extras must be a boolean value")
	}
	return nil
}

func (p *AutocompleteCardParams) ToURLValues() (url.Values, error) {
	if err := p.Validate(); err != nil {
		return nil, err
	}
	params := url.Values{}
	params.Add("q", p.Query)
	if p.Format != nil {
		params.Add("format", *p.Format)
	}
	if p.Pretty != nil {
		params.Add("pretty", strconv.FormatBool(*p.Pretty))
	}
	if p.IncludeExtras != nil {
		params.Add("include_extras", strconv.FormatBool(*p.IncludeExtras))
	}
	return params, nil
}

func NewAutocompleteCardParams(query string) *AutocompleteCardParams {
	return &AutocompleteCardParams{
		Query: query,
	}
}

func (p *AutocompleteCardParams) WithPretty(pretty bool) *AutocompleteCardParams {
	p.Pretty = &pretty
	return p
}

func (p *AutocompleteCardParams) WithIncludeExtras(includeExtras bool) *AutocompleteCardParams {
	p.IncludeExtras = &includeExtras
	return p
}
