# Symbol Service Documentation

The Symbol Service provides access to mana symbol information and imagery from the Scryfall API. This includes all mana symbols, hybrid symbols, special symbols, and their SVG representations.

## Available Methods

### GetSymbol

Retrieve information about a specific mana symbol.

```go
func (s *SymbolService) GetSymbol(ctx context.Context, symbol string) (*models.Symbol, error)
```

**Parameters:**
- `ctx`: Context for timeout and cancellation
- `symbol`: The symbol to look up (including braces)

**Example:**
```go
// Get information about the red mana symbol
symbol, err := client.Symbols.GetSymbol(ctx, "{R}")
if err != nil {
    log.Printf("Error: %v", err)
    return
}

fmt.Printf("Symbol: %s\n", symbol.Symbol)
fmt.Printf("SVG URI: %s\n", symbol.SvgUri)
fmt.Printf("Cost: %f\n", symbol.Cost)
fmt.Printf("CMC: %f\n", symbol.Cmc)
fmt.Printf("Colors: %v\n", symbol.Colors)
```

## Symbol Format

Mana symbols must be provided in Magic's standard format with curly braces:

### Basic Mana Symbols
```go
"{W}"  // White mana
"{U}"  // Blue mana  
"{B}"  // Black mana
"{R}"  // Red mana
"{G}"  // Green mana
"{C}"  // Colorless mana
```

### Generic Mana
```go
"{0}"  // Zero mana
"{1}"  // One generic mana
"{2}"  // Two generic mana
"{X}"  // X generic mana
"{Y}"  // Y generic mana
"{Z}"  // Z generic mana
```

### Hybrid Mana
```go
"{W/U}"  // White or Blue
"{U/B}"  // Blue or Black
"{B/R}"  // Black or Red
"{R/G}"  // Red or Green
"{G/W}"  // Green or White
"{2/W}"  // Two or White
"{2/U}"  // Two or Blue
```

### Phyrexian Mana
```go
"{W/P}"  // Phyrexian White
"{U/P}"  // Phyrexian Blue
"{B/P}"  // Phyrexian Black
"{R/P}"  // Phyrexian Red
"{G/P}"  // Phyrexian Green
```

### Special Symbols
```go
"{T}"     // Tap symbol
"{Q}"     // Untap symbol
"{E}"     // Energy counter
"{CHAOS}" // Chaos symbol
"{A}"     // Acorn symbol
"{TK}"    // Ticket symbol
```

## Symbol Information

Each symbol contains comprehensive information:

```go
type Symbol struct {
    Object    string   `json:"object"`     // Always "card_symbol"
    Symbol    string   `json:"symbol"`     // The symbol text
    SvgUri    string   `json:"svg_uri"`    // SVG image URL
    Cost      float64  `json:"cost"`       // Mana cost value
    Cmc       float64  `json:"cmc"`        // Converted mana cost
    Colors    []string `json:"colors"`     // Color identity
    Appears   bool     `json:"appears_in_mana_costs"` // Appears in costs
    Funny     bool     `json:"funny"`      // Un-set symbol
    Transposable bool  `json:"transposable"` // Can be substituted
    Represents_mana bool `json:"represents_mana"` // Represents mana
    English   string   `json:"english"`    // English description
}
```

## Complete Examples

