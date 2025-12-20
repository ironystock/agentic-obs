package docs

import (
	"strings"
	"testing"
)

func TestRenderHTML(t *testing.T) {
	html, err := RenderHTML("README")
	if err != nil {
		t.Fatalf("RenderHTML(README) error: %v", err)
	}

	if len(html) == 0 {
		t.Error("RenderHTML(README) returned empty content")
	}

	// Should contain HTML tags
	if !strings.Contains(html, "<") {
		t.Error("RenderHTML(README) should contain HTML tags")
	}
}

func TestRenderTerminal(t *testing.T) {
	output, err := RenderTerminal("README", 80)
	if err != nil {
		t.Fatalf("RenderTerminal(README, 80) error: %v", err)
	}

	if len(output) == 0 {
		t.Error("RenderTerminal(README, 80) returned empty content")
	}
}

func TestRenderMarkdownToHTML(t *testing.T) {
	md := "# Test\n\nThis is a **test**."
	html, err := RenderMarkdownToHTML(md)
	if err != nil {
		t.Fatalf("RenderMarkdownToHTML error: %v", err)
	}

	if !strings.Contains(html, "<h1") {
		t.Error("RenderMarkdownToHTML should render heading")
	}

	if !strings.Contains(html, "<strong>") {
		t.Error("RenderMarkdownToHTML should render bold text")
	}
}

func TestRenderMarkdownToTerminal(t *testing.T) {
	md := "# Test\n\nThis is a **test**."
	output, err := RenderMarkdownToTerminal(md, 80)
	if err != nil {
		t.Fatalf("RenderMarkdownToTerminal error: %v", err)
	}

	if len(output) == 0 {
		t.Error("RenderMarkdownToTerminal returned empty content")
	}
}

func TestRenderInvalidDoc(t *testing.T) {
	_, err := RenderHTML("NONEXISTENT")
	if err == nil {
		t.Error("RenderHTML(NONEXISTENT) should return error")
	}

	_, err = RenderTerminal("NONEXISTENT", 80)
	if err == nil {
		t.Error("RenderTerminal(NONEXISTENT) should return error")
	}
}
