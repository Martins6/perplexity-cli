package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"perplexity-cli/pkg/config"
	"perplexity-cli/pkg/perplexity"
	"perplexity-cli/pkg/session"
	"perplexity-cli/pkg/ui"
)

// sessionContinueCmd loads an existing session and continues the conversation
var sessionContinueCmd = &cobra.Command{
	Use:   "continue [id]",
	Short: "Continue a conversation session",
	Long: `Load an existing conversation session and prompt for a new message to continue the conversation.

The ID can be either:
- Short ID (e.g., a8x9k2) - the 6-7 character alphanumeric code shown in session list
- Full timestamp ID (e.g., 20240115-103045.123) - for backward compatibility

This command will:
1. Display the session history
2. Prompt you for a new message
3. Send your message with conversation context to Perplexity
4. Display the response with citations
5. Save the updated session

Examples:
  pplx session continue a8x9k2
  pplx session continue 20240115-103045.123`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		return continueSession(id)
	},
}

func init() {
	if sessionCmd != nil {
		sessionCmd.AddCommand(sessionContinueCmd)
	}
}

// continueSession handles the workflow for continuing a session
func continueSession(sessionID string) error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("configuration error: %w", err)
	}

	// Create session manager
	sessionManager, err := session.NewManager()
	if err != nil {
		return fmt.Errorf("failed to create session manager: %w", err)
	}

	// Load session (try short ID first, then full ID)
	var s *session.Session
	s, err = sessionManager.LoadByShortID(sessionID)
	if err != nil {
		s, err = sessionManager.Load(sessionID)
		if err != nil {
			return fmt.Errorf("session not found: %s\n\nRun 'pplx session list' to see available sessions", sessionID)
		}
	}

	// Display session history
	if err := session.DisplaySession(s); err != nil {
		return fmt.Errorf("failed to display session: %w", err)
	}

	// Display separator
	fmt.Println()
	ui.PrintSeparator(ui.Cyan)

	// Prompt for new message
	fmt.Print("You: ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	input = strings.TrimSpace(input)
	if input == "" {
		fmt.Println("No message provided. Exiting.")
		return nil
	}

	// Create API client
	clientConfig := perplexity.DefaultConfig(cfg.APIKey)
	clientConfig.Model = s.Metadata.Model
	client := perplexity.NewClientWithConfig(clientConfig)

	// Build API messages with conversation context
	messages := s.GetLastMessages(MaxContextMessages)
	apiMessages := make([]perplexity.Message, 0, len(messages)+1)

	// Convert session messages, stripping references from assistant messages
	for _, msg := range messages {
		content := msg.Content
		if msg.Role == "assistant" {
			content = perplexity.StripReferences(content)
		}
		apiMessages = append(apiMessages, perplexity.Message{
			Role:    msg.Role,
			Content: content,
		})
	}

	// Add the new user message
	apiMessages = append(apiMessages, perplexity.Message{
		Role:    "user",
		Content: input,
	})

	// Make API request
	req := &perplexity.ChatCompletionRequest{
		Model:           s.Metadata.Model,
		Messages:        apiMessages,
		MaxTokens:       cfg.MaxTokens,
		Temperature:     cfg.Temperature,
		TopP:            cfg.TopP,
		SearchMode:      cfg.SearchMode,
		ReasoningEffort: cfg.ReasoningEffort,
	}

	resp, err := client.CreateCompletionWithRequest(req)
	if err != nil {
		return fmt.Errorf("API request failed: %w", err)
	}

	// Parse response
	parsed := perplexity.ParseResponse(resp)

	// Display separator
	fmt.Println()
	ui.PrintSeparator(ui.Magenta)
	fmt.Print("PPLX: ")

	// Display response with markdown rendering
	formatted := perplexity.FormatWithReferences(parsed)
	rendered, err := ui.RenderMarkdown(formatted, cfg)
	if err != nil {
		fmt.Println(formatted)
	} else {
		fmt.Println(rendered)
	}
	fmt.Println()

	// Update session with new messages
	s.AddMessage("user", input)
	cleanContent := perplexity.StripReferences(parsed.Content)
	s.AddMessage("assistant", cleanContent)

	// Save updated session
	if err := sessionManager.Save(s); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to save session: %v\n", err)
	} else {
		session.Debugf("Session updated: %s", s.ID)
	}

	return nil
}