### Display Mana Symbol Info
```go
func displaySymbolInfo(client *scryfall.Client, symbolText string) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    symbol, err := client.Symbols.GetSymbol(ctx, symbolText)
    if err != nil {
        log.Printf("Error getting symbol %s: %v", symbolText, err)
        return
    }
    
    fmt.Printf("═══ %s ═══\n", symbol.Symbol)
    fmt.Printf("Description: %s\n", symbol.English)
    fmt.Printf("Mana Cost: %.1f\n", symbol.Cost)
    fmt.Printf("CMC: %.1f\n", symbol.Cmc)
    
    if len(symbol.Colors) > 0 {
        fmt.Printf("Colors: %s\n", strings.Join(symbol.Colors, ", "))
    } else {
        fmt.Printf("Colors: Colorless\n")
    }
    
    fmt.Printf("SVG Image: %s\n", symbol.SvgUri)
    fmt.Printf("Appears in costs: %t\n", symbol.Appears)
    
    if symbol.Funny {
        fmt.Printf("⚠️ Un-set symbol\n")
    }
    
    fmt.Println()
}

// Usage
symbols := []string{"{R}", "{W/U}", "{X}", "{T}", "{2/B}"}
for _, sym := range symbols {
    displaySymbolInfo(client, sym)
}
```

### Analyze Mana Cost
```go
func analyzeManaString(client *scryfall.Client, manaCost string) {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    // Parse mana cost string into individual symbols
    symbols := parseManaCost(manaCost) // You'd need to implement this
    
    var totalCost float64
    var totalCMC float64
    var colors []string
    colorMap := make(map[string]bool)
    
    fmt.Printf("🔮 Analyzing mana cost: %s\n\n", manaCost)
    
    for _, symbolText := range symbols {
        symbol, err := client.Symbols.GetSymbol(ctx, symbolText)
        if err != nil {
            fmt.Printf("❌ Error with symbol %s: %v\n", symbolText, err)
            continue
        }
        
        fmt.Printf("%-6s Cost: %3.1f  CMC: %3.1f  %s\n", 
            symbol.Symbol, symbol.Cost, symbol.Cmc, symbol.English)
        
        totalCost += symbol.Cost
        totalCMC += symbol.Cmc
        
        // Collect unique colors
        for _, color := range symbol.Colors {
            if !colorMap[color] {
                colorMap[color] = true
                colors = append(colors, color)
            }
        }
    }
    
    fmt.Printf("\n📊 Total Cost: %.1f\n", totalCost)
    fmt.Printf("📊 Total CMC: %.1f\n", totalCMC)
    fmt.Printf("🎨 Color Identity: %s\n", strings.Join(colors, ", "))
}

// Helper function to parse mana cost
func parseManaCost(manaCost string) []string {
    var symbols []string
    runes := []rune(manaCost)
    
    for i := 0; i < len(runes); i++ {
        if runes[i] == '{' {
            end := i + 1
            for end < len(runes) && runes[end] != '}' {
                end++
            }
            if end < len(runes) {
                symbols = append(symbols, string(runes[i:end+1]))
                i = end
            }
        }
    }
    
    return symbols
}
```

### Build Mana Symbol Reference
```go
func buildSymbolReference(client *scryfall.Client) {
    ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
    defer cancel()
    
    symbolCategories := map[string][]string{
        "Basic Mana": {"{W}", "{U}", "{B}", "{R}", "{G}", "{C}"},
        "Generic": {"{0}", "{1}", "{2}", "{3}", "{X}", "{Y}"},
        "Hybrid": {"{W/U}", "{U/B}", "{B/R}", "{R/G}", "{G/W}"},
        "Phyrexian": {"{W/P}", "{U/P}", "{B/P}", "{R/P}", "{G/P}"},
        "Special": {"{T}", "{Q}", "{E}", "{CHAOS}", "{A}"},
        "Two-Hybrid": {"{2/W}", "{2/U}", "{2/B}", "{2/R}", "{2/G}"},
    }
    
    for category, symbols := range symbolCategories {
        fmt.Printf("\n🔥 %s:\n", category)
        fmt.Printf("%-8s %-6s %-6s %-15s %s\n", "Symbol", "Cost", "CMC", "Colors", "Description")
        fmt.Printf("%-8s %-6s %-6s %-15s %s\n", "------", "----", "---", "------", "-----------")
        
        for _, symbolText := range symbols {
            symbol, err := client.Symbols.GetSymbol(ctx, symbolText)
            if err != nil {
                fmt.Printf("%-8s ❌ Error: %v\n", symbolText, err)
                continue
            }
            
            colors := strings.Join(symbol.Colors, ",")
            if colors == "" {
                colors = "None"
            }
            
            fmt.Printf("%-8s %-6.1f %-6.1f %-15s %s\n",
                symbol.Symbol, symbol.Cost, symbol.Cmc, colors, symbol.English)
        }
    }
}
```

