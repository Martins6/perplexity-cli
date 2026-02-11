# perplexity-cli

A command-line interface for interacting with the Perplexity Sonar API.

## Features

- One-shot queries - Send a single query and get a response
- Interactive mode - Have a conversation with Perplexity
- Session management - Save and search conversation history
- Citation handling - Automatic formatting of citations and references

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
```

Set your Perplexity API key:
```bash
export PPLX_API_KEY=your-api-key-here
```

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
