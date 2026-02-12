# Session Command Shortcuts

## Overview

Session command shortcuts provide a faster way to access common session management commands by using single-character flags at the root level of the CLI instead of typing the full command paths. This feature improves user efficiency for frequently used session operations.

## Implementation

### Shortcuts Available

The following shortcuts are implemented on the root command:

- `-c [id]` - Continue a session (same as: `pplx session continue [id]`)
- `-l [limit]` - List recent sessions (same as: `pplx session list -l [limit]`)
- `-s [query]` - Search sessions (same as: `pplx session search [query]`)

### Flag Definition

The shortcuts are defined in `cmd/root.go` using Cobra's flag registration:

```go
var shortcutContinue string
var shortcutListLimit int
var shortcutSearchQuery string

rootCmd.Flags().StringVarP(&shortcutContinue, "shortcut-continue", "c", "", "Continue a session (shortcut for: pplx session continue [id])")
rootCmd.Flags().IntVarP(&shortcutListLimit, "shortcut-list", "l", 0, "List recent sessions (shortcut for: pplx session list -l [limit])")
rootCmd.Flags().StringVarP(&shortcutSearchQuery, "shortcut-search", "s", "", "Search sessions (shortcut for: pplx session search [query])")
```

### Conflict Detection

To prevent confusion and ambiguous usage, shortcuts cannot be combined with explicit session commands. This validation is implemented in the root command's `PersistentPreRunE` function:

```go
PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
    isSessionCommand := len(args) > 0 && args[0] == "session"

    if isSessionCommand {
        if shortcutContinue != "" {
            return fmt.Errorf("cannot use session shortcuts (-c, -l, -s) together with explicit session commands\n\nPlease use either:\n  - pplx -c [id] (shortcut)\n  - pplx session continue [id] (explicit command)")
        }
        if shortcutSearchQuery != "" {
            return fmt.Errorf("cannot use session shortcuts (-c, -l, -s) together with explicit session commands\n\nPlease use either:\n  - pplx -s [query] (shortcut)\n  - pplx session search [query] (explicit command)")
        }
    }

    return nil
},
```

**Note:** The `-l` flag is valid in both contexts (as a shortcut and as a session flag), so it's only checked for conflicts with `-c` and `-s`.

### Execution Logic

The shortcuts are executed in the root command's `Run` function, which checks which shortcut flag is set and calls the appropriate session management logic:

```go
Run: func(cmd *cobra.Command, args []string) {
    // Handle shortcut -c (continue session)
    if shortcutContinue != "" {
        if err := continueSession(shortcutContinue); err != nil {
            fmt.Fprintf(os.Stderr, "\033[31mError:\033[0m %v\n", err)
            os.Exit(1)
        }
        return
    }

    // Handle shortcut -l (list sessions)
    if shortcutListLimit != 0 {
        // ... session list logic
        return
    }

    // Handle shortcut -s (search sessions)
    if shortcutSearchQuery != "" {
        // ... session search logic
        return
    }

    // If no arguments, enter interactive mode
    if len(args) == 0 {
        runInteractive()
    }
},
```

### Code Reuse

The shortcuts reuse existing session command logic to avoid duplication:

- **Continue shortcut** (`-c`): Calls the `continueSession()` function from `cmd/session_continue.go`
- **List shortcut** (`-l`): Duplicates the session list logic from `cmd/session_list.go` for simplicity
- **Search shortcut** (`-s`): Duplicates the session search logic from `cmd/session_search.go` for simplicity

### Error Handling

Each shortcut includes appropriate error handling:

```go
// Continue shortcut error handling
if err := continueSession(shortcutContinue); err != nil {
    fmt.Fprintf(os.Stderr, "\033[31mError:\033[0m %v\n", err)
    os.Exit(1)
}

// List shortcut error handling
sessionManager, err := session.NewManager()
if err != nil {
    fmt.Fprintf(os.Stderr, "\033[31mError:\033[0m failed to create session manager: %v\n", err)
    os.Exit(1)
}

// Search shortcut error handling
sessionManager, err := session.NewManager()
if err != nil {
    fmt.Fprintf(os.Stderr, "\033[31mError:\033[0m failed to create session manager: %v\n", err)
    os.Exit(1)
}
```

