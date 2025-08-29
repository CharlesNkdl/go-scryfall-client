# Examples and Use Cases

This document provides comprehensive examples of how to use the Go Scryfall Client library for various common tasks and use cases.

## Basic Examples

### Simple Card Lookup

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/CharlesNkdl/go-scryfall-client/scryfall"
    "github.com/CharlesNkdl/go-scryfall-client/scryfall/models/request/cards"
)

func main() {
    client := scryfall.NewClient()
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    // Look up a specific card
    card, err := client.Cards.GetByName(ctx, cards.NewExactCardParams("Lightning Bolt"))
    if err != nil {
        log.Fatalf("Error: %v", err)
    }

    fmt.Printf("Card: %s\n", card.Name)
    fmt.Printf("Type: %s\n", card.TypeLine)
    if card.ManaCost != nil {
        fmt.Printf("Mana Cost: %s\n", *card.ManaCost)
    }
    if card.OracleText != nil {
        fmt.Printf("Oracle Text: %s\n", *card.OracleText)
    }
}
```

### Fuzzy Search with Error Handling

```go
func fuzzyCardSearch(client *scryfall.Client, cardName string) {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    // Try fuzzy search to handle typos
    card, err := client.Cards.GetByName(ctx, cards.NewFuzzyCardParams(cardName))
    if err != nil {
        if apiErr, ok := err.(*errors.ApiError); ok {
            switch apiErr.ErrInfo.Status {
            case 404:
                fmt.Printf("❌ Card '%s' not found\n", cardName)
            case 429:
                fmt.Printf("⏱️ Rate limited, please wait\n")
            default:
                fmt.Printf("🚫 API Error: %s\n", apiErr.ErrInfo.Detail)
            }
        } else {
            fmt.Printf("🌐 Network Error: %v\n", err)
        }
        return
    }

    fmt.Printf("✅ Found: %s", card.Name)
    if card.Name != cardName {
        fmt.Printf(" (corrected from '%s')", cardName)
    }
    fmt.Println()
}

// Usage with various typos
fuzzyCardSearch(client, "Lighning Bolt")      // Finds "Lightning Bolt"
fuzzyCardSearch(client, "Ancestral Recal")   // Finds "Ancestral Recall"
fuzzyCardSearch(client, "Blak Lotus")        // Finds "Black Lotus"
```

## Deck Building Tools

### Deck Validator

```go
func validateDeck(client *scryfall.Client, deckList map[string]int) {
    ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
    defer cancel()

    fmt.Println("🔍 Validating deck...")
    
    var totalCards int
    var validCards int
    var totalValue float64

    for cardName, quantity := range deckList {
        totalCards += quantity
        
        card, err := client.Cards.GetByName(ctx, cards.NewExactCardParams(cardName))
        if err != nil {
            fmt.Printf("❌ %dx %s - Not found\n", quantity, cardName)
            continue
        }

        validCards += quantity
        fmt.Printf("✅ %dx %s", quantity, card.Name)

        // Show price if available
        if card.Prices != nil && card.Prices.USD != nil {
            price, _ := strconv.ParseFloat(*card.Prices.USD, 64)
            cardValue := price * float64(quantity)
            totalValue += cardValue
            fmt.Printf(" ($%.2f each, $%.2f total)", price, cardValue)
        }
        
        fmt.Println()
    }

    fmt.Printf("\n📊 Deck Summary:\n")
    fmt.Printf("   Total Cards: %d\n", totalCards)
    fmt.Printf("   Valid Cards: %d\n", validCards)
    fmt.Printf("   Missing Cards: %d\n", totalCards-validCards)
    if totalValue > 0 {
        fmt.Printf("   Estimated Value: $%.2f\n", totalValue)
    }
}

