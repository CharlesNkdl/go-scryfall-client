# Services Overview

The Go Scryfall Client provides several services to interact with different parts of the Scryfall API. Each service is specialized for a specific type of data or operation.

## Available Services

### [Card Service](cards.md)
Access to all card-related endpoints including search, retrieval by various identifiers, autocomplete, and random cards.

- **Namespace**: `client.Cards`
- **Primary Use**: Card data retrieval and search
- **Key Methods**: `GetByName`, `GetById`, `Search`, `GetRandom`, `Autocomplete`

### [Set Service](sets.md)
Retrieve information about Magic: The Gathering sets and releases.

- **Namespace**: `client.Sets`
- **Primary Use**: Set information and metadata
- **Key Methods**: `GetById`

### [Ruling Service](rulings.md)
Access card rulings and oracle text clarifications.

- **Namespace**: `client.Rulings`
- **Primary Use**: Official card rulings and clarifications
- **Key Methods**: `GetRulings`

### [Catalog Service](catalogs.md)
Access to various catalogs of Magic data like creature types, keywords, etc.

- **Namespace**: `client.Catalogs`
- **Primary Use**: Reference data and metadata
- **Key Methods**: `GetCreatureTypes`, `GetKeywords`, `GetArtifactTypes`, etc.

### [Symbol Service](symbols.md)
Retrieve mana symbol information and imagery.

- **Namespace**: `client.Symbols`
- **Primary Use**: Mana symbol data and SVG images
- **Key Methods**: `GetSymbol`

## Common Patterns

### Service Access

All services are accessed through the main client:

```go
client := scryfall.NewClient()

// Access different services
card, err := client.Cards.GetByName(ctx, "Lightning Bolt", false)
set, err := client.Sets.GetById(ctx, "khm")
rulings, err := client.Rulings.GetRulings(ctx, cardId)
catalog, err := client.Catalogs.GetCreatureTypes(ctx)
symbol, err := client.Symbols.GetSymbol(ctx, "{R}")
```

### Context Usage

All service methods require a context:

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

// Use context with any service method
result, err := client.ServiceName.Method(ctx, parameters...)
```

### Parameter Objects

Most service methods use parameter objects for type safety and validation:

```go
// Cards service with parameters
params := &cards.SearchParams{
    Query: "type:creature",
    Page:  1,
}
result, err := client.Cards.Search(ctx, params)

// Some methods have convenience wrappers
card, err := client.Cards.GetByName(ctx, "Lightning Bolt", false)
```

### Error Handling

All services return consistent error types:

```go
result, err := client.Cards.GetByName(ctx, "Nonexistent Card", false)
if err != nil {
    if apiErr, ok := err.(*errors.ApiError); ok {
        // API-specific error with status code and details
        fmt.Printf("API Error %d: %s\n", apiErr.ErrInfo.Status, apiErr.ErrInfo.Detail)
    } else {
        // Network or other error
        fmt.Printf("Error: %v\n", err)
    }
}
```

## Rate Limiting

All services automatically respect Scryfall's rate limits:

- Maximum 10 requests per second
- Built-in rate limiter shared across all services
- Automatic throttling prevents API violations

## Authentication

Scryfall API doesn't require authentication. All services work without API keys or tokens.

## Service-Specific Documentation

For detailed information about each service, including all available methods, parameters, and examples, see the individual service documentation files in this directory.