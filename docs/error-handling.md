# Error Handling Guide

This guide covers comprehensive error handling patterns and best practices when using the Go Scryfall Client library.

## Error Types

The library provides structured error handling through custom error types that wrap Scryfall API responses and network errors.

### ApiError

The primary error type for API-related errors:

```go
type ApiError struct {
    ErrInfo ScryfallError
}

type ScryfallError struct {
    Object   string                 `json:"object"`   // Always "error"
    Code     string                 `json:"code"`     // Error code
    Status   int                    `json:"status"`   // HTTP status code
    Detail   string                 `json:"detail"`   // Human-readable message
    Type     *string                `json:"type,omitempty"` // Error type
    Warnings []string               `json:"warnings,omitempty"` // API warnings
    Details  map[string]interface{} `json:"details,omitempty"` // Additional details
}

func (e *ApiError) Error() string {
    return e.ErrInfo.Detail
}
```

## Common HTTP Status Codes

### 404 Not Found
The most common error when searching for cards:

```go
card, err := client.Cards.GetByName(ctx, cards.NewExactCardParams("Nonexistent Card"))
if err != nil {
    if apiErr, ok := err.(*errors.ApiError); ok && apiErr.ErrInfo.Status == 404 {
        fmt.Println("Card not found")
        // Try fuzzy search as fallback
        card, err = client.Cards.GetByName(ctx, cards.NewFuzzyCardParams("Nonexistent Card"))
        if err != nil {
            fmt.Println("Card not found even with fuzzy search")
            return
        }
    }
}
```

### 429 Too Many Requests
Rate limiting errors:

```go
card, err := client.Cards.GetByName(ctx, cards.NewExactCardParams("Lightning Bolt"))
if err != nil {
    if apiErr, ok := err.(*errors.ApiError); ok && apiErr.ErrInfo.Status == 429 {
        fmt.Println("Rate limited - please wait before making more requests")
        // The client automatically handles rate limiting, but you might hit limits
        // with very aggressive usage patterns
        return
    }
}
```

### 400 Bad Request
Validation errors:

```go
params := &cards.SearchParams{
    Query: "", // Empty query will cause 400 error
}

results, err := client.Cards.Search(ctx, params)
if err != nil {
    if apiErr, ok := err.(*errors.ApiError); ok && apiErr.ErrInfo.Status == 400 {
        fmt.Printf("Bad request: %s\n", apiErr.ErrInfo.Detail)
        // Handle validation errors
        return
    }
}
```

### 500 Internal Server Error
Server-side errors:

```go
card, err := client.Cards.GetByName(ctx, cards.NewExactCardParams("Lightning Bolt"))
if err != nil {
    if apiErr, ok := err.(*errors.ApiError); ok && apiErr.ErrInfo.Status >= 500 {
        fmt.Println("Server error - try again later")
        return
    }
}
```

## Error Handling Patterns

### Basic Error Checking

Always check for errors and handle them appropriately:

```go
func getCard(client *scryfall.Client, cardName string) (*card.CardCore, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    card, err := client.Cards.GetByName(ctx, cards.NewExactCardParams(cardName))
    if err != nil {
        return nil, fmt.Errorf("failed to get card %s: %w", cardName, err)
    }

    return card, nil
}
```

### Comprehensive Error Handling

Handle different types of errors with specific responses:

```go
func robustCardLookup(client *scryfall.Client, cardName string) (*card.CardCore, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    card, err := client.Cards.GetByName(ctx, cards.NewExactCardParams(cardName))
    if err != nil {
        // Check if it's an API error
        if apiErr, ok := err.(*errors.ApiError); ok {
            switch apiErr.ErrInfo.Status {
            case 404:
                // Try fuzzy search as fallback
                fmt.Printf("Exact match not found for '%s', trying fuzzy search...\n", cardName)
                card, fuzzyErr := client.Cards.GetByName(ctx, cards.NewFuzzyCardParams(cardName))
                if fuzzyErr != nil {
                    return nil, fmt.Errorf("card not found with exact or fuzzy search: %s", cardName)
                }
                fmt.Printf("Found '%s' using fuzzy search\n", card.Name)
                return card, nil
                
            case 429:
                return nil, fmt.Errorf("rate limited - please wait before retrying")
                
            case 400:
                return nil, fmt.Errorf("bad request: %s", apiErr.ErrInfo.Detail)
                
            case 500, 502, 503, 504:
                return nil, fmt.Errorf("server error - try again later: %s", apiErr.ErrInfo.Detail)
                
            default:
                return nil, fmt.Errorf("API error %d: %s", apiErr.ErrInfo.Status, apiErr.ErrInfo.Detail)
            }
        }
        
        // Network or other error
        return nil, fmt.Errorf("network error: %w", err)
    }

    return card, nil
}
```

### Retry Logic

Implement retry logic for transient errors:

```go
func cardLookupWithRetry(client *scryfall.Client, cardName string, maxRetries int) (*card.CardCore, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
    defer cancel()

    var lastErr error
    
    for attempt := 1; attempt <= maxRetries; attempt++ {
        card, err := client.Cards.GetByName(ctx, cards.NewExactCardParams(cardName))
        if err == nil {
            return card, nil
        }

        lastErr = err

        // Check if error is retryable
        if apiErr, ok := err.(*errors.ApiError); ok {
            switch apiErr.ErrInfo.Status {
            case 404:
                // Not found errors shouldn't be retried
                return nil, fmt.Errorf("card not found: %s", cardName)
                
            case 400:
                // Bad request errors shouldn't be retried
                return nil, fmt.Errorf("bad request: %s", apiErr.ErrInfo.Detail)
                
            case 429:
                // Rate limit - wait and retry
                if attempt < maxRetries {
                    waitTime := time.Duration(attempt*attempt) * time.Second // Exponential backoff
                    fmt.Printf("Rate limited, waiting %v before attempt %d/%d\n", waitTime, attempt+1, maxRetries)
                    time.Sleep(waitTime)
                    continue
                }
                
            case 500, 502, 503, 504:
                // Server errors - wait and retry
                if attempt < maxRetries {
                    waitTime := time.Duration(attempt) * time.Second
                    fmt.Printf("Server error, waiting %v before attempt %d/%d\n", waitTime, attempt+1, maxRetries)
                    time.Sleep(waitTime)
                    continue
                }
            }
        }

        // Network errors - wait and retry
        if attempt < maxRetries {
            waitTime := time.Duration(attempt) * time.Second
            fmt.Printf("Network error, waiting %v before attempt %d/%d\n", waitTime, attempt+1, maxRetries)
            time.Sleep(waitTime)
        }
    }

    return nil, fmt.Errorf("failed after %d attempts, last error: %w", maxRetries, lastErr)
}
```

### Graceful Degradation

Handle errors gracefully while still providing useful functionality:

```go
func enrichedCardLookup(client *scryfall.Client, cardName string) {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    // Get basic card info
    card, err := client.Cards.GetByName(ctx, cards.NewExactCardParams(cardName))
    if err != nil {
        fmt.Printf("❌ Could not find card: %s\n", cardName)
        return
    }

    fmt.Printf("✅ Card: %s\n", card.Name)
    fmt.Printf("   Type: %s\n", card.TypeLine)
    
    if card.ManaCost != nil {
        fmt.Printf("   Mana Cost: %s\n", *card.ManaCost)
    }

    // Try to get rulings (optional - don't fail if this errors)
    rulings, err := client.Rulings.GetRulings(ctx, card.ID)
    if err != nil {
        fmt.Printf("⚠️  Could not get rulings: %v\n", err)
    } else if len(rulings) > 0 {
        fmt.Printf("   Rulings: %d available\n", len(rulings))
    }

    // Try to get set info (optional)
    set, err := client.Sets.GetById(ctx, card.Set)
    if err != nil {
        fmt.Printf("⚠️  Could not get set info: %v\n", err)
    } else {
        fmt.Printf("   Set: %s (%s)\n", set.Name, set.ReleasedAt)
    }
}
```

## Context Handling

### Context Cancellation

Handle context cancellation gracefully:

```go
func interruptibleSearch(client *scryfall.Client, query string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    // Set up signal handling for graceful shutdown
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

    // Run search in goroutine
    resultChan := make(chan *models.List[card.CardCore])
    errChan := make(chan error)

    go func() {
        params := &cards.SearchParams{Query: query}
        results, err := client.Cards.Search(ctx, params)
        if err != nil {
            errChan <- err
            return
        }
        resultChan <- results
    }()

    // Wait for result or interruption
    select {
    case results := <-resultChan:
        fmt.Printf("Found %d cards\n", len(results.Data))
        return nil
        
    case err := <-errChan:
        if errors.Is(err, context.Canceled) {
            fmt.Println("Search cancelled by user")
            return nil
        }
        return fmt.Errorf("search failed: %w", err)
        
    case <-sigChan:
        fmt.Println("Received interrupt signal, cancelling search...")
        cancel()
        return nil
    }
}
```

### Context Timeout

Handle timeout errors specifically:

```go
func timeoutAwareSearch(client *scryfall.Client, query string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // Short timeout
    defer cancel()

    params := &cards.SearchParams{Query: query}
    results, err := client.Cards.Search(ctx, params)
    if err != nil {
        if errors.Is(err, context.DeadlineExceeded) {
            fmt.Println("Search timed out - try a more specific query or increase timeout")
            return nil
        }
        return fmt.Errorf("search failed: %w", err)
    }

    fmt.Printf("Found %d cards\n", len(results.Data))
    return nil
}
```

## Validation Errors

### Parameter Validation

Handle validation errors from parameter objects:

```go
func safeCardSearch(client *scryfall.Client, query string) error {
    params := &cards.SearchParams{
        Query: query,
        Page:  1,
    }

    // Parameters are validated when converted to URL values
    _, err := params.ToURLValues()
    if err != nil {
        return fmt.Errorf("invalid search parameters: %w", err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    results, err := client.Cards.Search(ctx, params)
    if err != nil {
        return fmt.Errorf("search failed: %w", err)
    }

    fmt.Printf("Found %d cards\n", len(results.Data))
    return nil
}
```

