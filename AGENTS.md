# AGENTS.md - Coding Guidelines for perplexity-cli

## Project Overview

Go CLI application for interacting with the Perplexity Sonar API. Uses Cobra for CLI commands and Viper for configuration management.

## Build Commands

```bash
# Build binary
make build
# or
task build

# Run all tests
make test
# or
task test
# or
go test ./... -v

# Run single test
go test -run TestName ./pkg/perplexity
go test -run TestFunctionName ./pkg/session

# Run tests with coverage
make test-coverage
task test-coverage

# Format code
make fmt
task fmt

# Run linter (requires golangci-lint)
make lint
task lint

# Clean build artifacts
make clean
task clean

# Install binary locally
make install
task install

# Build and run
make run
task run
```

## Code Style Guidelines

### General

- **Go Version**: 1.25.2
- **Module**: `perplexity-cli`
- Run `go fmt ./...` before committing
- Run tests and linter before pushing

### Imports

Order imports as:
1. Standard library (no blank line after)
2. Third-party packages
3. Internal packages (e.g., `perplexity-cli/pkg/config`)

Example:
```go
import (
    "fmt"
    "os"
    
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    
    "perplexity-cli/pkg/config"
    "perplexity-cli/pkg/perplexity"
)
```

### Naming Conventions

- **Exported**: PascalCase (e.g., `NewClient`, `ChatCompletionRequest`)
- **Unexported**: camelCase (e.g., `config`, `httpClient`)
- **Acronyms**: All caps (e.g., `APIKey`, `HTTPClient`)
- **Test files**: Suffix with `_test.go`
- **Test functions**: Start with `Test` (e.g., `TestExtractCitations`)

### Types and Structs

- Add comments for all exported types, functions, methods, and constants
- Use struct tags for JSON serialization with `omitempty` when appropriate
- Group related fields in structs logically

Example:
```go
// Message represents a chat message in the conversation
type Message struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}
```

### Error Handling

- Wrap errors with context using `%w` verb
- Check errors immediately after function calls
- Use early returns to reduce nesting
- Print errors to stderr with red color formatting for CLI output

Example:
```go
if err != nil {
    return nil, fmt.Errorf("failed to create request: %w", err)
}
```

### Testing

- Use table-driven tests with descriptive test case names
- Use `t.Run()` for subtests
- Create temporary directories for file system tests and clean up with `defer`
- Use `t.Fatalf()` for fatal errors, `t.Errorf()` for non-fatal assertions
- Sleep small durations when testing time-dependent behavior

Example:
```go
func TestFunction(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {name: "valid case", input: "test", expected: "result"},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Function(tt.input)
            if result != tt.expected {
                t.Errorf("Function() = %v, expected %v", result, tt.expected)
            }
        })
    }
}
```

### Comments

- All exported identifiers must have a comment
- Comments should start with the identifier name
- Use complete sentences with proper punctuation
- Example: `// Client represents a Perplexity API client`

### Constants

- Use PascalCase for exported constants
- Group related constants together
- Prefer const over var for immutable values

Example:
```go
const DefaultAPIEndpoint = "https://api.perplexity.ai/chat/completions"
```

### Functions

- Keep functions focused and small
- Use early returns to minimize nesting
- Return errors as the last return value
- Use named return values sparingly

### Project Structure

```
cmd/           # Cobra command implementations
pkg/           # Library packages
  config/      # Configuration management
  perplexity/  # Perplexity API client
  session/     # Session management
  ui/          # UI utilities (colors, spinner)
main.go        # Entry point
```

## Dependencies

Key dependencies:
- `github.com/spf13/cobra` - CLI framework
- `github.com/spf13/viper` - Configuration management
- `github.com/fatih/color` - Terminal colors

Always run `go mod tidy` after adding/removing imports.
