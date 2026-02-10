package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"perplexity-cli/pkg/session"
)

// sessionSearchCmd searches through sessions
var sessionSearchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search through conversation sessions",
	Long: `Search through your saved conversation sessions for specific content.

The search is case-insensitive and matches against:
- Initial query
- All conversation messages

Examples:
  pplx session search "France"
  pplx session search "quantum computing"
  pplx session search "capital of"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := args[0]
		if query == "" {
			return fmt.Errorf("search query cannot be empty")
		}

		// Create session manager
		sessionManager, err := session.NewManager()
		if err != nil {
			return fmt.Errorf("failed to create session manager: %w", err)
		}

		// Search sessions
		results, err := sessionManager.Search(query)
		if err != nil {
			return fmt.Errorf("failed to search sessions: %w", err)
		}

		if len(results) == 0 {
			fmt.Println("No sessions found matching your query.")
			return nil
		}

		// Display results
		fmt.Printf("Found %d session(s) matching '%s':\n\n", len(results), query)

		for i, info := range results {
			fmt.Printf("%d. %s\n", i+1, session.FormatSessionTime(info.CreatedAt))
			fmt.Printf("   ID: %s\n", info.ID)
			fmt.Printf("   Query: %s\n", session.TruncateQuery(info.InitialQuery, 60))
			fmt.Printf("   Messages: %d\n", info.MessageCount)
			fmt.Println()
		}

		return nil
	},
}

func init() {
	// Add search command to the parent session command
	if sessionCmd != nil {
		sessionCmd.AddCommand(sessionSearchCmd)
	}
}
