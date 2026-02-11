package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

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
  pplx session search "France"`,
	Run: func(cmd *cobra.Command, args []string) {
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
