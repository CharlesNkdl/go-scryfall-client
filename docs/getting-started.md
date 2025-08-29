# Getting Started Guide

This guide will help you get up and running with the Go Scryfall Client library.

## Basic Usage

### Creating a Client

```go
package main

import (
    "context"
    "time"
    
    "github.com/CharlesNkdl/go-scryfall-client/scryfall"
)

func main() {
    // Create a new client with default configuration
    client := scryfall.NewClient()
    
    // The client automatically handles:
    // - Rate limiting (as required by Scryfall API)
    // - HTTP timeouts (20 seconds default)
    // - Proper User-Agent headers
    // - JSON response parsing
}
```

### Context Usage

Always use context for API calls to handle timeouts and cancellation:

```go
// Create a context with timeout
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

// Use the context in API calls
card, err := client.Cards.GetByName(ctx, "Lightning Bolt", false)
if err != nil {
    // Handle error
    return
}
```

## Common Operations

### Searching for Cards

#### By Exact Name

```go
card, err := client.Cards.GetByName(ctx, "Lightning Bolt", false)
if err != nil {
    log.Printf("Error: %v", err)
    return
}

fmt.Printf("Found card: %s\n", card.Name)
fmt.Printf("Mana cost: %s\n", *card.ManaCost)
fmt.Printf("Type: %s\n", card.TypeLine)
```

#### By Fuzzy Name (with typos)

```go
// This will find "Lightning Bolt" even with the typo
card, err := client.Cards.GetByName(ctx, "Lighning Bolt", true)
if err != nil {
    log.Printf("Error: %v", err)
    return
}

fmt.Printf("Found card: %s\n", card.Name)
```

#### By Card ID

```go
params := &cards.IdCardParams{
    Id: "60a25a91-40ba-4796-b9f2-cff5ac2e2c8a",
}

card, err := client.Cards.GetById(ctx, params)
if err != nil {
    log.Printf("Error: %v", err)
    return
}

fmt.Printf("Found card: %s\n", card.Name)
```

#### Advanced Search

```go
params := &cards.SearchParams{
    Query: "type:creature power:3",
    Page:  1,
}

searchResult, err := client.Cards.Search(ctx, params)
if err != nil {
    log.Printf("Error: %v", err)
    return
}

fmt.Printf("Found %d cards\n", len(searchResult.Data))
for _, card := range searchResult.Data {
    fmt.Printf("- %s\n", card.Name)
}
```

### Working with Sets

```go
// Get a specific set
set, err := client.Sets.GetById(ctx, "khm")
if err != nil {
    log.Printf("Error: %v", err)
    return
}

fmt.Printf("Set: %s (%s)\n", set.Name, set.Code)
fmt.Printf("Released: %s\n", set.ReleasedAt)
```

### Getting Card Rulings

```go
// Get rulings for a card
rulings, err := client.Rulings.GetRulings(ctx, "60a25a91-40ba-4796-b9f2-cff5ac2e2c8a")
if err != nil {
    log.Printf("Error: %v", err)
    return
}

for _, ruling := range rulings.Data {
    fmt.Printf("Ruling: %s\n", ruling.Comment)
    fmt.Printf("Source: %s\n", ruling.Source)
    fmt.Printf("Date: %s\n\n", ruling.Published_at)
}
```

## Error Handling

The library provides structured error handling:

```go
card, err := client.Cards.GetByName(ctx, "Nonexistent Card", false)
if err != nil {
    // Check if it's an API error
    if apiErr, ok := err.(*errors.ApiError); ok {
        fmt.Printf("API Error %d: %s\n", apiErr.ErrInfo.Status, apiErr.ErrInfo.Detail)
        
        // Handle specific errors
        switch apiErr.ErrInfo.Status {
        case 404:
            fmt.Println("Card not found")
        case 429:
            fmt.Println("Rate limited - wait and retry")
        default:
            fmt.Printf("API error: %s\n", apiErr.ErrInfo.Detail)
        }
    } else {
        // Network or other error
        fmt.Printf("Network error: %v\n", err)
    }
    return
}
```

## Rate Limiting

The client automatically respects Scryfall's rate limits:

- Maximum 10 requests per second
- Built-in rate limiter prevents violations
- No manual rate limiting needed

```go
// These calls will be automatically rate-limited
for i := 0; i < 20; i++ {
    card, err := client.Cards.GetRandom(ctx)
    if err != nil {
        log.Printf("Error: %v", err)
        continue
    }
    fmt.Printf("Random card %d: %s\n", i+1, card.Name)
    // No need to add delays - rate limiting is automatic
}
```

## Best Practices

### 1. Always Use Context

```go
// Good
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
card, err := client.Cards.GetByName(ctx, "Lightning Bolt", false)

// Bad - no context timeout
card, err := client.Cards.GetByName(context.Background(), "Lightning Bolt", false)
```

### 2. Handle Errors Properly

```go
// Good
card, err := client.Cards.GetByName(ctx, "Lightning Bolt", false)
if err != nil {
    log.Printf("Failed to get card: %v", err)
    return
}

// Bad - ignoring errors
card, _ := client.Cards.GetByName(ctx, "Lightning Bolt", false)
```

### 3. Use Specific Parameters

```go
// Good - specific search parameters
params := &cards.NamedCardParams{
    Exact: &cardName,
    Set:   &setCode,
}

// Less ideal - relying on defaults
card, err := client.Cards.GetByName(ctx, cardName, false)
```

## Next Steps

- Explore [Service Documentation](services/) for detailed API coverage
- Check [Examples](examples.md) for more complex use cases
- Read [Error Handling](error-handling.md) for comprehensive error management
- See [Advanced Features](advanced-features.md) for power user features