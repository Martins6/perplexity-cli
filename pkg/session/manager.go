package session

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Manager handles session CRUD operations
type Manager struct {
	sessionsDir string
}

// NewManager creates a new session manager
func NewManager() (*Manager, error) {
	sessionsDir := GetSessionsDir()
	if err := os.MkdirAll(sessionsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create sessions directory: %w", err)
	}
	return &Manager{sessionsDir: sessionsDir}, nil
}

// NewManagerWithDir creates a manager with a custom directory (for testing)
func NewManagerWithDir(dir string) *Manager {
	return &Manager{sessionsDir: dir}
}

// Save saves a session to disk atomically (write to temp then rename)
func (m *Manager) Save(session *Session) error {
	filename := filepath.Join(m.sessionsDir, session.ID+".json")

	// Ensure directory exists
	if err := os.MkdirAll(m.sessionsDir, 0755); err != nil {
		return fmt.Errorf("failed to create sessions directory: %w", err)
	}

	// Marshal session to JSON
	data, err := json.MarshalIndent(session, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	// Write to temporary file
	tempFile := filename + ".tmp"
	if err := os.WriteFile(tempFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write temp file: %w", err)
	}

	// Atomic rename
	if err := os.Rename(tempFile, filename); err != nil {
		os.Remove(tempFile) // Clean up temp file
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	Debugf("Saved session: %s", filename)
	return nil
}

// Load loads a session from disk by ID
func (m *Manager) Load(id string) (*Session, error) {
	filename := filepath.Join(m.sessionsDir, id+".json")
	return m.LoadFromFile(filename)
}

// LoadByShortID loads a session by its short ID
func (m *Manager) LoadByShortID(shortID string) (*Session, error) {
	entries, err := os.ReadDir(m.sessionsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read sessions directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filename := entry.Name()
		if !IsValidSessionFile(filename) {
			continue
		}

		filepath := filepath.Join(m.sessionsDir, filename)
		session, err := m.LoadFromFile(filepath)
		if err != nil {
			continue
		}

		if session.ShortID == shortID {
			return session, nil
		}
	}

	return nil, fmt.Errorf("session with short ID %s not found", shortID)
}

// LoadFromFile loads a session from a specific file path
func (m *Manager) LoadFromFile(filename string) (*Session, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read session file: %w", err)
	}

	var session Session
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %w", err)
	}

	// Handle old sessions without ShortID
	if session.ShortID == "" {
		// Generate ShortID from the session's creation time
		session.ShortID = GenerateShortID(session.Metadata.CreatedAt)

		// Save the session with the new ShortID
		if err := m.Save(&session); err != nil {
			// Log warning but continue - session will work without saving
			Debugf("Failed to save migrated session %s: %v", session.ID, err)
		} else {
			Debugf("Migrated session %s with ShortID %s", session.ID, session.ShortID)
		}
	}

	return &session, nil
}

// Delete deletes a session by ID
func (m *Manager) Delete(id string) error {
	filename := GenerateSessionFilename(id)
	if err := os.Remove(filename); err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}
	return nil
}

// List returns a list of all sessions sorted by time (newest first)
func (m *Manager) List() ([]SessionInfo, error) {
	entries, err := os.ReadDir(m.sessionsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read sessions directory: %w", err)
	}

	var sessions []SessionInfo
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filename := entry.Name()
		if !IsValidSessionFile(filename) {
			continue
		}

		// Load session to get info
		filepath := filepath.Join(m.sessionsDir, filename)
		session, err := m.LoadFromFile(filepath)
		if err != nil {
			Debugf("Failed to load session %s: %v", filename, err)
			continue
		}

		sessions = append(sessions, session.ToInfo())
	}

	// Sort by created time, newest first
	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].CreatedAt.After(sessions[j].CreatedAt)
	})

	return sessions, nil
}

// ListRecent returns the n most recent sessions
func (m *Manager) ListRecent(n int) ([]SessionInfo, error) {
	sessions, err := m.List()
	if err != nil {
		return nil, err
	}

	if len(sessions) > n {
		return sessions[:n], nil
	}
	return sessions, nil
}

// Search searches for sessions containing the query string
func (m *Manager) Search(query string) ([]SessionInfo, error) {
	queryLower := strings.ToLower(query)
	sessions, err := m.List()
	if err != nil {
		return nil, err
	}

	var results []SessionInfo
	for _, info := range sessions {
		// Load full session to search content
		session, err := m.Load(info.ID)
		if err != nil {
			continue
		}

		found := false

		// Search by short ID (exact match, case-insensitive)
		if strings.EqualFold(session.ShortID, query) {
			found = true
		}

		// Search in initial query
		if !found && strings.Contains(strings.ToLower(session.Metadata.InitialQuery), queryLower) {
			found = true
		}

		// Search in all messages
		if !found {
			for _, msg := range session.Messages {
				if strings.Contains(strings.ToLower(msg.Content), queryLower) {
					found = true
					break
				}
			}
		}

		if found {
			results = append(results, info)
		}
	}

	return results, nil
}

// GetSessionDir returns the sessions directory path
func (m *Manager) GetSessionDir() string {
	return m.sessionsDir
}

// CreateAndSave creates a new session and saves it
func (m *Manager) CreateAndSave(model, initialQuery string) (*Session, error) {
	session := NewSession(model, initialQuery)
	if err := m.Save(session); err != nil {
		return nil, err
	}
	return session, nil
}

// Update appends a message to a session and saves it
func (m *Manager) Update(id, role, content string) error {
	session, err := m.Load(id)
	if err != nil {
		return fmt.Errorf("failed to load session for update: %w", err)
	}

	session.AddMessage(role, content)
	return m.Save(session)
}

// GetSessionFilename returns the full path for a session file
func (m *Manager) GetSessionFilename(id string) string {
	return filepath.Join(m.sessionsDir, id+".json")
}

// SessionExists checks if a session exists
func (m *Manager) SessionExists(id string) bool {
	filename := GenerateSessionFilename(id)
	_, err := os.Stat(filename)
	return err == nil
}

// GetLatestSession returns the most recent session
func (m *Manager) GetLatestSession() (*Session, error) {
	sessions, err := m.List()
	if err != nil {
		return nil, err
	}

	if len(sessions) == 0 {
		return nil, fmt.Errorf("no sessions found")
	}

	return m.Load(sessions[0].ID)
}

// CreateSessionFromPerplexityMessages creates a session from perplexity messages
func (m *Manager) CreateSessionFromPerplexityMessages(model string, messages []struct {
	Role    string
	Content string
}) (*Session, error) {
	var initialQuery string
	for _, msg := range messages {
		if msg.Role == "user" {
			initialQuery = msg.Content
			break
		}
	}

	session := NewSession(model, initialQuery)
	for _, msg := range messages {
		session.AddMessage(msg.Role, msg.Content)
	}

	if err := m.Save(session); err != nil {
		return nil, err
	}

	return session, nil
}

// GetStats returns statistics about sessions
func (m *Manager) GetStats() (total int, oldest, newest time.Time, err error) {
	sessions, err := m.List()
	if err != nil {
		return 0, time.Time{}, time.Time{}, err
	}

	if len(sessions) == 0 {
		return 0, time.Time{}, time.Time{}, nil
	}

	total = len(sessions)
	oldest = sessions[len(sessions)-1].CreatedAt
	newest = sessions[0].CreatedAt

	return total, oldest, newest, nil
}
