# perplexity-cli

A command-line interface for interacting with the Perplexity Sonar API.

## Features

- One-shot queries - Send a single query and get a response
- Interactive mode - Have a conversation with Perplexity with beautiful markdown rendering
- Session management - Save and search conversation history
- Citation handling - Automatic formatting of citations and references
- Markdown rendering - Rich formatting with syntax highlighting in interactive mode

## Installation

### Prerequisites

- Go 1.25.2 or later
- [Task](https://taskfile.dev) - Modern build tool (optional but recommended)
  ```bash
  go install github.com/go-task/task/v3/cmd/task@latest
  ```

### Quick Install

If you have Task installed, simply run:

```bash
task install
```

This will build the binary and install it to `$GOPATH/bin` (if set) or `/usr/local/bin`.

### Manual Install

```bash
task install
# or
task build && cp build/pplx $GOPATH/bin/
```

### Build from Source

```bash
# Using Task (recommended)
task build

# Using Make
task build

# Manual
go build -o build/pplx main.go
```

## Available Tasks

Run `task --list` to see all available tasks:

```bash
$ task --list
task: Available tasks for this project:
* all:                 Run clean, build, and test
* build:               Build binary for current platform
* clean:               Clean build artifacts
* default:             Build the pplx binary
* deps:                Download and tidy dependencies
* fmt:                 Format Go code
* install:             Install binary to $GOPATH/bin or /usr/local/bin (always cleans first)
* lint:                Run linter (requires golangci-lint)
* run:                 Build and run the application
* test:                Run all tests
* test-coverage:       Run tests with coverage report
* verify:              Build and run tests (CI verification)
```

## Usage

```bash
# One-shot query
pplx run "What is the capital of France?"

# Interactive mode
pplx

# List recent sessions
pplx session -l 10

# Search sessions
pplx session search "France"
```

## Configuration

Configuration can be provided via:
- Environment variables (prefix: `PPLX_`)
- Config file: `~/.pplx/config.yaml`

Example config:
```yaml
api_key: your-api-key-here
model: sonar
max_tokens: 1024
temperature: 0.7
top_p: 0.9
search_context_size: medium
search_mode: true
reasoning_effort: low

# Markdown rendering (interactive mode only)
use_glow: true         # Enable/disable markdown rendering
glow_style: auto       # Style: auto, dark, light, or custom JSON path
glow_width: 0          # Word wrap width (0 = terminal width)
```

Set your Perplexity API key:
```bash
export PPLX_API_KEY=your-api-key-here
```

## Markdown Rendering

The interactive mode supports beautiful markdown rendering with syntax highlighting using the Charmbracelet Glamour library.

### Features
- **Syntax highlighting** - Code blocks are automatically highlighted
- **Rich formatting** - Headers, lists, quotes, and tables are properly rendered
- **Multiple themes** - Choose from dark, light, or auto (detects terminal background)
- **Customizable** - Adjust word wrap width or use custom JSON stylesheets
- **Graceful fallback** - Automatically uses plain text if not a terminal

### Configuration Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `use_glow` | bool | true | Enable/disable markdown rendering in interactive mode |
| `glow_style` | string | auto | Style theme: auto, dark, light, or custom JSON path |
| `glow_width` | int | 0 | Word wrap width (0 = use terminal width) |

### Examples

Enable markdown rendering with dark theme:
```yaml
use_glow: true
glow_style: dark
glow_width: 0
```

Disable markdown rendering (plain text only):
```yaml
use_glow: false
```

Use a custom word wrap width:
```yaml
use_glow: true
glow_style: auto
glow_width: 80
```

**Note:** Markdown rendering is only available in interactive mode. One-shot queries (`pplx run`) always display plain text output.

## Development

```bash
# Run tests
task test

# Format code
task fmt

# Run linter
task lint

# Full verification (build + test)
task verify
```

## Why Taskfile?

Task provides several advantages over traditional Makefiles:

- **Better incremental builds** - Tasks are skipped when sources haven't changed
- **Task dependencies** - Tasks can depend on other tasks (e.g., `install` runs `clean` and `build` first)
- **Cross-platform** - Works consistently across macOS, Linux, and Windows
- **Built-in help** - `task --list` shows all available tasks with descriptions
- **Modern YAML syntax** - Easier to read and maintain than Makefile syntax

Both Taskfile and Makefile are supported and maintained. Use whichever you prefer!

## License

[License](LICENSE)
