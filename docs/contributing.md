# Contributing Guidelines

Thank you for your interest in contributing to the Go Scryfall Client! This document provides guidelines and information for contributors.

## Getting Started

### Prerequisites

- Go 1.21 or later
- Git
- Basic understanding of the Scryfall API
- Familiarity with Magic: The Gathering terminology

### Setting Up the Development Environment

1. **Fork the repository**
   ```bash
   # Fork on GitHub, then clone your fork
   git clone https://github.com/YOUR_USERNAME/go-scryfall-client.git
   cd go-scryfall-client
   ```

2. **Set up upstream remote**
   ```bash
   git remote add upstream https://github.com/CharlesNkdl/go-scryfall-client.git
   ```

3. **Install dependencies**
   ```bash
   go mod tidy
   ```

4. **Run tests to ensure everything works**
   ```bash
   go test ./...
   ```

## Development Workflow

### 1. Create a Feature Branch

```bash
git checkout main
git pull upstream main
git checkout -b feature/your-feature-name
```

### 2. Make Your Changes

Follow the coding standards and guidelines outlined below.

### 3. Write Tests

All new code should include comprehensive tests. See [Testing Guidelines](#testing-guidelines) below.

### 4. Run Tests and Linting

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test ./... -cover

# Run specific package tests
go test ./scryfall/services -v

# Check formatting
go fmt ./...

# Run vet for static analysis
go vet ./...
```

### 5. Commit Your Changes

```bash
git add .
git commit -m "feat: add support for new endpoint"
```

See [Commit Message Guidelines](#commit-message-guidelines) for formatting.

### 6. Push and Create Pull Request

```bash
git push origin feature/your-feature-name
```

Then create a pull request on GitHub.

## Coding Standards

### Go Style Guide

Follow the standard Go style guide and conventions:

- Use `gofmt` for formatting
- Follow Go naming conventions (PascalCase for exported, camelCase for unexported)
- Write clear, descriptive variable and function names
- Use Go doc comments for all exported functions and types

### Code Organization

```
scryfall/
├── client.go              # Main client
├── errors/                # Error types
├── models/                # Data models
│   ├── card/             # Card-related models
│   ├── set/              # Set-related models
│   └── request/          # Request parameter models
└── services/             # API service implementations
```

### Naming Conventions

- **Files**: Use snake_case (e.g., `card_service.go`)
- **Types**: Use PascalCase (e.g., `CardService`)
- **Functions**: Use PascalCase for exported, camelCase for unexported
- **Variables**: Use camelCase
- **Constants**: Use PascalCase or UPPER_CASE for package-level constants

### Documentation

All exported functions, types, and packages must have doc comments:

```go
// CardService provides access to card-related Scryfall API endpoints.
type CardService struct {
    Client HTTPClient
}

// GetByName retrieves a card by its name using exact or fuzzy matching.
// It returns the card data or an error if the card is not found or if there's an API error.
func (s *CardService) GetByName(ctx context.Context, params *NamedCardParams) (*card.CardCore, error) {
    // Implementation...
}
```

## API Guidelines

### Adding New Endpoints

When adding support for new Scryfall API endpoints:

1. **Create parameter struct** in appropriate `request` package:
   ```go
   type NewEndpointParams struct {
       RequiredField string  `json:"required_field"`
       OptionalField *string `json:"optional_field,omitempty"`
   }

   func (p *NewEndpointParams) Validate() error {
       if p.RequiredField == "" {
           return fmt.Errorf("required_field cannot be empty")
       }
       return nil
   }

   func (p *NewEndpointParams) ToURLValues() (url.Values, error) {
       if err := p.Validate(); err != nil {
           return nil, err
       }
       // Convert to URL values
   }
   ```

2. **Add service method**:
   ```go
   func (s *ServiceName) NewMethod(ctx context.Context, params *NewEndpointParams) (*ResponseType, error) {
       urlValues, err := params.ToURLValues()
       if err != nil {
           return nil, fmt.Errorf("failed to get URL values: %w", err)
       }
       
       path := fmt.Sprintf("/endpoint?%s", urlValues.Encode())
       req, err := s.Client.NewRequest(ctx, "GET", path)
       if err != nil {
           return nil, err
       }
       
       var result ResponseType
       if err := s.Client.Do(req, &result); err != nil {
           return nil, err
       }
       
       return &result, nil
   }
   ```

3. **Add comprehensive tests** (see Testing Guidelines)

4. **Update documentation** in the `docs/` directory

### Error Handling

- Always return meaningful errors with context
- Use `fmt.Errorf` with `%w` verb to wrap errors
- Validate parameters before making API calls
- Handle all HTTP status codes appropriately

```go
func (s *CardService) GetByName(ctx context.Context, params *NamedCardParams) (*card.CardCore, error) {
    if err := params.Validate(); err != nil {
        return nil, fmt.Errorf("invalid parameters: %w", err)
    }
    
    // API call implementation...
    
    if err != nil {
        return nil, fmt.Errorf("failed to get card by name: %w", err)
    }
    
    return result, nil
}
```

## Testing Guidelines

### Test Structure

Use table-driven tests for comprehensive coverage:

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
        {
            name: "successful exact search",
            params: &NamedCardParams{
                Exact: StringPtr("Lightning Bolt"),
            },
            mockResponse: card.CardCore{
                Name: "Lightning Bolt",
                ID:   "test-id",
            },
            expectedResult: &card.CardCore{
                Name: "Lightning Bolt", 
                ID:   "test-id",
            },
        },
        {
            name: "card not found",
            params: &NamedCardParams{
                Exact: StringPtr("Nonexistent Card"),
            },
            mockError: &errors.ApiError{
                ErrInfo: models.ScryfallError{
                    Status: 404,
                    Detail: "Card not found",
                },
            },
            expectedError: "Card not found",
        },
        // More test cases...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockClient := &MockHTTPClient{
                DoFunc: func(req *http.Request, v interface{}) error {
                    if tt.mockError != nil {
                        return tt.mockError
                    }
                    // Set response data
                    return nil
                },
            }

            service := &CardService{Client: mockClient}
            result, err := service.GetByName(context.Background(), tt.params)

            if tt.expectedError != "" {
                assert.Error(t, err)
                assert.Contains(t, err.Error(), tt.expectedError)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.expectedResult, result)
            }
        })
    }
}
```

### Mock Client

Use the provided mock client for testing:

```go
type MockHTTPClient struct {
    DoFunc        func(req *http.Request, v interface{}) error
    NewRequestFunc func(ctx context.Context, method, path string) (*http.Request, error)
}

func (m *MockHTTPClient) Do(req *http.Request, v interface{}) error {
    if m.DoFunc != nil {
        return m.DoFunc(req, v)
    }
    return nil
}

func (m *MockHTTPClient) NewRequest(ctx context.Context, method, path string) (*http.Request, error) {
    if m.NewRequestFunc != nil {
        return m.NewRequestFunc(ctx, method, path)
    }
    return http.NewRequestWithContext(ctx, method, "https://api.scryfall.com"+path, nil)
}
```

### Test Coverage

- Aim for >90% test coverage
- Test all public methods
- Test error conditions
- Test parameter validation
- Test edge cases

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test ./... -cover

# Run specific tests
go test ./scryfall/services -run TestCardService_GetByName

# Verbose output
go test ./... -v

# Race condition detection
go test ./... -race
```

## Model Guidelines

### Adding New Models

When adding new data models:

1. **Follow JSON tag conventions**:
   ```go
   type NewModel struct {
       ID       string  `json:"id"`
       Name     string  `json:"name"`
       Optional *string `json:"optional,omitempty"`
   }
   ```

2. **Use appropriate Go types**:
   - `*string` for optional string fields
   - `*int` for optional integer fields
   - `[]Type` for arrays
   - Custom enum types for constrained values

3. **Add validation methods when needed**:
   ```go
   func (m *NewModel) Validate() error {
       if m.ID == "" {
           return fmt.Errorf("id cannot be empty")
       }
       return nil
   }
   ```

4. **Document all fields**:
   ```go
   type NewModel struct {
       // ID is the unique identifier for this model
       ID string `json:"id"`
       
       // Name is the display name
       Name string `json:"name"`
   }
   ```

### Enum Types

Use typed constants for enums:

```go
type Rarity string

const (
    RarityCommon   Rarity = "common"
    RarityUncommon Rarity = "uncommon" 
    RarityRare     Rarity = "rare"
    RarityMythic   Rarity = "mythic"
)
```

## Documentation Guidelines

### Code Documentation

- All exported functions and types must have doc comments
- Start doc comments with the name of the item being documented
- Include examples for complex functions
- Document parameters and return values

### README and Docs

When updating documentation:

1. **Keep README.md concise** - focus on quick start
2. **Use docs/ directory** for detailed documentation
3. **Update relevant service docs** when adding endpoints
4. **Include examples** for new features
5. **Update CHANGELOG.md** for significant changes

### Documentation Structure

```
docs/
├── installation.md         # Installation guide
├── getting-started.md      # Basic usage
├── services/              # Service-specific docs
│   ├── cards.md
│   ├── sets.md
│   └── ...
├── examples.md            # Code examples
├── models.md             # Data model reference
├── error-handling.md     # Error handling guide
├── advanced-features.md  # Advanced usage
└── contributing.md       # This file
```

## Commit Message Guidelines

Use conventional commit format:

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

### Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `test`: Adding or updating tests
- `refactor`: Code refactoring
- `style`: Code style changes (formatting, etc.)
- `chore`: Maintenance tasks

### Examples

```
feat(cards): add support for card autocomplete endpoint

fix(client): handle rate limiting correctly

docs(services): update card service documentation

test(sets): add comprehensive set service tests
```

## Pull Request Guidelines

### Before Submitting

- [ ] Tests pass locally
- [ ] Code is properly formatted (`go fmt`)
- [ ] No linting errors (`go vet`)
- [ ] Documentation is updated
- [ ] Commit messages follow guidelines
- [ ] Branch is up to date with main

### PR Description Template

```markdown
## Description
Brief description of the changes.

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Documentation update
- [ ] Refactoring
- [ ] Other (please describe)

## Testing
- [ ] Tests added/updated
- [ ] All tests pass
- [ ] Manual testing performed

## Documentation
- [ ] Documentation updated
- [ ] Examples added/updated

## Checklist
- [ ] Code follows project style guidelines
- [ ] Self-review completed
- [ ] No breaking changes (or documented)
```

## Release Process

### Versioning

We follow [Semantic Versioning](https://semver.org/):

- **MAJOR**: Breaking changes
- **MINOR**: New features (backward compatible)
- **PATCH**: Bug fixes (backward compatible)

### Release Checklist

1. Update CHANGELOG.md
2. Update version in relevant files
3. Create release PR
4. Tag release after merge
5. Update documentation

## Getting Help

### Resources

- [Scryfall API Documentation](https://scryfall.com/docs/api)
- [Go Documentation](https://golang.org/doc/)
- [Project Issues](https://github.com/CharlesNkdl/go-scryfall-client/issues)

### Communication

- **GitHub Issues**: Bug reports and feature requests
- **GitHub Discussions**: Questions and general discussion
- **Pull Requests**: Code contributions

### Questions?

If you have questions about contributing:

1. Check existing issues and documentation
2. Create a GitHub issue with the "question" label
3. Be specific about what you're trying to achieve

## Code of Conduct

### Our Standards

- Be respectful and inclusive
- Focus on constructive feedback
- Help create a welcoming environment
- Respect differing viewpoints

### Unacceptable Behavior

- Harassment or discrimination
- Trolling or insulting comments
- Personal attacks
- Publishing private information

### Reporting

Report unacceptable behavior by creating an issue or contacting maintainers directly.

## Recognition

Contributors will be recognized in:

- CONTRIBUTORS.md file
- Release notes for significant contributions
- GitHub contributors list

Thank you for contributing to the Go Scryfall Client! Your efforts help make this library better for the entire Magic: The Gathering development community.