# Set Service Documentation

The Set Service provides access to Magic: The Gathering set information through the Scryfall API.

## Available Methods

### GetById

Retrieve information about a specific Magic set by its set code.

```go
func (s *SetService) GetById(ctx context.Context, id string) (*set.Set, error)
```

**Parameters:**
- `ctx`: Context for timeout and cancellation
- `id`: Set code (3-4 character identifier)

**Example:**
```go
// Get information about Kaldheim
set, err := client.Sets.GetById(ctx, "khm")
if err != nil {
    log.Printf("Error: %v", err)
    return
}

fmt.Printf("Set Name: %s\n", set.Name)
fmt.Printf("Set Code: %s\n", set.Code)
fmt.Printf("Release Date: %s\n", set.ReleasedAt)
fmt.Printf("Card Count: %d\n", set.CardCount)
fmt.Printf("Set Type: %s\n", set.SetType)

if set.Block != nil {
    fmt.Printf("Block: %s\n", *set.Block)
}

if set.ParentSetCode != nil {
    fmt.Printf("Parent Set: %s\n", *set.ParentSetCode)
}
```

## Set Information

The returned `Set` object contains comprehensive information about the Magic set:

### Basic Information
```go
type Set struct {
    ID          string  `json:"id"`           // Scryfall ID
    Code        string  `json:"code"`         // Set code (khm, lea, etc.)
    Name        string  `json:"name"`         // Display name
    ReleasedAt  string  `json:"released_at"`  // Release date (YYYY-MM-DD)
    SetType     string  `json:"set_type"`     // Set type
    CardCount   int     `json:"card_count"`   // Number of cards in set
    Digital     bool    `json:"digital"`      // Digital-only set
    FoilOnly    bool    `json:"foil_only"`    // Foil-only set
}
```

### Optional Information
```go
// Block information (older sets)
Block *string `json:"block,omitempty"`
ParentSetCode *string `json:"parent_set_code,omitempty"`

// URLs and identifiers
ScryfallURI string `json:"scryfall_uri"`    // Scryfall page URL
URI         string `json:"uri"`             // API endpoint
IconSVGURI  string `json:"icon_svg_uri"`    // Set symbol SVG
SearchURI   string `json:"search_uri"`      // Search cards in set
```

## Common Set Codes

Here are some commonly used set codes:

### Recent Standard Sets
```go
"khm"  // Kaldheim
"znr"  // Zendikar Rising
"m21"  // Core Set 2021
"iko"  // Ikoria: Lair of Behemoths
"thb"  // Theros Beyond Death
"eld"  // Throne of Eldraine
"m20"  // Core Set 2020
"war"  // War of the Spark
"rna"  // Ravnica Allegiance
"grn"  // Guilds of Ravnica
```

### Classic Sets
```go
"lea"  // Limited Edition Alpha
"leb"  // Limited Edition Beta
"2ed"  // Unlimited Edition
"3ed"  // Revised Edition
"4ed"  // Fourth Edition
"ice"  // Ice Age
"mir"  // Mirage
"tmp"  // Tempest
"usg"  // Urza's Saga
"mmq"  // Mercadian Masques
```

### Special Sets
```go
"cmr"  // Commander Legends
"2xm"  // Double Masters
"jmp"  // Jumpstart
"m25"  // Masters 25
"ima"  // Iconic Masters
"mm3"  // Modern Masters 2017
"ema"  // Eternal Masters
```

## Use Cases

### Get Set Information for Deck Building
```go
func getSetInfo(client *scryfall.Client, setCodes []string) {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    for _, code := range setCodes {
        set, err := client.Sets.GetById(ctx, code)
        if err != nil {
            fmt.Printf("❌ %s: Error - %v\n", code, err)
            continue
        }
        
        fmt.Printf("📦 %s (%s)\n", set.Name, set.Code)
        fmt.Printf("   Released: %s\n", set.ReleasedAt)
        fmt.Printf("   Cards: %d\n", set.CardCount)
        fmt.Printf("   Type: %s\n", set.SetType)
        
        if set.Block != nil {
            fmt.Printf("   Block: %s\n", *set.Block)
        }
        
        fmt.Println()
    }
}

// Usage
setCodes := []string{"khm", "znr", "iko", "thb", "eld"}
getSetInfo(client, setCodes)
```

### Validate Set Codes
```go
func validateSetCode(client *scryfall.Client, setCode string) bool {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    _, err := client.Sets.GetById(ctx, setCode)
    return err == nil
}

// Usage
if validateSetCode(client, "khm") {
    fmt.Println("Valid set code")
} else {
    fmt.Println("Invalid set code")
}
```

### Check Set Release Timeline
```go
func checkSetTimeline(client *scryfall.Client, setCodes []string) {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    type SetInfo struct {
        Name        string
        Code        string
        ReleasedAt  string
    }
    
    var sets []SetInfo
    
    for _, code := range setCodes {
        set, err := client.Sets.GetById(ctx, code)
        if err != nil {
            continue
        }
        
        sets = append(sets, SetInfo{
            Name:       set.Name,
            Code:       set.Code,
            ReleasedAt: set.ReleasedAt,
        })
    }
    
    // Sort by release date if needed
    fmt.Println("Set Release Timeline:")
    for _, set := range sets {
        fmt.Printf("%s: %s (%s)\n", set.ReleasedAt, set.Name, set.Code)
    }
}
```

## Error Handling

The Set Service returns standard API errors:

```go
set, err := client.Sets.GetById(ctx, "invalid")
if err != nil {
    if apiErr, ok := err.(*errors.ApiError); ok {
        switch apiErr.ErrInfo.Status {
        case 404:
            fmt.Println("Set not found")
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

1. **Use standard set codes** - Always use the official 3-4 character set codes
2. **Cache set information** - Set data rarely changes, so consider caching
3. **Handle invalid codes** - Always check for 404 errors with user-provided set codes
4. **Use context timeouts** - Set appropriate timeouts for your use case
5. **Consider rate limiting** - The client handles this automatically, but be aware of it

## Related Services

- Use **Card Service** to get cards from a specific set: `client.Cards.Search(ctx, &cards.SearchParams{Query: "set:khm"})`
- Use **Card Service** with set parameters: `cards.NewExactCardParams("Lightning Bolt").WithSet("lea")`