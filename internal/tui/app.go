package tui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ironystock/agentic-obs/config"
	"github.com/ironystock/agentic-obs/internal/docs"
	"github.com/ironystock/agentic-obs/internal/storage"
)

// ViewType represents the different views in the TUI
type ViewType int

const (
	ViewStatus ViewType = iota
	ViewConfig
	ViewHistory
	ViewDocs
)

// Number of views for tab cycling
const numViews = 4

// App represents the TUI application
type App struct {
	db         *storage.DB
	cfg        *config.Config
	appName    string
	appVersion string
}

// New creates a new TUI application
func New(db *storage.DB, cfg *config.Config, appName, appVersion string) *App {
	return &App{
		db:         db,
		cfg:        cfg,
		appName:    appName,
		appVersion: appVersion,
	}
}

// Run starts the TUI application
func (a *App) Run() error {
	model := newModel(a.db, a.cfg, a.appName, a.appVersion)
	p := tea.NewProgram(model, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

// Model represents the TUI model state
type Model struct {
	db         *storage.DB
	cfg        *config.Config
	appName    string
	appVersion string

	// View state
	currentView ViewType
	width       int
	height      int
	ready       bool

	// Status data
	startTime       time.Time
	lastRefresh     time.Time
	screenshotCount int
	historyCount    int

	// History data
	actions       []storage.ActionRecord
	historyOffset int

	// Docs data
	docsList       []docs.Doc
	selectedDocIdx int
	docsContent    string
	docsOffset     int

	// UI components
	spinner spinner.Model

	// Error state
	lastError error
}

// newModel creates a new Model
func newModel(db *storage.DB, cfg *config.Config, appName, appVersion string) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(colorAccent)

	return Model{
		db:          db,
		cfg:         cfg,
		appName:     appName,
		appVersion:  appVersion,
		currentView: ViewStatus,
		startTime:   time.Now(),
		spinner:     s,
	}
}

// Messages
type tickMsg time.Time
type statusUpdateMsg struct {
	screenshotCount int
	historyCount    int
}
type historyUpdateMsg struct {
	actions []storage.ActionRecord
}
type docsListMsg struct {
	docs []docs.Doc
}
type docsContentMsg struct {
	content string
}
type errMsg struct {
	err error
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		tickCmd(),
		fetchStatusCmd(m.db),
		fetchHistoryCmd(m.db, historyFetchLimit),
	)
}

