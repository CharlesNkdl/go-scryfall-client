# Ruling Service Documentation

The Ruling Service provides access to official card rulings and Oracle text clarifications from Wizards of the Coast.

## Available Methods

### GetRulings

Retrieve all rulings for a specific card by its Scryfall ID.

```go
func (s *RulingService) GetRulings(ctx context.Context, cardID string) ([]models.Ruling, error)
```

**Parameters:**
- `ctx`: Context for timeout and cancellation
- `cardID`: Scryfall ID of the card

**Example:**
```go
// Get rulings for a specific card
cardID := "60a25a91-40ba-4796-b9f2-cff5ac2e2c8a" // Lightning Bolt
rulings, err := client.Rulings.GetRulings(ctx, cardID)
if err != nil {
    log.Printf("Error: %v", err)
    return
}

fmt.Printf("Rulings for card ID %s:\n", cardID)
for i, ruling := range rulings {
    fmt.Printf("\nRuling %d:\n", i+1)
    fmt.Printf("  Source: %s\n", ruling.Source)
    fmt.Printf("  Published: %s\n", ruling.Published_at)
    fmt.Printf("  Comment: %s\n", ruling.Comment)
}
```

## Ruling Information

Each ruling contains the following information:

```go
type Ruling struct {
    Object       string `json:"object"`        // Always "ruling"
    Oracle_id    string `json:"oracle_id"`     // Oracle ID (same for all printings)
    Source       string `json:"source"`        // Source of the ruling
    Published_at string `json:"published_at"`  // Publication date (YYYY-MM-DD)
    Comment      string `json:"comment"`       // The ruling text
}
```

### Ruling Sources

Rulings can come from different sources:

- **"wotc"** - Official Wizards of the Coast rulings
- **"scryfall"** - Scryfall-specific clarifications

## Getting Card ID for Rulings

To get rulings, you need the card's Scryfall ID. Here are common ways to obtain it:

### From Card Search
```go
// First, get the card
card, err := client.Cards.GetByName(ctx, cards.NewExactCardParams("Lightning Bolt"))
if err != nil {
    log.Printf("Error finding card: %v", err)
    return
}

// Then get its rulings
rulings, err := client.Rulings.GetRulings(ctx, card.ID)
if err != nil {
    log.Printf("Error getting rulings: %v", err)
    return
}

fmt.Printf("Rulings for %s:\n", card.Name)
for _, ruling := range rulings {
    fmt.Printf("- %s\n", ruling.Comment)
}
```

### From Oracle ID
If you have the Oracle ID (same for all printings of a card), you can still get rulings:

```go
// Oracle IDs are the same across all printings of a card
// So any printing's ID will work for rulings
card, err := client.Cards.GetByName(ctx, cards.NewExactCardParams("Lightning Bolt"))
if err != nil {
    return
}

rulings, err := client.Rulings.GetRulings(ctx, card.ID)
// These rulings apply to ALL printings of Lightning Bolt
```

## Complete Example

Here's a complete example that finds a card and displays its rulings:

```go
func displayCardRulings(client *scryfall.Client, cardName string) {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    // Find the card first
    fmt.Printf("Looking up card: %s\n", cardName)
    card, err := client.Cards.GetByName(ctx, cards.NewExactCardParams(cardName))
    if err != nil {
        fmt.Printf("❌ Card not found: %v\n", err)
        return
    }
    
    fmt.Printf("✅ Found: %s\n", card.Name)
    fmt.Printf("   Set: %s\n", card.Set)
    fmt.Printf("   ID: %s\n", card.ID)
    
    // Get rulings for the card
    rulings, err := client.Rulings.GetRulings(ctx, card.ID)
    if err != nil {
        fmt.Printf("❌ Error getting rulings: %v\n", err)
        return
    }
    
    if len(rulings) == 0 {
        fmt.Printf("📝 No rulings found for this card.\n")
        return
    }
    
    fmt.Printf("\n📜 Rulings (%d total):\n", len(rulings))
    for i, ruling := range rulings {
        fmt.Printf("\n%d. [%s] %s\n", i+1, ruling.Source, ruling.Published_at)
        fmt.Printf("   %s\n", ruling.Comment)
    }
}

// Usage
displayCardRulings(client, "Necropotence")
displayCardRulings(client, "Black Lotus")
displayCardRulings(client, "Oko, Thief of Crowns")
```

## Advanced Use Cases

