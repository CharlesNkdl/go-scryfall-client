# Card Service Documentation

The Card Service provides access to all card-related endpoints in the Scryfall API. It supports searching, retrieving cards by various identifiers, autocomplete functionality, and getting random cards.

## Available Methods

### GetById

Retrieve a specific card by its Scryfall ID.

```go
func (s *CardService) GetById(ctx context.Context, params *cardreq.IdCardParams) (*card.CardCore, error)
```

**Parameters:**
- `ctx`: Context for timeout and cancellation
- `params`: ID card parameters

**Example:**
```go
params := &cards.IdCardParams{
    ID: "60a25a91-40ba-4796-b9f2-cff5ac2e2c8a",
    Format: StringPtr("json"), // Optional: json, text, image
    Version: StringPtr("small"), // Optional: small, normal, large, png, art_crop, border_crop
    Face: StringPtr("front"), // Optional: front, back (for double-faced cards)
    Pretty: BoolPtr(true), // Optional: pretty-print JSON
}

card, err := client.Cards.GetById(ctx, params)
if err != nil {
    log.Printf("Error: %v", err)
    return
}

fmt.Printf("Card: %s\n", card.Name)
fmt.Printf("Type: %s\n", card.TypeLine)
if card.ManaCost != nil {
    fmt.Printf("Mana Cost: %s\n", *card.ManaCost)
}
```

### GetByName

Search for cards by name using exact or fuzzy matching.

```go
func (s *CardService) GetByName(ctx context.Context, params *cardreq.NamedCardParams) (*card.CardCore, error)
```

**Parameters:**
- `ctx`: Context for timeout and cancellation
- `params`: Named card parameters

**Convenience Constructors:**
```go
// Exact name search
params := cards.NewExactCardParams("Lightning Bolt")

// Fuzzy name search (allows typos)
params := cards.NewFuzzyCardParams("Lighning Bolt")

// With additional parameters
params := cards.NewExactCardParams("Lightning Bolt").
    WithSet("lea").
    WithFormat("json").
    WithPretty(true)
```

**Examples:**

#### Exact Name Search
```go
params := cards.NewExactCardParams("Lightning Bolt")
card, err := client.Cards.GetByName(ctx, params)
if err != nil {
    log.Printf("Error: %v", err)
    return
}
fmt.Printf("Found: %s\n", card.Name)
```

#### Fuzzy Name Search
```go
// This will find "Lightning Bolt" even with the typo
params := cards.NewFuzzyCardParams("Lighning Bolt")
card, err := client.Cards.GetByName(ctx, params)
if err != nil {
    log.Printf("Error: %v", err)
    return
}
fmt.Printf("Found: %s (corrected from fuzzy search)\n", card.Name)
```

#### With Set Specification
```go
params := cards.NewExactCardParams("Lightning Bolt").WithSet("lea")
card, err := client.Cards.GetByName(ctx, params)
if err != nil {
    log.Printf("Error: %v", err)
    return
}
fmt.Printf("Found: %s from set %s\n", card.Name, card.Set)
```

### Search

Perform advanced card searches using Scryfall's powerful search syntax.

```go
func (s *CardService) Search(ctx context.Context, params *cardreq.SearchParams) (*models.List[card.CardCore], error)
```

**Parameters:**
- `ctx`: Context for timeout and cancellation
- `params`: Search parameters

**Example:**
```go
params := &cards.SearchParams{
    Query: "type:creature power:3 cmc:2",
    Page: 1,
    Order: StringPtr("name"),
    Dir: StringPtr("asc"),
    Unique: StringPtr("cards"),
    Format: StringPtr("json"),
    Pretty: BoolPtr(true),
}

searchResult, err := client.Cards.Search(ctx, params)
if err != nil {
    log.Printf("Error: %v", err)
    return
}

fmt.Printf("Found %d cards (total: %d)\n", len(searchResult.Data), searchResult.TotalCards)
for _, card := range searchResult.Data {
    fmt.Printf("- %s (%s)\n", card.Name, card.TypeLine)
}

// Handle pagination
if searchResult.HasMore {
    fmt.Printf("More results available. Next page URL: %s\n", searchResult.NextPage)
}
```

**Search Query Examples:**
```go
// Creatures with power 3
"type:creature power:3"

// Red spells costing 1 mana
"color:red cmc:1"

// Cards from specific set
"set:khm"

// Legendary creatures
"type:legendary type:creature"

// Cards containing specific text
"oracle:\"draw a card\""

// Cards by artist
"artist:\"Rebecca Guay\""

// Expensive cards
"usd>50"

// Complex queries
"type:creature power>=5 toughness<=3 cmc<=4"
```

### Autocomplete

Get autocomplete suggestions for card names.

```go
func (s *CardService) Autocomplete(ctx context.Context, params *cardreq.AutocompleteCardParams) (*models.Catalog, error)
```

**Example:**
```go
params := &cards.AutocompleteCardParams{
    Query: "light",
    IncludeExtras: BoolPtr(false),
}

suggestions, err := client.Cards.Autocomplete(ctx, params)
if err != nil {
    log.Printf("Error: %v", err)
    return
}

fmt.Printf("Autocomplete suggestions for 'light':\n")
for _, suggestion := range suggestions.Data {
    fmt.Printf("- %s\n", suggestion)
}
```

### GetRandom

Get a random card, optionally filtered by a search query.