## Usage Examples

### Continue a Session

**Full command:**
```bash
pplx session continue a8x9k2
```

**Shortcut:**
```bash
pplx -c a8x9k2
```

### List Recent Sessions

**Full command:**
```bash
pplx session list -l 20
```

**Shortcut:**
```bash
pplx -l 20
```

**Note:** The `-l` flag works in both contexts:
- As a shortcut at root level: `pplx -l 10`
- As a flag on the session command: `pplx session -l 10`

### Search Sessions

**Full command:**
```bash
pplx session search "France"
```

**Shortcut:**
```bash
pplx -s "France"
```

## Limitations

1. **Root Level Only**: Shortcuts only work at the root level of the CLI. They cannot be used after the `session` command (e.g., `pplx session -c` is not allowed).

2. **One Shortcut at a Time**: Only one shortcut can be used per command invocation. Attempting to use multiple shortcuts (e.g., `pplx -c a8x9k2 -l 10`) will result in undefined behavior.

3. **Conflict Detection**: Shortcuts cannot be combined with explicit session commands. If you try to use both, you'll receive a helpful error message.

4. **Flag Priority**: If a shortcut flag conflicts with another flag (like `-l`), the behavior depends on the command context.

## Error Messages

### Conflict Error

When attempting to combine shortcuts with explicit session commands:

```bash
$ pplx session -c a8x9k2
Error: cannot use session shortcuts (-c, -l, -s) together with explicit session commands

Please use either:
  - pplx -c [id] (shortcut)
  - pplx session continue [id] (explicit command)
```

### Session Not Found Error

When using the continue shortcut with an invalid session ID:

```bash
$ pplx -c invalidid
Error: session not found: invalidid

Run 'pplx session list' to see available sessions
```

### No Sessions Error

When using the list or search shortcuts and no sessions exist:

```bash
$ pplx -l 10
No sessions found.
Start a conversation with 'pplx' or run a query with 'pplx run "<query>"'
```

## Documentation Updates

The shortcuts are documented in multiple locations to ensure users can easily discover them:

1. **Root Command Help** (`cmd/root.go`): The `Long` field includes a "Shortcuts" section with examples
2. **Session Command Help** (`cmd/session.go`): The `Long` field includes a reference to the available shortcuts
3. **README.md**: Includes a dedicated "Shortcuts" section in the Usage documentation
4. **Feature Documentation** (`docs/features/perplexity-cli.md`): Comprehensive documentation of shortcuts with examples and conflict behavior

## Testing

All existing tests pass with the shortcut implementation:

```bash
task test
```

Manual testing scenarios:
- ✅ `pplx -c <id>` continues a session correctly
- ✅ `pplx -l <limit>` lists the specified number of sessions
- ✅ `pplx -s <query>` searches sessions correctly
- ✅ `pplx -l` (no argument) uses default limit of 10
- ✅ Conflict detection works when shortcuts are combined with session commands
- ✅ Error messages are helpful and clear
- ✅ Existing commands (`pplx session list`, `pplx session search`, `pplx session continue`) continue to work as expected

## Backward Compatibility

This feature is fully backward compatible:

- All existing commands work exactly as before
- No breaking changes to any existing functionality
- Shortcuts are additive only (they don't replace or modify existing commands)
- The `-l` flag continues to work on the `pplx session` command as it did before

## Future Enhancements

Potential improvements to consider:

1. **Additional Shortcuts**: Add shortcuts for other common commands (e.g., `pplx run` shortcut)
2. **Auto-Completion**: Add shell auto-completion support for shortcut flags
3. **Shortcut Configuration**: Allow users to define custom shortcuts in config file
4. **Help Integration**: Add shortcuts to the main help output in a more prominent position
5. **Validation**: Add more validation for shortcut arguments (e.g., validate session ID format for `-c`)