// Usage
deck := map[string]int{
    "Lightning Bolt":    4,
    "Counterspell":      4,
    "Island":           20,
    "Mountain":         20,
    "Nonexistent Card":  1, // This will show as not found
}
validateDeck(client, deck)
```

### Random Deck Generator

```go
func generateRandomDeck(client *scryfall.Client, format string, size int) {
    ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
    defer cancel()

    fmt.Printf("🎲 Generating random %s deck (%d cards)...\n\n", format, size)
    
    // Search constraints for different formats
    var query string
    switch format {
    case "standard":
        query = "legal:standard"
    case "modern":
        query = "legal:modern"
    case "legacy":
        query = "legal:legacy"
    default:
        query = "legal:vintage"
    }

    deck := make(map[string]int)
    
    for len(deck) < size {
        params := &cards.RandomCardParams{
            Query: &query,
        }
        
        card, err := client.Cards.GetRandom(ctx, params)
        if err != nil {
            fmt.Printf("Error getting random card: %v\n", err)
            continue
        }

        // Skip basic lands for variety
        if strings.Contains(strings.ToLower(card.TypeLine), "basic") {
            continue
        }

        // Add to deck or increase quantity
        if deck[card.Name] < 4 { // Max 4 copies
            deck[card.Name]++
            fmt.Printf("Added: %s (%s)\n", card.Name, card.Set)
        }
    }

    fmt.Printf("\n📋 Random %s Deck:\n", format)
    for cardName, quantity := range deck {
        fmt.Printf("%dx %s\n", quantity, cardName)
    }
}

// Usage
generateRandomDeck(client, "modern", 60)
```

## Advanced Search Examples

### Complex Card Search

```go
func advancedCardSearch(client *scryfall.Client) {
    ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
    defer cancel()

    searches := map[string]string{
        "Expensive Red Cards": "color:red usd>100",
        "Cheap Creatures":     "type:creature usd<1",
        "Draw Spells":         "oracle:\"draw a card\" -type:creature",
        "Artifacts in Modern": "type:artifact legal:modern",
        "Legendary Dragons":   "type:legendary type:dragon",
        "Planeswalkers":       "type:planeswalker cmc<=4",
        "High Power/Low CMC":  "power>=4 cmc<=3 type:creature",
    }

    for description, query := range searches {
        fmt.Printf("🔍 %s:\n", description)
        
        params := &cards.SearchParams{
            Query: query,
            Page:  1,
        }
        
        results, err := client.Cards.Search(ctx, params)
        if err != nil {
            fmt.Printf("   ❌ Error: %v\n\n", err)
            continue
        }

        fmt.Printf("   Found %d cards (showing first 5):\n", results.TotalCards)
        
        limit := 5
        if len(results.Data) < limit {
            limit = len(results.Data)
        }
        
        for i := 0; i < limit; i++ {
            card := results.Data[i]
            fmt.Printf("   %d. %s (%s)\n", i+1, card.Name, card.Set)
        }
        
        fmt.Println()
    }
}
```

### Set Analysis

```go
func analyzeSet(client *scryfall.Client, setCode string) {
    ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
    defer cancel()

    // Get set information
    set, err := client.Sets.GetById(ctx, setCode)
    if err != nil {
        log.Printf("Error getting set info: %v", err)
        return
    }

    fmt.Printf("📦 Set Analysis: %s (%s)\n", set.Name, set.Code)
    fmt.Printf("Released: %s\n", set.ReleasedAt)
    fmt.Printf("Total Cards: %d\n\n", set.CardCount)

    // Analyze cards by type
    query := fmt.Sprintf("set:%s", setCode)
    params := &cards.SearchParams{
        Query: query,
        Page:  1,
    }

    results, err := client.Cards.Search(ctx, params)
    if err != nil {
        log.Printf("Error searching set: %v", err)
        return
    }

    // Count by type
    typeCounts := make(map[string]int)
    colorCounts := make(map[string]int)
    rarityCounts := make(map[string]int)

    for _, card := range results.Data {
        // Count primary type
        typeParts := strings.Fields(card.TypeLine)
        if len(typeParts) > 0 {
            primaryType := typeParts[0]
            typeCounts[primaryType]++
        }

        // Count colors
        if len(card.Colors) == 0 {
            colorCounts["Colorless"]++
        } else if len(card.Colors) == 1 {
            colorCounts[string(card.Colors[0])]++
        } else {
            colorCounts["Multicolor"]++
        }

        // Count rarities
        rarityCounts[string(card.Rarity)]++
    }

    fmt.Printf("📊 Type Distribution:\n")
    for cardType, count := range typeCounts {
        fmt.Printf("   %s: %d\n", cardType, count)
    }

    fmt.Printf("\n🎨 Color Distribution:\n")
    for color, count := range colorCounts {
        fmt.Printf("   %s: %d\n", color, count)
    }

    fmt.Printf("\n✨ Rarity Distribution:\n")
    for rarity, count := range rarityCounts {
        fmt.Printf("   %s: %d\n", strings.Title(rarity), count)
    }
}

