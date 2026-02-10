package perplexity

import (
	"fmt"
	"regexp"
	"strings"
)

var citationRegex = regexp.MustCompile(`\[(\d+)\]`)

// ExtractCitations finds all citations in the content and returns them with their positions
func ExtractCitations(content string) []Citation {
	matches := citationRegex.FindAllStringSubmatchIndex(content, -1)
	citations := make([]Citation, 0, len(matches))

	for _, match := range matches {
		if len(match) >= 4 {
			// match[0], match[1] are the positions of the full match
			// match[2], match[3] are the positions of the first submatch (the number)
			numStr := content[match[2]:match[3]]
			var num int
			fmt.Sscanf(numStr, "%d", &num)

			// Convert to 0-based index for search_results array
			citations = append(citations, Citation{
				Number: num,
				Index:  num - 1,
			})
		}
	}

	return citations
}

// ParseResponse parses the API response and extracts citations
func ParseResponse(resp *ChatCompletionResponse) *ParsedResponse {
	if len(resp.Choices) == 0 {
		return &ParsedResponse{
			Content:       "",
			Citations:     []Citation{},
			SearchResults: resp.SearchResults,
		}
	}

	content := resp.Choices[0].Message.Content
	citations := ExtractCitations(content)

	return &ParsedResponse{
		Content:       content,
		Citations:     citations,
		SearchResults: resp.SearchResults,
	}
}

// FormatWithReferences formats the response content with a references section
func FormatWithReferences(parsed *ParsedResponse) string {
	if len(parsed.SearchResults) == 0 || len(parsed.Citations) == 0 {
		return parsed.Content
	}

	var sb strings.Builder
	sb.WriteString(parsed.Content)
	sb.WriteString("\n\n## References:\n")

	// Track which references we've already included
	seen := make(map[int]bool)
	refNum := 1

	for _, citation := range parsed.Citations {
		if citation.Index >= 0 && citation.Index < len(parsed.SearchResults) {
			if !seen[citation.Index] {
				seen[citation.Index] = true
				result := parsed.SearchResults[citation.Index]
				sb.WriteString(fmt.Sprintf("[%d] %s - %s\n", refNum, result.Title, result.URL))
				refNum++
			}
		}
	}

	return sb.String()
}

// StripReferences removes the references section from content before sending to API
// This prevents the model from receiving formatted references as context
func StripReferences(content string) string {
	// Find the references section and remove it
	refMarkers := []string{"\n## References:", "\n# References:", "\nReferences:", "\n##references:", "\n#references:"}

	lowerContent := strings.ToLower(content)
	for _, marker := range refMarkers {
		if idx := strings.Index(lowerContent, strings.ToLower(marker)); idx != -1 {
			return strings.TrimSpace(content[:idx])
		}
	}

	return content
}
