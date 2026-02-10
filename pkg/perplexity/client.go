package perplexity

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const DefaultAPIEndpoint = "https://api.perplexity.ai/chat/completions"

// Client represents a Perplexity API client
type Client struct {
	config     *ClientConfig
	httpClient *http.Client
	endpoint   string
}

// NewClient creates a new Perplexity API client
func NewClient(apiKey string) *Client {
	return NewClientWithConfig(DefaultConfig(apiKey))
}

// NewClientWithConfig creates a client with custom configuration
func NewClientWithConfig(config *ClientConfig) *Client {
	return &Client{
		config:     config,
		httpClient: &http.Client{Timeout: config.Timeout},
		endpoint:   DefaultAPIEndpoint,
	}
}

// SetEndpoint allows changing the API endpoint (useful for testing)
func (c *Client) SetEndpoint(endpoint string) {
	c.endpoint = endpoint
}

// CreateCompletion sends a chat completion request to the Perplexity API
func (c *Client) CreateCompletion(messages []Message) (*ChatCompletionResponse, error) {
	if c.config.APIKey == "" {
		return nil, fmt.Errorf("API key is required")
	}

	req := ChatCompletionRequest{
		Model:    c.config.Model,
		Messages: messages,
	}

	return c.CreateCompletionWithRequest(&req)
}

// CreateCompletionWithRequest sends a chat completion request with full configuration
func (c *Client) CreateCompletionWithRequest(req *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	if c.config.APIKey == "" {
		return nil, fmt.Errorf("API key is required. Set PPLX_API_KEY environment variable.")
	}

	// Set default model if not specified
	if req.Model == "" {
		req.Model = c.config.Model
	}

	// Marshal request body
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequest("POST", c.endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.config.APIKey)

	// Make request with retries
	var httpResp *http.Response
	for attempt := 0; attempt <= c.config.MaxRetries; attempt++ {
		httpResp, err = c.httpClient.Do(httpReq)
		if err == nil {
			break
		}
		if attempt < c.config.MaxRetries {
			time.Sleep(time.Duration(attempt+1) * time.Second)
		}
	}
	if err != nil {
		return nil, fmt.Errorf("failed to make request after %d attempts: %w", c.config.MaxRetries+1, err)
	}
	defer httpResp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for HTTP errors
	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", httpResp.StatusCode, string(respBody))
	}

	// Parse response
	var resp ChatCompletionResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &resp, nil
}

// Ask sends a simple query and returns the formatted response with references
func (c *Client) Ask(query string) (string, error) {
	messages := []Message{
		{Role: "user", Content: query},
	}

	resp, err := c.CreateCompletion(messages)
	if err != nil {
		return "", err
	}

	parsed := ParseResponse(resp)
	return FormatWithReferences(parsed), nil
}

// SetModel changes the default model for the client
func (c *Client) SetModel(model string) {
	c.config.Model = model
}

// GetModel returns the current default model
func (c *Client) GetModel() string {
	return c.config.Model
}
