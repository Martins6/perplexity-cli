# Perplexity CLI Implementation Plan

## Problem Statement
Build a Go CLI application using Cobra framework to interact with the Perplexity Sonar API. The CLI should support:
1. One-shot queries (`pplx run "question"`)
2. Interactive conversation mode (`pplx`)
3. Session persistence in `~/.pplx/sessions/`
4. Session search and listing functionality

## Solution Overview
Create a well-structured Go CLI using Cobra with the following architecture:
- **Perplexity API Client**: Custom lightweight HTTP client for the chat completions endpoint
- **Session Management**: JSON-based conversation storage with timestamp naming
- **Interactive REPL**: Line-by-line interface with conversation context tracking
- **Citation Parser**: Extract and format reference numbers with URLs
- **Search/List Commands**: Grep-based session file utilities

## User Preferences
- Default model: `sonar`
- API key: Environment variable `PPLX_API_KEY`
- Citation format: Bracket numbers at end of sentences (e.g., "Paris is the capital of France. [1]")
- Persistence: Auto-save sessions immediately after each exchange

## Implementation Steps

### Step 1: Project Setup and Dependencies
**Files to Create:**
- `go.mod` (via `go mod init`)
- `go.sum` (via `go get`)
- `main.go` (entry point)
- `cmd/root.go` (root command)

**Dependencies:** cobra, viper

### Step 2: Create Core API Client
**Files to Create:**
- `pkg/perplexity/client.go` - HTTP client with API key auth
- `pkg/perplexity/types.go` - Request/response structures
- `pkg/perplexity/citations.go` - Citation extraction and formatting

**Key Details:**
- API endpoint: `https://api.perplexity.ai/chat/completions`
- Parse content for bracket citations `[1]`, `[2]`, etc.
- Extract `search_results` array for reference URLs
- Format references at end of output with `# References` heading

### Step 3: Session Management System
**Files to Create:**
- `pkg/session/manager.go` - Session CRUD operations
- `pkg/session/types.go` - Session JSON structure
- `pkg/session/utils.go` - Timestamp generation, directory management

**Key Details:**
- Storage path: `~/.pplx/sessions/`
- Filename format: `20060102-150405.000.json` (timestamp with milliseconds)
- Auto-save after every API exchange in interactive mode

### Step 4: Implement `pplx run` Command
**Files to Create:**
- `cmd/run.go` - `pplx run` subcommand

**Command Syntax:** `pplx run "query string" [--model sonar]`

**Key Details:**
- Accept query as positional argument
- Make single API call with user query
- Display content with citations inline
- Add `# References` section at end

### Step 5: Implement Interactive Mode (`pplx`)
**Files to Create:**
- `cmd/interactive.go` or integrate in root.go

**Key Details:**
- Detect no arguments → enter interactive mode
- Display prompt: `You: `
- Exit commands: `/q`, `/quit`, `ctrl+c`
- Maintain conversation history in memory
- **CRITICAL**: Don't send reference content to API (strip references section)
- Auto-save session after each exchange
- Format: `PPLX: <content with citations>\n\n## References:\n[1] <title> - <url>`

### Step 6: Implement Session Search Command
**Files to Create:**
- `cmd/session.go` - Parent session command
- `cmd/session_search.go` - Search subcommand

**Command Syntax:** `pplx session search "search term"`

**Key Details:**
- Search through `~/.pplx/sessions/*.json`
- Return matching filenames with context
- Case-insensitive search

### Step 7: Implement Session List Command
**Files to Create:**
- `cmd/session_list.go` - List subcommand

**Command Syntax:** `pplx session -l 10` or `pplx session list -n 10`

**Key Details:**
- Sort by filename (timestamp) descending
- Display: `<timestamp> - <initial_query_truncated>`
- Default limit: 10 sessions

### Step 8: Configuration and Environment Setup
**Files to Create:**
- `pkg/config/config.go` - Configuration management

**Key Details:**
- Read API key from `PPLX_API_KEY` environment variable
- Support config file at `~/.pplx/config.yaml` for defaults
- Config options: default_model, max_tokens, temperature

### Step 9: Error Handling and User Experience
**Files to Modify:**
- All command files for error messages
- Add `--help` documentation

**Key Details:**
- Graceful API error handling
- Clear error messages for missing API key
- Signal handling for clean shutdown
- Colorized output (optional)

### Step 10: Testing and Build
**Files to Create:**
- `*_test.go` files for core packages
- `Makefile` for build automation

## Project Structure
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
│   └── config/
│       └── config.go     # Configuration management
├── main.go               # Entry point
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## Key Challenges
1. **Citation parsing complexity** - Implement regex pattern `\[\d+\]` to find citations
2. **Conversation context management** - Maintain separate arrays: one for display (with citations) and one for API (clean content)
3. **Session file corruption** - Write to temp file then atomic rename
4. **Signal handling in REPL** - Use `signal.Notify` for graceful shutdown
5. **Large conversation history** - Limit context window to last N messages

## Verification Steps
1. Unit tests for citation parsing, session file operations, API client
2. Integration tests for one-shot queries and interactive mode
3. Manual verification of session JSON structure
4. Test search and list commands
