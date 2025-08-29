# Go Scryfall Client

A comprehensive, type-safe Go client library for the [Scryfall API](https://scryfall.com/docs/api) - the most powerful Magic: The Gathering card database and search engine.

## Features

- ✅ **Complete API Coverage** - All major Scryfall endpoints (cards, sets, rulings, catalogs, symbols)
- ✅ **Type Safety** - Strongly-typed models with comprehensive data structures
- ✅ **Built-in Rate Limiting** - Automatic compliance with Scryfall's API limits (10 req/sec)
- ✅ **Context Support** - Full `context.Context` support for timeouts and cancellation
- ✅ **Error Handling** - Structured error types with detailed HTTP status information
- ✅ **Concurrent Safe** - Thread-safe operations across goroutines
- ✅ **Comprehensive Testing** - >95% test coverage with mock-based testing
- ✅ **Extensive Documentation** - Complete documentation with examples and use cases

## Quick Start

### Installation

```bash
go get github.com/CharlesNkdl/go-scryfall-client
```

### Basic Usage

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
    // Create client
    client := scryfall.NewClient()
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    // Search for a card
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

### Advanced Search

```go
// Complex search with multiple criteria
params := &cards.SearchParams{
    Query: "type:creature power:3 cmc:2 legal:modern",
    Page:  1,
}

results, err := client.Cards.Search(ctx, params)
if err != nil {
    log.Fatalf("Error: %v", err)
}

fmt.Printf("Found %d creatures\n", len(results.Data))
for _, card := range results.Data {
    fmt.Printf("- %s (%s)\n", card.Name, card.Set)
}
```

## 📚 Documentation

**[Complete Documentation](docs/README.md)** - Comprehensive guides and references

### Quick Links
- **[Installation Guide](docs/installation.md)** - Setup and verification
- **[Getting Started](docs/getting-started.md)** - Basic usage patterns
- **[API Reference](docs/services/README.md)** - All services and methods
- **[Examples](docs/examples.md)** - Real-world code examples
- **[Error Handling](docs/error-handling.md)** - Comprehensive error management

### Service Documentation
- **[Card Service](docs/services/cards.md)** - Search, retrieve, and manage card data
- **[Set Service](docs/services/sets.md)** - Magic set information
- **[Ruling Service](docs/services/rulings.md)** - Official card rulings
- **[Catalog Service](docs/services/catalogs.md)** - Reference data (types, keywords, etc.)
- **[Symbol Service](docs/services/symbols.md)** - Mana symbol information

## API Coverage

| Service | Endpoints | Status |
|---------|-----------|--------|
| **Cards** | Search, Get by ID/Name, Random, Autocomplete | ✅ Complete |
| **Sets** | Get by ID/Code | ✅ Complete |
| **Rulings** | Get card rulings | ✅ Complete |
| **Catalogs** | All catalog types | ✅ Complete |
| **Symbols** | Mana symbol lookup | ✅ Complete |

## Usage Examples

### Deck Validation
```go
func validateDeck(client *scryfall.Client, cardNames []string) {
    for _, name := range cardNames {
        card, err := client.Cards.GetByName(ctx, cards.NewExactCardParams(name))
        if err != nil {
            fmt.Printf("❌ %s: Not found\n", name)
            continue
        }
        fmt.Printf("✅ %s: Valid\n", card.Name)
    }
}
```

### Price Tracking
```go
func checkPrice(client *scryfall.Client, cardName string) {
    card, err := client.Cards.GetByName(ctx, cards.NewExactCardParams(cardName))
    if err != nil {
        return
    }
    
    if card.Prices != nil && card.Prices.USD != nil {
        fmt.Printf("%s: $%s\n", card.Name, *card.Prices.USD)
    }
}
```

### Random Card Discovery
```go
func getRandomCreature(client *scryfall.Client) {
    params := &cards.RandomCardParams{
        Query: StringPtr("type:creature"),
    }
    
    card, err := client.Cards.GetRandom(ctx, params)
    if err != nil {
        return
    }
    
    fmt.Printf("Random creature: %s\n", card.Name)
}
```

## Error Handling

The library provides structured error handling with detailed information:

```go
card, err := client.Cards.GetByName(ctx, cards.NewExactCardParams("Nonexistent Card"))
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

## Rate Limiting

The client automatically handles Scryfall's rate limits (10 requests per second):

```go
// These requests will be automatically throttled
for i := 0; i < 20; i++ {
    card, err := client.Cards.GetRandom(ctx)
    // No manual delays needed
}
```

## Contributing

We welcome contributions! Please see our [Contributing Guidelines](docs/contributing.md) for details.

### Development Setup

```bash
# Clone the repository
git clone https://github.com/CharlesNkdl/go-scryfall-client.git
cd go-scryfall-client

# Install dependencies
go mod tidy

# Run tests
go test ./...

# Run tests with coverage
go test ./... -cover
```

## Testing

The library includes comprehensive tests with >95% coverage:

```bash
# Run all tests
go test ./...

# Run specific service tests
go test ./scryfall/services -v

# Run with race detection
go test ./... -race

# Generate coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Scryfall](https://scryfall.com/) for providing the excellent Magic: The Gathering API
- The Go community for the tools and libraries that make this possible
- Magic: The Gathering players and developers who help improve this library

## Related Projects

- [Scryfall API Documentation](https://scryfall.com/docs/api)
- [Other Scryfall API clients](https://scryfall.com/docs/api-clients)

---

**Note**: This is an unofficial library and is not affiliated with Scryfall or Wizards of the Coast.
