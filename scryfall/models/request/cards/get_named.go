package cards

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type NamedCardParams struct {
	Exact *string `json:"exact,omitempty" validate:"omitempty"`
	Fuzzy *string `json:"fuzzy,omitempty" validate:"omitempty"`

	Set     *string `json:"set,omitempty" validate:"omitempty,min=3,max=5"`
	Format  *string `json:"format,omitempty" validate:"omitempty,oneof=json text image"`
	Face    *string `json:"face,omitempty" validate:"omitempty,oneof=front back"`
	Version *string `json:"version,omitempty" validate:"omitempty,oneof=small normal large png art_crop border_crop"`
	Pretty  *bool   `json:"pretty,omitempty"`
}

func (p *NamedCardParams) Validate() error {
	if p.Exact == nil && p.Fuzzy == nil {
		return fmt.Errorf("either 'exact' or 'fuzzy' parameter is required")
	}
	if p.Exact != nil && p.Fuzzy != nil {
		return fmt.Errorf("'exact' and 'fuzzy' parameters are mutually exclusive")
	}
	if p.Exact != nil && strings.TrimSpace(*p.Exact) == "" {
		return fmt.Errorf("exact name cannot be empty")
	}

	if p.Fuzzy != nil && strings.TrimSpace(*p.Fuzzy) == "" {
		return fmt.Errorf("fuzzy name cannot be empty")
	}
	if p.Set != nil {
		setLen := len(strings.TrimSpace(*p.Set))
		if setLen < 3 || setLen > 5 {
			return fmt.Errorf("set code must be between 3 and 5 characters")
		}
	}
	if p.Format != nil {
		switch *p.Format {
		case "json", "text", "image":
		default:
			return fmt.Errorf("invalid format: must be 'json', 'text', or 'image'")
		}
	}
	if p.Version != nil {
		switch *p.Version {
		case "small", "normal", "large", "png", "art_crop", "border_crop":
		default:
			return fmt.Errorf("invalid version: must be 'small', 'normal', 'large', 'png', 'art_crop', or 'border_crop'")
		}
	}
	if p.Face != nil {
		switch *p.Face {
		case "front", "back":
		default:
			return fmt.Errorf("invalid face: must be 'front' or 'back'")
		}
	}
	if p.Face != nil && *p.Face == "back" {
		if p.Format == nil || *p.Format != "image" {
			return fmt.Errorf("face=back can only be used with format=image")
		}
	}
	return nil
}

func (p *NamedCardParams) ToURLValues() (url.Values, error) {
	if err := p.Validate(); err != nil {
		return nil, err
	}
	params := url.Values{}
	if p.Exact != nil {
		params.Add("exact", *p.Exact)
	}
	if p.Fuzzy != nil {
		params.Add("fuzzy", *p.Fuzzy)
	}
	if p.Set != nil {
		params.Add("set", strings.ToLower(strings.TrimSpace(*p.Set)))
	}
	if p.Format != nil {
		params.Add("format", *p.Format)
	}
	if p.Face != nil {
		params.Add("face", *p.Face)
	}
	if p.Version != nil {
		params.Add("version", *p.Version)
	}
	if p.Pretty != nil {
		params.Add("pretty", strconv.FormatBool(*p.Pretty))
	}
	return params, nil
}

func (p *NamedCardParams) BuildPath() string {
	return "/cards/named"
}

func (p *NamedCardParams) BuildFullURL(baseURL string) (string, error) {
	queryParams, err := p.ToURLValues()
	if err != nil {
		return "", err
	}

	fullURL := baseURL + p.BuildPath()

	if len(queryParams) > 0 {
		fullURL += "?" + queryParams.Encode()
	}

	return fullURL, nil
}

func NewExactCardParams(exactName string) *NamedCardParams {
	return &NamedCardParams{
		Exact: &exactName,
	}
}

func NewFuzzyCardParams(fuzzyName string) *NamedCardParams {
	return &NamedCardParams{
		Fuzzy: &fuzzyName,
	}
}

func (p *NamedCardParams) WithSet(setCode string) *NamedCardParams {
	normalizedSet := strings.ToLower(strings.TrimSpace(setCode))
	p.Set = &normalizedSet
	return p
}

func (p *NamedCardParams) WithFormat(format string) *NamedCardParams {
	p.Format = &format
	return p
}

func (p *NamedCardParams) WithVersion(version string) *NamedCardParams {
	p.Version = &version
	return p
}

func (p *NamedCardParams) WithFace(face string) *NamedCardParams {
	p.Face = &face
	return p
}

func (p *NamedCardParams) WithPretty(pretty bool) *NamedCardParams {
	p.Pretty = &pretty
	return p
}

func (p *NamedCardParams) GetSearchType() string {
	if p.Exact != nil {
		return "exact"
	}
	if p.Fuzzy != nil {
		return "fuzzy"
	}
	return "none"
}

func (p *NamedCardParams) GetSearchTerm() string {
	if p.Exact != nil {
		return *p.Exact
	}
	if p.Fuzzy != nil {
		return *p.Fuzzy
	}
	return ""
}
