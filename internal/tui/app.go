package tui

import (
	"context"
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ironystock/agentic-obs/config"
	"github.com/ironystock/agentic-obs/internal/storage"
)

// ViewType represents the different views in the TUI
type ViewType int

const (
	ViewStatus ViewType = iota
	ViewConfig
	ViewHistory
)

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

	// UI components
	spinner spinner.Model

	// Error state
	lastError error
}

// newModel creates a new Model
func newModel(db *storage.DB, cfg *config.Config, appName, appVersion string) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

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
type errMsg struct {
	err error
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		tickCmd(),
		fetchStatusCmd(m.db),
		fetchHistoryCmd(m.db, 50),
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
			return m, fetchHistoryCmd(m.db, 50)
		case "tab", "right":
			m.currentView = (m.currentView + 1) % 3
			if m.currentView == ViewHistory {
				return m, fetchHistoryCmd(m.db, 50)
			}
		case "shift+tab", "left":
			m.currentView = (m.currentView + 2) % 3
			if m.currentView == ViewHistory {
				return m, fetchHistoryCmd(m.db, 50)
			}
		case "r":
			// Refresh current view
			return m, tea.Batch(fetchStatusCmd(m.db), fetchHistoryCmd(m.db, 50))
		case "j", "down":
			if m.currentView == ViewHistory && m.historyOffset < len(m.actions)-10 {
				m.historyOffset++
			}
		case "k", "up":
			if m.currentView == ViewHistory && m.historyOffset > 0 {
				m.historyOffset--
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
	}

	// Compose final view
	tabs := m.renderTabs()
	help := m.renderHelp()

	return lipgloss.JoinVertical(lipgloss.Left, tabs, content, help)
}

// renderTabs renders the tab bar
func (m Model) renderTabs() string {
	tabStyle := lipgloss.NewStyle().Padding(0, 2)
	activeStyle := tabStyle.Copy().Bold(true).Foreground(lipgloss.Color("205")).Background(lipgloss.Color("236"))
	inactiveStyle := tabStyle.Copy().Foreground(lipgloss.Color("250"))

	tabs := []string{"Status", "Config", "History"}
	var rendered []string

	for i, tab := range tabs {
		if ViewType(i) == m.currentView {
			rendered = append(rendered, activeStyle.Render(tab))
		} else {
			rendered = append(rendered, inactiveStyle.Render(tab))
		}
	}

	tabBar := lipgloss.JoinHorizontal(lipgloss.Top, rendered...)
	return lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(lipgloss.Color("240")).
		Width(m.width).
		Render(tabBar)
}

// renderStatusView renders the status view
func (m Model) renderStatusView() string {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("250"))
	valueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	okStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("82"))
	errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(1, 2).
		Width(m.width - 4)

	// Server info
	uptime := time.Since(m.startTime).Round(time.Second)
	serverInfo := fmt.Sprintf("%s\n\n%s %s\n%s %s\n%s %s\n%s %s",
		titleStyle.Render("Server Info"),
		labelStyle.Render("Name:"),
		valueStyle.Render(m.appName),
		labelStyle.Render("Version:"),
		valueStyle.Render(m.appVersion),
		labelStyle.Render("Uptime:"),
		valueStyle.Render(uptime.String()),
		labelStyle.Render("Last Refresh:"),
		valueStyle.Render(m.lastRefresh.Format("15:04:05")),
	)

	// OBS info
	obsStatus := okStyle.Render("Configured")
	obsInfo := fmt.Sprintf("%s\n\n%s %s\n%s %s:%s",
		titleStyle.Render("OBS Connection"),
		labelStyle.Render("Status:"),
		obsStatus,
		labelStyle.Render("Address:"),
		valueStyle.Render(m.cfg.OBSHost),
		valueStyle.Render(m.cfg.OBSPort),
	)

	// Stats
	statsInfo := fmt.Sprintf("%s\n\n%s %s\n%s %s",
		titleStyle.Render("Statistics"),
		labelStyle.Render("Screenshot Sources:"),
		valueStyle.Render(fmt.Sprintf("%d", m.screenshotCount)),
		labelStyle.Render("Action History:"),
		valueStyle.Render(fmt.Sprintf("%d records", m.historyCount)),
	)

	// Error display
	errorInfo := ""
	if m.lastError != nil {
		errorInfo = fmt.Sprintf("\n\n%s %s",
			errorStyle.Render("Error:"),
			errorStyle.Render(m.lastError.Error()),
		)
	}

	content := lipgloss.JoinVertical(lipgloss.Left,
		boxStyle.Render(serverInfo),
		"",
		boxStyle.Render(obsInfo),
		"",
		boxStyle.Render(statsInfo),
		errorInfo,
	)

	return content
}