### Batch Ruling Lookup
```go
func getRulingsForDeck(client *scryfall.Client, cardNames []string) {
    ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
    defer cancel()
    
    for _, cardName := range cardNames {
        // Get card
        card, err := client.Cards.GetByName(ctx, cards.NewExactCardParams(cardName))
        if err != nil {
            fmt.Printf("❌ %s: Card not found\n", cardName)
            continue
        }
        
        // Get rulings
        rulings, err := client.Rulings.GetRulings(ctx, card.ID)
        if err != nil {
            fmt.Printf("❌ %s: Error getting rulings\n", cardName)
            continue
        }
        
        if len(rulings) > 0 {
            fmt.Printf("📜 %s: %d ruling(s)\n", cardName, len(rulings))
            for _, ruling := range rulings {
                fmt.Printf("   - %s\n", ruling.Comment)
            }
        } else {
            fmt.Printf("📝 %s: No rulings\n", cardName)
        }
        fmt.Println()
    }
}

// Usage
complexCards := []string{
    "Necropotence",
    "Oko, Thief of Crowns", 
    "Teferi, Time Raveler",
    "Arcum's Astrolabe",
}
getRulingsForDeck(client, complexCards)
```

### Filter Rulings by Source
```go
func getWOTCRulings(client *scryfall.Client, cardName string) {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    card, err := client.Cards.GetByName(ctx, cards.NewExactCardParams(cardName))
    if err != nil {
        log.Printf("Error: %v", err)
        return
    }
    
    rulings, err := client.Rulings.GetRulings(ctx, card.ID)
    if err != nil {
        log.Printf("Error: %v", err)
        return
    }
    
    fmt.Printf("Official WotC rulings for %s:\n", card.Name)
    for _, ruling := range rulings {
        if ruling.Source == "wotc" {
            fmt.Printf("- [%s] %s\n", ruling.Published_at, ruling.Comment)
        }
    }
}
```

### Check for Recent Rulings
```go
func hasRecentRulings(client *scryfall.Client, cardName string, sinceDate string) bool {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    card, err := client.Cards.GetByName(ctx, cards.NewExactCardParams(cardName))
    if err != nil {
        return false
    }
    
    rulings, err := client.Rulings.GetRulings(ctx, card.ID)
    if err != nil {
        return false
    }
    
    for _, ruling := range rulings {
        if ruling.Published_at >= sinceDate {
            return true
        }
    }
    return false
}

// Usage
if hasRecentRulings(client, "Oko, Thief of Crowns", "2020-01-01") {
    fmt.Println("Card has recent rulings")
}
```

## Error Handling

```go
rulings, err := client.Rulings.GetRulings(ctx, cardID)
if err != nil {
    if apiErr, ok := err.(*errors.ApiError); ok {
        switch apiErr.ErrInfo.Status {
        case 404:
            fmt.Println("Card not found or has no rulings")
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

1. **Get card ID first** - Always retrieve the card to get its ID before getting rulings
2. **Handle empty results** - Not all cards have rulings
3. **Check ruling sources** - Filter by source if you only want official rulings
4. **Use context timeouts** - Ruling requests can be slower for complex cards
5. **Consider caching** - Rulings don't change frequently
6. **Batch efficiently** - Use appropriate delays between requests for bulk operations

## Related Operations

### Combined Card and Ruling Display
```go
func fullCardInfo(client *scryfall.Client, cardName string) {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    // Get card
    card, err := client.Cards.GetByName(ctx, cards.NewExactCardParams(cardName))
    if err != nil {
        log.Printf("Error: %v", err)
        return
    }
    
    // Display card info
    fmt.Printf("═══ %s ═══\n", card.Name)
    fmt.Printf("Type: %s\n", card.TypeLine)
    if card.ManaCost != nil {
        fmt.Printf("Mana Cost: %s\n", *card.ManaCost)
    }
    if card.OracleText != nil {
        fmt.Printf("Oracle Text: %s\n", *card.OracleText)
    }
    
    // Get and display rulings
    rulings, err := client.Rulings.GetRulings(ctx, card.ID)
    if err != nil {
        fmt.Printf("Could not get rulings: %v\n", err)
        return
    }
    
    if len(rulings) > 0 {
        fmt.Printf("\n📜 Rulings:\n")
        for _, ruling := range rulings {
            fmt.Printf("• %s\n", ruling.Comment)
        }
    }
}
```