// Update handles messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "1":
			m.currentView = ViewStatus
		case "2":
			m.currentView = ViewConfig
		case "3":
			m.currentView = ViewHistory
			return m, fetchHistoryCmd(m.db, historyFetchLimit)
		case "4":
			m.currentView = ViewDocs
			return m, fetchDocsListCmd()
		case "tab", "right":
			m.currentView = (m.currentView + 1) % numViews
			return m, m.fetchViewDataCmd()
		case "shift+tab", "left":
			m.currentView = (m.currentView + numViews - 1) % numViews
			return m, m.fetchViewDataCmd()
		case "r":
			// Refresh current view
			return m, tea.Batch(fetchStatusCmd(m.db), fetchHistoryCmd(m.db, historyFetchLimit))
		case "j", "down":
			if m.currentView == ViewHistory && m.historyOffset < len(m.actions)-scrollMargin {
				m.historyOffset++
			}
			if m.currentView == ViewDocs {
				if m.docsContent == "" && len(m.docsList) > 0 {
					// In list view: navigate selection
					if m.selectedDocIdx < len(m.docsList)-1 {
						m.selectedDocIdx++
					}
				} else if m.docsOffset < maxDocsScrollOffset(m.docsContent, m.height) {
					// In content view: scroll
					m.docsOffset++
				}
			}
		case "k", "up":
			if m.currentView == ViewHistory && m.historyOffset > 0 {
				m.historyOffset--
			}
			if m.currentView == ViewDocs {
				if m.docsContent == "" && len(m.docsList) > 0 {
					// In list view: navigate selection
					if m.selectedDocIdx > 0 {
						m.selectedDocIdx--
					}
				} else if m.docsOffset > 0 {
					// In content view: scroll
					m.docsOffset--
				}
			}
		case "pgdown", "ctrl+d", "d":
			// Page down - scroll by half page
			if m.currentView == ViewDocs && m.docsContent != "" {
				pageSize := (m.height - uiChromeHeight) / 2
				maxOffset := maxDocsScrollOffset(m.docsContent, m.height)
				m.docsOffset += pageSize
				if m.docsOffset > maxOffset {
					m.docsOffset = maxOffset
				}
			}
		case "pgup", "ctrl+u", "u":
			// Page up - scroll by half page
			if m.currentView == ViewDocs && m.docsContent != "" {
				pageSize := (m.height - uiChromeHeight) / 2
				m.docsOffset -= pageSize
				if m.docsOffset < 0 {
					m.docsOffset = 0
				}
			}
		case "home", "g":
			// Go to top
			if m.currentView == ViewDocs && m.docsContent != "" {
				m.docsOffset = 0
			}
		case "end", "G":
			// Go to bottom
			if m.currentView == ViewDocs && m.docsContent != "" {
				m.docsOffset = maxDocsScrollOffset(m.docsContent, m.height)
			}
		case "enter":
			if m.currentView == ViewDocs && m.docsContent == "" && len(m.docsList) > 0 {
				return m, fetchDocsContentCmd(m.docsList[m.selectedDocIdx].Name, m.width)
			}
		case "esc", "backspace":
			if m.currentView == ViewDocs && m.docsContent != "" {
				m.docsContent = ""
				m.docsOffset = 0
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true

	case tickMsg:
		m.lastRefresh = time.Time(msg)
		return m, tea.Batch(tickCmd(), fetchStatusCmd(m.db))

	case statusUpdateMsg:
		m.screenshotCount = msg.screenshotCount
		m.historyCount = msg.historyCount
		m.lastError = nil

	case historyUpdateMsg:
		m.actions = msg.actions
		m.lastError = nil

	case docsListMsg:
		m.docsList = msg.docs
		m.lastError = nil

	case docsContentMsg:
		m.docsContent = msg.content
		m.docsOffset = 0
		m.lastError = nil

	case errMsg:
		m.lastError = msg.err

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

// View renders the UI
func (m Model) View() string {
	if !m.ready {
		return m.spinner.View() + " Loading..."
	}

	// Build the UI
	var content string
	switch m.currentView {
	case ViewStatus:
		content = m.renderStatusView()
	case ViewConfig:
		content = m.renderConfigView()
	case ViewHistory:
		content = m.renderHistoryView()
	case ViewDocs:
		content = m.renderDocsView()
	}

	// Compose final view
	header := m.renderHeader()
	tabs := m.renderTabs()
	help := m.renderHelp()

	return lipgloss.JoinVertical(lipgloss.Left, header, tabs, content, help)
}

// renderHeader renders the header bar with app info and connection status
func (m Model) renderHeader() string {
	// App name and version
	appInfo := styleHeader.Render(fmt.Sprintf("%s v%s", m.appName, m.appVersion))

	// Connection status indicator
	status := styleSuccess.Render(StatusConnected + " Connected")

	// Current time
	timeStr := styleMuted.Render(time.Now().Format("15:04:05"))

	// Calculate spacing
	leftPart := appInfo
	rightPart := status + "  " + timeStr
	spacing := m.width - lipgloss.Width(leftPart) - lipgloss.Width(rightPart) - headerSpacing
	if spacing < 1 {
		spacing = 1
	}

	headerContent := leftPart + fmt.Sprintf("%*s", spacing, "") + rightPart

	return styleHeaderBox.Copy().Width(m.width - headerWidthOffset).Render(headerContent)
}

// renderTabs renders the tab bar with emoji icons
func (m Model) renderTabs() string {
	tabs := []struct {
		icon string
		name string
	}{
		{TabIconStatus, "Status"},
		{TabIconConfig, "Config"},
		{TabIconHistory, "History"},
		{TabIconDocs, "Docs"},
	}

	var rendered []string
	for i, tab := range tabs {
		label := tab.icon + " " + tab.name
		if ViewType(i) == m.currentView {
			rendered = append(rendered, styleTabActive.Render(label))
		} else {
			rendered = append(rendered, styleTabInactive.Render(label))
		}
		// Add separator between tabs (except after last)
		if i < len(tabs)-1 {
			rendered = append(rendered, styleTabSeparator.Render(" │ "))
		}
	}

	tabBar := lipgloss.JoinHorizontal(lipgloss.Top, rendered...)
	separator := styleDim.Render(fmt.Sprintf("%s", repeatChar("─", m.width)))

	return lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.NewStyle().Padding(0, 1).Render(tabBar),
		separator,
	)
}

