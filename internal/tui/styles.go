package tui

import "github.com/charmbracelet/lipgloss"

// Color palette - Pink/magenta accent theme
var (
	colorAccent    = lipgloss.Color("205") // Pink/magenta (primary)
	colorSuccess   = lipgloss.Color("82")  // Green
	colorError     = lipgloss.Color("196") // Red
	colorWarning   = lipgloss.Color("214") // Orange
	colorMuted     = lipgloss.Color("241") // Gray (help text)
	colorSubtle    = lipgloss.Color("245") // Light gray (labels)
	colorText      = lipgloss.Color("252") // Near white (values)
	colorBorder    = lipgloss.Color("240") // Border gray
	colorHighlight = lipgloss.Color("236") // Background highlight
	colorDim       = lipgloss.Color("239") // Dimmed text (timestamps)
)

// Shared styles
var (
	// Text styles
	styleTitle   = lipgloss.NewStyle().Bold(true).Foreground(colorAccent)
	styleLabel   = lipgloss.NewStyle().Foreground(colorSubtle).Width(20)
	styleValue   = lipgloss.NewStyle().Foreground(colorText)
	styleSuccess = lipgloss.NewStyle().Foreground(colorSuccess)
	styleError   = lipgloss.NewStyle().Foreground(colorError)
	styleWarning = lipgloss.NewStyle().Foreground(colorWarning)
	styleMuted   = lipgloss.NewStyle().Foreground(colorMuted)
	styleDim     = lipgloss.NewStyle().Foreground(colorDim)

	// Container styles
	styleBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorBorder).
			Padding(1, 2)

	// Tab styles
	styleTabActive = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorAccent).
			Background(colorHighlight).
			Padding(0, 2)
	styleTabInactive = lipgloss.NewStyle().
				Foreground(colorSubtle).
				Padding(0, 2)
	styleTabSeparator = lipgloss.NewStyle().
				Foreground(colorBorder)

	// Header styles
	styleHeader = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorAccent)
	styleHeaderBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorAccent).
			Padding(0, 1)

	// Table styles
	styleTableHeader = lipgloss.NewStyle().
				Bold(true).
				Foreground(colorSubtle)
	styleTableRow = lipgloss.NewStyle().
			Foreground(colorText)
	styleTableRowAlt = lipgloss.NewStyle().
				Foreground(colorText).
				Background(lipgloss.Color("235"))

	// Help bar styles
	styleHelpKey = lipgloss.NewStyle().
			Foreground(colorAccent)
	styleHelpText = lipgloss.NewStyle().
			Foreground(colorMuted)
)

// Status indicator characters
const (
	StatusConnected    = "●"
	StatusDisconnected = "○"
	StatusConnecting   = "◐"
)

// Tab icons
const (
	TabIconStatus  = "📊"
	TabIconConfig  = "⚙️"
	TabIconHistory = "📜"
	TabIconDocs    = "📖"
)

// Layout constants
const (
	// Box and container offsets
	boxWidthOffset    = 4  // Offset for box width from terminal width
	headerWidthOffset = 2  // Offset for header box width
	headerSpacing     = 6  // Spacing offset in header
	tablePadding      = 12 // Padding/border offset for table width
	columnSpacing     = 6  // Spacing between table columns

	// Table column widths
	colWidthTimestamp = 19 // "2006-01-02 15:04:05"
	colWidthStatus    = 6  // "OK" or "FAIL"
	colWidthDuration  = 10 // "12345ms"
	colWidthToolMin   = 15 // Minimum tool column width
	ellipsisLen       = 3  // Length of "..."

	// View dimensions
	uiChromeHeight = 18 // Height used by header, tabs, help bar
	minVisibleRows = 5  // Minimum rows to show in history
	scrollMargin   = 10 // Margin for scroll detection

	// Data limits
	historyFetchLimit = 50 // Number of history records to fetch
)
