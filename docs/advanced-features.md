# Advanced Features

This guide covers advanced features and patterns for power users of the Go Scryfall Client library.

## Rate Limiting

The client includes built-in rate limiting to comply with Scryfall's API requirements (maximum 10 requests per second).

### Understanding Rate Limiting

```go
// The client automatically handles rate limiting
client := scryfall.NewClient()

// These requests will be automatically throttled
for i := 0; i < 20; i++ {
    card, err := client.Cards.GetRandom(ctx)
    if err != nil {
        log.Printf("Error: %v", err)
        continue
    }
    fmt.Printf("Random card %d: %s\n", i+1, card.Name)
    // No manual delays needed - rate limiter handles this
}
```

### Custom Rate Limiting

You can customize the rate limiter if needed:

```go
import "golang.org/x/time/rate"

// Create client with custom rate limiter
client := scryfall.NewClient()
// More conservative: 5 requests per second with burst of 3
client.RateLimiter = rate.NewLimiter(5, 3)

// More aggressive: 10 requests per second with burst of 10 
// (this is the default, don't exceed Scryfall's limits)
client.RateLimiter = rate.NewLimiter(10, 10)
```

### Handling Rate Limit Errors

Even with built-in rate limiting, you might encounter 429 errors under high load:

```go
func robustRequest(client *scryfall.Client, cardName string) (*card.CardCore, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
    defer cancel()

    maxRetries := 3
    baseDelay := time.Second

    for attempt := 0; attempt < maxRetries; attempt++ {
        card, err := client.Cards.GetByName(ctx, cards.NewExactCardParams(cardName))
        if err == nil {
            return card, nil
        }

        if apiErr, ok := err.(*errors.ApiError); ok && apiErr.ErrInfo.Status == 429 {
            if attempt < maxRetries-1 {
                delay := baseDelay * time.Duration(1<<attempt) // Exponential backoff
                fmt.Printf("Rate limited, waiting %v before retry %d/%d\n", delay, attempt+2, maxRetries)
                time.Sleep(delay)
                continue
            }
        }

        return nil, err
    }

    return nil, fmt.Errorf("failed after %d attempts", maxRetries)
}
```

## Concurrent Operations

### Safe Concurrent Usage

The client is safe for concurrent use across goroutines:

```go
func concurrentCardLookup(client *scryfall.Client, cardNames []string) {
    ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
    defer cancel()

    type result struct {
        Name string
        Card *card.CardCore
        Err  error
    }

    results := make(chan result, len(cardNames))
    var wg sync.WaitGroup

    // Launch goroutines for concurrent lookups
    for _, name := range cardNames {
        wg.Add(1)
        go func(cardName string) {
            defer wg.Done()
            
            card, err := client.Cards.GetByName(ctx, cards.NewExactCardParams(cardName))
            results <- result{
                Name: cardName,
                Card: card,
                Err:  err,
            }
        }(name)
    }

    // Close results channel when all goroutines complete
    go func() {
        wg.Wait()
        close(results)
    }()

    // Process results
    for res := range results {
        if res.Err != nil {
            fmt.Printf("❌ %s: %v\n", res.Name, res.Err)
        } else {
            fmt.Printf("✅ %s: Found\n", res.Name)
        }
    }
}

// Usage
cardNames := []string{
    "Lightning Bolt", "Counterspell", "Giant Growth", 
    "Dark Ritual", "Healing Salve", "Ancestral Recall",
}
concurrentCardLookup(client, cardNames)
```

### Worker Pool Pattern

For high-throughput operations, use a worker pool:

```go
func workerPoolCardLookup(client *scryfall.Client, cardNames []string, numWorkers int) {
    ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
    defer cancel()

    // Input channel
    jobs := make(chan string, len(cardNames))
    
    // Result channel
    type result struct {
        Name string
        Card *card.CardCore
        Err  error
    }
    results := make(chan result, len(cardNames))

    // Start workers
    var wg sync.WaitGroup
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func(workerID int) {
            defer wg.Done()
            fmt.Printf("Worker %d started\n", workerID)
            
            for cardName := range jobs {
                card, err := client.Cards.GetByName(ctx, cards.NewExactCardParams(cardName))
                results <- result{
                    Name: cardName,
                    Card: card,
                    Err:  err,
                }
            }
            
            fmt.Printf("Worker %d finished\n", workerID)
        }(i)
    }

    // Send jobs
    go func() {
        for _, name := range cardNames {
            jobs <- name
        }
        close(jobs)
    }()

    // Close results when all workers finish
    go func() {
        wg.Wait()
        close(results)
    }()

    // Process results
    successCount := 0
    errorCount := 0
    
    for res := range results {
        if res.Err != nil {
            fmt.Printf("❌ %s: %v\n", res.Name, res.Err)
            errorCount++
        } else {
            fmt.Printf("✅ %s: %s (%s)\n", res.Name, res.Card.Name, res.Card.Set)
            successCount++
        }
    }

    fmt.Printf("\n📊 Summary: %d successful, %d errors\n", successCount, errorCount)
}

// Usage with 3 workers
workerPoolCardLookup(client, cardNames, 3)
```

## Advanced Search Patterns

### Complex Query Building

Build complex search queries programmatically:

```go
type SearchBuilder struct {
    conditions []string
}

func NewSearchBuilder() *SearchBuilder {
    return &SearchBuilder{}
}

func (sb *SearchBuilder) AddCondition(condition string) *SearchBuilder {
    sb.conditions = append(sb.conditions, condition)
    return sb
}

func (sb *SearchBuilder) Colors(colors []string) *SearchBuilder {
    if len(colors) > 0 {
        colorStr := strings.Join(colors, "")
        sb.conditions = append(sb.conditions, fmt.Sprintf("color:%s", colorStr))
    }
    return sb
}

func (sb *SearchBuilder) CMC(operator string, value int) *SearchBuilder {
    sb.conditions = append(sb.conditions, fmt.Sprintf("cmc%s%d", operator, value))
    return sb
}

func (sb *SearchBuilder) Type(cardType string) *SearchBuilder {
    sb.conditions = append(sb.conditions, fmt.Sprintf("type:%s", cardType))
    return sb
}

func (sb *SearchBuilder) Legal(format string) *SearchBuilder {
    sb.conditions = append(sb.conditions, fmt.Sprintf("legal:%s", format))
    return sb
}

func (sb *SearchBuilder) Oracle(text string) *SearchBuilder {
    sb.conditions = append(sb.conditions, fmt.Sprintf("oracle:\"%s\"", text))
    return sb
}

func (sb *SearchBuilder) Build() string {
    return strings.Join(sb.conditions, " ")
}

// Usage
query := NewSearchBuilder().
    Type("creature").
    Colors([]string{"R", "G"}).
    CMC("<=", 4).
    Legal("modern").
    Oracle("trample").
    Build()

fmt.Printf("Query: %s\n", query)
// Output: type:creature color:RG cmc<=4 legal:modern oracle:"trample"

params := &cards.SearchParams{Query: query}
results, err := client.Cards.Search(ctx, params)
```

### Pagination Handling

Handle large search results with pagination:

```go
func searchAllPages(client *scryfall.Client, query string) ([]*card.CardCore, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
    defer cancel()

    var allCards []*card.CardCore
    page := 1

    for {
        params := &cards.SearchParams{
            Query: query,
            Page:  page,
        }

        results, err := client.Cards.Search(ctx, params)
        if err != nil {
            return nil, fmt.Errorf("error on page %d: %w", page, err)
        }

        // Add cards from this page
        for _, card := range results.Data {
            cardCopy := card // Important: copy the card to avoid reference issues
            allCards = append(allCards, &cardCopy)
        }

        fmt.Printf("Processed page %d: %d cards (total so far: %d)\n", 
            page, len(results.Data), len(allCards))

        // Check if there are more pages
        if !results.HasMore {
            break
        }

        page++
        
        // Safety check to prevent infinite loops
        if page > 100 {
            return nil, fmt.Errorf("too many pages, stopping at page %d", page)
        }

        // Small delay to be nice to the API
        time.Sleep(100 * time.Millisecond)
    }

    return allCards, nil
}

// Usage
allCreatures, err := searchAllPages(client, "type:creature legal:modern")
if err != nil {
    log.Printf("Error: %v", err)
    return
}

fmt.Printf("Found %d total creatures legal in Modern\n", len(allCreatures))
```

