# Session Continue Command

## Overview

The `pplx session continue` command allows users to load an existing conversation session and continue the conversation interactively. Users can provide a new message that is sent to the Perplexity API with the full conversation context, receiving a response that is displayed with citations and saved to the session.

## Command Usage

```bash
pplx session continue [id]
```

The ID can be either:
- Short ID (e.g., `a8x9k2`) - the 6-7 character alphanumeric code shown in session list
- Full timestamp ID (e.g., `20240115-103045.123`) - for backward compatibility

## Implementation

### Files Created
- `cmd/session_continue.go` - New command file implementing the `session continue` subcommand

### Key Features

1. **Session Loading**: 
   - Attempts to load session by short ID first
   - Falls back to full timestamp ID if short ID lookup fails
   - Returns helpful error message if session not found, suggesting `pplx session list`

2. **Session History Display**:
   - Uses `session.DisplaySession()` to show conversation context
   - Displays session metadata (ID, model, created time)

3. **Interactive Input**:
   - Prompts user with `You:` using colored output
   - Uses `bufio.NewReader` for input handling
   - Exits gracefully if no message is provided

4. **Context Management**:
   - Retrieves last N messages using `session.GetLastMessages(MaxContextMessages)`
   - Strips "## References:" section from assistant messages using `perplexity.StripReferences()`
   - Converts session messages to API message format
   - Appends new user message to context

5. **API Integration**:
   - Creates API client with session's model configuration
   - Builds request with conversation context and new message
   - Includes all configurable parameters (max tokens, temperature, top P, search mode, reasoning effort)

6. **Response Display**:
   - Displays separator line before response (magenta color)
   - Parses response using `perplexity.ParseResponse()`
   - Formats with references using `perplexity.FormatWithReferences()`
   - Renders markdown using `ui.RenderMarkdown()` with graceful fallback to plain text

7. **Session Persistence**:
   - Adds new user message to session
   - Adds assistant response (with references stripped)
   - Saves updated session
   - Warns if save fails but doesn't fail the operation

### Error Handling

- Configuration validation before proceeding
- Helpful error messages for session not found
- API errors wrapped with context
- Non-fatal warning for save failures
- Empty input detection and graceful exit

### Command Documentation

The command includes comprehensive help text accessible via `pplx session continue --help` that explains:
- ID format options (short vs full)
- Command workflow steps
- Usage examples

## Technical Details

### Context Window Management

To avoid token limits, only the last 20 messages (configurable via `MaxContextMessages`) are included in the API request. This balances context preservation with API efficiency.

### Reference Stripping

Critical implementation detail: assistant messages must have the "## References:" section stripped before being sent back to the API. Without this, the API may be confused by formatted references. This is handled using `perplexity.StripReferences()`.

### Session ID Fallback Pattern

The command follows the same pattern as `session_show.go`:
1. Try `sessionManager.LoadByShortID(sessionID)`
2. If that fails, try `sessionManager.Load(sessionID)` (full ID)
3. If both fail, return error with helpful message

### Colored Output

The command uses the UI color utilities:
- Cyan color for user input separator
- Magenta color for response separator
- Markdown rendering for formatted response

## Dependencies

- Go 1.25.2
- Existing session management (`pkg/session`)
- Perplexity API client (`pkg/perplexity`)
- UI utilities (`pkg/ui`)
- Configuration management (`pkg/config`)

## Related Features

- Session listing (`pplx session list`) - provides session IDs
- Session show (`pplx session show`) - displays session details
- Interactive mode - similar interactive prompt behavior
- One-shot queries (`pplx run`) - single query without session persistence