// repeatChar repeats a character n times using strings.Repeat (O(n) complexity)
func repeatChar(char string, n int) string {
	if n <= 0 {
		return ""
	}
	return strings.Repeat(char, n)
}

// renderStatusView renders the status view with aligned key-value pairs
func (m Model) renderStatusView() string {
	box := styleBox.Copy().Width(m.width - boxWidthOffset)

	// Server info with aligned labels
	uptime := time.Since(m.startTime).Round(time.Second)
	serverInfo := lipgloss.JoinVertical(lipgloss.Left,
		styleTitle.Render("Server Info"),
		"",
		styleLabel.Render("Name")+styleValue.Render(m.appName),
		styleLabel.Render("Version")+styleValue.Render(m.appVersion),
		styleLabel.Render("Uptime")+styleValue.Render(uptime.String()),
		styleLabel.Render("Last Refresh")+styleValue.Render(m.lastRefresh.Format("15:04:05")),
	)

	// OBS info with status indicator
	obsStatusText := styleSuccess.Render(StatusConnected + " Connected")
	obsInfo := lipgloss.JoinVertical(lipgloss.Left,
		styleTitle.Render("OBS Connection"),
		"",
		styleLabel.Render("Status")+obsStatusText,
		styleLabel.Render("Address")+styleValue.Render(m.cfg.OBSHost+":"+m.cfg.OBSPort),
	)

	// Stats
	statsInfo := lipgloss.JoinVertical(lipgloss.Left,
		styleTitle.Render("Statistics"),
		"",
		styleLabel.Render("Screenshot Sources")+styleValue.Render(fmt.Sprintf("%d", m.screenshotCount)),
		styleLabel.Render("Action History")+styleValue.Render(fmt.Sprintf("%d records", m.historyCount)),
	)

	// Error display
	errorInfo := ""
	if m.lastError != nil {
		errorInfo = "\n\n" + styleError.Render("Error: "+m.lastError.Error())
	}

	content := lipgloss.JoinVertical(lipgloss.Left,
		box.Render(serverInfo),
		"",
		box.Render(obsInfo),
		"",
		box.Render(statsInfo),
		errorInfo,
	)

	return content
}

// renderConfigView renders the config view with consistent formatting
func (m Model) renderConfigView() string {
	box := styleBox.Copy().Width(m.width - boxWidthOffset)

	boolStr := func(b bool) string {
		if b {
			return styleSuccess.Render("Enabled")
		}
		return styleError.Render("Disabled")
	}

	// OBS config
	password := "(none)"
	if m.cfg.OBSPassword != "" {
		password = "****"
	}
	obsConfig := lipgloss.JoinVertical(lipgloss.Left,
		styleTitle.Render("OBS WebSocket"),
		"",
		styleLabel.Render("Host")+styleValue.Render(m.cfg.OBSHost),
		styleLabel.Render("Port")+styleValue.Render(m.cfg.OBSPort),
		styleLabel.Render("Password")+styleValue.Render(password),
	)

	// HTTP config
	httpConfig := lipgloss.JoinVertical(lipgloss.Left,
		styleTitle.Render("HTTP Server"),
		"",
		styleLabel.Render("Status")+boolStr(m.cfg.WebServer.Enabled),
		styleLabel.Render("Host")+styleValue.Render(m.cfg.WebServer.Host),
		styleLabel.Render("Port")+styleValue.Render(fmt.Sprintf("%d", m.cfg.WebServer.Port)),
	)

	// Tool groups
	toolGroups := lipgloss.JoinVertical(lipgloss.Left,
		styleTitle.Render("Tool Groups"),
		"",
		styleLabel.Render("Core")+boolStr(m.cfg.ToolGroups.Core),
		styleLabel.Render("Visual")+boolStr(m.cfg.ToolGroups.Visual),
		styleLabel.Render("Layout")+boolStr(m.cfg.ToolGroups.Layout),
		styleLabel.Render("Audio")+boolStr(m.cfg.ToolGroups.Audio),
		styleLabel.Render("Sources")+boolStr(m.cfg.ToolGroups.Sources),
		styleLabel.Render("Design")+boolStr(m.cfg.ToolGroups.Design),
	)

	return lipgloss.JoinVertical(lipgloss.Left,
		box.Render(obsConfig),
		"",
		box.Render(httpConfig),
		"",
		box.Render(toolGroups),
	)
}

