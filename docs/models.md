# Models and Data Structures

This document describes the data models and structures used throughout the Go Scryfall Client library.

## Overview

The library uses strongly-typed Go structs to represent all data from the Scryfall API. These models are organized into logical packages and provide comprehensive coverage of Magic: The Gathering card data.

## Core Models

### CardCore

The main card data structure representing a Magic card:

```go
type CardCore struct {
    // Identifiers
    ID                string   `json:"id"`                      // Scryfall ID
    ArenaID           *int     `json:"arena_id,omitempty"`      // Arena ID
    MTGOID            *int     `json:"mtgo_id,omitempty"`       // MTGO ID
    MTGOFoilID        *int     `json:"mtgo_foil_id,omitempty"`  // MTGO Foil ID
    MultiverseIDs     []int    `json:"multiverse_ids,omitempty"` // Multiverse IDs
    TcgPlayerID       *int     `json:"tcgplayer_id,omitempty"`  // TCGPlayer ID
    TcgPlayerEtchedID *int     `json:"tcgplayer_etched_id,omitempty"` // TCGPlayer Etched ID
    CardmarketID      *int     `json:"cardmarket_id,omitempty"` // Cardmarket ID
    
    // Basic Properties
    Object            string   `json:"object"`             // Always "card"
    Layout            Layout   `json:"layout"`             // Card layout
    OracleID          *string  `json:"oracle_id,omitempty"` // Oracle ID (same across printings)
    Lang              string   `json:"lang"`               // Language code
    
    // URIs
    PrintsSearchURI   string   `json:"prints_search_uri"`  // Search for all printings
    RulingsURI        string   `json:"rulings_uri"`        // Rulings endpoint
    ScryfallURI       string   `json:"scryfall_uri"`       // Scryfall web page
    URI               string   `json:"uri"`                // API endpoint
    
    // Embedded Data
    GameplayFields    Gameplay `json:"gameplay_fields"`    // Gameplay-related data
    PrintFields       Print    `json:"print_fields"`       // Print-specific data
    CardFaces         []Face   `json:"card_faces,omitempty"` // Multi-faced cards
}
```

### Gameplay Fields

Contains all gameplay-related information:

```go
type Gameplay struct {
    AllParts       []Related      `json:"all_parts"`         // Related cards
    CardFaces      []Face         `json:"card_faces,omitempty"` // Card faces
    CMC            *float32       `json:"cmc,omitempty"`     // Converted mana cost
    ColorIdentity  []models.Color `json:"color_identity,omitempty"` // Color identity
    ColorIndicator []models.Color `json:"color_indicator,omitempty"` // Color indicator
    Colors         []models.Color `json:"colors,omitempty"`  // Card colors
    Defense        *string        `json:"defense,omitempty"` // Defense value
    EDHRecRank     *int           `json:"edhrec_rank,omitempty"` // EDHREC rank
    HandModifier   *string        `json:"hand_modifier,omitempty"` // Hand size modifier
    Keywords       []string       `json:"keywords"`          // Keyword abilities
    Legalities     []Legality     `json:"legalities"`        // Format legalities
    LifeModifier   *string        `json:"life_modifier,omitempty"` // Life modifier
    Loyalty        *string        `json:"loyalty,omitempty"` // Loyalty value
    ManaCost       *string        `json:"mana_cost,omitempty"` // Mana cost
    Name           string         `json:"name"`              // Card name
    OracleText     *string        `json:"oracle_text,omitempty"` // Oracle text
    Power          *string        `json:"power,omitempty"`   // Power value
    ProducedMana   []models.Color `json:"produced_mana,omitempty"` // Mana production
    Reserved       bool           `json:"reserved"`          // Reserved list
    Toughness      *string        `json:"toughness,omitempty"` // Toughness value
    TypeLine       string         `json:"type_line"`         // Type line
}
```

### Print Fields

Contains print-specific information:

