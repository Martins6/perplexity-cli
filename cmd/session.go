package cmd

import (
	"github.com/spf13/cobra"
)

// sessionCmd is the parent command for session management
var sessionCmd = &cobra.Command{
	Use:   "session",
	Short: "Manage conversation sessions",
	Long: `Manage your saved conversation sessions.

This command provides subcommands for listing, searching, and managing
your conversation history stored in ~/.pplx/sessions/.

Examples:
  # List recent sessions
  pplx session -l 10

  # Search sessions
  pplx session search "France"

  # List all sessions
  pplx session list`,
}

func init() {
	rootCmd.AddCommand(sessionCmd)
}
