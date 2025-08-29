# Documentation Index

Welcome to the comprehensive documentation for the Go Scryfall Client library. This documentation covers everything you need to know to effectively use this library for Magic: The Gathering applications.

## Quick Start

- **[Installation Guide](installation.md)** - Get up and running quickly
- **[Getting Started](getting-started.md)** - Basic usage and common patterns

## Core Documentation

### API Reference
- **[Services Overview](services/README.md)** - All available services
- **[Card Service](services/cards.md)** - Card search, retrieval, and data
- **[Set Service](services/sets.md)** - Magic set information
- **[Ruling Service](services/rulings.md)** - Official card rulings
- **[Catalog Service](services/catalogs.md)** - Reference data catalogs
- **[Symbol Service](services/symbols.md)** - Mana symbol information

### Data Models
- **[Models and Data Structures](models.md)** - Complete model reference

### Error Handling
- **[Error Handling Guide](error-handling.md)** - Comprehensive error management

## Advanced Topics

- **[Advanced Features](advanced-features.md)** - Rate limiting, caching, concurrency
- **[Examples and Use Cases](examples.md)** - Real-world code examples

## Development

- **[Contributing Guidelines](contributing.md)** - How to contribute to the project
- **[Architecture Overview](architecture.md)** - Library design and implementation

## API Coverage

This library provides Go bindings for the following Scryfall API endpoints:

### Cards
- ✅ Get card by ID
- ✅ Get card by name (exact and fuzzy)
- ✅ Search cards
- ✅ Card autocomplete
- ✅ Random card
- ✅ Get card by set/collector number

### Sets
- ✅ Get set by ID/code

### Rulings
- ✅ Get card rulings

### Catalogs
- ✅ Get all catalog types (creature types, keywords, etc.)

### Symbology  
- ✅ Get mana symbol information

## Quick Navigation

### By User Type

#### **New Users**
1. [Installation](installation.md)
2. [Getting Started](getting-started.md)
3. [Card Service](services/cards.md)
4. [Examples](examples.md)

#### **Experienced Users**
1. [Services Overview](services/README.md)
2. [Advanced Features](advanced-features.md)
3. [Error Handling](error-handling.md)
4. [Models Reference](models.md)

#### **Contributors**
1. [Contributing Guidelines](contributing.md)
2. [Architecture Overview](architecture.md)
3. [Testing Guidelines](contributing.md#testing-guidelines)

### By Use Case

#### **Building a Deck Tool**
- [Card Service](services/cards.md) - Search and validate cards
- [Set Service](services/sets.md) - Validate set codes
- [Examples: Deck Building](examples.md#deck-building-tools)
- [Error Handling](error-handling.md) - Handle missing cards

#### **Creating a Card Database**
- [Advanced Features: Pagination](advanced-features.md#pagination-handling)
- [Advanced Features: Caching](advanced-features.md#caching-strategies)
- [Models Reference](models.md) - Complete data structures
- [Examples: Data Analysis](examples.md#data-analysis-examples)

#### **Building a Web API**
- [Error Handling](error-handling.md) - Proper error responses
- [Advanced Features: Concurrency](advanced-features.md#concurrent-operations)
- [Examples: Web Integration](examples.md#integration-examples)

#### **Price Tracking**
- [Card Service: Search](services/cards.md#search)
- [Models: Prices](models.md#prices)
- [Examples: Price Tracking](examples.md#price-tracking)

#### **Learning Magic Rules**
- [Ruling Service](services/rulings.md) - Official rulings
- [Catalog Service](services/catalogs.md) - Keywords and abilities
- [Symbol Service](services/symbols.md) - Mana symbols

## Common Tasks

### Basic Card Lookup
```go
client := scryfall.NewClient()
card, err := client.Cards.GetByName(ctx, cards.NewExactCardParams("Lightning Bolt"))
```
👉 [Full Example](getting-started.md#searching-for-cards)

### Advanced Search
```go
params := &cards.SearchParams{
    Query: "type:creature power:3 cmc:2",
}
results, err := client.Cards.Search(ctx, params)
```
👉 [Search Documentation](services/cards.md#search)

### Error Handling
```go
if apiErr, ok := err.(*errors.ApiError); ok {
    switch apiErr.ErrInfo.Status {
    case 404:
        // Card not found
    case 429:
        // Rate limited
    }
}
```
👉 [Error Handling Guide](error-handling.md)

### Get Card Rulings
```go
rulings, err := client.Rulings.GetRulings(ctx, cardID)
```
👉 [Ruling Service](services/rulings.md)

## Library Features

### ✅ Type Safety
- Strongly-typed models for all API responses
- Compile-time validation of parameters
- Custom enum types for constrained values

### ✅ Error Handling
- Structured error types with detailed information
- HTTP status code mapping
- Comprehensive error context

### ✅ Rate Limiting
- Built-in rate limiting (10 requests/second)
- Automatic compliance with Scryfall API limits
- Configurable rate limiter

### ✅ Context Support
- All operations support `context.Context`
- Timeout and cancellation support
- Request tracing ready

### ✅ Comprehensive Testing
- >95% test coverage
- Mock-based testing
- Table-driven test patterns

### ✅ Concurrent Safe
- Thread-safe client operations
- Safe for use across goroutines
- Worker pool patterns supported

## API Limits and Guidelines

### Rate Limiting
- **Maximum**: 10 requests per second
- **Built-in**: Automatic rate limiting
- **Burst**: Small burst allowance for occasional spikes

### Request Guidelines
- Use specific searches when possible
- Cache results when appropriate
- Implement retry logic for transient errors
- Use context timeouts

### Best Practices
- Always handle errors appropriately
- Use fuzzy search as fallback for exact search
- Validate parameters before API calls
- Monitor request rates in production

## Getting Help

### Documentation Issues
If you find errors or gaps in the documentation:
1. Check the [GitHub Issues](https://github.com/CharlesNkdl/go-scryfall-client/issues)
2. Create a new issue with the "documentation" label
3. Be specific about what's unclear or missing

### Usage Questions
For help using the library:
1. Check the [Examples](examples.md) for similar use cases
2. Review the specific [Service Documentation](services/)
3. Create a GitHub issue with the "question" label

### Bug Reports
If you encounter bugs:
1. Check [existing issues](https://github.com/CharlesNkdl/go-scryfall-client/issues)
2. Create a detailed bug report
3. Include code examples and error messages

### Feature Requests
For new features:
1. Check the [Scryfall API documentation](https://scryfall.com/docs/api) for endpoint availability
2. Create a feature request issue
3. Describe the use case and expected behavior

## Related Resources

### Scryfall
- [Scryfall API Documentation](https://scryfall.com/docs/api)
- [Scryfall Website](https://scryfall.com/)
- [API Status Page](https://scryfall.statuspage.io/)

### Go Resources
- [Go Documentation](https://golang.org/doc/)
- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go Modules](https://golang.org/ref/mod)

### Magic: The Gathering
- [Official Rules](https://magic.wizards.com/en/rules)
- [Comprehensive Rules](https://magic.wizards.com/en/game-info/gameplay/rules-and-formats/rules)
- [Glossary](https://mtg.fandom.com/wiki/Category:Magic_glossary)

---

*This documentation is for the Go Scryfall Client library. For the official Scryfall API documentation, visit [scryfall.com/docs/api](https://scryfall.com/docs/api).*