### Search Result Filtering

Filter and process search results:

```go
func advancedCreatureAnalysis(client *scryfall.Client) {
    ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
    defer cancel()

    // Search for creatures
    params := &cards.SearchParams{
        Query: "type:creature legal:modern",
    }

    results, err := client.Cards.Search(ctx, params)
    if err != nil {
        log.Printf("Error: %v", err)
        return
    }

    // Analyze creatures by power/toughness
    powerDistribution := make(map[string]int)
    colorDistribution := make(map[string]int)
    expensiveCards := []*card.CardCore{}

    for _, creature := range results.Data {
        // Power distribution
        power := "Unknown"
        if creature.Power != nil {
            power = *creature.Power
        }
        powerDistribution[power]++

        // Color distribution
        if len(creature.Colors) == 0 {
            colorDistribution["Colorless"]++
        } else if len(creature.Colors) == 1 {
            colorDistribution[string(creature.Colors[0])]++
        } else {
            colorDistribution["Multicolor"]++
        }

        // Find expensive cards
        if creature.Prices != nil && creature.Prices.USD != nil {
            if price, err := strconv.ParseFloat(*creature.Prices.USD, 64); err == nil && price > 10.0 {
                expensiveCards = append(expensiveCards, &creature)
            }
        }
    }

    // Display analysis
    fmt.Printf("🔍 Creature Analysis (%d cards analyzed):\n\n", len(results.Data))

    fmt.Printf("💪 Power Distribution:\n")
    for power, count := range powerDistribution {
        fmt.Printf("   %s: %d\n", power, count)
    }

    fmt.Printf("\n🎨 Color Distribution:\n")
    for color, count := range colorDistribution {
        fmt.Printf("   %s: %d\n", color, count)
    }

    fmt.Printf("\n💰 Expensive Creatures (>$10):\n")
    sort.Slice(expensiveCards, func(i, j int) bool {
        priceI, _ := strconv.ParseFloat(*expensiveCards[i].Prices.USD, 64)
        priceJ, _ := strconv.ParseFloat(*expensiveCards[j].Prices.USD, 64)
        return priceI > priceJ
    })

    for i, creature := range expensiveCards {
        if i >= 10 { // Show top 10
            break
        }
        fmt.Printf("   %s: $%s\n", creature.Name, *creature.Prices.USD)
    }
}
```

## Caching Strategies

### In-Memory Caching

Implement caching to reduce API calls:

```go
import (
    "sync"
    "time"
)

type CachedClient struct {
    client *scryfall.Client
    cache  map[string]*cacheEntry
    mutex  sync.RWMutex
    ttl    time.Duration
}

type cacheEntry struct {
    card      *card.CardCore
    timestamp time.Time
}

func NewCachedClient(client *scryfall.Client, ttl time.Duration) *CachedClient {
    return &CachedClient{
        client: client,
        cache:  make(map[string]*cacheEntry),
        ttl:    ttl,
    }
}

func (cc *CachedClient) GetByName(ctx context.Context, cardName string) (*card.CardCore, error) {
    // Check cache first
    cc.mutex.RLock()
    if entry, exists := cc.cache[cardName]; exists {
        if time.Since(entry.timestamp) < cc.ttl {
            cc.mutex.RUnlock()
            fmt.Printf("Cache hit for %s\n", cardName)
            return entry.card, nil
        }
    }
    cc.mutex.RUnlock()

    // Cache miss or expired - fetch from API
    fmt.Printf("Cache miss for %s, fetching from API\n", cardName)
    card, err := cc.client.Cards.GetByName(ctx, cards.NewExactCardParams(cardName))
    if err != nil {
        return nil, err
    }

    // Store in cache
    cc.mutex.Lock()
    cc.cache[cardName] = &cacheEntry{
        card:      card,
        timestamp: time.Now(),
    }
    cc.mutex.Unlock()

    return card, nil
}

func (cc *CachedClient) ClearExpired() {
    cc.mutex.Lock()
    defer cc.mutex.Unlock()

    now := time.Now()
    for key, entry := range cc.cache {
        if now.Sub(entry.timestamp) >= cc.ttl {
            delete(cc.cache, key)
        }
    }
}

func (cc *CachedClient) Stats() (int, int) {
    cc.mutex.RLock()
    defer cc.mutex.RUnlock()

    total := len(cc.cache)
    expired := 0
    now := time.Now()
    
    for _, entry := range cc.cache {
        if now.Sub(entry.timestamp) >= cc.ttl {
            expired++
        }
    }

    return total, expired
}

// Usage
cachedClient := NewCachedClient(client, 10*time.Minute)

// These will be cached
card1, _ := cachedClient.GetByName(ctx, "Lightning Bolt")
card2, _ := cachedClient.GetByName(ctx, "Lightning Bolt") // Cache hit

// Clean up periodically
go func() {
    ticker := time.NewTicker(5 * time.Minute)
    for range ticker.C {
        cachedClient.ClearExpired()
        total, expired := cachedClient.Stats()
        fmt.Printf("Cache stats: %d total, %d expired\n", total, expired)
    }
}()
```

### File-Based Caching

For persistent caching across application runs:

```go
import (
    "encoding/json"
    "os"
    "path/filepath"
)

type FileCachedClient struct {
    client   *scryfall.Client
    cacheDir string
    ttl      time.Duration
}

func NewFileCachedClient(client *scryfall.Client, cacheDir string, ttl time.Duration) *FileCachedClient {
    os.MkdirAll(cacheDir, 0755)
    return &FileCachedClient{
        client:   client,
        cacheDir: cacheDir,
        ttl:      ttl,
    }
}

func (fc *FileCachedClient) getCacheFilePath(cardName string) string {
    // Use a safe filename
    safeName := strings.ReplaceAll(cardName, " ", "_")
    safeName = strings.ReplaceAll(safeName, "/", "_")
    return filepath.Join(fc.cacheDir, safeName+".json")
}

func (fc *FileCachedClient) GetByName(ctx context.Context, cardName string) (*card.CardCore, error) {
    cacheFile := fc.getCacheFilePath(cardName)
    
    // Check if cached file exists and is recent
    if info, err := os.Stat(cacheFile); err == nil {
        if time.Since(info.ModTime()) < fc.ttl {
            // Load from cache
            data, err := os.ReadFile(cacheFile)
            if err == nil {
                var card card.CardCore
                if json.Unmarshal(data, &card) == nil {
                    fmt.Printf("File cache hit for %s\n", cardName)
                    return &card, nil
                }
            }
        }
    }

    // Cache miss - fetch from API
    fmt.Printf("File cache miss for %s, fetching from API\n", cardName)
    card, err := fc.client.Cards.GetByName(ctx, cards.NewExactCardParams(cardName))
    if err != nil {
        return nil, err
    }

    // Save to cache
    data, err := json.Marshal(card)
    if err == nil {
        os.WriteFile(cacheFile, data, 0644)
    }

    return card, nil
}

func (fc *FileCachedClient) CleanExpired() error {
    return filepath.Walk(fc.cacheDir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        
        if !info.IsDir() && time.Since(info.ModTime()) >= fc.ttl {
            fmt.Printf("Removing expired cache file: %s\n", path)
            return os.Remove(path)
        }
        
        return nil
    })
}

// Usage
fileCachedClient := NewFileCachedClient(client, "./card_cache", 24*time.Hour)

// These will be cached to disk
card, err := fileCachedClient.GetByName(ctx, "Lightning Bolt")
```