```go
func (s *CardService) GetRandom(ctx context.Context, params *cardreq.RandomCardParams) (*card.CardCore, error)
```

**Example:**
```go
// Completely random card
params := &cards.RandomCardParams{}
card, err := client.Cards.GetRandom(ctx, params)
if err != nil {
    log.Printf("Error: %v", err)
    return
}
fmt.Printf("Random card: %s\n", card.Name)

// Random creature
params = &cards.RandomCardParams{
    Query: StringPtr("type:creature"),
}
card, err = client.Cards.GetRandom(ctx, params)
if err != nil {
    log.Printf("Error: %v", err)
    return
}
fmt.Printf("Random creature: %s\n", card.Name)
```

### GetByCodeNumberLang

Get cards from a specific set by collector number and language.

```go
func (s *CardService) GetByCodeNumberLang(ctx context.Context, params *cardreq.CardsByCodeNumberLangParams) (*models.List[card.CardCore], error)
```

**Example:**
```go
params := &cards.CardsByCodeNumberLangParams{
    SetCode: "khm",
    CollectorNumber: "123",
    Lang: StringPtr("en"), // Optional language code
}

cards, err := client.Cards.GetByCodeNumberLang(ctx, params)
if err != nil {
    log.Printf("Error: %v", err)
    return
}

for _, card := range cards.Data {
    fmt.Printf("Card: %s (Set: %s, Number: %s)\n", 
        card.Name, card.Set, card.CollectorNumber)
}
```

## Parameter Types

### IdCardParams

Used with `GetById`:

```go
type IdCardParams struct {
    ID      string   // Required: Scryfall ID
    Format  *string  // Optional: json, text, image
    Version *string  // Optional: small, normal, large, png, art_crop, border_crop
    Face    *string  // Optional: front, back
    Pretty  *bool    // Optional: pretty-print JSON
}
```

### NamedCardParams

Used with `GetByName`:

```go
type NamedCardParams struct {
    Exact   *string  // Exact card name
    Fuzzy   *string  // Fuzzy card name (allows typos)
    Set     *string  // Optional: set code
    Format  *string  // Optional: json, text, image
    Version *string  // Optional: version type
    Face    *string  // Optional: card face
    Pretty  *bool    // Optional: pretty-print JSON
}
```

**Validation Rules:**
- Must specify either `Exact` or `Fuzzy`, not both
- Set code must be 3-4 characters if specified
- Face "back" requires image format

### SearchParams

Used with `Search`:

```go
type SearchParams struct {
    Query          string   // Required: search query
    Page           int      // Optional: page number (default: 1)
    Order          *string  // Optional: name, set, released, rarity, color, usd, tix, eur, cmc, power, toughness, edhrec, artist
    Dir            *string  // Optional: auto, asc, desc
    Unique         *string  // Optional: cards, art, prints
    IncludeExtras  *bool    // Optional: include extras
    IncludeMultilingual *bool // Optional: include multilingual
    IncludeVariations *bool  // Optional: include variations
    Format         *string  // Optional: json, csv
    Pretty         *bool    // Optional: pretty-print JSON
}
```

### AutocompleteCardParams

Used with `Autocomplete`:

```go
type AutocompleteCardParams struct {
    Query         string // Required: partial card name
    IncludeExtras *bool  // Optional: include extra cards
}
```

### RandomCardParams

Used with `GetRandom`:

```go
type RandomCardParams struct {
    Query *string // Optional: filter random cards by search query
}
```

## Helper Functions

For convenience, use these helper functions to create pointers to basic types:

```go
func StringPtr(s string) *string {
    return &s
}

func BoolPtr(b bool) *bool {
    return &b
}

func IntPtr(i int) *int {
    return &i
}
```

## Error Handling

All card service methods return `*errors.ApiError` for API-related errors:

```go
card, err := client.Cards.GetByName(ctx, params)
if err != nil {
    if apiErr, ok := err.(*errors.ApiError); ok {
        switch apiErr.ErrInfo.Status {
        case 404:
            fmt.Println("Card not found")
        case 429:
            fmt.Println("Rate limited")
        default:
            fmt.Printf("API error: %s\n", apiErr.ErrInfo.Detail)
        }
    } else {
        fmt.Printf("Network error: %v\n", err)
    }
}
```

## Best Practices

1. **Use exact searches when possible** for better performance
2. **Specify set codes** when looking for specific printings
3. **Handle pagination** for search results with many cards
4. **Use context timeouts** for all API calls
5. **Cache results** when appropriate to reduce API calls
6. **Validate parameters** before making requests

## Common Use Cases

### Building a Deck Checker
```go
func checkDeck(client *scryfall.Client, cardNames []string) {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    for _, name := range cardNames {
        params := cards.NewExactCardParams(name)
        card, err := client.Cards.GetByName(ctx, params)
        if err != nil {
            fmt.Printf("❌ %s: Not found\n", name)
            continue
        }
        fmt.Printf("✅ %s: Found (%s)\n", name, card.Set)
    }
}
```

### Random Card Discovery
```go
func discoverRandomCreatures(client *scryfall.Client, count int) {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    params := &cards.RandomCardParams{
        Query: StringPtr("type:creature"),
    }
    
    for i := 0; i < count; i++ {
        card, err := client.Cards.GetRandom(ctx, params)
        if err != nil {
            log.Printf("Error getting random card: %v", err)
            continue
        }
        fmt.Printf("%d. %s - %s\n", i+1, card.Name, card.TypeLine)
    }
}
```