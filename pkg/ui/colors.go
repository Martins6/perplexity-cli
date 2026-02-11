package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

var (
	// Colors for different output types
	Green   = color.New(color.FgGreen).SprintFunc()
	Red     = color.New(color.FgRed).SprintFunc()
	Yellow  = color.New(color.FgYellow).SprintFunc()
	Blue    = color.New(color.FgBlue).SprintFunc()
	Cyan    = color.New(color.FgCyan).SprintFunc()
	Magenta = color.New(color.FgMagenta).SprintFunc()
	White   = color.New(color.FgWhite).SprintFunc()
	Bold    = color.New(color.Bold).SprintFunc()

	// Specific use colors
	PromptColor  = color.New(color.FgCyan, color.Bold).SprintFunc()
	ErrorColor   = color.New(color.FgRed, color.Bold).SprintFunc()
	SuccessColor = color.New(color.FgGreen).SprintFunc()
	InfoColor    = color.New(color.FgBlue).SprintFunc()
	WarningColor = color.New(color.FgYellow).SprintFunc()
	HeaderColor  = color.New(color.FgWhite, color.Bold, color.Underline).SprintFunc()
)

// DisableColors disables all color output (useful for pipes/non-tty)
func DisableColors() {
	color.NoColor = true
}

// EnableColors enables color output
func EnableColors() {
	color.NoColor = false
}

// IsTerminal returns true if stdout is a terminal
func IsTerminal() bool {
	return color.NoColor == false || os.Getenv("FORCE_COLOR") != ""
}

// PrintError prints an error message in red
func PrintError(format string, args ...interface{}) {
	color.Red(format, args...)
}

// PrintSuccess prints a success message in green
func PrintSuccess(format string, args ...interface{}) {
	color.Green(format, args...)
}

// PrintInfo prints an info message in blue
func PrintInfo(format string, args ...interface{}) {
	color.Blue(format, args...)
}

// PrintWarning prints a warning message in yellow
func PrintWarning(format string, args ...interface{}) {
	color.Yellow(format, args...)
}

// PrintPrompt prints a prompt in cyan
func PrintPrompt(text string) string {
	return PromptColor(text)
}

// PrintSeparator prints a colored separator line
func PrintSeparator(colorFunc func(...interface{}) string) {
	fmt.Println(colorFunc(strings.Repeat("─", 60)))
}
