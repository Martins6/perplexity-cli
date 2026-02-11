package session

import (
	"fmt"

	"perplexity-cli/pkg/config"
	"perplexity-cli/pkg/perplexity"
	"perplexity-cli/pkg/ui"
)

// DisplaySession displays a full session conversation with formatting
func DisplaySession(s *Session) error {
	cfg, err := config.Load()
	if err != nil {
		cfg = config.DefaultConfig()
	}

	// Display header with metadata
	fmt.Println()
	ui.PrintSeparator(ui.HeaderColor)
	fmt.Printf("Session: [%s] %s\n", s.ShortID, FormatSessionTime(s.Metadata.CreatedAt))
	fmt.Printf("Model: %s\n", s.Metadata.Model)
	fmt.Printf("Messages: %d\n", len(s.Messages))
	ui.PrintSeparator(ui.HeaderColor)
	fmt.Println()

	// Display each message
	for i, msg := range s.Messages {
		if msg.Role == "user" {
			// Print separator before user message (except first)
			if i > 0 {
				ui.PrintSeparator(ui.Cyan)
			}
			fmt.Print("You: ")
			fmt.Println(msg.Content)
		} else if msg.Role == "assistant" {
			fmt.Println()
			ui.PrintSeparator(ui.Magenta)
			fmt.Print("PPLX: ")

			// Check if content has citations and format accordingly
			formatted := formatMessageWithCitations(msg.Content)
			rendered, err := ui.RenderMarkdown(formatted, cfg)
			if err != nil {
				fmt.Println(formatted)
			} else {
				fmt.Println(rendered)
			}
			fmt.Println()
		}
	}

	return nil
}

// formatMessageWithCitations checks for citations and formats them
func formatMessageWithCitations(content string) string {
	// Extract citations from content
	citations := perplexity.ExtractCitations(content)

	if len(citations) == 0 {
		return content
	}

	// Content already has citations inline (e.g., [1], [2])
	// Since we don't have search results stored in the session,
	// we can't show the full references section
	// Just return the content as-is with citations preserved

	return content
}

// DisplaySessionSummary displays a brief summary of the session
func DisplaySessionSummary(s *Session) {
	fmt.Printf("[%s] %s - %s (%d messages)\n",
		s.ShortID,
		FormatSessionTime(s.Metadata.CreatedAt),
		TruncateQuery(s.Metadata.InitialQuery, 50),
		len(s.Messages))
}
