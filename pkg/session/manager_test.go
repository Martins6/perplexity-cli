package session

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestNewSession(t *testing.T) {
	session := NewSession("sonar", "What is 2+2?")

	if session.ID == "" {
		t.Error("NewSession() created session with empty ID")
	}

	if session.Metadata.Model != "sonar" {
		t.Errorf("NewSession() model = %s, expected sonar", session.Metadata.Model)
	}

	if session.Metadata.InitialQuery != "What is 2+2?" {
		t.Errorf("NewSession() initialQuery = %s, expected 'What is 2+2?'", session.Metadata.InitialQuery)
	}

	if len(session.Messages) != 0 {
		t.Errorf("NewSession() should have 0 messages, got %d", len(session.Messages))
	}

	if session.Metadata.CreatedAt.IsZero() {
		t.Error("NewSession() created_at is zero")
	}
}

func TestSessionAddMessage(t *testing.T) {
	session := NewSession("sonar", "Test")

	session.AddMessage("user", "Hello")
	if len(session.Messages) != 1 {
		t.Errorf("AddMessage() should have 1 message, got %d", len(session.Messages))
	}

	if session.Messages[0].Role != "user" {
		t.Errorf("AddMessage() role = %s, expected user", session.Messages[0].Role)
	}

	if session.Messages[0].Content != "Hello" {
		t.Errorf("AddMessage() content = %s, expected Hello", session.Messages[0].Content)
	}

	session.AddMessage("assistant", "Hi there")
	if len(session.Messages) != 2 {
		t.Errorf("AddMessage() should have 2 messages, got %d", len(session.Messages))
	}
}

func TestSessionGetLastMessages(t *testing.T) {
	session := NewSession("sonar", "Test")

	// Add 5 messages
	for i := 0; i < 5; i++ {
		session.AddMessage("user", string(rune('A'+i)))
	}

	// Get last 3
	last3 := session.GetLastMessages(3)
	if len(last3) != 3 {
		t.Errorf("GetLastMessages(3) returned %d messages, expected 3", len(last3))
	}

	// Get last 10 (more than available)
	last10 := session.GetLastMessages(10)
	if len(last10) != 5 {
		t.Errorf("GetLastMessages(10) returned %d messages, expected 5", len(last10))
	}
}

func TestGenerateSessionID(t *testing.T) {
	now := time.Date(2024, 1, 15, 10, 30, 45, 123456789, time.UTC)
	id := generateSessionID(now)

	expected := "20240115-103045.123"
	if id != expected {
		t.Errorf("generateSessionID() = %s, expected %s", id, expected)
	}
}

func TestManagerSaveAndLoad(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "session-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	manager := NewManagerWithDir(tempDir)

	// Create and save session
	session := NewSession("sonar", "Test query")
	session.AddMessage("user", "Hello")
	session.AddMessage("assistant", "Hi!")

	if err := manager.Save(session); err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	// Load session
	loaded, err := manager.Load(session.ID)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if loaded.ID != session.ID {
		t.Errorf("Loaded session ID = %s, expected %s", loaded.ID, session.ID)
	}

	if loaded.Metadata.Model != "sonar" {
		t.Errorf("Loaded session model = %s, expected sonar", loaded.Metadata.Model)
	}

	if len(loaded.Messages) != 2 {
		t.Errorf("Loaded session has %d messages, expected 2", len(loaded.Messages))
	}
}

func TestManagerList(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "session-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	manager := NewManagerWithDir(tempDir)

	// Create multiple sessions
	for i := 0; i < 3; i++ {
		session := NewSession("sonar", fmt.Sprintf("Query %d", i))
		time.Sleep(10 * time.Millisecond) // Ensure different timestamps
		if err := manager.Save(session); err != nil {
			t.Fatalf("Save() failed: %v", err)
		}
	}

	// List sessions
	sessions, err := manager.List()
	if err != nil {
		t.Fatalf("List() failed: %v", err)
	}

	if len(sessions) != 3 {
		t.Errorf("List() returned %d sessions, expected 3", len(sessions))
	}

	// Check order (newest first)
	for i := 0; i < len(sessions)-1; i++ {
		if sessions[i].CreatedAt.Before(sessions[i+1].CreatedAt) {
			t.Error("List() sessions not sorted by date (newest first)")
		}
	}
}

func TestManagerSearch(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "session-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	manager := NewManagerWithDir(tempDir)

	// Create sessions with different queries
	session1 := NewSession("sonar", "What is Paris?")
	session1.AddMessage("user", "Tell me about France")
	manager.Save(session1)

	session2 := NewSession("sonar", "What is London?")
	session2.AddMessage("user", "Tell me about England")
	manager.Save(session2)

	time.Sleep(10 * time.Millisecond)

	session3 := NewSession("sonar", "What is Tokyo?")
	session3.AddMessage("user", "Tell me about Japan")
	manager.Save(session3)

	// Search for "France"
	results, err := manager.Search("France")
	if err != nil {
		t.Fatalf("Search() failed: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Search('France') returned %d results, expected 1", len(results))
	}

	// Search for "is" (case insensitive)
	results, err = manager.Search("is")
	if err != nil {
		t.Fatalf("Search() failed: %v", err)
	}

	if len(results) != 3 {
		t.Errorf("Search('is') returned %d results, expected 3", len(results))
	}
}

func TestTruncateQuery(t *testing.T) {
	tests := []struct {
		query    string
		maxLen   int
		expected string
	}{
		{"Short query", 100, "Short query"},
		{"This is a very long query that needs truncation", 20, "This is a very lo..."},
		{"Exact", 5, "Exact"},
		{"Too short max", 2, "To"},
	}

	for _, tt := range tests {
		result := TruncateQuery(tt.query, tt.maxLen)
		if result != tt.expected {
			t.Errorf("TruncateQuery(%q, %d) = %q, expected %q", tt.query, tt.maxLen, result, tt.expected)
		}
	}
}

func TestIsValidSessionFile(t *testing.T) {
	tests := []struct {
		filename string
		expected bool
	}{
		{"20240115-103045.123.json", true},
		{"20240115-103045.json", true},
		{"session.json", true}, // Changed: we now accept any .json file
		{"invalid", false},
		{"test.txt", false},
		{".json", false},
		{"20240115-103045.123", false},
	}

	for _, tt := range tests {
		result := IsValidSessionFile(tt.filename)
		if result != tt.expected {
			t.Errorf("IsValidSessionFile(%q) = %v, expected %v", tt.filename, result, tt.expected)
		}
	}
}

func TestParseSessionID(t *testing.T) {
	tests := []struct {
		filename string
		expected string
	}{
		{"/path/to/20240115-103045.123.json", "20240115-103045.123"},
		{"20240115-103045.123.json", "20240115-103045.123"},
		{"20240115-103045.json", "20240115-103045"},
		{"session.json", "session"},
		{"session", "session"},
	}

	for _, tt := range tests {
		result := ParseSessionID(tt.filename)
		if result != tt.expected {
			t.Errorf("ParseSessionID(%q) = %q, expected %q", tt.filename, result, tt.expected)
		}
	}
}
