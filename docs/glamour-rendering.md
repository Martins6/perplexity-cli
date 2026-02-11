# Glamour Markdown Rendering

This document describes the markdown rendering integration using Charmbracelet Glamour in the perplexity-cli tool.

## Overview

The interactive mode of `pplx` uses [Charmbracelet Glamour](https://github.com/charmbracelet/glamour) to render beautiful, styled markdown output with syntax highlighting, proper formatting, and theme support.

## Features

### Syntax Highlighting
Code blocks are automatically highlighted with appropriate colors based on the programming language:
```markdown
\`\`\`go
func main() {
    fmt.Println("Hello, World!")
}
\`\`\`
```

### Rich Formatting
All markdown elements are properly rendered:
- Headers (#, ##, ###)
- Bold (**text**) and italic (*text*)
- Lists (unordered and ordered)
- Blockquotes (> quote)
- Tables (if supported)
- Links and images
- Horizontal rules (---)

### Multiple Themes
Choose from built-in themes or use custom stylesheets:
- **auto** - Automatically detects terminal background color (default)
- **dark** - Dark theme
- **light** - Light theme
- **Custom** - Path to a custom JSON stylesheet

### Terminal Detection
The renderer automatically detects whether output is a terminal:
- If output is a terminal: Uses markdown rendering
- If output is piped/redirected: Falls back to plain text
- Configurable via `use_glow` option

## Configuration

### Environment Variables

```bash
export PPLX_USE_GLOW=true
export PPLX_GLOW_STYLE=dark
export PPLX_GLOW_WIDTH=80
```

### Config File (~/.pplx/config.yaml)

```yaml
use_glow: true
glow_style: auto
glow_width: 0
```

## Configuration Options

### use_glow
- **Type:** `bool`
- **Default:** `true`
- **Description:** Enable or disable markdown rendering in interactive mode
- **Example:** `use_glow: false` (disables rendering, uses plain text)

### glow_style
- **Type:** `string`
- **Default:** `"auto"`
- **Description:** Theme style for markdown rendering
- **Options:**
  - `"auto"` - Automatically detects terminal background (recommended)
  - `"dark"` - Dark theme
  - `"light"` - Light theme
  - `"/path/to/style.json"` - Custom stylesheet path
- **Example:** `glow_style: dark`

### glow_width
- **Type:** `int`
- **Default:** `0`
- **Description:** Word wrap width in characters
- **Behavior:**
  - `0` - Use terminal width (default)
  - `>0` - Fixed width, content wraps at this column
- **Example:** `glow_width: 80`

## Usage Examples

### Interactive Mode

Markdown rendering is automatically enabled in interactive mode:

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

Use the dark theme:

```yaml
# ~/.pplx/config.yaml
use_glow: true
glow_style: dark
```

Use the light theme:

```yaml
use_glow: true
glow_style: light
```

### Custom Word Wrap

Set a fixed width of 80 characters:

```yaml
use_glow: true
glow_style: auto
glow_width: 80
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

## Implementation Details

### How It Works

1. When a response is received from the Perplexity API, it's formatted with citations
2. The formatted content is passed through `ui.RenderMarkdown(content, cfg)`
3. If `use_glow` is true and output is a terminal:
   - A Glamour TermRenderer is created with the configured options
   - The markdown is rendered to ANSI codes
   - The rendered output is displayed
4. Otherwise, the original plain text is displayed

### Graceful Fallback

The implementation includes several fallback mechanisms:

1. **Config disabled** - If `use_glow: false`, returns plain text
2. **Not a terminal** - If output is piped/redirected, returns plain text
3. **Renderer error** - If Glamour fails to initialize, returns plain text
4. **Nil config** - If config is nil, returns plain text

### Terminal Detection

Terminal detection is performed using the `ui.IsTerminal()` function:

```go
func IsTerminal() bool {
    return color.NoColor == false || os.Getenv("FORCE_COLOR") != ""
}
```

You can force markdown rendering even when piping:

```bash
FORCE_COLOR=1 pplx | less -R
```

## Performance Considerations

- **Small overhead** - Glamour rendering adds minimal overhead (~1-5ms per response)
- **Memory** - Temporary memory usage during rendering (~1-2MB)
- **Impact:** Negligible for typical responses; acceptable in interactive mode

## Troubleshooting

### Rendering Looks Broken

**Problem:** Markdown is not displaying correctly or showing garbled text.

**Solutions:**
1. Check your terminal supports ANSI escape codes
2. Try a different style: `glow_style: dark`
3. Disable rendering: `use_glow: false`
4. Check `TERM` environment variable: `echo $TERM`

### Colors Look Wrong

**Problem:** Colors don't look good on your terminal.

**Solutions:**
1. Use "auto" style to detect terminal background: `glow_style: auto`
2. Manually select dark or light based on your terminal: `glow_style: dark`
3. Try forcing a different terminal type

### Rendering Is Too Slow

**Problem:** Responses take longer to display.

**Solutions:**
1. Disable rendering: `use_glow: false`
2. Check your system resources
3. Try reducing word wrap width: `glow_width: 60`

### Piped Output Shows ANSI Codes

**Problem:** When piping to a file, you see ANSI escape codes.

**Solution:** This is expected behavior. Plain text is only used when NOT a terminal. To get plain text in a pipe, disable rendering:

```bash
PPLX_USE_GLOW=false pplx > output.txt
```

## Creating Custom Styles

You can create custom JSON stylesheets for Glamour. For details, see the [Glamour documentation](https://github.com/charmbracelet/glamour#custom-styles).

Example custom style:

```json
{
  "document": {
    "prefix": "",
    "suffix": "",
    "block_prefix": "\n",
    "block_suffix": "\n",
    "color": "229",
    "indent": 2
  },
  "block_quote": {
    "prefix": " │ ",
    "indent": 2,
    "color": "240",
    "indent_token": " │ "
  },
  "paragraph": {
    "prefix": "",
    "suffix": "",
    "indent": 0
  },
  "list": {
    "prefix": "• ",
    "indent": 2
  },
  "heading": {
    "prefix": "# ",
    "color": "33",
    "bold": true
  },
  "code_block": {
    "theme": "Monokai Extended",
    "chroma": {
      "background": "#262626"
    },
    "width": 80,
    "indent": 2,
    "prefix": "  ",
    "block_suffix": "\n"
  }
}
```

Use a custom style:

```yaml
glow_style: /path/to/custom-style.json
```

## Session Display

Markdown rendering is also applied when viewing saved sessions:

```bash
$ pplx session show abc123
```

The response messages will be rendered with markdown formatting, just like in interactive mode.

## References

- [Charmbracelet Glamour](https://github.com/charmbracelet/glamour) - Markdown rendering library
- [Glamour Styles Gallery](https://github.com/charmbracelet/glamour/tree/master/styles) - Available built-in styles
- [Chroma Syntax Highlighting](https://github.com/alecthomas/chroma) - Syntax highlighter used by Glamour