```go
type Print struct {
    Artist              *string        `json:"artist,omitempty"`
    ArtistIDs           []string       `json:"artist_ids,omitempty"`
    BorderColor         BorderColor    `json:"border_color"`
    CardBackID          *string        `json:"card_back_id,omitempty"`
    CollectorNumber     string         `json:"collector_number"`
    ContentWarning      *bool          `json:"content_warning,omitempty"`
    Digital             bool           `json:"digital"`
    Finishes            []Finish       `json:"finishes"`
    FlavorName          *string        `json:"flavor_name,omitempty"`
    FlavorText          *string        `json:"flavor_text,omitempty"`
    Frame               Frame          `json:"frame"`
    FrameEffects        []FrameEffect  `json:"frame_effects,omitempty"`
    FullArt             bool           `json:"full_art"`
    Games               []Game         `json:"games"`
    HighresImage        bool           `json:"highres_image"`
    IllustrationID      *string        `json:"illustration_id,omitempty"`
    ImageStatus         ImageStatus    `json:"image_status"`
    ImageUris           *ImageUris     `json:"image_uris,omitempty"`
    Oversized           bool           `json:"oversized"`
    Prices              *Prices        `json:"prices,omitempty"`
    PrintedName         *string        `json:"printed_name,omitempty"`
    PrintedText         *string        `json:"printed_text,omitempty"`
    PrintedTypeLine     *string        `json:"printed_type_line,omitempty"`
    Promo               bool           `json:"promo"`
    PromoTypes          []string       `json:"promo_types,omitempty"`
    PurchaseUris        *PurchaseUris  `json:"purchase_uris,omitempty"`
    Rarity              Rarity         `json:"rarity"`
    RelatedUris         *RelatedUris   `json:"related_uris,omitempty"`
    ReleasedAt          string         `json:"released_at"`
    Reprint             bool           `json:"reprint"`
    ScryfallSetURI      string         `json:"scryfall_set_uri"`
    SetName             string         `json:"set_name"`
    SetSearchURI        string         `json:"set_search_uri"`
    SetType             string         `json:"set_type"`
    SetURI              string         `json:"set_uri"`
    Set                 string         `json:"set"`
    SetID               string         `json:"set_id"`
    StorySpotlight      bool           `json:"story_spotlight"`
    Textless            bool           `json:"textless"`
    Variation           bool           `json:"variation"`
    VariationOf         *string        `json:"variation_of,omitempty"`
    SecurityStamp       *SecurityStamp `json:"security_stamp,omitempty"`
    Watermark           *string        `json:"watermark,omitempty"`
}
```

## Enumeration Types

### Color

Magic card colors:

```go
type Color string

const (
    ColorWhite Color = "W"  // White
    ColorBlue  Color = "U"  // Blue
    ColorBlack Color = "B"  // Black
    ColorRed   Color = "R"  // Red
    ColorGreen Color = "G"  // Green
)
```

### Layout

Card layouts:

```go
type Layout string

const (
    LayoutNormal        Layout = "normal"
    LayoutSplit         Layout = "split"
    LayoutFlip          Layout = "flip"
    LayoutTransform     Layout = "transform"
    LayoutModalDFC      Layout = "modal_dfc"
    LayoutMeld          Layout = "meld"
    LayoutLeveler       Layout = "leveler"
    LayoutSaga          Layout = "saga"
    LayoutAdventure     Layout = "adventure"
    LayoutPlanar        Layout = "planar"
    LayoutScheme        Layout = "scheme"
    LayoutVanguard      Layout = "vanguard"
    LayoutToken         Layout = "token"
    LayoutDoubleFacedToken Layout = "double_faced_token"
    LayoutEmblem        Layout = "emblem"
    LayoutAugment       Layout = "augment"
    LayoutHost          Layout = "host"
    LayoutArtSeries     Layout = "art_series"
    LayoutReversibleCard Layout = "reversible_card"
)
```

### Rarity

Card rarities:

```go
type Rarity string

const (
    RarityCommon    Rarity = "common"
    RarityUncommon  Rarity = "uncommon"
    RarityRare      Rarity = "rare"
    RaritySpecial   Rarity = "special"
    RarityMythic    Rarity = "mythic"
    RarityBonus     Rarity = "bonus"
)
```

### Frame

Card frame types:

```go
type Frame string

const (
    Frame1993   Frame = "1993"
    Frame1997   Frame = "1997"
    Frame2003   Frame = "2003"
    Frame2015   Frame = "2015"
    FrameFuture Frame = "future"
)
```

### BorderColor

Card border colors:

```go
type BorderColor string

const (
    BorderColorBlack  BorderColor = "black"
    BorderColorWhite  BorderColor = "white"
    BorderColorBorderless BorderColor = "borderless"
    BorderColorSilver BorderColor = "silver"
    BorderColorGold   BorderColor = "gold"
)
```

## Specialized Models

### Face

For multi-faced cards (double-faced, split, etc.):

