package perplexity

import (
	"testing"
)

func TestExtractCitations(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected []Citation
	}{
		{
			name:     "No citations",
			content:  "The capital of France is Paris.",
			expected: []Citation{},
		},
		{
			name:    "Single citation",
			content: "The capital of France is Paris. [1]",
			expected: []Citation{
				{Number: 1, Index: 0},
			},
		},
		{
			name:    "Multiple citations",
			content: "Paris is the capital of France. [1] It has been the capital since 987. [2]",
			expected: []Citation{
				{Number: 1, Index: 0},
				{Number: 2, Index: 1},
			},
		},
		{
			name:    "Repeated citation",
			content: "Paris is beautiful. [1] Paris has great food. [1]",
			expected: []Citation{
				{Number: 1, Index: 0},
				{Number: 1, Index: 0},
			},
		},
		{
			name:    "Citation at beginning",
			content: "[1] Paris is the capital of France.",
			expected: []Citation{
				{Number: 1, Index: 0},
			},
		},
		{
			name:    "Large citation numbers",
			content: "This is supported by multiple studies. [10] [25] [100]",
			expected: []Citation{
				{Number: 10, Index: 9},
				{Number: 25, Index: 24},
				{Number: 100, Index: 99},
			},
		},
		{
			name:    "Invalid bracket format",
			content: "Paris is great [not a citation] and nice [1]",
			expected: []Citation{
				{Number: 1, Index: 0},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractCitations(tt.content)

			if len(result) != len(tt.expected) {
				t.Errorf("ExtractCitations() returned %d citations, expected %d", len(result), len(tt.expected))
				return
			}

			for i, citation := range result {
				if citation.Number != tt.expected[i].Number {
					t.Errorf("Citation %d: Number = %d, expected %d", i, citation.Number, tt.expected[i].Number)
				}
				if citation.Index != tt.expected[i].Index {
					t.Errorf("Citation %d: Index = %d, expected %d", i, citation.Index, tt.expected[i].Index)
				}
			}
		})
	}
}

func TestParseResponse(t *testing.T) {
	resp := &ChatCompletionResponse{
		ID:      "test-123",
		Model:   "sonar",
		Created: 1234567890,
		Object:  "chat.completion",
		Usage:   Usage{TotalTokens: 50},
		Choices: []Choice{
			{
				Index:        0,
				FinishReason: "stop",
				Message: Message{
					Role:    "assistant",
					Content: "The capital of France is Paris. [1] It is known for the Eiffel Tower. [2]",
				},
			},
		},
		SearchResults: []SearchResult{
			{Title: "Wikipedia - Paris", URL: "https://en.wikipedia.org/wiki/Paris"},
			{Title: "Eiffel Tower Official", URL: "https://www.toureiffel.paris/en"},
		},
	}

	parsed := ParseResponse(resp)

	if len(parsed.Citations) != 2 {
		t.Errorf("ParseResponse() returned %d citations, expected 2", len(parsed.Citations))
	}

	if len(parsed.SearchResults) != 2 {
		t.Errorf("ParseResponse() returned %d search results, expected 2", len(parsed.SearchResults))
	}

	if parsed.Content == "" {
		t.Error("ParseResponse() returned empty content")
	}
}

func TestFormatWithReferences(t *testing.T) {
	tests := []struct {
		name     string
		parsed   *ParsedResponse
		expected string
	}{
		{
			name: "With citations and references",
			parsed: &ParsedResponse{
				Content: "The capital of France is Paris. [1]",
				Citations: []Citation{
					{Number: 1, Index: 0},
				},
				SearchResults: []SearchResult{
					{Title: "Wikipedia", URL: "https://wikipedia.org/Paris"},
				},
			},
			expected: "The capital of France is Paris. [1]\n\n## References:\n[1] Wikipedia - https://wikipedia.org/Paris\n",
		},
		{
			name: "No citations",
			parsed: &ParsedResponse{
				Content:       "The capital of France is Paris.",
				Citations:     []Citation{},
				SearchResults: []SearchResult{},
			},
			expected: "The capital of France is Paris.",
		},
		{
			name: "No search results",
			parsed: &ParsedResponse{
				Content:       "The capital of France is Paris. [1]",
				Citations:     []Citation{{Number: 1, Index: 0}},
				SearchResults: []SearchResult{},
			},
			expected: "The capital of France is Paris. [1]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatWithReferences(tt.parsed)
			if result != tt.expected {
				t.Errorf("FormatWithReferences() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

func TestStripReferences(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name:     "No references section",
			content:  "The capital of France is Paris.",
			expected: "The capital of France is Paris.",
		},
		{
			name:     "With ## References section",
			content:  "The capital of France is Paris.\n\n## References:\n[1] Wikipedia - https://wikipedia.org/Paris",
			expected: "The capital of France is Paris.",
		},
		{
			name:     "With # References section",
			content:  "The capital of France is Paris.\n\n# References:\n[1] Wikipedia",
			expected: "The capital of France is Paris.",
		},
		{
			name:     "With References section (no header level)",
			content:  "The capital of France is Paris.\n\nReferences:\n[1] Wikipedia",
			expected: "The capital of France is Paris.",
		},
		{
			name:     "Multiple references",
			content:  "Paris is great. [1]\n\n## References:\n[1] Source\n[2] Another",
			expected: "Paris is great. [1]",
		},
		{
			name:     "Lowercase references header",
			content:  "Paris is great.\n\n## references:\n[1] Source",
			expected: "Paris is great.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StripReferences(tt.content)
			if result != tt.expected {
				t.Errorf("StripReferences() = %q, expected %q", result, tt.expected)
			}
		})
	}
}
