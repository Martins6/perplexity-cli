package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"perplexity-cli/pkg/session"
)

var (
	listLimit int
)

// sessionListCmd lists recent sessions
var sessionListCmd = &cobra.Command{
	Use:   "list",
	Short: "List recent conversation sessions",
	Long: `List your recent conversation sessions sorted by date.

By default, shows the 10 most recent sessions. Use the -l/--limit flag
to change the number of sessions displayed.

Examples:
  pplx session list
  pplx session list -l 5
  pplx session list --limit 20`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Create session manager
		sessionManager, err := session.NewManager()
		if err != nil {
			return fmt.Errorf("failed to create session manager: %w", err)
		}

		// Get recent sessions
		sessions, err := sessionManager.ListRecent(listLimit)
		if err != nil {
			return fmt.Errorf("failed to list sessions: %w", err)
		}

		if len(sessions) == 0 {
			fmt.Println("No sessions found.")
			fmt.Println("Start a conversation with 'pplx' or run a query with 'pplx run \"<query>\"'")
			return nil
		}

		// Display results
		fmt.Printf("Recent sessions (showing %d of %d total):\n\n", len(sessions), len(sessions))

		for i, info := range sessions {
			// Format: 1. [shortid] Jan 02, 2006 15:04:05
			fmt.Printf("%d. [%s] %s\n", i+1, info.ShortID, session.FormatSessionTime(info.CreatedAt))
			fmt.Printf("   %s\n", session.TruncateQuery(info.InitialQuery, 60))
			fmt.Printf("   (%d messages)\n", info.MessageCount)

			// Add spacing between entries
			if i < len(sessions)-1 {
				fmt.Println()
			}
		}

		return nil
	},
}

func init() {
	// Add list command to the parent session command
	if sessionCmd != nil {
		sessionCmd.AddCommand(sessionListCmd)

		// Add the -l/--limit flag
		// Also add it to the parent session command for 'pplx session -l 10' syntax
		sessionListCmd.Flags().IntVarP(&listLimit, "limit", "l", 10, "Number of sessions to display")
		sessionCmd.Flags().IntVarP(&listLimit, "limit", "l", 10, "Number of sessions to display")
	}
}