### Named Card Parameter Validation

```go
func validateNamedCardParams(cardName string, setCode string) error {
    params := &cards.NamedCardParams{
        Exact: &cardName,
        Set:   &setCode,
    }

    if err := params.Validate(); err != nil {
        return fmt.Errorf("invalid parameters: %w", err)
    }

    return nil
}

// Usage
if err := validateNamedCardParams("Lightning Bolt", "lea"); err != nil {
    fmt.Printf("Validation error: %v\n", err)
    return
}
```

## Logging and Monitoring

### Structured Logging

Use structured logging for better error tracking:

```go
import "log/slog"

func loggedCardLookup(client *scryfall.Client, cardName string) (*card.CardCore, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    logger := slog.Default().With("operation", "card_lookup", "card_name", cardName)
    logger.Info("Starting card lookup")

    card, err := client.Cards.GetByName(ctx, cards.NewExactCardParams(cardName))
    if err != nil {
        if apiErr, ok := err.(*errors.ApiError); ok {
            logger.Error("API error during card lookup",
                "status", apiErr.ErrInfo.Status,
                "code", apiErr.ErrInfo.Code,
                "detail", apiErr.ErrInfo.Detail)
        } else {
            logger.Error("Network error during card lookup", "error", err)
        }
        return nil, err
    }

    logger.Info("Card lookup successful", "card_id", card.ID, "set", card.Set)
    return card, nil
}
```

### Error Metrics

Track error rates for monitoring:

```go
type ErrorMetrics struct {
    TotalRequests int
    Errors        map[int]int // Status code -> count
    NetworkErrors int
}

func (m *ErrorMetrics) RecordError(err error) {
    m.TotalRequests++
    
    if apiErr, ok := err.(*errors.ApiError); ok {
        if m.Errors == nil {
            m.Errors = make(map[int]int)
        }
        m.Errors[apiErr.ErrInfo.Status]++
    } else {
        m.NetworkErrors++
    }
}

func (m *ErrorMetrics) ErrorRate() float64 {
    if m.TotalRequests == 0 {
        return 0
    }
    
    totalErrors := m.NetworkErrors
    for _, count := range m.Errors {
        totalErrors += count
    }
    
    return float64(totalErrors) / float64(m.TotalRequests)
}

// Usage
var metrics ErrorMetrics

card, err := client.Cards.GetByName(ctx, cards.NewExactCardParams("Lightning Bolt"))
if err != nil {
    metrics.RecordError(err)
    return err
}

fmt.Printf("Error rate: %.2f%%\n", metrics.ErrorRate()*100)
```

## Testing Error Handling

### Mock Error Responses

```go
func TestErrorHandling(t *testing.T) {
    // Create a mock client that returns specific errors
    mockClient := &MockHTTPClient{
        DoFunc: func(req *http.Request, v interface{}) error {
            return &errors.ApiError{
                ErrInfo: models.ScryfallError{
                    Status: 404,
                    Code:   "not_found",
                    Detail: "Card not found",
                },
            }
        },
    }

    service := &services.CardService{Client: mockClient}
    
    _, err := service.GetByName(context.Background(), cards.NewExactCardParams("Nonexistent"))
    
    if err == nil {
        t.Fatal("Expected error, got nil")
    }

    apiErr, ok := err.(*errors.ApiError)
    if !ok {
        t.Fatalf("Expected ApiError, got %T", err)
    }

    if apiErr.ErrInfo.Status != 404 {
        t.Errorf("Expected status 404, got %d", apiErr.ErrInfo.Status)
    }
}
```

## Best Practices

### 1. Always Check Errors
```go
// Good
card, err := client.Cards.GetByName(ctx, params)
if err != nil {
    return fmt.Errorf("failed to get card: %w", err)
}

// Bad
card, _ := client.Cards.GetByName(ctx, params)
```

### 2. Use Appropriate Context Timeouts
```go
// Good - reasonable timeout
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

// Bad - no timeout
ctx := context.Background()
```

### 3. Handle Specific Error Types
```go
// Good - specific handling
if apiErr, ok := err.(*errors.ApiError); ok {
    switch apiErr.ErrInfo.Status {
    case 404:
        // Handle not found
    case 429:
        // Handle rate limit
    }
}

// Bad - generic handling
if err != nil {
    log.Printf("Error: %v", err)
}
```

### 4. Provide Fallbacks
```go
// Good - try fuzzy search as fallback
card, err := client.Cards.GetByName(ctx, cards.NewExactCardParams(name))
if err != nil {
    if apiErr, ok := err.(*errors.ApiError); ok && apiErr.ErrInfo.Status == 404 {
        card, err = client.Cards.GetByName(ctx, cards.NewFuzzyCardParams(name))
    }
}
```

### 5. Wrap Errors with Context
```go
// Good - provides context
if err != nil {
    return fmt.Errorf("failed to lookup card %s: %w", cardName, err)
}

// Bad - loses context
if err != nil {
    return err
}
```

By following these error handling patterns, your applications will be more robust and provide better user experiences when dealing with various error conditions that can occur when interacting with the Scryfall API.