# Perplexity CLI

## Overview

The Perplexity CLI is a Go-based command-line interface tool that provides seamless interaction with the Perplexity Sonar API. Built with Cobra and Viper, it offers both one-shot query execution and interactive conversation modes with session persistence, citation handling, and session management capabilities.

## Architecture

### Core Components

**API Client (`pkg/perplexity/`)**
- HTTP client with API key authentication
- Chat completions endpoint integration
- Request/response type definitions
- Citation extraction and formatting

**Session Management (`pkg/session/`)**
- JSON-based conversation storage
- Timestamp-based file naming
- Auto-save functionality
- CRUD operations for sessions

**Configuration (`pkg/config/`)**
- Environment variable support (PPLX_API_KEY)
- Configuration file management (~/.pplx/config.yaml)
- Default model and settings

**CLI Commands (`cmd/`)**
- Root command and interactive mode
- One-shot query execution (`pplx run`)
- Session management (`pplx session`)
- Session search and listing

### Project Structure

```
perplexity-cli/
├── cmd/
│   ├── root.go           # Root command & interactive mode
│   ├── run.go            # pplx run command
│   └── session.go        # pplx session command group
│       ├── list.go       # pplx session list
│       └── search.go     # pplx session search
├── pkg/
│   ├── perplexity/
│   │   ├── client.go     # API client
│   │   ├── types.go      # Request/response structs
│   │   └── citations.go  # Citation parsing
│   ├── session/
│   │   ├── manager.go    # Session CRUD
│   │   ├── types.go      # Session JSON structures
│   │   └── utils.go      # Helper functions
│   ├── config/
│   │   └── config.go     # Configuration management
│   └── ui/
│       ├── colors.go     # Terminal color formatting
│       └── spinner.go    # Loading indicators
├── main.go               # Entry point
├── Taskfile.yaml         # Task-based build system
├── Makefile              # Traditional build system
└── README.md
```

## Features

### One-Shot Queries

Execute a single query and receive a response with citations.

**Usage:**
```bash
pplx run "What is the capital of France?" [--model sonar]
```

**Output Format:**
```
PPLX: Paris is the capital of France. [1]

## References:
[1] France - https://en.wikipedia.org/wiki/France
```

**Key Details:**
- Query provided as positional argument
- Accepts optional `--model` flag (default: `sonar`)
- Displays content with inline citations in bracket format `[1]`, `[2]`
- Appends formatted references section at the end
- Single API call with no session persistence

### Interactive Mode

Enter a REPL-style conversation with context retention.

**Usage:**
```bash
pplx
```

**Features:**
- Detects no arguments and enters interactive mode
- Prompt: `You: `
- Exit commands: `/q`, `/quit`, `ctrl+c`
- Maintains conversation history in memory
- Auto-saves session after each exchange
- Displays formatted responses with citations

**Critical Behavior:**
- Does not send reference content to API
- Strips references section from conversation history before API calls
- Maintains separate arrays: display context (with citations) and API context (clean)

**Output Format:**
```
You: What is the capital of France?

PPLX: Paris is the capital of France. [1]

## References:
[1] France - https://en.wikipedia.org/wiki/France

You: What about Germany?

PPLX: The capital of Germany is Berlin. [2]

## References:
[2] Germany - https://en.wikipedia.org/wiki/Germany
```

### Session Management

Sessions are stored as JSON files in `~/.pplx/sessions/` with timestamp-based naming.

**Storage Path:** `~/.pplx/sessions/`

**Filename Format:** `20060102-150405.000.json` (RFC3339 timestamp with milliseconds)

**Session Structure:**
```json
{
  "id": "20060102-150405.000",
  "created_at": "2006-01-02T15:04:05.000Z",
  "messages": [
    {
      "role": "user",
      "content": "query text"
    },
    {
      "role": "assistant",
      "content": "response text with [1] citations"
    }
  ],
  "model": "sonar"
}
```

**Auto-Save:**
- Session files created immediately on interactive mode start
- Updated after each API exchange in interactive mode
- Written to temp file then atomic rename to prevent corruption

### Session List

Display recent sessions sorted by timestamp.

**Usage:**
```bash
pplx session list -n 10
# or
pplx session -l 10
```

**Output:**
```
20060102-150405.000 - What is the capital of France?
20060102-140030.123 - Explain quantum computing
20060102-130015.456 - How does photosynthesis work?
```

**Features:**
- Sort by filename (timestamp) descending
- Default limit: 10 sessions
- Displays timestamp and truncated initial query
- Uses `-n` flag for custom limit

### Session Search

Search through session files for specific terms.

**Usage:**
```bash
pplx session search "photosynthesis"
```

**Output:**
```
20060102-130015.456 - How does photosynthesis work?
  User: How does photosynthesis work?
  Assistant: Photosynthesis is the process by which plants convert light energy...
```

**Features:**
- Case-insensitive search through all session files
- Searches both user and assistant message content
- Returns matching filename with context
- Uses grep for efficient file system search

