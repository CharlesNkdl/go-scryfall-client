package cards

import (
	"fmt"
	"net/url"
	"strconv"
)

type IdCardParams struct {
	Id      string  `json:"id" validate:"required,min=1"`
	Format  *string `json:"format,omitempty" validate:"omitempty,oneof=json text image"`
	Version *string `json:"version,omitempty" validate:"omitempty,oneof=small normal large png border_crop art_crop"`
	Face    *string `json:"face,omitempty" validate:"omitempty,oneof=front back"`
	Pretty  *bool   `json:"pretty,omitempty"`
}

func (p *IdCardParams) Validate() error {
	if p.Id == "" {
		return fmt.Errorf("id must be a non-empty string")
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
		case "small", "normal", "large", "png", "border_crop", "art_crop":
		default:
			return fmt.Errorf("invalid version: must be 'small', 'normal', 'large', 'png', 'border_crop', or 'art_crop'")
		}
	}
	if p.Face != nil {
		switch *p.Face {
		case "front", "back":
		default:
			return fmt.Errorf("invalid face: must be 'front' or 'back'")
		}
	}
	if p.Pretty != nil && !*p.Pretty && *p.Pretty {
		return fmt.Errorf("pretty must be a boolean value")
	}
	return nil
}

func (p *IdCardParams) ToURLValues() (url.Values, error) {
	if err := p.Validate(); err != nil {
		return nil, err
	}
	params := url.Values{}

	if p.Format != nil {
		params.Add("format", *p.Format)
	}
	if p.Version != nil {
		params.Add("version", *p.Version)
	}
	if p.Face != nil {
		params.Add("face", *p.Face)
	}
	if p.Pretty != nil {
		params.Add("pretty", strconv.FormatBool(*p.Pretty))
	}
	return params, nil
}