// renderHistoryView renders the history view with dynamic columns
func (m Model) renderHistoryView() string {
	box := styleBox.Copy().Width(m.width - boxWidthOffset)

	if len(m.actions) == 0 {
		return box.Render(styleTitle.Render("Action History") + "\n\n" + styleMuted.Render("No actions recorded yet"))
	}

	// Calculate column widths dynamically based on terminal width
	availableWidth := m.width - tablePadding
	colTool := availableWidth - colWidthTimestamp - colWidthStatus - colWidthDuration - columnSpacing
	if colTool < colWidthToolMin {
		colTool = colWidthToolMin
	}

	// Helper to pad styled text to display width
	padStyled := func(s string, width int) string {
		displayWidth := lipgloss.Width(s)
		if displayWidth >= width {
			return s
		}
		return s + strings.Repeat(" ", width-displayWidth)
	}

	// Header - pad AFTER styling
	header := padStyled(styleTableHeader.Render("Timestamp"), colWidthTimestamp) + "  " +
		padStyled(styleTableHeader.Render("Tool"), colTool) + "  " +
		padStyled(styleTableHeader.Render("Status"), colWidthStatus) + "  " +
		styleTableHeader.Render("Duration")

	// Separator
	separator := styleDim.Render(repeatChar("─", availableWidth))

	// Rows
	var rows []string
	rows = append(rows, header)
	rows = append(rows, separator)

	// Calculate visible range
	maxVisible := m.height - uiChromeHeight
	if maxVisible < minVisibleRows {
		maxVisible = minVisibleRows
	}
	start := m.historyOffset
	end := start + maxVisible
	if end > len(m.actions) {
		end = len(m.actions)
	}

	for _, action := range m.actions[start:end] {
		// Format status with consistent width
		var status string
		if action.Success {
			status = styleSuccess.Render("OK  ")
		} else {
			status = styleError.Render("FAIL")
		}

		// Truncate tool name if too long
		toolName := action.ToolName
		if len(toolName) > colTool {
			toolName = toolName[:colTool-ellipsisLen] + "..."
		}

		// Format duration
		duration := fmt.Sprintf("%dms", action.DurationMs)

		// Build row with proper padding
		row := padStyled(styleDim.Render(action.CreatedAt.Format("2006-01-02 15:04:05")), colWidthTimestamp) + "  " +
			padStyled(toolName, colTool) + "  " +
			padStyled(status, colWidthStatus) + "  " +
			duration

		rows = append(rows, row)
	}

	// Scroll indicator
	scrollInfo := ""
	if len(m.actions) > maxVisible {
		scrollInfo = fmt.Sprintf("\n\n%s",
			styleMuted.Render(fmt.Sprintf("Showing %d-%d of %d (↑/↓ or j/k to scroll)", start+1, end, len(m.actions))))
	}

	content := styleTitle.Render("Action History") + "\n\n" + lipgloss.JoinVertical(lipgloss.Left, rows...) + scrollInfo
	return box.Render(content)
}

