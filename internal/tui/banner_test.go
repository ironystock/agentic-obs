package tui

import (
	"strings"
	"testing"
)

func TestBannerBox(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    []string
	}{
		{
			name:    "short version",
			version: "1.0.0",
			want:    []string{"agentic-obs", "1.0.0", "MCP"},
		},
		{
			name:    "dev version",
			version: "dev",
			want:    []string{"agentic-obs", "dev", "MCP"},
		},
		{
			name:    "long version truncated",
			version: "1.0.0-beta.1+build.12345",
			want:    []string{"agentic-obs", "1.0.0-beta.1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BannerBox(tt.version)
			for _, want := range tt.want {
				if !strings.Contains(got, want) {
					t.Errorf("BannerBox(%q) should contain %q, got:\n%s", tt.version, want, got)
				}
			}
		})
	}
}

func TestPadRight(t *testing.T) {
	tests := []struct {
		s      string
		length int
		want   string
	}{
		{"test", 8, "test    "},
		{"longstring", 4, "long"},
		{"exact", 5, "exact"},
		{"", 3, "   "},
	}

	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			got := padRight(tt.s, tt.length)
			if got != tt.want {
				t.Errorf("padRight(%q, %d) = %q, want %q", tt.s, tt.length, got, tt.want)
			}
		})
	}
}

func TestBannerConstants(t *testing.T) {
	// Verify banner constants are non-empty
	if BannerSmall == "" {
		t.Error("BannerSmall should not be empty")
	}
	if BannerMedium == "" {
		t.Error("BannerMedium should not be empty")
	}
	if BannerLarge == "" {
		t.Error("BannerLarge should not be empty")
	}

	// Verify they contain expected content
	if !strings.Contains(BannerSmall, "agentic-obs") {
		t.Error("BannerSmall should contain 'agentic-obs'")
	}
	if !strings.Contains(BannerMedium, "MCP") {
		t.Error("BannerMedium should contain 'MCP'")
	}
	if !strings.Contains(BannerLarge, "MCP") {
		t.Error("BannerLarge should contain 'MCP'")
	}
}