```go
type Face struct {
    Artist          *string        `json:"artist,omitempty"`
    ArtistID        *string        `json:"artist_id,omitempty"`
    CMC             *float32       `json:"cmc,omitempty"`
    ColorIndicator  []models.Color `json:"color_indicator,omitempty"`
    Colors          []models.Color `json:"colors,omitempty"`
    Defense         *string        `json:"defense,omitempty"`
    FlavorText      *string        `json:"flavor_text,omitempty"`
    IllustrationID  *string        `json:"illustration_id,omitempty"`
    ImageUris       *ImageUris     `json:"image_uris,omitempty"`
    Layout          *Layout        `json:"layout,omitempty"`
    Loyalty         *string        `json:"loyalty,omitempty"`
    ManaCost        string         `json:"mana_cost"`
    Name            string         `json:"name"`
    Object          string         `json:"object"`
    OracleID        *string        `json:"oracle_id,omitempty"`
    OracleText      *string        `json:"oracle_text,omitempty"`
    Power           *string        `json:"power,omitempty"`
    PrintedName     *string        `json:"printed_name,omitempty"`
    PrintedText     *string        `json:"printed_text,omitempty"`
    PrintedTypeLine *string        `json:"printed_type_line,omitempty"`
    Toughness       *string        `json:"toughness,omitempty"`
    TypeLine        *string        `json:"type_line,omitempty"`
    Watermark       *string        `json:"watermark,omitempty"`
}
```

### ImageUris

Card image URLs:

```go
type ImageUris struct {
    Small      string `json:"small"`       // Small image (146×204)
    Normal     string `json:"normal"`      // Normal image (488×680)
    Large      string `json:"large"`       // Large image (672×936)
    PNG        string `json:"png"`         // PNG format
    ArtCrop    string `json:"art_crop"`    // Cropped art
    BorderCrop string `json:"border_crop"` // Cropped with border
}
```

### Prices

Card pricing information:

```go
type Prices struct {
    USD       *string `json:"usd,omitempty"`        // USD price
    USDFoil   *string `json:"usd_foil,omitempty"`   // USD foil price
    USDEtched *string `json:"usd_etched,omitempty"` // USD etched price
    EUR       *string `json:"eur,omitempty"`        // EUR price
    EURFoil   *string `json:"eur_foil,omitempty"`   // EUR foil price
    Tix       *string `json:"tix,omitempty"`        // MTGO tickets
}
```

### Legality

Format legality information:

```go
type Legality struct {
    Format   string `json:"format"`   // Format name
    Legality string `json:"legality"` // legal, not_legal, restricted, banned
}
```

## Collection Types

### List[T]

Generic paginated list response:

```go
type List[T any] struct {
    Object     string   `json:"object"`               // Always "list"
    Data       []T      `json:"data"`                 // Array of items
    HasMore    bool     `json:"has_more"`             // More pages available
    NextPage   *string  `json:"next_page,omitempty"`  // Next page URL
    TotalCards *int     `json:"total_cards,omitempty"` // Total count
    Warnings   []string `json:"warnings,omitempty"`   // API warnings
}
```

### Catalog

Catalog of string values:

```go
type Catalog struct {
    Object string   `json:"object"` // Always "catalog"
    URI    string   `json:"uri"`    // API endpoint
    Data   []string `json:"data"`   // Array of strings
}
```

## Set Models

### Set

Magic set information:

```go
type Set struct {
    ID             string  `json:"id"`              // Scryfall ID
    Code           string  `json:"code"`            // Set code
    Name           string  `json:"name"`            // Set name
    URI            string  `json:"uri"`             // API endpoint
    ScryfallURI    string  `json:"scryfall_uri"`    // Scryfall page
    SearchURI      string  `json:"search_uri"`      // Card search
    ReleasedAt     string  `json:"released_at"`     // Release date
    SetType        string  `json:"set_type"`        // Set type
    CardCount      int     `json:"card_count"`      // Number of cards
    ParentSetCode  *string `json:"parent_set_code,omitempty"` // Parent set
    Digital        bool    `json:"digital"`         // Digital only
    FoilOnly       bool    `json:"foil_only"`       // Foil only
    NonfoilOnly    bool    `json:"nonfoil_only"`    // Non-foil only
    Block          *string `json:"block,omitempty"` // Block name
    BlockCode      *string `json:"block_code,omitempty"` // Block code
    IconSVGURI     string  `json:"icon_svg_uri"`    // Set symbol SVG
}
```

## Ruling Models

### Ruling

Official card ruling:

```go
type Ruling struct {
    Object       string `json:"object"`        // Always "ruling"
    OracleID     string `json:"oracle_id"`     // Oracle ID
    Source       string `json:"source"`        // Ruling source
    PublishedAt  string `json:"published_at"`  // Publication date
    Comment      string `json:"comment"`       // Ruling text
}
```

## Symbol Models

### Symbol

Mana symbol information:

