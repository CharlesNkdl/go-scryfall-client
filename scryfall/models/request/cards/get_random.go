package cards

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type RandomCardParams struct {
	Query   *string `json:"q,omitempty" validate:"omitempty"`
	Format  *string `json:"format,omitempty" validate:"omitempty,oneof=json text image"`
	Face    *string `json:"face,omitempty" validate:"omitempty,oneof=front back"`
	Version *string `json:"version,omitempty" validate:"omitempty,oneof=small normal large png art_crop border_crop"`
	Pretty  *bool   `json:"pretty,omitempty"`
}

func (p *RandomCardParams) Validate() error {
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
	if p.Query != nil && strings.TrimSpace(*p.Query) == "" {
		return fmt.Errorf("search query 'q' cannot be empty when provided")
	}

	return nil
}

func (p *RandomCardParams) ToURLValues() (url.Values, error) {
	if err := p.Validate(); err != nil {
		return nil, err
	}
	params := url.Values{}
	if p.Query != nil {
		params.Add("q", strings.TrimSpace(*p.Query))
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

func (p *RandomCardParams) BuildPath() string {
	return "/cards/random"
}

func (p *RandomCardParams) BuildFullURL(baseURL string) (string, error) {
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

func NewRandomCardParams() *RandomCardParams {
	return &RandomCardParams{}
}

func NewRandomCardParamsWithQuery(query string) *RandomCardParams {
	trimmedQuery := strings.TrimSpace(query)
	return &RandomCardParams{
		Query: &trimmedQuery,
	}
}

func (p *RandomCardParams) WithQuery(query string) *RandomCardParams {
	trimmedQuery := strings.TrimSpace(query)
	p.Query = &trimmedQuery
	return p
}

func (p *RandomCardParams) WithFormat(format string) *RandomCardParams {
	p.Format = &format
	return p
}

func (p *RandomCardParams) WithVersion(version string) *RandomCardParams {
	p.Version = &version
	return p
}

func (p *RandomCardParams) WithFace(face string) *RandomCardParams {
	p.Face = &face
	return p
}

func (p *RandomCardParams) WithPretty(pretty bool) *RandomCardParams {
	p.Pretty = &pretty
	return p
}

func (p *RandomCardParams) HasQuery() bool {
	return p.Query != nil && strings.TrimSpace(*p.Query) != ""
}

func (p *RandomCardParams) GetQuery() string {
	if p.Query != nil {
		return strings.TrimSpace(*p.Query)
	}
	return ""
}

func (p *RandomCardParams) ClearQuery() *RandomCardParams {
	p.Query = nil
	return p
}
