package tui

// Banner contains ASCII art for the agentic-obs application.
// Used in TUI header and CLI output.

// BannerSmall is a compact single-line banner
const BannerSmall = `agentic-obs`

// BannerMedium is a stylized text banner for terminal display
const BannerMedium = `
   ▲
  ╱ ╲   agentic-obs
 ●───●  OBS Studio Control via MCP
`

// BannerLarge is a full ASCII art banner for splash screens
const BannerLarge = `
       ▲
      ╱ ╲
     ╱   ╲
    ●─────●

  ┌─────────────────────────────────┐
  │  agentic-obs                    │
  │  OBS Studio Control via MCP     │
  └─────────────────────────────────┘
`

// BannerBox returns a boxed banner with version info
func BannerBox(version string) string {
	return `
┌─────────────────────────────────────┐
│  agentic-obs ` + padRight(version, 22) + ` │
│  OBS Studio Control via MCP         │
└─────────────────────────────────────┘`
}

// padRight pads a string to the specified length
func padRight(s string, length int) string {
	if len(s) >= length {
		return s[:length]
	}
	return s + spaces(length-len(s))
}

// spaces returns a string of n spaces
func spaces(n int) string {
	if n <= 0 {
		return ""
	}
	b := make([]byte, n)
	for i := range b {
		b[i] = ' '
	}
	return string(b)
}
