package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"perplexity-cli/pkg/config"
	"perplexity-cli/pkg/perplexity"
	"perplexity-cli/pkg/session"
	"perplexity-cli/pkg/ui"
)

const (
	MaxContextMessages = 20
	ExitCommand        = "/q"
	AltExitCommand     = "/quit"
)

// InteractiveSession manages an interactive conversation
type InteractiveSession struct {
	client         *perplexity.Client
	sessionManager *session.Manager
	session        *session.Session
	config         *config.Config
	reader         *bufio.Reader
	quit           chan bool
	firstMessage   bool
}

// NewInteractiveSession creates a new interactive session
func NewInteractiveSession(cfg *config.Config) (*InteractiveSession, error) {
	clientConfig := perplexity.DefaultConfig(cfg.APIKey)
	clientConfig.Model = cfg.Model

	sessionManager, err := session.NewManager()
	if err != nil {
		return nil, fmt.Errorf("failed to create session manager: %w", err)
	}

	return &InteractiveSession{
		client:         perplexity.NewClientWithConfig(clientConfig),
		sessionManager: sessionManager,
		config:         cfg,
		reader:         bufio.NewReader(os.Stdin),
		quit:           make(chan bool),
		firstMessage:   true,
	}, nil
}

// Run starts the interactive REPL
func (is *InteractiveSession) Run() error {
	// Set up signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nReceived interrupt signal. Saving session...")
		is.saveSession()
		is.quit <- true
	}()

	fmt.Println("Welcome to PPLX Interactive Mode!")
	fmt.Printf("Model: %s | Type '%s' or press Ctrl+C to exit\n\n", is.config.Model, ExitCommand)

	// Main loop
	for {
		select {
		case <-is.quit:
			fmt.Println("Goodbye!")
			return nil
		default:
			if err := is.processInput(); err != nil {
				if err.Error() == "exit" {
					is.saveSession()
					fmt.Println("Goodbye!")
					return nil
				}
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			}
		}
	}
}

// processInput processes a single user input
func (is *InteractiveSession) processInput() error {
	// Display separator for consecutive messages
	if !is.firstMessage {
		ui.PrintSeparator(ui.Cyan)
	}

	// Display prompt
	fmt.Print("You: ")

	// Read input
	input, err := is.reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	// Trim whitespace and newlines
	input = strings.TrimSpace(input)

	// Check for exit command
	if input == ExitCommand || input == AltExitCommand {
		return fmt.Errorf("exit")
	}

	// Skip empty input
	if input == "" {
		return nil
	}

	// Initialize session on first message
	if is.session == nil {
		is.session = session.NewSession(is.config.Model, input)
	}

	// Build message history for API (limit to last N messages)
	apiMessages := is.buildAPIMessages(input)

	// Make API request
	req := &perplexity.ChatCompletionRequest{
		Model:           is.config.Model,
		Messages:        apiMessages,
		MaxTokens:       is.config.MaxTokens,
		Temperature:     is.config.Temperature,
		TopP:            is.config.TopP,
		SearchMode:      is.config.SearchMode,
		ReasoningEffort: is.config.ReasoningEffort,
	}

	resp, err := is.client.CreateCompletionWithRequest(req)
	if err != nil {
		return fmt.Errorf("API request failed: %w", err)
	}

	// Parse response
	parsed := perplexity.ParseResponse(resp)

	// Add messages to session
	is.session.AddMessage("user", input)

	// Add assistant response (without references for clean context)
	cleanContent := perplexity.StripReferences(parsed.Content)
	is.session.AddMessage("assistant", cleanContent)

	// Display response
	fmt.Println()
	ui.PrintSeparator(ui.Magenta)
	fmt.Print("PPLX: ")
	formatted := perplexity.FormatWithReferences(parsed)
	fmt.Println(formatted)
	fmt.Println()

	// Mark first message as complete
	is.firstMessage = false

	// Auto-save session
	if err := is.saveSession(); err != nil {
		session.Debugf("Failed to auto-save session: %v", err)
	}

	return nil
}

// buildAPIMessages builds the message array for API request
// It strips references from previous assistant messages before sending to API
func (is *InteractiveSession) buildAPIMessages(newInput string) []perplexity.Message {
	if is.session == nil {
		return []perplexity.Message{
			{Role: "user", Content: newInput},
		}
	}

	// Get last N messages for context
	messages := is.session.GetLastMessages(MaxContextMessages)
	apiMessages := make([]perplexity.Message, 0, len(messages)+1)

	// Convert session messages to perplexity messages, stripping references
	for _, msg := range messages {
		content := msg.Content

		// Strip references section from assistant messages
		if msg.Role == "assistant" {
			content = perplexity.StripReferences(content)
		}

		apiMessages = append(apiMessages, perplexity.Message{
			Role:    msg.Role,
			Content: content,
		})
	}

	// Add the new user message if not already in history
	if len(messages) == 0 || messages[len(messages)-1].Content != newInput {
		apiMessages = append(apiMessages, perplexity.Message{
			Role:    "user",
			Content: newInput,
		})
	}

	return apiMessages
}

// saveSession saves the current session
func (is *InteractiveSession) saveSession() error {
	if is.session == nil {
		return nil
	}

	if err := is.sessionManager.Save(is.session); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to save session: %v\n", err)
		return err
	}

	session.Debugf("Session auto-saved: %s", is.session.ID)
	return nil
}

// runInteractive is the entry point called from root.go
func runInteractive() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "Configuration error: %v\n", err)
		fmt.Fprintln(os.Stderr, "Please set PPLX_API_KEY environment variable or add it to ~/.pplx/config.yaml")
		os.Exit(1)
	}

	// Create and run interactive session
	interactive, err := NewInteractiveSession(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating interactive session: %v\n", err)
		os.Exit(1)
	}

	if err := interactive.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
