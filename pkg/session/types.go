package session

import (
	"time"

	"perplexity-cli/pkg/perplexity"
)

// SessionMessage represents a single message in a conversation
type SessionMessage struct {
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

// SessionMetadata contains metadata about the session
type SessionMetadata struct {
	Model        string    `json:"model"`
	InitialQuery string    `json:"initial_query"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Session represents a conversation session stored as JSON
type Session struct {
	ID       string           `json:"id"`
	Messages []SessionMessage `json:"messages"`
	Metadata SessionMetadata  `json:"metadata"`
}

// NewSession creates a new session with the given model and initial query
func NewSession(model, initialQuery string) *Session {
	now := time.Now()
	return &Session{
		ID:       generateSessionID(now),
		Messages: make([]SessionMessage, 0),
		Metadata: SessionMetadata{
			Model:        model,
			InitialQuery: initialQuery,
			CreatedAt:    now,
			UpdatedAt:    now,
		},
	}
}

// AddMessage adds a message to the session
func (s *Session) AddMessage(role, content string) {
	s.Messages = append(s.Messages, SessionMessage{
		Role:      role,
		Content:   content,
		Timestamp: time.Now(),
	})
	s.Metadata.UpdatedAt = time.Now()
}

// AddPerplexityMessages converts and adds perplexity messages
func (s *Session) AddPerplexityMessages(messages []perplexity.Message) {
	for _, msg := range messages {
		s.AddMessage(msg.Role, msg.Content)
	}
}

// ToPerplexityMessages converts session messages to perplexity format
func (s *Session) ToPerplexityMessages() []perplexity.Message {
	messages := make([]perplexity.Message, len(s.Messages))
	for i, msg := range s.Messages {
		messages[i] = perplexity.Message{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}
	return messages
}

// GetLastMessages returns the last n messages (for context window limiting)
func (s *Session) GetLastMessages(n int) []SessionMessage {
	if n >= len(s.Messages) {
		return s.Messages
	}
	return s.Messages[len(s.Messages)-n:]
}

// SessionInfo represents summary information for listing sessions
type SessionInfo struct {
	ID           string    `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	InitialQuery string    `json:"initial_query"`
	MessageCount int       `json:"message_count"`
}

// ToInfo converts a Session to SessionInfo
func (s *Session) ToInfo() SessionInfo {
	return SessionInfo{
		ID:           s.ID,
		CreatedAt:    s.Metadata.CreatedAt,
		InitialQuery: s.Metadata.InitialQuery,
		MessageCount: len(s.Messages),
	}
}
