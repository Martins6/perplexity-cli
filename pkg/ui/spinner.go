package ui

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Spinner represents a loading animation
type Spinner struct {
	frames []string
	index  int
	active bool
	delay  time.Duration
	mu     sync.Mutex
	stop   chan bool
	text   string
}

// NewSpinner creates a new spinner with default frames
func NewSpinner(text string) *Spinner {
	return &Spinner{
		frames: []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		delay:  100 * time.Millisecond,
		text:   text,
		stop:   make(chan bool),
	}
}

// Start begins the spinner animation
func (s *Spinner) Start() {
	s.mu.Lock()
	if s.active {
		s.mu.Unlock()
		return
	}
	s.active = true
	s.mu.Unlock()

	// Set up signal handling to stop spinner on interrupt
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		ticker := time.NewTicker(s.delay)
		defer ticker.Stop()

		for {
			select {
			case <-s.stop:
				s.clear()
				return
			case <-sigChan:
				s.Stop()
				return
			case <-ticker.C:
				s.mu.Lock()
				if !s.active {
					s.mu.Unlock()
					return
				}
				frame := s.frames[s.index]
				s.index = (s.index + 1) % len(s.frames)
				s.mu.Unlock()

				// Print spinner frame
				fmt.Printf("\r%s %s", Cyan(frame), s.text)
			}
		}
	}()
}

// Stop stops the spinner animation
func (s *Spinner) Stop() {
	s.mu.Lock()
	if !s.active {
		s.mu.Unlock()
		return
	}
	s.active = false
	s.mu.Unlock()

	close(s.stop)
	time.Sleep(10 * time.Millisecond) // Small delay to ensure clear happens
}

// clear removes the spinner from the terminal
func (s *Spinner) clear() {
	fmt.Printf("\r%s\r", "                                                                                                                                                                                                                                                                ")
}

// WithDelay sets a custom delay between frames
func (s *Spinner) WithDelay(delay time.Duration) *Spinner {
	s.delay = delay
	return s
}

// WithFrames sets custom spinner frames
func (s *Spinner) WithFrames(frames []string) *Spinner {
	s.frames = frames
	return s
}

// SimpleProgress just prints a progress message
func SimpleProgress(text string) func() {
	fmt.Printf("%s %s...", Cyan("→"), text)
	return func() {
		fmt.Println(" " + Green("✓"))
	}
}