// renderDocsView renders the documentation view
func (m Model) renderDocsView() string {
	box := styleBox.Copy().Width(m.width - boxWidthOffset)

	// If viewing a specific document
	if m.docsContent != "" {
		lines := strings.Split(m.docsContent, "\n")
		maxVisible := m.height - uiChromeHeight
		if maxVisible < minVisibleRows {
			maxVisible = minVisibleRows
		}

		start := m.docsOffset
		end := start + maxVisible
		if end > len(lines) {
			end = len(lines)
		}
		if start > end {
			start = 0
		}

		visibleLines := lines[start:end]
		content := strings.Join(visibleLines, "\n")

		scrollInfo := ""
		if len(lines) > maxVisible {
			scrollInfo = fmt.Sprintf("\n\n%s",
				styleMuted.Render(fmt.Sprintf("Lines %d-%d of %d │ ↑/↓ j/k scroll │ d/u page │ g/G top/bottom │ Esc back", start+1, end, len(lines))))
		}

		return box.Render(content + scrollInfo)
	}

	// Show docs list
	if len(m.docsList) == 0 {
		return box.Render(styleTitle.Render("Documentation") + "\n\n" + styleMuted.Render("Loading..."))
	}

	var items []string
	items = append(items, styleTitle.Render("Documentation"))
	items = append(items, "")
	items = append(items, styleMuted.Render("Select a document and press Enter to view:"))
	items = append(items, "")

	for i, doc := range m.docsList {
		prefix := "  "
		if i == m.selectedDocIdx {
			prefix = styleAccent.Render("▶ ")
		}
		title := doc.Title
		if i == m.selectedDocIdx {
			title = styleAccent.Render(doc.Title)
		}
		desc := styleMuted.Render(" - " + doc.Description)
		items = append(items, prefix+title+desc)
	}

	items = append(items, "")
	items = append(items, styleMuted.Render("Use ↑/↓ to navigate, Enter to open, Esc to close"))

	return box.Render(lipgloss.JoinVertical(lipgloss.Left, items...))
}

// renderHelp renders the help bar with enhanced formatting
func (m Model) renderHelp() string {
	help := fmt.Sprintf("%s Tab • %s Refresh • %s Scroll • %s Quit",
		styleHelpKey.Render("[1-4]"),
		styleHelpKey.Render("[r]"),
		styleHelpKey.Render("[↑/↓]"),
		styleHelpKey.Render("[q]"),
	)

	return styleHelpText.Copy().Padding(0, 1).Render(help)
}

// Command functions
func tickCmd() tea.Cmd {
	return tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func fetchStatusCmd(db *storage.DB) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()

		// Get screenshot count
		sources, err := db.ListScreenshotSources(ctx)
		if err != nil {
			return errMsg{err}
		}

		// Get history count
		stats, err := db.GetActionStats(ctx)
		if err != nil {
			return errMsg{err}
		}

		historyCount := 0
		if total, ok := stats["total_actions"].(int64); ok {
			historyCount = int(total)
		}

		return statusUpdateMsg{
			screenshotCount: len(sources),
			historyCount:    historyCount,
		}
	}
}

func fetchHistoryCmd(db *storage.DB, limit int) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		actions, err := db.GetRecentActions(ctx, limit)
		if err != nil {
			return errMsg{err}
		}
		return historyUpdateMsg{actions: actions}
	}
}

func fetchDocsListCmd() tea.Cmd {
	return func() tea.Msg {
		docsList, err := docs.List()
		if err != nil {
			return errMsg{err}
		}
		return docsListMsg{docs: docsList}
	}
}

func fetchDocsContentCmd(name string, width int) tea.Cmd {
	return func() tea.Msg {
		content, err := docs.RenderTerminal(name, width-boxWidthOffset-4)
		if err != nil {
			return errMsg{err}
		}
		return docsContentMsg{content: content}
	}
}

// fetchViewDataCmd returns the appropriate fetch command for the current view
func (m Model) fetchViewDataCmd() tea.Cmd {
	switch m.currentView {
	case ViewHistory:
		return fetchHistoryCmd(m.db, historyFetchLimit)
	case ViewDocs:
		return fetchDocsListCmd()
	default:
		return nil
	}
}

// maxDocsScrollOffset calculates the maximum scroll offset for docs
func maxDocsScrollOffset(content string, height int) int {
	lines := strings.Count(content, "\n") + 1
	maxVisible := height - uiChromeHeight
	if maxVisible < minVisibleRows {
		maxVisible = minVisibleRows
	}
	offset := lines - maxVisible
	if offset < 0 {
		return 0
	}
	return offset
}

// styleAccent is used for highlighting in docs view
var styleAccent = lipgloss.NewStyle().Foreground(colorAccent)
