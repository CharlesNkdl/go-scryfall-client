# Installation Guide

This guide covers how to install and set up the Go Scryfall Client library.

## Prerequisites

- Go 1.21 or later
- Internet connection for API access

## Installation

### Using go get

```bash
go get github.com/CharlesNkdl/go-scryfall-client
```

### Using go mod

Add to your `go.mod`:

```go
require github.com/CharlesNkdl/go-scryfall-client v0.1.0
```

Then run:

```bash
go mod tidy
```

## Quick Setup Verification

Create a simple test file to verify the installation:

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/CharlesNkdl/go-scryfall-client/scryfall"
)

func main() {
    // Create a new client
    client := scryfall.NewClient()
    
    // Create a context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    // Test the connection with a simple card search
    card, err := client.Cards.GetByName(ctx, "Lightning Bolt", false)
    if err != nil {
        log.Fatalf("Failed to fetch card: %v", err)
    }
    
    fmt.Printf("Successfully fetched card: %s\n", card.Name)
    fmt.Printf("Installation verified!\n")
}
```

Run the test:

```bash
go run main.go
```

If you see "Installation verified!" the library is correctly installed and working.

## Dependencies

The library has minimal dependencies:

- `golang.org/x/time/rate` - For API rate limiting compliance

All dependencies are automatically managed by Go modules.

## Next Steps

- Read the [Getting Started Guide](getting-started.md)
- Explore [Service Documentation](services/)
- Check out [Examples](examples.md)