### Download Symbol Images
```go
func downloadSymbolSVG(client *scryfall.Client, symbolText string, filename string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    symbol, err := client.Symbols.GetSymbol(ctx, symbolText)
    if err != nil {
        return fmt.Errorf("error getting symbol: %w", err)
    }
    
    // Download SVG (you'd need to implement HTTP download)
    fmt.Printf("📥 Download %s SVG from: %s\n", symbol.Symbol, symbol.SvgUri)
    fmt.Printf("   Save to: %s\n", filename)
    
    return nil
}

// Usage
downloadSymbolSVG(client, "{R}", "red_mana.svg")
downloadSymbolSVG(client, "{W/U}", "azorius_hybrid.svg")
```

### Validate Mana Cost String
```go
func validateManaCost(client *scryfall.Client, manaCost string) bool {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    symbols := parseManaCost(manaCost)
    
    for _, symbolText := range symbols {
        _, err := client.Symbols.GetSymbol(ctx, symbolText)
        if err != nil {
            fmt.Printf("❌ Invalid symbol: %s\n", symbolText)
            return false
        }
    }
    
    return true
}

// Usage
testCosts := []string{
    "{3}{R}{R}",      // Valid
    "{W/U}{T}",       // Valid
    "{Z}{Invalid}",   // Invalid
}

for _, cost := range testCosts {
    if validateManaCost(client, cost) {
        fmt.Printf("✅ %s is valid\n", cost)
    } else {
        fmt.Printf("❌ %s is invalid\n", cost)
    }
}
```

## URL Encoding

The service automatically handles URL encoding for special characters:

```go
// These all work correctly
symbols := []string{
    "{W/U}",    // Contains slash
    "{2/W}",    // Contains slash  
    "{CHAOS}",  // Long symbol name
    "{∞}",      // Unicode character
    "{½}",      // Unicode fraction
}

for _, sym := range symbols {
    symbol, err := client.Symbols.GetSymbol(ctx, sym)
    // URL encoding is handled automatically
}
```

## Error Handling

```go
symbol, err := client.Symbols.GetSymbol(ctx, "{INVALID}")
if err != nil {
    if apiErr, ok := err.(*errors.ApiError); ok {
        switch apiErr.ErrInfo.Status {
        case 404:
            fmt.Println("Symbol not found")
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

1. **Include braces** - Always use `{symbol}` format, not just `symbol`
2. **Handle Unicode** - Some symbols use Unicode characters
3. **Cache symbols** - Symbol data doesn't change, cache locally
4. **Validate input** - Check symbol format before API calls
5. **Use SVG images** - SVG format scales well for any size
6. **Consider alternatives** - Some symbols have multiple representations

## Common Symbol Collections

### Deck Building Symbols
```go
basicSymbols := []string{"{W}", "{U}", "{B}", "{R}", "{G}"}
hybridSymbols := []string{"{W/U}", "{U/B}", "{B/R}", "{R/G}", "{G/W}"}
utilitySymbols := []string{"{T}", "{Q}", "{X}"}
```

### Un-set Symbols
```go
funSymbols := []string{"{A}", "{TK}", "{CHAOS}"}
```

### Modern Symbols
```go
newSymbols := []string{"{E}", "{C}", "{Q}"}
```

## Related Operations

- Use with **Card Service** to understand mana costs in search results
- Combine with **Catalog Service** for complete Magic reference tools
- Use for deck analysis and mana curve calculations
- Integrate with UI applications for mana symbol display