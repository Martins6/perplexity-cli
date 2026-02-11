package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"perplexity-cli/pkg/session"
)

// sessionShowCmd displays a specific session by ID
var sessionShowCmd = &cobra.Command{
	Use:   "show [id]",
	Short: "Show a specific conversation session",
	Long: `Display a full conversation session by its ID.

The ID can be either:
- Short ID (e.g., a8x9k2) - the 6-7 character alphanumeric code shown in brackets
- Full timestamp ID (e.g., 20240115-103045.123) - for backward compatibility

Examples:
  pplx session show a8x9k2
  pplx session show 20240115-103045.123`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		// Create session manager
		sessionManager, err := session.NewManager()
		if err != nil {
			return fmt.Errorf("failed to create session manager: %w", err)
		}

		// Try to load by short ID first, then by full ID
		var s *session.Session
		s, err = sessionManager.LoadByShortID(id)
		if err != nil {
			// Try loading as full timestamp ID
			s, err = sessionManager.Load(id)
			if err != nil {
				return fmt.Errorf("session not found: %s", id)
			}
		}

		// Display the session using display utility
		return session.DisplaySession(s)
	},
}

func init() {
	if sessionCmd != nil {
		sessionCmd.AddCommand(sessionShowCmd)
	}
}
