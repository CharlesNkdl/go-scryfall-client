# Go Scryfall Client Test Suite

This document describes the comprehensive test suite for the Go Scryfall Client library.

## Running Tests

### Run All Tests
```bash
go test ./...
```

### Run Tests with Verbose Output
```bash
go test ./... -v
```

### Run Tests for Specific Packages
```bash
# Service tests
go test ./scryfall/services -v

# Client tests
go test ./scryfall -v

# Request parameter tests
go test ./scryfall/models/request/cards -v
```

### Run Tests with Coverage
```bash
go test ./... -cover
```

## Test Structure

### Service Tests (`scryfall/services/`)
- **CardService**: Tests for all card endpoints (GetById, GetByName, Search, Autocomplete, GetRandom, GetByCodeNumberLang)
- **SetService**: Tests for set retrieval functionality
- **RulingService**: Tests for card ruling retrieval
- **CatalogService**: Tests for various catalog endpoints
- **SymbolService**: Tests for mana symbol lookup with URL encoding

### Client Tests (`scryfall/`)
- **Client Configuration**: Tests client initialization and configuration
- **Request Creation**: Tests HTTP request creation with proper headers
- **Rate Limiting**: Tests that rate limiting is properly implemented
- **Error Handling**: Tests proper handling of HTTP errors and API errors
- **Context Support**: Tests context cancellation and timeouts

### Request Parameter Tests (`scryfall/models/request/cards/`)
- **Validation**: Tests parameter validation for all request types
- **URL Encoding**: Tests proper conversion to URL parameters
- **Edge Cases**: Tests boundary conditions and error scenarios

## Test Features

### Mock HTTP Client
The test suite uses a comprehensive mock HTTP client that allows testing without making actual API calls:
- Configurable responses and errors
- Request validation
- Path and method verification

### Table-Driven Tests
All tests use Go's table-driven test pattern for comprehensive coverage:
- Multiple test cases per function
- Clear separation of test data and logic
- Easy to add new test cases

### Validation Testing
Extensive testing of parameter validation following Scryfall API specifications:
- Required parameter validation
- Format validation (json, text, image)
- Range validation (page numbers, string lengths)
- Mutual exclusivity validation (exact vs fuzzy search)

### Error Scenario Testing
Tests cover various error scenarios:
- Network errors
- HTTP status code errors (404, 429, 500)
- JSON parsing errors
- Context cancellation
- Rate limiting

### API Compliance
Tests ensure compliance with Scryfall API requirements:
- Proper User-Agent headers
- Rate limiting respect
- Correct URL encoding
- Parameter validation according to API specs

## Best Practices Demonstrated

1. **Separation of Concerns**: Mock client separates testing logic from actual HTTP calls
2. **Comprehensive Coverage**: Tests cover happy paths, error paths, and edge cases
3. **Go Conventions**: Follows Go testing conventions and naming patterns
4. **API Validation**: Tests ensure requests match API specifications
5. **Performance Awareness**: Rate limiting tests ensure API guidelines are followed

## Adding New Tests

When adding new endpoints or modifying existing ones:

1. Add service tests in `scryfall/services/`
2. Add request parameter tests in appropriate request package
3. Update mock client if new interfaces are needed
4. Ensure validation tests cover all parameter combinations
5. Test both success and error scenarios