// Usage
analyzeSet(client, "khm") // Analyze Kaldheim
```

## Data Analysis Examples

### Price Tracking

```go
func trackCardPrices(client *scryfall.Client, cardNames []string) {
    ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
    defer cancel()

    fmt.Printf("💰 Price Tracking Report:\n\n")

    for _, cardName := range cardNames {
        card, err := client.Cards.GetByName(ctx, cards.NewExactCardParams(cardName))
        if err != nil {
            fmt.Printf("❌ %s: Not found\n", cardName)
            continue
        }

        fmt.Printf("💳 %s (%s)\n", card.Name, card.Set)
        
        if card.Prices != nil {
            if card.Prices.USD != nil {
                fmt.Printf("   USD: $%s\n", *card.Prices.USD)
            }
            if card.Prices.EUR != nil {
                fmt.Printf("   EUR: €%s\n", *card.Prices.EUR)
            }
            if card.Prices.Tix != nil {
                fmt.Printf("   MTGO: %s tix\n", *card.Prices.Tix)
            }
        } else {
            fmt.Printf("   No price data available\n")
        }
        
        fmt.Println()
    }
}

// Usage
expensiveCards := []string{
    "Black Lotus",
    "Ancestral Recall", 
    "Mox Ruby",
    "Time Walk",
    "Timetwister",
}
trackCardPrices(client, expensiveCards)
```

### Format Legality Checker

```go
func checkFormatLegality(client *scryfall.Client, cardNames []string, format string) {
    ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
    defer cancel()

    fmt.Printf("⚖️ %s Legality Check:\n\n", strings.Title(format))

    for _, cardName := range cardNames {
        card, err := client.Cards.GetByName(ctx, cards.NewExactCardParams(cardName))
        if err != nil {
            fmt.Printf("❌ %s: Not found\n", cardName)
            continue
        }

        var legal bool
        var status string

        for _, legality := range card.Legalities {
            if strings.ToLower(legality.Format) == strings.ToLower(format) {
                legal = legality.Legality == "legal"
                status = legality.Legality
                break
            }
        }

        icon := "✅"
        if !legal {
            icon = "❌"
        }

        fmt.Printf("%s %s: %s\n", icon, card.Name, status)
    }
}

// Usage
modernDeck := []string{
    "Lightning Bolt",
    "Counterspell", 
    "Black Lotus",    // Illegal in Modern
    "Tarmogoyf",
}
checkFormatLegality(client, modernDeck, "modern")
```

## Utility Functions

### Auto-complete Helper

```go
func cardAutoComplete(client *scryfall.Client, partial string) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    params := &cards.AutocompleteCardParams{
        Query: partial,
    }

    suggestions, err := client.Catalogs.Autocomplete(ctx, params)
    if err != nil {
        fmt.Printf("Error getting suggestions: %v\n", err)
        return
    }

    fmt.Printf("💭 Suggestions for '%s':\n", partial)
    for i, suggestion := range suggestions.Data {
        if i >= 10 { // Limit to 10 suggestions
            break
        }
        fmt.Printf("   %d. %s\n", i+1, suggestion)
    }
}

// Usage
cardAutoComplete(client, "light")  // Shows cards starting with "light"
```

### Card Image Downloader

```go
func downloadCardImage(client *scryfall.Client, cardName string, imageType string) {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    card, err := client.Cards.GetByName(ctx, cards.NewExactCardParams(cardName))
    if err != nil {
        fmt.Printf("Error finding card: %v\n", err)
        return
    }

    if card.ImageUris == nil {
        fmt.Printf("No images available for %s\n", card.Name)
        return
    }

    var imageURL string
    switch imageType {
    case "small":
        imageURL = card.ImageUris.Small
    case "normal":
        imageURL = card.ImageUris.Normal
    case "large":
        imageURL = card.ImageUris.Large
    case "art":
        imageURL = card.ImageUris.ArtCrop
    case "border":
        imageURL = card.ImageUris.BorderCrop
    default:
        imageURL = card.ImageUris.Normal
    }

    fmt.Printf("🖼️ Image URL for %s (%s): %s\n", card.Name, imageType, imageURL)
    // You would implement actual download logic here
}

