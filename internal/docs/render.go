package docs

import (
	"bytes"
	"fmt"

	"github.com/charmbracelet/glamour"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// htmlRenderer is a configured goldmark markdown renderer for HTML output.
var htmlRenderer = goldmark.New(
	goldmark.WithExtensions(extension.GFM), // GitHub Flavored Markdown
	goldmark.WithParserOptions(
		parser.WithAutoHeadingID(), // Auto-generate heading IDs for anchors
	),
	goldmark.WithRendererOptions(
		html.WithHardWraps(),
		html.WithXHTML(),
		html.WithUnsafe(), // Allow raw HTML in markdown
	),
)

// RenderHTML converts markdown content to HTML.
func RenderHTML(name string) (string, error) {
	markdown, err := Content(name)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := htmlRenderer.Convert([]byte(markdown), &buf); err != nil {
		return "", fmt.Errorf("failed to render markdown: %w", err)
	}

	return buf.String(), nil
}

// RenderTerminal converts markdown content to styled terminal output.
// width specifies the terminal width for word wrapping.
func RenderTerminal(name string, width int) (string, error) {
	markdown, err := Content(name)
	if err != nil {
		return "", err
	}

	if width <= 0 {
		width = 80
	}

	// Use dracula style which matches our dark theme
	renderer, err := glamour.NewTermRenderer(
		glamour.WithStylePath("dracula"),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create terminal renderer: %w", err)
	}

	out, err := renderer.Render(markdown)
	if err != nil {
		return "", fmt.Errorf("failed to render markdown for terminal: %w", err)
	}

	return out, nil
}

// RenderMarkdownToHTML converts arbitrary markdown string to HTML.
func RenderMarkdownToHTML(markdown string) (string, error) {
	var buf bytes.Buffer
	if err := htmlRenderer.Convert([]byte(markdown), &buf); err != nil {
		return "", fmt.Errorf("failed to render markdown: %w", err)
	}
	return buf.String(), nil
}

// RenderMarkdownToTerminal converts arbitrary markdown string to terminal output.
func RenderMarkdownToTerminal(markdown string, width int) (string, error) {
	if width <= 0 {
		width = 80
	}

	renderer, err := glamour.NewTermRenderer(
		glamour.WithStylePath("dracula"),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create terminal renderer: %w", err)
	}

	out, err := renderer.Render(markdown)
	if err != nil {
		return "", fmt.Errorf("failed to render markdown: %w", err)
	}

	return out, nil
}
