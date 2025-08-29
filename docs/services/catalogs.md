# Catalog Service Documentation

The Catalog Service provides access to various catalogs of Magic: The Gathering reference data, such as creature types, keywords, planeswalker types, and more.

## Available Methods

### GetCatalog

Retrieve a specific catalog of Magic data by catalog type.

```go
func (s *CatalogService) GetCatalog(ctx context.Context, catalogType string) (*models.Catalog, error)
```

**Parameters:**
- `ctx`: Context for timeout and cancellation
- `catalogType`: Type of catalog to retrieve

**Example:**
```go
// Get all creature types
catalog, err := client.Catalogs.GetCatalog(ctx, "creature-types")
if err != nil {
    log.Printf("Error: %v", err)
    return
}

fmt.Printf("Creature Types (%d total):\n", len(catalog.Data))
for i, creatureType := range catalog.Data {
    fmt.Printf("%d. %s\n", i+1, creatureType)
}
```

## Available Catalog Types

### creature-types
All creature types in Magic.

```go
catalog, err := client.Catalogs.GetCatalog(ctx, "creature-types")
// Returns: Human, Elf, Goblin, Dragon, Angel, etc.
```

### planeswalker-types
All planeswalker types (subtypes).

```go
catalog, err := client.Catalogs.GetCatalog(ctx, "planeswalker-types")
// Returns: Jace, Chandra, Garruk, Liliana, etc.
```

### land-types
All land types in Magic.

```go
catalog, err := client.Catalogs.GetCatalog(ctx, "land-types")
// Returns: Plains, Island, Swamp, Mountain, Forest, etc.
```

### artifact-types
All artifact types in Magic.

```go
catalog, err := client.Catalogs.GetCatalog(ctx, "artifact-types")
// Returns: Equipment, Vehicle, Treasure, etc.
```

### enchantment-types
All enchantment types in Magic.

```go
catalog, err := client.Catalogs.GetCatalog(ctx, "enchantment-types")
// Returns: Aura, Cartouche, Curse, Saga, etc.
```

### spell-types
All spell types in Magic.

```go
catalog, err := client.Catalogs.GetCatalog(ctx, "spell-types")
// Returns: Adventure, Arcane, Trap, etc.
```

### powers
All power values that appear on creatures.

```go
catalog, err := client.Catalogs.GetCatalog(ctx, "powers")
// Returns: 0, 1, 2, 3, *, X, etc.
```

### toughnesses
All toughness values that appear on creatures.

```go
catalog, err := client.Catalogs.GetCatalog(ctx, "toughnesses")
// Returns: 0, 1, 2, 3, *, X, etc.
```

### loyalties
All loyalty values that appear on planeswalkers.

```go
catalog, err := client.Catalogs.GetCatalog(ctx, "loyalties")
// Returns: 1, 2, 3, 4, 5, 6, X, etc.
```

### watermarks
All watermarks that appear on cards.

```go
catalog, err := client.Catalogs.GetCatalog(ctx, "watermarks")
// Returns: set symbols, guild symbols, etc.
```

### keyword-abilities
All keyword abilities in Magic.

```go
catalog, err := client.Catalogs.GetCatalog(ctx, "keyword-abilities")
// Returns: Flying, Trample, Haste, Vigilance, etc.
```

### keyword-actions
All keyword actions in Magic.

```go
catalog, err := client.Catalogs.GetCatalog(ctx, "keyword-actions")
// Returns: Activate, Attach, Cast, Counter, etc.
```

### ability-words
All ability words in Magic.

```go
catalog, err := client.Catalogs.GetCatalog(ctx, "ability-words")
// Returns: Threshold, Metalcraft, Delirium, etc.
```

## Catalog Data Structure

```go
type Catalog struct {
    Object string   `json:"object"`     // Always "catalog"
    URI    string   `json:"uri"`        // API endpoint for this catalog
    Data   []string `json:"data"`       // Array of catalog entries
}
```

## Complete Examples

### Browse All Creature Types
```go
func browseCreatureTypes(client *scryfall.Client) {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    catalog, err := client.Catalogs.GetCatalog(ctx, "creature-types")
    if err != nil {
        log.Printf("Error: %v", err)
        return
    }
    
    fmt.Printf("🧙 Magic Creature Types (%d total):\n\n", len(catalog.Data))
    
    for i, creatureType := range catalog.Data {
        fmt.Printf("%-3d %-20s", i+1, creatureType)
        if (i+1)%3 == 0 {
            fmt.Println()
        }
    }
    fmt.Println()
}
```

### Validate Card Types
```go
func validateCreatureType(client *scryfall.Client, typeName string) bool {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    catalog, err := client.Catalogs.GetCatalog(ctx, "creature-types")
    if err != nil {
        return false
    }
    
    for _, validType := range catalog.Data {
        if strings.EqualFold(validType, typeName) {
            return true
        }
    }
    return false
}

// Usage
if validateCreatureType(client, "Dragon") {
    fmt.Println("✅ Valid creature type")
} else {
    fmt.Println("❌ Invalid creature type")
}
```

