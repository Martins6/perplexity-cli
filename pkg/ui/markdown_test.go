package ui

import (
	"testing"

	"github.com/fatih/color"
	"perplexity-cli/pkg/config"
)

func TestRenderMarkdown(t *testing.T) {
	tests := []struct {
		name    string
		content string
		cfg     *config.Config
		wantErr bool
		setup   func() (restore func())
	}{
		{
			name:    "rendering enabled with terminal",
			content: "# Header\n\nThis is a paragraph with **bold** text.\n\n```go\nfmt.Println(\"Hello\")\n```",
			cfg: &config.Config{
				UseGlow:   true,
				GlowStyle: "auto",
				GlowWidth: 80,
			},
			wantErr: false,
			setup: func() (restore func()) {
				oldNoColor := color.NoColor
				color.NoColor = false
				return func() { color.NoColor = oldNoColor }
			},
		},
		{
			name:    "rendering disabled",
			content: "# Header\n\nThis is a paragraph.",
			cfg: &config.Config{
				UseGlow:   false,
				GlowStyle: "auto",
				GlowWidth: 80,
			},
			wantErr: false,
			setup: func() (restore func()) {
				oldNoColor := color.NoColor
				color.NoColor = false
				return func() { color.NoColor = oldNoColor }
			},
		},
		{
			name:    "not a terminal",
			content: "# Header\n\nThis is a paragraph.",
			cfg: &config.Config{
				UseGlow:   true,
				GlowStyle: "auto",
				GlowWidth: 80,
			},
			wantErr: false,
			setup: func() (restore func()) {
				oldNoColor := color.NoColor
				color.NoColor = true
				return func() { color.NoColor = oldNoColor }
			},
		},
		{
			name:    "dark style",
			content: "# Header\n\nThis is a paragraph with *italic* text.",
			cfg: &config.Config{
				UseGlow:   true,
				GlowStyle: "dark",
				GlowWidth: 0,
			},
			wantErr: false,
			setup: func() (restore func()) {
				oldNoColor := color.NoColor
				color.NoColor = false
				return func() { color.NoColor = oldNoColor }
			},
		},
		{
			name:    "light style",
			content: "# Header\n\nThis is a paragraph with *italic* text.",
			cfg: &config.Config{
				UseGlow:   true,
				GlowStyle: "light",
				GlowWidth: 0,
			},
			wantErr: false,
			setup: func() (restore func()) {
				oldNoColor := color.NoColor
				color.NoColor = false
				return func() { color.NoColor = oldNoColor }
			},
		},
		{
			name:    "custom word wrap width",
			content: "# Header\n\nThis is a paragraph with **bold** text that should wrap at 40 characters.",
			cfg: &config.Config{
				UseGlow:   true,
				GlowStyle: "auto",
				GlowWidth: 40,
			},
			wantErr: false,
			setup: func() (restore func()) {
				oldNoColor := color.NoColor
				color.NoColor = false
				return func() { color.NoColor = oldNoColor }
			},
		},
		{
			name:    "empty content",
			content: "",
			cfg: &config.Config{
				UseGlow:   true,
				GlowStyle: "auto",
				GlowWidth: 80,
			},
			wantErr: false,
			setup: func() (restore func()) {
				oldNoColor := color.NoColor
				color.NoColor = false
				return func() { color.NoColor = oldNoColor }
			},
		},
		{
			name:    "complex markdown with lists and code blocks",
			content: "# Features\n\n- Feature 1\n- Feature 2\n- Feature 3\n\n## Example Code\n\n```go\nfunc main() {\n    fmt.Println(\"Hello, World!\")\n}\n```\n\n> This is a blockquote",
			cfg: &config.Config{
				UseGlow:   true,
				GlowStyle: "auto",
				GlowWidth: 80,
			},
			wantErr: false,
			setup: func() (restore func()) {
				oldNoColor := color.NoColor
				color.NoColor = false
				return func() { color.NoColor = oldNoColor }
			},
		},
		{
			name:    "plain text without markdown",
			content: "This is just plain text without any markdown formatting.",
			cfg: &config.Config{
				UseGlow:   true,
				GlowStyle: "auto",
				GlowWidth: 80,
			},
			wantErr: false,
			setup: func() (restore func()) {
				oldNoColor := color.NoColor
				color.NoColor = false
				return func() { color.NoColor = oldNoColor }
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			restore := tt.setup()
			defer restore()

			result, err := RenderMarkdown(tt.content, tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("RenderMarkdown() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.cfg != nil && !tt.cfg.UseGlow || !IsTerminal() {
				if result != tt.content {
					t.Errorf("RenderMarkdown() should return original content when disabled or not terminal")
				}
			} else {
				if result == "" && tt.content != "" {
					t.Errorf("RenderMarkdown() returned empty result for non-empty content")
				}
			}
		})
	}
}

func TestRenderMarkdownFallback(t *testing.T) {
	t.Run("fallback on nil config", func(t *testing.T) {
		oldNoColor := color.NoColor
		defer func() { color.NoColor = oldNoColor }()
		color.NoColor = false

		content := "# Header\n\nContent"
		result, err := RenderMarkdown(content, nil)

		if err != nil {
			t.Errorf("Expected no error with nil config, got: %v", err)
		}
		if result == "" {
			t.Errorf("Expected result to be returned, got empty string")
		}
	})
}
