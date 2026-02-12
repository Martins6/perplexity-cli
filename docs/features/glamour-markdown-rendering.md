# Glamour Markdown Rendering

## Overview

The Glamour Markdown Rendering feature provides beautiful, styled markdown output in the interactive mode of the Perplexity CLI tool. It uses the [Charmbracelet Glamour](https://github.com/charmbracelet/glamour) library to render markdown content with syntax highlighting, proper formatting, and theme support.

## Motivation

Prior to this feature, the interactive `pplx` command and one-shot `pplx run` command displayed markdown responses from the Perplexity Sonar API as plain text with minimal formatting. This made it difficult to read and appreciate the rich markdown content, including headers, code blocks, lists, and other markdown elements that Perplexity often returns.

## Implementation Details

### Core Components

1. **Markdown Rendering Utility** (`pkg/ui/markdown.go`)
   - `RenderMarkdown(content string, cfg *config.Config) (string, error)` - Main rendering function for interactive mode
   - `RenderMarkdownAlways(content string, cfg *config.Config) string` - Always-on rendering for one-shot queries
   - Creates Glamour TermRenderer with configurable options
   - Graceful fallback to plain text when disabled or not a terminal

2. **Configuration Options** (`pkg/config/config.go`)
   - `UseGlow bool` - Enable/disable markdown rendering (default: true)
   - `GlowStyle string` - Select style theme (default: "auto")
   - `GlowWidth int` - Set word wrap width (default: 0 for terminal width)

3. **Integration Points**
   - Interactive mode (`cmd/interactive.go`) - Renders assistant responses
   - One-shot queries (`cmd/run.go`) - Renders response content and citations with markdown
   - Session display (`pkg/session/display.go`) - Renders saved session messages

### Key Features

#### Syntax Highlighting
Code blocks are automatically highlighted with appropriate colors based on programming language:
```markdown
\`\`\`go
func main() {
    fmt.Println("Hello, World!")
}
\`\`\`
```

#### Rich Formatting
All markdown elements are properly rendered:
- Headers (#, ##, ###)
- Bold (**text**) and italic (*text*)
- Lists (unordered and ordered)
- Blockquotes (> quote)
- Tables (if supported)
- Links and images

#### Multiple Themes
- **auto** - Automatically detects terminal background color (default)
- **dark** - Dark theme
- **light** - Light theme
- **Custom** - Path to a custom JSON stylesheet

#### Terminal Detection
The renderer automatically detects whether output is a terminal:
- If output is a terminal: Uses markdown rendering
- If output is piped/redirected: Falls back to plain text
- Configurable via `use_glow` option

### Graceful Degradation

The implementation includes several fallback mechanisms:

1. **Config disabled** - If `use_glow: false`, returns plain text
2. **Not a terminal** - If output is piped/redirected, returns plain text
3. **Renderer error** - If Glamour fails to initialize, returns plain text
4. **Nil config** - If config is nil, returns plain text

### Performance Considerations

- **Small overhead** - Glamour rendering adds minimal overhead (~1-5ms per response)
- **Memory** - Temporary memory usage during rendering (~1-2MB)
- **Impact:** Negligible for typical responses; acceptable in interactive mode

## Configuration

### Environment Variables

```bash
export PPLX_USE_GLOW=true
export PPLX_GLOW_STYLE=dark
export PPLX_GLOW_WIDTH=80
```

### Config File (~/.pplx/config.yaml)

```yaml
# Markdown rendering (interactive mode only)
use_glow: true         # Enable/disable markdown rendering
glow_style: auto       # Style: auto, dark, light, or custom JSON path
glow_width: 0          # Word wrap width (0 = terminal width)
```

## Usage Examples

### One-Shot Queries

Markdown rendering is automatically enabled for all `pplx run` commands:

```bash
$ pplx run "Explain quantum computing"
**Quantum computing** is a new approach to calculation that uses principles of quantum mechanics...

### Key Concepts
- **Qubits**: The basic unit of quantum information...
- **Superposition**: Qubits hold a combination of all possible states...

## References:
[1] What is quantum computing? - https://example.com/quantum
```

The `pplx run` command always renders markdown when output is to a terminal, regardless of the `use_glow` configuration option. This ensures consistent, beautiful output for one-shot queries.

### Interactive Mode

Markdown rendering is automatically enabled in interactive mode and respects the `use_glow` configuration:

```bash
$ pplx
Welcome to PPLX Interactive Mode!
Model: sonar | Type '/q' or press Ctrl+C to exit

You: Show me a Go function example
PPLX: Here's an example:

# Go Function Example

\`\`\`go
package main

import "fmt"

func greet(name string) {
    fmt.Printf("Hello, %s!\n", name)
}

func main() {
    greet("World")
}
\`\`\`

This function takes a name parameter and prints a greeting.
```

### Custom Themes

Use dark theme:

```yaml
use_glow: true
glow_style: dark
```

Use light theme:

```yaml
use_glow: true
glow_style: light
```

### Disable Markdown Rendering

If you prefer plain text output:

```yaml
use_glow: false
```

Or via environment variable:

```bash
PPLX_USE_GLOW=false pplx
```

## Testing

The feature includes comprehensive unit tests in `pkg/ui/markdown_test.go`:

- Test rendering with enabled glamour
- Test fallback when glamour is disabled
- Test fallback when not a terminal
- Test with different markdown content (headers, code blocks, lists, etc.)
- Test with different styles (auto, dark, light)
- Test error handling (invalid style, etc.)

Run tests with:
```bash
go test ./pkg/ui/...
```

## Dependencies

- `github.com/charmbracelet/glamour` - Markdown rendering library
- `github.com/fatih/color` - Terminal color support (existing dependency)

## Scope

- **Both modes** - Markdown rendering is enabled for both interactive mode (`pplx`) and one-shot queries (`pplx run`)
- **Always-on for run command** - The `pplx run` command always uses glamour rendering when in a terminal, regardless of `use_glow` configuration
- **Configurable for interactive mode** - Interactive mode respects the `use_glow` configuration option
- **Optional dependency** - Gracefully degrades if glamour is not available or not a terminal
- **Zero breaking changes** - Existing functionality is preserved

## Related Documentation

- [Glamour Rendering Guide](../glamour-rendering.md) - Detailed usage and configuration guide
- [Charmbracelet Glamour](https://github.com/charmbracelet/glamour) - Official library documentation
