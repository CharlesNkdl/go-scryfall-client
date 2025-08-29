# Architecture Overview

This document provides a comprehensive overview of the Go Scryfall Client library architecture, design decisions, and implementation details.

## High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                        Client Layer                         │
├─────────────────────────────────────────────────────────────┤
│  scryfall.Client                                           │
│  ├── HTTP Client (with timeout & rate limiting)            │
│  ├── Base URL configuration                                │
│  └── Service instances                                     │
└─────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────┐
│                       Service Layer                        │
├─────────────────────────────────────────────────────────────┤
│  CardService    │ SetService   │ RulingService              │
│  CatalogService │ SymbolService                             │
│                                                             │
│  Each service:                                              │
│  ├── Wraps specific API endpoints                          │
│  ├── Handles request construction                          │
│  ├── Manages parameter validation                          │
│  └── Processes responses                                   │
└─────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────┐
│                      Request Layer                         │
├─────────────────────────────────────────────────────────────┤
│  Parameter Models (request/*)                              │
│  ├── Type-safe parameter objects                          │
│  ├── Validation logic                                     │
│  ├── URL encoding                                         │
│  └── Builder patterns                                     │
└─────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────┐
│                       Model Layer                          │
├─────────────────────────────────────────────────────────────┤
│  Data Models (models/*)                                    │
│  ├── Strongly-typed structs                               │
│  ├── JSON serialization tags                              │
│  ├── Enum types with constants                            │
│  └── Embedded composition                                 │
└─────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────┐
│                       HTTP Layer                           │
├─────────────────────────────────────────────────────────────┤
│  net/http + golang.org/x/time/rate                        │
│  ├── Request creation and execution                       │
│  ├── Response parsing                                     │
│  ├── Error handling                                       │
│  └── Rate limiting                                        │
└─────────────────────────────────────────────────────────────┘
```

## Core Components

### 1. Client (`scryfall.Client`)

The main entry point that orchestrates all operations:

```go
type Client struct {
    httpClient *http.Client     // HTTP client with timeout
    BaseUrl    string           // API base URL
    RateLimiter *rate.Limiter   // Rate limiting (10 req/sec)
    
    // Service instances
    Cards    *services.CardService
    Sets     *services.SetService
    Rulings  *services.RulingService
    Catalogs *services.CatalogService
    Symbols  *services.SymbolService
}
```

**Responsibilities:**
- HTTP client configuration and management
- Rate limiting enforcement
- Service instance creation and injection
- Common request/response handling

### 2. Service Layer

Each service encapsulates related API endpoints:

#### CardService
```go
type CardService struct {
    Client HTTPClient
}

// Methods:
// - GetById(ctx, params) - Get card by Scryfall ID
// - GetByName(ctx, params) - Get card by name (exact/fuzzy)
// - Search(ctx, params) - Advanced card search
// - Autocomplete(ctx, params) - Name autocomplete
// - GetRandom(ctx, params) - Random card
// - GetByCodeNumberLang(ctx, params) - Get by set/collector number
```

#### Other Services
- **SetService**: Set information retrieval
- **RulingService**: Card ruling retrieval  
- **CatalogService**: Reference data catalogs
- **SymbolService**: Mana symbol information

**Design Principles:**
- Single responsibility per service
- Dependency injection via `HTTPClient` interface
- Context-aware operations
- Consistent error handling

### 3. Request Parameter Models

Type-safe parameter objects with validation:

```go
type NamedCardParams struct {
    Exact   *string  // Exact card name
    Fuzzy   *string  // Fuzzy card name
    Set     *string  // Set code
    Format  *string  // Response format
    Version *string  // Image version
    Face    *string  // Card face
    Pretty  *bool    // Pretty JSON
}

func (p *NamedCardParams) Validate() error {
    // Validation logic
}

func (p *NamedCardParams) ToURLValues() (url.Values, error) {
    // URL encoding logic
}
```

**Key Features:**
- Pointer fields for optional parameters
- Comprehensive validation
- Builder pattern support
- URL encoding abstraction

### 4. Data Models

Strongly-typed representations of API responses:

```go
type CardCore struct {
    // Core identifiers
    ID       string `json:"id"`
    OracleID *string `json:"oracle_id,omitempty"`
    
    // Embedded gameplay and print data
    GameplayFields Gameplay `json:"gameplay_fields"`
    PrintFields    Print    `json:"print_fields"`
    
    // Multi-faced card support
    CardFaces []Face `json:"card_faces,omitempty"`
}
```

**Design Patterns:**
- Composition over inheritance
- Embedded structs for logical grouping
- Pointer types for optional fields
- Custom types for enums

## Design Decisions

### 1. Interface-Based Architecture

The library uses interfaces for testability and flexibility:

```go
type HTTPClient interface {
    NewRequest(ctx context.Context, method, path string) (*http.Request, error)
    Do(req *http.Request, v interface{}) error
}
```

**Benefits:**
- Easy mocking for testing
- Potential for custom HTTP client implementations
- Clear separation of concerns
- Better unit test coverage

### 2. Embedded Struct Composition

Card data is organized using embedded structs:

```go
type CardCore struct {
    // Basic fields
    ID       string   `json:"id"`
    OracleID *string  `json:"oracle_id,omitempty"`
    
    // Composed data
    GameplayFields Gameplay `json:"gameplay_fields"`
    PrintFields    Print    `json:"print_fields"`
}

type Gameplay struct {
    Name       string         `json:"name"`
    ManaCost   *string        `json:"mana_cost,omitempty"`
    TypeLine   string         `json:"type_line"`
    OracleText *string        `json:"oracle_text,omitempty"`
    Colors     []models.Color `json:"colors,omitempty"`
    // ... more gameplay fields
}
```

**Rationale:**
- Logical grouping of related fields
- Easier to understand and maintain
- Mirrors Scryfall API structure
- Enables selective field access

### 3. Generic List Type

Using Go generics for type-safe collections:

```go
type List[T any] struct {
    Object     string   `json:"object"`
    Data       []T      `json:"data"`
    HasMore    bool     `json:"has_more"`
    NextPage   *string  `json:"next_page,omitempty"`
    TotalCards *int     `json:"total_cards,omitempty"`
    Warnings   []string `json:"warnings,omitempty"`
}
```

**Advantages:**
- Type safety at compile time
- Consistent pagination handling
- Reduced code duplication
- Better IDE support

### 4. Context-First API Design

All service methods require `context.Context`:

```go
func (s *CardService) GetByName(ctx context.Context, params *NamedCardParams) (*card.CardCore, error)
```

**Benefits:**
- Timeout and cancellation support
- Request tracing capabilities
- Better resource management
- Cloud-native ready

### 5. Rate Limiting Integration

Built-in rate limiting using `golang.org/x/time/rate`:

```go
type Client struct {
    RateLimiter *rate.Limiter // 10 requests per second
    // ...
}

func (c *Client) Do(req *http.Request, v interface{}) error {
    if err := c.RateLimiter.Wait(req.Context()); err != nil {
        return err
    }
    // ... execute request
}
```

**Why Built-in:**
- Scryfall API requirements (10 req/sec max)
- Prevents accidental rate limit violations
- Transparent to users
- Configurable if needed

## Error Handling Strategy

### 1. Structured Error Types

Custom error types for different scenarios:

```go
type ApiError struct {
    ErrInfo ScryfallError
}

type ScryfallError struct {
    Object   string `json:"object"`
    Code     string `json:"code"`
    Status   int    `json:"status"`
    Detail   string `json:"detail"`
    Type     *string `json:"type,omitempty"`
    Warnings []string `json:"warnings,omitempty"`
}
```

### 2. Error Wrapping

Using `fmt.Errorf` with `%w` for error chains:

```go
func (s *CardService) GetByName(ctx context.Context, params *NamedCardParams) (*card.CardCore, error) {
    urlValues, err := params.ToURLValues()
    if err != nil {
        return nil, fmt.Errorf("failed to get URL values: %w", err)
    }
    // ...
}
```

### 3. Status Code Handling

Comprehensive HTTP status code handling:

```go
func (c *Client) Do(req *http.Request, v interface{}) error {
    // Rate limiting check
    if err := c.RateLimiter.Wait(req.Context()); err != nil {
        return err
    }
    
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    // Handle rate limiting
    if resp.StatusCode == http.StatusTooManyRequests {
        // Return structured error
    }
    
    // Handle other error statuses
    if resp.StatusCode < 200 || resp.StatusCode >= 300 {
        // Parse and return API error
    }
    
    // Parse successful response
    return json.NewDecoder(resp.Body).Decode(v)
}
```

## Testing Architecture

### 1. Mock-Based Testing

Using interface mocks for unit testing:

```go
type MockHTTPClient struct {
    DoFunc        func(req *http.Request, v interface{}) error
    NewRequestFunc func(ctx context.Context, method, path string) (*http.Request, error)
}
```

### 2. Table-Driven Tests

Comprehensive test coverage using table-driven patterns:

```go
func TestCardService_GetByName(t *testing.T) {
    tests := []struct {
        name           string
        params         *NamedCardParams
        mockResponse   interface{}
        mockError      error
        expectedResult *card.CardCore
        expectedError  string
    }{
        // Test cases...
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation...
        })
    }
}
```

### 3. Test Organization

```
package_test.go
├── Unit tests for public methods
├── Error condition testing
├── Parameter validation testing
└── Integration-style tests with mocks
```

## Performance Considerations

### 1. Memory Management

- Pointer types for optional fields reduce memory allocation
- Streaming JSON parsing where possible
- Careful handling of large search results

### 2. Connection Pooling

HTTP client configured with connection pooling:

```go
client := &http.Client{
    Timeout: 20 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns:       100,
        MaxIdleConnsPerHost: 10,
        IdleConnTimeout:    90 * time.Second,
    },
}
```

### 3. Rate Limiting Efficiency

Token bucket algorithm for smooth rate limiting:
- Allows burst requests up to limit
- Prevents request queuing under normal load
- Graceful degradation under high load

## Extensibility Points

### 1. Custom HTTP Clients

Users can potentially provide custom HTTP clients:

```go
// If library supports this pattern
client := scryfall.NewClient()
client.HTTPClient = customHTTPClient
```

### 2. Additional Services

New services can be added following existing patterns:

```go
type NewService struct {
    Client HTTPClient
}

func (s *NewService) NewMethod(ctx context.Context, params *NewParams) (*Result, error) {
    // Implementation following established patterns
}
```

### 3. Custom Models

Users can define custom models for specific use cases:

```go
type CustomCard struct {
    card.CardCore
    CustomField string `json:"custom_field"`
}
```

## Security Considerations

### 1. Input Validation

All user inputs are validated:
- Parameter validation in request objects
- URL encoding for special characters
- Length limits on string fields

### 2. Safe JSON Parsing

Using standard library JSON parsing:
- No eval() or unsafe operations
- Structured parsing into known types
- Error handling for malformed responses

### 3. Rate Limiting as Protection

Built-in rate limiting protects both:
- Client applications from being banned
- Scryfall API from overload

## Maintenance and Evolution

### 1. Backwards Compatibility

- New fields added as optional (pointers)
- Existing methods maintain signatures
- Deprecation process for removed features

### 2. API Version Handling

Currently targets Scryfall API version 1.0:
- Can be extended for multiple API versions
- Version-specific models if needed
- Migration helpers for version changes

### 3. Dependency Management

Minimal external dependencies:
- `golang.org/x/time/rate` for rate limiting
- Standard library for HTTP and JSON
- No heavyweight frameworks

## Future Architecture Considerations

### 1. Caching Layer

Potential addition of caching interface:

```go
type CacheInterface interface {
    Get(key string) (interface{}, bool)
    Set(key string, value interface{}, ttl time.Duration)
}
```

### 2. Metrics and Observability

Potential hooks for metrics:

```go
type MetricsCollector interface {
    RecordRequest(endpoint string, duration time.Duration, status int)
    RecordError(endpoint string, error error)
}
```

### 3. Plugin System

Potential plugin architecture for custom behaviors:

```go
type Plugin interface {
    BeforeRequest(req *http.Request) error
    AfterResponse(resp *http.Response) error
}
```

This architecture provides a solid foundation for a Magic: The Gathering API client while maintaining flexibility for future enhancements and user customizations.