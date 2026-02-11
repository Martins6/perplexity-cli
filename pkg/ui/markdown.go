package ui

import (
	"perplexity-cli/pkg/config"

	"github.com/charmbracelet/glamour"
)

func RenderMarkdown(content string, cfg *config.Config) (string, error) {
	if cfg == nil || !cfg.UseGlow || !IsTerminal() {
		return content, nil
	}

	opts := []glamour.TermRendererOption{
		glamour.WithEmoji(),
		glamour.WithPreservedNewLines(),
	}

	if cfg.GlowStyle == "" || cfg.GlowStyle == "auto" {
		opts = append(opts, glamour.WithAutoStyle())
	} else {
		opts = append(opts, glamour.WithStandardStyle(cfg.GlowStyle))
	}

	if cfg.GlowWidth > 0 {
		opts = append(opts, glamour.WithWordWrap(cfg.GlowWidth))
	}

	renderer, err := glamour.NewTermRenderer(opts...)
	if err != nil {
		return content, nil
	}

	return renderer.Render(content)
}