### Compare Keywords Across Time
```go
func analyzeKeywords(client *scryfall.Client) {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    // Get keyword abilities
    abilities, err := client.Catalogs.GetCatalog(ctx, "keyword-abilities")
    if err != nil {
        log.Printf("Error getting abilities: %v", err)
        return
    }
    
    // Get keyword actions
    actions, err := client.Catalogs.GetCatalog(ctx, "keyword-actions")
    if err != nil {
        log.Printf("Error getting actions: %v", err)
        return
    }
    
    // Get ability words
    words, err := client.Catalogs.GetCatalog(ctx, "ability-words")
    if err != nil {
        log.Printf("Error getting ability words: %v", err)
        return
    }
    
    fmt.Printf("📊 Magic Keywords Analysis:\n")
    fmt.Printf("   Keyword Abilities: %d\n", len(abilities.Data))
    fmt.Printf("   Keyword Actions:   %d\n", len(actions.Data))
    fmt.Printf("   Ability Words:     %d\n", len(words.Data))
    fmt.Printf("   Total Keywords:    %d\n", 
        len(abilities.Data) + len(actions.Data) + len(words.Data))
    
    fmt.Printf("\n🎯 Recent Keyword Abilities (last 10):\n")
    start := len(abilities.Data) - 10
    if start < 0 {
        start = 0
    }
    for i := start; i < len(abilities.Data); i++ {
        fmt.Printf("   - %s\n", abilities.Data[i])
    }
}
```

### Build Reference Guide
```go
func buildReferenceGuide(client *scryfall.Client) {
    ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
    defer cancel()
    
    catalogs := map[string]string{
        "creature-types":     "🧙 Creature Types",
        "planeswalker-types": "✨ Planeswalker Types", 
        "land-types":         "🏔️ Land Types",
        "artifact-types":     "⚙️ Artifact Types",
        "enchantment-types":  "✨ Enchantment Types",
        "keyword-abilities":  "🎯 Keyword Abilities",
        "ability-words":      "💫 Ability Words",
    }
    
    for catalogType, title := range catalogs {
        catalog, err := client.Catalogs.GetCatalog(ctx, catalogType)
        if err != nil {
            fmt.Printf("❌ Error getting %s: %v\n", catalogType, err)
            continue
        }
        
        fmt.Printf("\n%s (%d):\n", title, len(catalog.Data))
        
        // Show first 10 entries
        limit := 10
        if len(catalog.Data) < limit {
            limit = len(catalog.Data)
        }
        
        for i := 0; i < limit; i++ {
            fmt.Printf("   %s\n", catalog.Data[i])
        }
        
        if len(catalog.Data) > limit {
            fmt.Printf("   ... and %d more\n", len(catalog.Data)-limit)
        }
    }
}
```

### Search Helper
```go
func findInCatalog(client *scryfall.Client, catalogType, searchTerm string) []string {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    catalog, err := client.Catalogs.GetCatalog(ctx, catalogType)
    if err != nil {
        return nil
    }
    
    var matches []string
    searchLower := strings.ToLower(searchTerm)
    
    for _, item := range catalog.Data {
        if strings.Contains(strings.ToLower(item), searchLower) {
            matches = append(matches, item)
        }
    }
    
    return matches
}

// Usage
matches := findInCatalog(client, "creature-types", "human")
fmt.Printf("Creature types containing 'human': %v\n", matches)
// Output: [Human, Zombie, etc.]
```

## Error Handling

```go
catalog, err := client.Catalogs.GetCatalog(ctx, "invalid-catalog")
if err != nil {
    if apiErr, ok := err.(*errors.ApiError); ok {
        switch apiErr.ErrInfo.Status {
        case 404:
            fmt.Println("Catalog type not found")
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

1. **Cache catalog data** - Catalogs change infrequently, cache them locally
2. **Use exact catalog names** - Catalog types are case-sensitive
3. **Handle large catalogs** - Some catalogs (like creature-types) are very large
4. **Implement search** - For user interfaces, implement search within catalogs
5. **Combine catalogs** - Use multiple catalogs to build comprehensive references
6. **Validate user input** - Use catalogs to validate user-provided type names

## Available Catalog Types Summary

| Catalog Type | Description | Typical Count |
|-------------|-------------|---------------|
| `creature-types` | All creature types | 200+ |
| `planeswalker-types` | Planeswalker subtypes | 40+ |
| `land-types` | All land types | 15+ |
| `artifact-types` | All artifact types | 10+ |
| `enchantment-types` | All enchantment types | 10+ |
| `spell-types` | All spell types | 5+ |
| `powers` | All power values | 20+ |
| `toughnesses` | All toughness values | 20+ |
| `loyalties` | All loyalty values | 10+ |
| `watermarks` | All watermarks | 50+ |
| `keyword-abilities` | All keyword abilities | 100+ |
| `keyword-actions` | All keyword actions | 30+ |
| `ability-words` | All ability words | 30+ |

## Related Services

- Use with **Card Service** to validate search queries
- Combine with **Search** to find cards of specific types
- Use for building deck construction tools and validators