## Custom HTTP Configuration

### Custom HTTP Client

Customize the underlying HTTP client:

```go
func createCustomClient() *scryfall.Client {
    // Create custom HTTP client with specific settings
    httpClient := &http.Client{
        Timeout: 60 * time.Second,
        Transport: &http.Transport{
            MaxIdleConns:        100,
            MaxIdleConnsPerHost: 10,
            IdleConnTimeout:     90 * time.Second,
            DisableCompression:  false,
        },
    }

    // Create Scryfall client
    client := scryfall.NewClient()
    client.HTTPClient = httpClient // If the library allows this customization
    
    return client
}
```

### Proxy Support

Configure proxy usage:

```go
func createProxyClient(proxyURL string) (*scryfall.Client, error) {
    proxyURL, err := url.Parse(proxyURL)
    if err != nil {
        return nil, fmt.Errorf("invalid proxy URL: %w", err)
    }

    httpClient := &http.Client{
        Transport: &http.Transport{
            Proxy: http.ProxyURL(proxyURL),
        },
        Timeout: 30 * time.Second,
    }

    client := scryfall.NewClient()
    client.HTTPClient = httpClient // If the library allows this customization
    
    return client, nil
}

// Usage
client, err := createProxyClient("http://proxy.example.com:8080")
if err != nil {
    log.Fatalf("Failed to create proxy client: %v", err)
}
```

## Performance Optimization

### Batch Operations

Process multiple cards efficiently:

```go
func batchCardLookup(client *scryfall.Client, cardNames []string, batchSize int) {
    ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
    defer cancel()

    for i := 0; i < len(cardNames); i += batchSize {
        end := i + batchSize
        if end > len(cardNames) {
            end = len(cardNames)
        }

        batch := cardNames[i:end]
        fmt.Printf("Processing batch %d-%d of %d\n", i+1, end, len(cardNames))

        // Process batch concurrently
        var wg sync.WaitGroup
        for _, cardName := range batch {
            wg.Add(1)
            go func(name string) {
                defer wg.Done()
                
                card, err := client.Cards.GetByName(ctx, cards.NewExactCardParams(name))
                if err != nil {
                    fmt.Printf("❌ %s: %v\n", name, err)
                } else {
                    fmt.Printf("✅ %s: %s\n", name, card.Set)
                }
            }(cardName)
        }

        wg.Wait()

        // Small delay between batches
        if end < len(cardNames) {
            time.Sleep(500 * time.Millisecond)
        }
    }
}

// Usage
cardNames := []string{/* large list of card names */}
batchCardLookup(client, cardNames, 10) // Process 10 cards at a time
```

### Memory Management

For large datasets, manage memory usage:

```go
func processLargeSearchResults(client *scryfall.Client, query string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
    defer cancel()

    page := 1
    processedCount := 0

    for {
        params := &cards.SearchParams{
            Query: query,
            Page:  page,
        }

        results, err := client.Cards.Search(ctx, params)
        if err != nil {
            return fmt.Errorf("error on page %d: %w", page, err)
        }

        // Process cards from this page immediately
        for _, card := range results.Data {
            // Process the card
            processCard(card)
            processedCount++
            
            // Optional: trigger GC periodically for very large datasets
            if processedCount%1000 == 0 {
                runtime.GC()
                fmt.Printf("Processed %d cards, triggered GC\n", processedCount)
            }
        }

        // results.Data will be garbage collected after this iteration
        
        if !results.HasMore {
            break
        }

        page++
        
        // Progress reporting
        if page%10 == 0 {
            fmt.Printf("Completed page %d, processed %d cards\n", page, processedCount)
        }
    }

    fmt.Printf("Finished processing %d cards\n", processedCount)
    return nil
}

func processCard(card card.CardCore) {
    // Process individual card here
    // Keep processing minimal to reduce memory usage
}
```

These advanced features allow you to build sophisticated applications with the Go Scryfall Client while maintaining good performance and reliability.