```go
type Symbol struct {
    Object          string   `json:"object"`           // Always "card_symbol"
    Symbol          string   `json:"symbol"`           // Symbol text
    SvgURI          string   `json:"svg_uri"`          // SVG image URL
    Cost            float64  `json:"cost"`             // Mana cost
    CMC             float64  `json:"cmc"`              // CMC contribution
    Colors          []string `json:"colors"`           // Color identity
    AppearsInManaCosts bool  `json:"appears_in_mana_costs"` // In mana costs
    Funny           bool     `json:"funny"`            // Un-set symbol
    Transposable    bool     `json:"transposable"`     // Can substitute
    RepresentsMana  bool     `json:"represents_mana"`  // Represents mana
    English         string   `json:"english"`          // Description
}
```

## Error Models

### ScryfallError

API error information:

```go
type ScryfallError struct {
    Object   string                 `json:"object"`   // Always "error"
    Code     string                 `json:"code"`     // Error code
    Status   int                    `json:"status"`   // HTTP status
    Detail   string                 `json:"detail"`   // Error message
    Type     *string                `json:"type,omitempty"` // Error type
    Warnings []string               `json:"warnings,omitempty"` // Warnings
    Details  map[string]interface{} `json:"details,omitempty"` // Extra details
}
```

### ApiError

Wrapped API error:

```go
type ApiError struct {
    ErrInfo ScryfallError
}

func (e *ApiError) Error() string {
    return e.ErrInfo.Detail
}
```

## Usage Examples

### Working with Card Data

```go
// Access gameplay fields
card, err := client.Cards.GetByName(ctx, cards.NewExactCardParams("Lightning Bolt"))
if err != nil {
    return err
}

fmt.Printf("Name: %s\n", card.Name)
fmt.Printf("Type: %s\n", card.TypeLine)
if card.ManaCost != nil {
    fmt.Printf("Mana Cost: %s\n", *card.ManaCost)
}
if card.OracleText != nil {
    fmt.Printf("Oracle Text: %s\n", *card.OracleText)
}

// Check colors
if len(card.Colors) > 0 {
    fmt.Printf("Colors: %v\n", card.Colors)
}

// Check legalities
for _, legality := range card.Legalities {
    if legality.Format == "modern" {
        fmt.Printf("Modern: %s\n", legality.Legality)
    }
}
```

### Working with Multi-faced Cards

```go
card, err := client.Cards.GetByName(ctx, cards.NewExactCardParams("Delver of Secrets"))
if err != nil {
    return err
}

if len(card.CardFaces) > 0 {
    for i, face := range card.CardFaces {
        fmt.Printf("Face %d: %s\n", i+1, face.Name)
        if face.OracleText != nil {
            fmt.Printf("  Text: %s\n", *face.OracleText)
        }
    }
}
```

### Working with Prices

```go
if card.Prices != nil {
    if card.Prices.USD != nil {
        fmt.Printf("USD: $%s\n", *card.Prices.USD)
    }
    if card.Prices.USDFoil != nil {
        fmt.Printf("USD Foil: $%s\n", *card.Prices.USDFoil)
    }
    if card.Prices.EUR != nil {
        fmt.Printf("EUR: €%s\n", *card.Prices.EUR)
    }
    if card.Prices.Tix != nil {
        fmt.Printf("MTGO: %s tix\n", *card.Prices.Tix)
    }
}
```

### Working with Images

```go
if card.ImageUris != nil {
    fmt.Printf("Small: %s\n", card.ImageUris.Small)
    fmt.Printf("Normal: %s\n", card.ImageUris.Normal)
    fmt.Printf("Large: %s\n", card.ImageUris.Large)
    fmt.Printf("Art Crop: %s\n", card.ImageUris.ArtCrop)
}
```

## Type Safety Features

### Pointer Helpers

For optional fields, use helper functions to create pointers:

```go
func StringPtr(s string) *string {
    return &s
}

func IntPtr(i int) *int {
    return &i
}

func BoolPtr(b bool) *bool {
    return &b
}

// Usage
params := &cards.IdCardParams{
    ID:     "card-id",
    Format: StringPtr("json"),
    Pretty: BoolPtr(true),
}
```

### Type Assertions

When working with interface{} values:

```go
if card.Prices.USD != nil {
    if price, err := strconv.ParseFloat(*card.Prices.USD, 64); err == nil {
        fmt.Printf("Price: $%.2f\n", price)
    }
}
```

## Best Practices

1. **Check for nil pointers** - Many fields are optional and may be nil
2. **Use type constants** - Use provided constants for enums instead of strings
3. **Handle different layouts** - Different card layouts have different data available
4. **Parse prices carefully** - Prices are strings and may be empty or "N/A"
5. **Validate data** - Always validate data before using it in business logic

The models provide a comprehensive, type-safe way to work with Magic: The Gathering data from the Scryfall API, ensuring your applications can handle the full complexity of Magic cards while maintaining Go's type safety guarantees.