## Configuration

### Environment Variables

**Required:**
- `PPLX_API_KEY`: Perplexity API key for authentication

**Optional:**
- `PPLX_MODEL`: Default model (default: `sonar`)
- `PPLX_MAX_TOKENS`: Maximum tokens for responses
- `PPLX_TEMPERATURE`: Response temperature (0.0-2.0)

### Configuration File

Location: `~/.pplx/config.yaml`

**Options:**
```yaml
default_model: sonar
max_tokens: 1024
temperature: 0.7
```

**Priority:**
1. Command-line flags (highest)
2. Environment variables
3. Configuration file
4. Default values (lowest)

## API Integration

### Endpoint

**URL:** `https://api.perplexity.ai/chat/completions`

**Method:** POST

**Headers:**
```
Authorization: Bearer <API_KEY>
Content-Type: application/json
```

**Request Body:**
```json
{
  "model": "sonar",
  "messages": [
    {"role": "user", "content": "query text"}
  ],
  "stream": false
}
```

**Response:**
```json
{
  "id": "chatcmpl-xxx",
  "model": "sonar",
  "choices": [
    {
      "message": {
        "role": "assistant",
        "content": "Response text [1] [2]"
      }
    }
  ],
  "citations": [
    "https://example.com/reference1",
    "https://example.com/reference2"
  ]
}
```

### Citation Handling

**Extraction:**
- Parse content for pattern `\[\d+\]` using regex
- Extract reference URLs from `citations` array
- Map citation numbers to URLs by index

**Display Format:**
```
Content text with [1] citations inline.

## References:
[1] Title - https://example.com/reference1
[2] Title - https://example.com/reference2
```

**Reference Cleanup:**
- Strip `## References:` section and all lines after before sending to API
- Maintain display context with citations for output
- Use clean context for API conversation history

## Error Handling

### Graceful Degradation

**Missing API Key:**
```
Error: PPLX_API_KEY environment variable not set
Please set your Perplexity API key: export PPLX_API_KEY=your_key_here
```

**API Errors:**
- Parse error messages from API responses
- Display formatted error with context
- Return appropriate exit codes

**Network Issues:**
- Timeout handling with configurable duration
- Retry logic for transient failures
- Connection error messaging

### Signal Handling

**Graceful Shutdown:**
- Handle `SIGINT` (Ctrl+C) for clean exit
- Save current session state before termination
- Clean up temporary files
- Display exit message

## Implementation Details

### Citation Parser

**Pattern Matching:**
```go
regex := regexp.MustCompile(`\[(\d+)\]`)
matches := regex.FindAllStringSubmatch(content, -1)
```

**Reference Extraction:**
- Extract citation numbers from content
- Map to URLs from API response
- Format as numbered list with titles

### Context Management

**Display Context:**
- Full message history with citations
- Used for output formatting
- Maintains reference numbers

**API Context:**
- Clean message history without references
- Stripped of `## References:` section
- Sent to API for subsequent requests

### Session Persistence

**Write Strategy:**
1. Create temp file in same directory
2. Write session data to temp file
3. Atomic rename to target filename
4. Prevents corruption on write failures

**Atomic Rename:**
```go
tempFile := filepath.Join(dir, filename+".tmp")
os.Rename(tempFile, targetFile)
```

## Testing

### Unit Tests

**API Client:**
- Request serialization
- Response parsing
- Citation extraction
- Error handling

**Session Management:**
- File creation and deletion
- JSON serialization/deserialization
- Timestamp generation
- Search functionality

**Configuration:**
- Environment variable parsing
- Config file loading
- Default value application

### Integration Tests

**One-Shot Queries:**
- Query execution
- Output formatting
- Citation display

**Interactive Mode:**
- Conversation flow
- Context retention
- Auto-save functionality

**Session Commands:**
- List sorting and limiting
- Search accuracy
- Context display

## Dependencies

**Core Framework:**
- `github.com/spf13/cobra` - CLI command framework
- `github.com/spf13/viper` - Configuration management

**Utilities:**
- `github.com/fatih/color` - Terminal color formatting

**Go Version:** 1.25.2

## Build and Installation

**Using Taskfile:**
```bash
task build
task install
```

**Using Makefile:**
```bash
make build
make install
```

**Binary Name:** `pplx`

**Install Locations:**
- GOPATH/bin (if GOPATH set)
- /usr/local/bin (fallback)

## Development

**Run Tests:**
```bash
task test
# or
make test
```

**Format Code:**
```bash
task fmt
# or
make fmt
```

**Run Linter:**
```bash
task lint
# or
make lint
```

**Build and Run:**
```bash
task run
# or
make run
```

## Future Enhancements

Potential improvements:
- Stream response support for real-time output
- Session export to different formats (Markdown, JSON)
- Conversation export/import functionality
- Multiple model selection from API
- Custom citation formatting options
- Session tagging and organization
- Conversation context window management
- Configurable output verbosity
