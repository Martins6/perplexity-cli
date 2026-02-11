package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"perplexity-cli/pkg/config"
	"perplexity-cli/pkg/perplexity"
	"perplexity-cli/pkg/ui"
)

var runCmd = &cobra.Command{
	Use:   "run [query]",
	Short: "Send a one-shot query to Perplexity",
	Long: `Send a single query to the Perplexity Sonar API and display the response.

Examples:
  pplx run "What is the capital of France?"
  pplx run "Explain quantum computing" --model sonar-pro
  echo "What is 2+2?" | pplx run`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get query from argument or stdin
		var query string
		if len(args) > 0 {
			query = args[0]
		} else {
			// Try to read from stdin
			stat, _ := os.Stdin.Stat()
			if (stat.Mode() & os.ModeCharDevice) == 0 {
				// Data is being piped
				var buf strings.Builder
				buffer := make([]byte, 1024)
				for {
					n, err := os.Stdin.Read(buffer)
					if n > 0 {
						buf.Write(buffer[:n])
					}
					if err != nil {
						break
					}
				}
				query = strings.TrimSpace(buf.String())
			}
		}

		if query == "" {
			return fmt.Errorf("query is required\n\nUsage: pplx run \"<query>\"\n   or: echo \"<query>\" | pplx run")
		}

		// Load configuration
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}

		if err := cfg.Validate(); err != nil {
			return fmt.Errorf("configuration error: %w", err)
		}

		// Get model from flag or config
		model, _ := cmd.Flags().GetString("model")
		if model == "" {
			model = cfg.Model
		}

		// Create API client
		clientConfig := perplexity.DefaultConfig(cfg.APIKey)
		clientConfig.Model = model
		client := perplexity.NewClientWithConfig(clientConfig)

		// Prepare messages
		messages := []perplexity.Message{
			{Role: "user", Content: query},
		}

		// Make API request with spinner
		req := &perplexity.ChatCompletionRequest{
			Model:           model,
			Messages:        messages,
			MaxTokens:       cfg.MaxTokens,
			Temperature:     cfg.Temperature,
			TopP:            cfg.TopP,
			SearchMode:      cfg.SearchMode,
			ReasoningEffort: cfg.ReasoningEffort,
		}

		// Show progress spinner
		spinner := ui.NewSpinner("Thinking...")
		spinner.Start()

		resp, err := client.CreateCompletionWithRequest(req)
		spinner.Stop()

		if err != nil {
			return fmt.Errorf("API request failed: %w", err)
		}

		// Parse response
		parsed := perplexity.ParseResponse(resp)

		// Display formatted response with references
		formatted := perplexity.FormatWithReferences(parsed)
		fmt.Println(formatted)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
