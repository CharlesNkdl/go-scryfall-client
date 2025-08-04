package cards

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

var supportedLanguages = map[string]bool{
	"en":  true, // English
	"es":  true, // Spanish
	"fr":  true, // French
	"de":  true, // German
	"it":  true, // Italian
	"pt":  true, // Portuguese
	"ja":  true, // Japanese
	"ko":  true, // Korean
	"ru":  true, // Russian
	"zhs": true, // Simplified Chinese
	"zht": true, // Traditional Chinese
	"he":  true, // Hebrew
	"la":  true, // Latin
	"grc": true, // Ancient Greek
	"ar":  true, // Arabic
	"sa":  true, // Sanskrit
	"ph":  true, // Phyrexian
	"qya": true, // Quenya
}

type SetCardParams struct {
	Code   string  `json:"code" validate:"required,min=3,max=5"`
	Number string  `json:"number" validate:"required"`
	Lang   *string `json:"lang,omitempty" validate:"omitempty,oneof=en es fr de it pt ja ko ru zhs zht he la grc ar sa ph qya"`

	Format  *string `json:"format,omitempty" validate:"omitempty,oneof=json text image"`
	Face    *string `json:"face,omitempty" validate:"omitempty,oneof=front back"`
	Version *string `json:"version,omitempty" validate:"omitempty,oneof=small normal large png art_crop border_crop"`
	Pretty  *bool   `json:"pretty,omitempty"`
}

func (p *SetCardParams) Validate() error {
	if len(p.Code) < 3 || len(p.Code) > 5 {
		return fmt.Errorf("code must be between 3 and 5 characters")
	}
	if strings.TrimSpace(p.Code) == "" {
		return fmt.Errorf("code cannot be empty")
	}
	if strings.TrimSpace(p.Number) == "" {
		return fmt.Errorf("number cannot be empty")
	}
	if p.Lang != nil {
		if !supportedLanguages[*p.Lang] {
			return fmt.Errorf("unsupported language: %s. Supported languages: en, es, fr, de, it, pt, ja, ko, ru, zhs, zht, he, la, grc, ar, sa, ph, qya", *p.Lang)
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

func (p *SetCardParams) ToURLValues() (url.Values, error) {
	if err := p.Validate(); err != nil {
		return nil, err
	}
	params := url.Values{}
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

func (p *SetCardParams) BuildPath() string {
	if p.Lang != nil && *p.Lang != "" {
		return fmt.Sprintf("/cards/%s/%s/%s", p.Code, p.Number, *p.Lang)
	}
	return fmt.Sprintf("/cards/%s/%s", p.Code, p.Number)
}

func (p *SetCardParams) BuildFullURL(baseURL string) (string, error) {
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

func NewSetCardParams(code, number string) *SetCardParams {
	return &SetCardParams{
		Code:   strings.ToLower(code),
		Number: number,
	}
}

func IsValidLanguage(lang string) bool {
	return supportedLanguages[lang]
}

func GetSupportedLanguages() []string {
	langs := make([]string, 0, len(supportedLanguages))
	for lang := range supportedLanguages {
		langs = append(langs, lang)
	}
	return langs
}

func (p *SetCardParams) WithLang(lang string) *SetCardParams {
	p.Lang = &lang
	return p
}

func (p *SetCardParams) WithFormat(format string) *SetCardParams {
	p.Format = &format
	return p
}

func (p *SetCardParams) WithVersion(version string) *SetCardParams {
	p.Version = &version
	return p
}

func (p *SetCardParams) WithFace(face string) *SetCardParams {
	p.Face = &face
	return p
}

func (p *SetCardParams) WithPretty(pretty bool) *SetCardParams {
	p.Pretty = &pretty
	return p
}
