package perplexity

import "time"

// Message represents a chat message in the conversation
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// SearchResult represents a web search result from Perplexity
type SearchResult struct {
	Title string `json:"title"`
	URL   string `json:"url"`
	Date  string `json:"date,omitempty"`
}

// Usage represents token usage information
type Usage struct {
	PromptTokens      int    `json:"prompt_tokens"`
	CompletionTokens  int    `json:"completion_tokens"`
	TotalTokens       int    `json:"total_tokens"`
	SearchContextSize string `json:"search_context_size,omitempty"`
	CitationTokens    int    `json:"citation_tokens,omitempty"`
	NumSearchQueries  int    `json:"num_search_queries,omitempty"`
	ReasoningTokens   int    `json:"reasoning_tokens,omitempty"`
}

// Choice represents a completion choice from the API
type Choice struct {
	Index        int     `json:"index"`
	FinishReason string  `json:"finish_reason"`
	Message      Message `json:"message"`
}

// ChatCompletionRequest represents the request body for chat completions
type ChatCompletionRequest struct {
	Model                  string    `json:"model"`
	Messages               []Message `json:"messages"`
	MaxTokens              int       `json:"max_tokens,omitempty"`
	Temperature            float64   `json:"temperature,omitempty"`
	TopP                   float64   `json:"top_p,omitempty"`
	SearchMode             string    `json:"search_mode,omitempty"`
	ReasoningEffort        string    `json:"reasoning_effort,omitempty"`
	Stream                 bool      `json:"stream,omitempty"`
	ReturnImages           bool      `json:"return_images,omitempty"`
	ReturnRelatedQuestions bool      `json:"return_related_questions,omitempty"`
}

// ChatCompletionResponse represents the response from the API
type ChatCompletionResponse struct {
	ID            string         `json:"id"`
	Model         string         `json:"model"`
	Created       int64          `json:"created"`
	Object        string         `json:"object"`
	Usage         Usage          `json:"usage"`
	Choices       []Choice       `json:"choices"`
	SearchResults []SearchResult `json:"search_results,omitempty"`
}

// ClientConfig holds configuration for the API client
type ClientConfig struct {
	APIKey     string
	Model      string
	Timeout    time.Duration
	MaxRetries int
}

// DefaultConfig returns a default configuration
func DefaultConfig(apiKey string) *ClientConfig {
	return &ClientConfig{
		APIKey:     apiKey,
		Model:      "sonar",
		Timeout:    30 * time.Second,
		MaxRetries: 3,
	}
}

// Citation represents a parsed citation with its reference
type Citation struct {
	Number int
	Index  int // 0-based index into SearchResults
}

// ParsedResponse holds the parsed content with citation information
type ParsedResponse struct {
	Content       string
	Citations     []Citation
	SearchResults []SearchResult
}
