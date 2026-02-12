package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"perplexity-cli/pkg/session"
)

var cfgFile string
var shortcutContinue string
var shortcutListLimit int
var shortcutSearchQuery string

var rootCmd = &cobra.Command{
	Use:   "pplx",
	Short: "A CLI for the Perplexity Sonar API",
	Long: `PPLX is a command-line interface for interacting with the Perplexity Sonar API.

It supports one-shot queries, interactive conversations, and session management.

Examples:
  # One-shot query
  pplx run "What is the capital of France?"

  # Interactive mode
  pplx

  # List recent sessions
  pplx session -l 10

  # Search sessions
  pplx session search "France"

	Shortcuts:
  pplx -c [id]           Continue a session (same as: pplx session continue [id])
  pplx -l [limit]        List recent sessions (same as: pplx session list -l [limit])
  pplx -s [query]        Search sessions (same as: pplx session search [query])

Note: Shortcuts only work at the root level and cannot be combined with session commands.`,
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
	Run: func(cmd *cobra.Command, args []string) {
		// Handle shortcut -sc (continue session)
		if shortcutContinue != "" {
			if err := continueSession(shortcutContinue); err != nil {
				fmt.Fprintf(os.Stderr, "\033[31mError:\033[0m %v\n", err)
				os.Exit(1)
			}
			return
		}

		// Handle shortcut -sl (list sessions)
		if shortcutListLimit != 0 {
			sessionManager, err := session.NewManager()
			if err != nil {
				fmt.Fprintf(os.Stderr, "\033[31mError:\033[0m failed to create session manager: %v\n", err)
				os.Exit(1)
			}

			limit := shortcutListLimit
			if limit <= 0 {
				limit = 10
			}

			sessions, err := sessionManager.ListRecent(limit)
			if err != nil {
				fmt.Fprintf(os.Stderr, "\033[31mError:\033[0m failed to list sessions: %v\n", err)
				os.Exit(1)
			}

			if len(sessions) == 0 {
				fmt.Println("No sessions found.")
				fmt.Println("Start a conversation with 'pplx' or run a query with 'pplx run \"<query>\"'")
				return
			}

			fmt.Printf("Recent sessions (showing %d of %d total):\n\n", len(sessions), len(sessions))

			for i, info := range sessions {
				fmt.Printf("%d. [%s] %s\n", i+1, info.ShortID, session.FormatSessionTime(info.CreatedAt))
				fmt.Printf("   %s\n", session.TruncateQuery(info.InitialQuery, 60))
				fmt.Printf("   (%d messages)\n", info.MessageCount)

				if i < len(sessions)-1 {
					fmt.Println()
				}
			}
			return
		}

		// Handle shortcut -ss (search sessions)
		if shortcutSearchQuery != "" {
			sessionManager, err := session.NewManager()
			if err != nil {
				fmt.Fprintf(os.Stderr, "\033[31mError:\033[0m failed to create session manager: %v\n", err)
				os.Exit(1)
			}

			results, err := sessionManager.Search(shortcutSearchQuery)
			if err != nil {
				fmt.Fprintf(os.Stderr, "\033[31mError:\033[0m failed to search sessions: %v\n", err)
				os.Exit(1)
			}

			if len(results) == 0 {
				fmt.Println("No sessions found matching your query.")
				return
			}

			fmt.Printf("Found %d session(s) matching '%s':\n\n", len(results), shortcutSearchQuery)

			for i, info := range results {
				fmt.Printf("%d. [%s] %s\n", i+1, info.ShortID, session.FormatSessionTime(info.CreatedAt))
				fmt.Printf("   Query: %s\n", session.TruncateQuery(info.InitialQuery, 60))
				fmt.Printf("   Messages: %d\n", info.MessageCount)
				fmt.Println()
			}
			return
		}

		// If no arguments, enter interactive mode
		if len(args) == 0 {
			runInteractive()
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "\033[31mError:\033[0m %v\n", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pplx/config.yaml)")
	rootCmd.PersistentFlags().String("model", "sonar", "Perplexity model to use")
	viper.BindPFlag("model", rootCmd.PersistentFlags().Lookup("model"))
	rootCmd.Flags().StringVarP(&shortcutContinue, "shortcut-continue", "c", "", "Continue a session (shortcut for: pplx session continue [id])")
	rootCmd.Flags().IntVarP(&shortcutListLimit, "shortcut-list", "l", 0, "List recent sessions (shortcut for: pplx session list -l [limit])")
	rootCmd.Flags().StringVarP(&shortcutSearchQuery, "shortcut-search", "s", "", "Search sessions (shortcut for: pplx session search [query])")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
			os.Exit(1)
		}

		viper.AddConfigPath(home + "/.pplx")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	viper.SetEnvPrefix("PPLX")
	viper.AutomaticEnv()
}

// runInteractive is defined in interactive.go