// renderConfigView renders the config view
func (m Model) renderConfigView() string {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("250"))
	enabledStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("82"))
	disabledStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(1, 2).
		Width(m.width - 4)

	boolStr := func(b bool) string {
		if b {
			return enabledStyle.Render("Enabled")
		}
		return disabledStyle.Render("Disabled")
	}

	// OBS config
	obsConfig := fmt.Sprintf("%s\n\n%s %s\n%s %s\n%s %s",
		titleStyle.Render("OBS WebSocket"),
		labelStyle.Render("Host:"),
		m.cfg.OBSHost,
		labelStyle.Render("Port:"),
		m.cfg.OBSPort,
		labelStyle.Render("Password:"),
		func() string {
			if m.cfg.OBSPassword != "" {
				return "****"
			}
			return "(none)"
		}(),
	)

	// HTTP config
	httpConfig := fmt.Sprintf("%s\n\n%s %s\n%s %s\n%s %d",
		titleStyle.Render("HTTP Server"),
		labelStyle.Render("Status:"),
		boolStr(m.cfg.WebServer.Enabled),
		labelStyle.Render("Host:"),
		m.cfg.WebServer.Host,
		labelStyle.Render("Port:"),
		m.cfg.WebServer.Port,
	)

	// Tool groups
	toolGroups := fmt.Sprintf("%s\n\n%s %s\n%s %s\n%s %s\n%s %s\n%s %s\n%s %s",
		titleStyle.Render("Tool Groups"),
		labelStyle.Render("Core:"),
		boolStr(m.cfg.ToolGroups.Core),
		labelStyle.Render("Visual:"),
		boolStr(m.cfg.ToolGroups.Visual),
		labelStyle.Render("Layout:"),
		boolStr(m.cfg.ToolGroups.Layout),
		labelStyle.Render("Audio:"),
		boolStr(m.cfg.ToolGroups.Audio),
		labelStyle.Render("Sources:"),
		boolStr(m.cfg.ToolGroups.Sources),
		labelStyle.Render("Design:"),
		boolStr(m.cfg.ToolGroups.Design),
	)

	return lipgloss.JoinVertical(lipgloss.Left,
		boxStyle.Render(obsConfig),
		"",
		boxStyle.Render(httpConfig),
		"",
		boxStyle.Render(toolGroups),
	)
}

// renderHistoryView renders the history view
func (m Model) renderHistoryView() string {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("250"))
	successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("82"))
	failStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("245"))

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(1, 2).
		Width(m.width - 4)

	if len(m.actions) == 0 {
		return boxStyle.Render(titleStyle.Render("Action History") + "\n\n" + dimStyle.Render("No actions recorded yet"))
	}

	// Header
	header := fmt.Sprintf("%-20s %-20s %-8s %-10s",
		headerStyle.Render("Timestamp"),
		headerStyle.Render("Tool"),
		headerStyle.Render("Status"),
		headerStyle.Render("Duration"),
	)

	// Rows
	var rows []string
	rows = append(rows, header)
	rows = append(rows, dimStyle.Render("─────────────────────────────────────────────────────────────────"))

	// Calculate visible range
	maxVisible := m.height - 15
	if maxVisible < 5 {
		maxVisible = 5
	}
	start := m.historyOffset
	end := start + maxVisible
	if end > len(m.actions) {
		end = len(m.actions)
	}

	for _, action := range m.actions[start:end] {
		status := successStyle.Render("OK")
		if !action.Success {
			status = failStyle.Render("FAIL")
		}

		row := fmt.Sprintf("%-20s %-20s %-8s %10dms",
			dimStyle.Render(action.CreatedAt.Format("2006-01-02 15:04:05")),
			action.ToolName,
			status,
			action.DurationMs,
		)
		rows = append(rows, row)
	}

	// Scroll indicator
	scrollInfo := ""
	if len(m.actions) > maxVisible {
		scrollInfo = fmt.Sprintf("\n\n%s",
			dimStyle.Render(fmt.Sprintf("Showing %d-%d of %d (j/k to scroll)", start+1, end, len(m.actions))))
	}

	content := titleStyle.Render("Action History") + "\n\n" + lipgloss.JoinVertical(lipgloss.Left, rows...) + scrollInfo
	return boxStyle.Render(content)
}

// renderHelp renders the help bar
func (m Model) renderHelp() string {
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Padding(0, 1)

	keyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205"))

	help := fmt.Sprintf("%s/%s/%s tabs • %s refresh • %s/%s scroll • %s quit",
		keyStyle.Render("1"),
		keyStyle.Render("2"),
		keyStyle.Render("3"),
		keyStyle.Render("r"),
		keyStyle.Render("j"),
		keyStyle.Render("k"),
		keyStyle.Render("q"),
	)

	return helpStyle.Render(help)
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
