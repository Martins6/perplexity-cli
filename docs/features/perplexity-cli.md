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
│       ├── search.go     # pplx session search
│       └── show.go       # pplx session show
├── pkg/
│   ├── perplexity/
│   │   ├── client.go     # API client
│   │   ├── types.go      # Request/response structs
│   │   └── citations.go  # Citation parsing
│   ├── session/
│   │   ├── manager.go    # Session CRUD + LoadByShortID
│   │   ├── types.go      # Session JSON structures + ShortID field
│   │   ├── utils.go      # Helper functions + GenerateShortID
│   │   └── display.go    # Conversation display utilities
│   ├── config/
│   │   └── config.go     # Configuration management
│   └── ui/
│       └── colors.go     # Terminal color formatting
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
- Sessions are NOT saved for one-shot queries (unlike interactive mode)
- No session files are created in `~/.pplx/sessions/`

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
Welcome to PPLX Interactive Mode!
Model: sonar | Type '/q' or press Ctrl+C to exit

You: What is the capital of France?

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
PPLX: Paris is the capital of France. [1]

## References:
[1] France - https://en.wikipedia.org/wiki/France

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
You: What about Germany?

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
PPLX: The capital of Germany is Berlin. [2]

## References:
[2] Germany - https://en.wikipedia.org/wiki/Germany
```

**Visual Enhancements:**
- Colored separator lines appear between consecutive messages
- Different colors for YOU (blue) vs PPLX (cyan) separators
- No separator before the first "You:" prompt
- Improves readability and visual distinction in conversation flow

### Session Management

Sessions are stored as JSON files in `~/.pplx/sessions/` with timestamp-based naming.

**Storage Path:** `~/.pplx/sessions/`

**Filename Format:** `20060102-150405.000.json` (RFC3339 timestamp with milliseconds)

**Session Structure:**
```json
{
  "id": "20060102-150405.000",
  "short_id": "a8x9k2",
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

**Short ID:**
- `short_id` field contains Base62-encoded timestamp
- Generated from `created_at` timestamp for consistency
- Present in both Session and SessionInfo structures
- Used for quick reference in list/search/show commands

**Auto-Save:**
- Session files created immediately on interactive mode start
- Updated after each API exchange in interactive mode
- Written to temp file then atomic rename to prevent corruption
- New sessions include `short_id` field automatically

**Backward Compatibility:**
- Existing sessions without `short_id` are handled gracefully
- Short ID is generated from `created_at` timestamp on first load
- Session is automatically saved with new short ID (lazy migration)
- No batch migration needed; transparent to user

### Session List

Display recent sessions sorted by timestamp with short IDs.

**Usage:**
```bash
pplx session list -n 10
# or
pplx session -l 10
```

**Output:**
```
1. [a8x9k2] Jan 02, 2006 15:04:05
   What is the capital of France?
   (3 messages)

2. [b7j8l1] Jan 02, 2006 14:00:30
   Explain quantum computing
   (5 messages)

3. [c6i7k0] Jan 02, 2006 13:00:15
   How does photosynthesis work?
   (4 messages)
```

**Features:**
- Sort by filename (timestamp) descending
- Default limit: 10 sessions
- Displays short Base62 ID in brackets for easy reference
- Shows timestamp, truncated initial query, and message count
- Uses `-n` flag for custom limit

**Short ID Format:**
- Base62-encoded Unix timestamp (milliseconds)
- 6-7 alphanumeric characters (0-9a-zA-Z)
- Collision-free due to monotonic timestamps
- Easy to copy-paste without special characters

### Session Search

Search through session files for specific terms or short IDs.

**Usage:**
```bash
pplx session search "photosynthesis"
# or search by short ID
pplx session search a8x9k2
```

**Output:**
```
[a8x9k2] Jan 02, 2006 15:04:05
What is the capital of France?

User: How does photosynthesis work?
Assistant: Photosynthesis is the process by which plants convert light energy...
```

**Features:**
- Case-insensitive search through all session files
- Searches both user and assistant message content
- Matches against short IDs as well as message content
- Returns matching short ID with context and timestamp
- Uses grep for efficient file system search

**ID Lookup:**
- Short IDs are alphanumeric (6-7 chars): `a8x9k2`
- Timestamp IDs have format: `20060102-150405.000`
- Command tries short ID match first, then falls back to full timestamp ID

### Session Show

Display the full conversation of a specific session with formatted messages and citations.

**Usage:**
```bash
pplx session show a8x9k2
# or use full timestamp ID
pplx session show 20060102-150405.000
```

**Output:**
```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
You: What is the capital of France?
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
PPLX: Paris is the capital of France. [1]

## References:
[1] France - https://en.wikipedia.org/wiki/France
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
You: What about Germany?
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
PPLX: The capital of Germany is Berlin. [2]

## References:
[2] Germany - https://en.wikipedia.org/wiki/Germany
```

**Features:**
- Accepts short Base62 ID or full timestamp ID
- Displays entire conversation chronologically
- Color-coded separators (blue for user, cyan for assistant)
- Shows inline citations with reconstructed reference sections
- Displays timestamps and message metadata
- Handles not found errors gracefully

**Display Format:**
- Each message wrapped in separator lines for clarity
- User messages labeled with "You:" (blue separator)
- Assistant messages labeled with "PPLX:" (cyan separator)
- Citations displayed inline in bracket format
- Full reference section shown after each assistant response
- Original message content preserved exactly as stored

**ID Resolution:**
- Tries short ID lookup first (LoadByShortID)
- Falls back to full timestamp ID lookup if needed
- Provides clear error if session not found

### Session Command Shortcuts

For faster access to common session commands, use these root-level flags:

**Available Shortcuts:**
- `-c [id]` - Continue a session (same as: `pplx session continue [id]`)
- `-l [limit]` - List recent sessions (same as: `pplx session list -l [limit]`)
- `-s [query]` - Search sessions (same as: `pplx session search [query]`)

**Examples:**

**Continue a session:**
```bash
# Full command
pplx session continue a8x9k2

# Shortcut
pplx -c a8x9k2
```

**List recent sessions:**
```bash
# Full command
pplx session list -l 20

# Shortcut
pplx -l 20
```

**Search sessions:**
```bash
# Full command
pplx session search "France"

# Shortcut
pplx -s "France"
```

**Conflict Behavior:**
Shortcuts cannot be combined with explicit session commands. If you try to use both together, you'll receive an error:

```bash
$ pplx session -c a8x9k2
Error: cannot use session shortcuts (-c, -l, -s) together with explicit session commands

Please use either:
  - pplx -c [id] (shortcut)
  - pplx session continue [id] (explicit command)
```

**Limitations:**
- Shortcuts only work at the root level (e.g., `pplx -c` works, `pplx session -c` does not)
- Only one shortcut can be used at a time
- All other flags work normally with shortcuts (e.g., `pplx -c a8x9k2 --model sonar-pro`)

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

**Note:** The configuration file is optional. The CLI will work correctly using only the `PPLX_API_KEY` environment variable. If the config file does not exist, it will not cause an error.

**Options:**
```yaml
default_model: sonar
max_tokens: 1024
temperature: 0.7
```

**Priority:**
1. Command-line flags (highest)
2. Environment variables
3. Configuration file (optional)
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

### Short ID Generation

**Base62 Encoding:**
```go
func GenerateShortID(t time.Time) string {
    timestamp := t.UnixMilli()
    chars := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
    result := make([]byte, 7)
    for i := 6; i >= 0; i-- {
        result[i] = chars[timestamp%62]
        timestamp /= 62
    }
    return string(result)
}
```

**Properties:**
- Encodes Unix timestamp in milliseconds to Base62
- Uses 62-character alphabet: `0-9a-zA-Z`
- Produces 6-7 character strings (7 for full millisecond precision)
- Monotonic increasing (newer timestamps = lexicographically larger)
- Guaranteed unique (timestamps never repeat)
- No collision checking required

**Example IDs:**
- `a8x9k2` - January 2, 2006 15:04:05
- `b7j8l1` - January 2, 2006 14:00:30
- `c6i7k0` - January 2, 2006 13:00:15

**Migration:**
```go
func (m *Manager) Load(id string) (*Session, error) {
    session, err := m.loadSession(id)
    if err != nil {
        return nil, err
    }
    
    // Lazy migration: generate short ID if missing
    if session.ShortID == "" {
        session.ShortID = GenerateShortID(session.CreatedAt)
        m.Save(session)
    }
    
    return session, nil
}
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
- Short ID generation and uniqueness
- LoadByShortID functionality
- Search functionality
- Show command formatting

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
- List sorting and limiting with short ID display
- Search accuracy with short ID matching
- Show command formatting with citations
- Backward compatibility with old sessions

**Short ID Generation:**
- Base62 encoding correctness
- Timestamp uniqueness guarantees
- Monotonic ordering verification
- ID collision scenarios (should not occur)

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
task build
task install
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
task test
```

**Format Code:**
```bash
task fmt
# or
task fmt
```

**Run Linter:**
```bash
task lint
# or
task lint
```

**Build and Run:**
```bash
task run
# or
task run
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