// Usage
downloadCardImage(client, "Lightning Bolt", "large")
```

## Integration Examples

### Web API Integration

```go
type CardResponse struct {
    Name     string  `json:"name"`
    Type     string  `json:"type"`
    ManaCost *string `json:"mana_cost"`
    Text     *string `json:"oracle_text"`
    ImageURL string  `json:"image_url"`
}

func cardAPIHandler(client *scryfall.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        cardName := r.URL.Query().Get("name")
        if cardName == "" {
            http.Error(w, "Missing 'name' parameter", http.StatusBadRequest)
            return
        }

        ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
        defer cancel()

        card, err := client.Cards.GetByName(ctx, cards.NewExactCardParams(cardName))
        if err != nil {
            if apiErr, ok := err.(*errors.ApiError); ok && apiErr.ErrInfo.Status == 404 {
                http.Error(w, "Card not found", http.StatusNotFound)
            } else {
                http.Error(w, "Internal server error", http.StatusInternalServerError)
            }
            return
        }

        response := CardResponse{
            Name:     card.Name,
            Type:     card.TypeLine,
            ManaCost: card.ManaCost,
            Text:     card.OracleText,
        }

        if card.ImageUris != nil {
            response.ImageURL = card.ImageUris.Normal
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
    }
}

// Usage in main()
func main() {
    client := scryfall.NewClient()
    http.HandleFunc("/card", cardAPIHandler(client))
    log.Println("Starting server on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

### Command Line Tool

```go
func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: card-lookup <card-name>")
        os.Exit(1)
    }

    cardName := strings.Join(os.Args[1:], " ")
    client := scryfall.NewClient()
    
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    // Try exact search first
    card, err := client.Cards.GetByName(ctx, cards.NewExactCardParams(cardName))
    if err != nil {
        // Try fuzzy search
        card, err = client.Cards.GetByName(ctx, cards.NewFuzzyCardParams(cardName))
        if err != nil {
            fmt.Printf("❌ Card '%s' not found\n", cardName)
            os.Exit(1)
        }
        fmt.Printf("🔍 Found similar card: %s\n", card.Name)
    }

    // Display card information
    fmt.Printf("\n═══ %s ═══\n", card.Name)
    fmt.Printf("Type: %s\n", card.TypeLine)
    if card.ManaCost != nil {
        fmt.Printf("Mana Cost: %s\n", *card.ManaCost)
    }
    if card.OracleText != nil {
        fmt.Printf("Oracle Text: %s\n", *card.OracleText)
    }
    fmt.Printf("Set: %s (%s)\n", card.SetName, card.Set)
}
```

## Error Handling Patterns

### Retry Logic

```go
func robustCardLookup(client *scryfall.Client, cardName string, maxRetries int) (*card.CardCore, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
    defer cancel()

    for attempt := 1; attempt <= maxRetries; attempt++ {
        card, err := client.Cards.GetByName(ctx, cards.NewExactCardParams(cardName))
        if err == nil {
            return card, nil
        }

        if apiErr, ok := err.(*errors.ApiError); ok {
            switch apiErr.ErrInfo.Status {
            case 429: // Rate limited
                if attempt < maxRetries {
                    waitTime := time.Duration(attempt) * time.Second
                    fmt.Printf("⏱️ Rate limited, waiting %v before retry %d/%d\n", waitTime, attempt+1, maxRetries)
                    time.Sleep(waitTime)
                    continue
                }
            case 404: // Not found
                return nil, fmt.Errorf("card not found: %s", cardName)
            default:
                return nil, fmt.Errorf("API error: %s", apiErr.ErrInfo.Detail)
            }
        }

        if attempt < maxRetries {
            fmt.Printf("🔄 Attempt %d failed, retrying...\n", attempt)
            time.Sleep(time.Second)
        }
    }

    return nil, fmt.Errorf("failed after %d attempts", maxRetries)
}

// Usage
card, err := robustCardLookup(client, "Lightning Bolt", 3)
if err != nil {
    log.Printf("Failed to get card: %v", err)
    return
}
```

These examples demonstrate the versatility and power of the Go Scryfall Client library for building Magic: The Gathering applications, tools, and integrations.