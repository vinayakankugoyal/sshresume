package teahandler

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/vinayakankugoyal/sshresume/pkg/config"
)

const (
	appWidth  = 78
	appHeight = 30
)

type model struct {
	Tabs      []string
	TabFiles  []string
	activeTab int
	width     int
	height    int
	name      string
	email     string
	github    string
	linkedin  string
	viewport  viewport.Model
	ready     bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Initialize viewport on first update
	if !m.ready {
		contentBoxWidth := appWidth - docStyle.GetHorizontalFrameSize() - 2
		m.viewport = viewport.New(contentBoxWidth, appHeight-10)
		m.ready = true
		m = m.updateViewportContent()
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "right", "l", "n", "tab":
			if m.activeTab < len(m.Tabs)-1 {
				m.activeTab++
				m = m.updateViewportContent()
			}
			return m, nil
		case "left", "h", "p", "shift+tab":
			if m.activeTab > 0 {
				m.activeTab--
				m = m.updateViewportContent()
			}
			return m, nil
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if m.ready {
			contentBoxWidth := appWidth - docStyle.GetHorizontalFrameSize() - 2
			m.viewport.Width = contentBoxWidth
			m.viewport.Height = appHeight - 10
		}
	}

	// Pass other messages to viewport for scrolling (including arrow keys for scrolling)
	if m.ready {
		m.viewport, cmd = m.viewport.Update(msg)
	}
	return m, cmd
}

func (m model) updateViewportContent() model {
	// Read and render markdown content for the active tab
	mdContent, err := os.ReadFile(m.TabFiles[m.activeTab])
	if err != nil {
		m.viewport.SetContent("Error reading file: " + err.Error())
		return m
	}

	// Use fixed app width for content, accounting for outer padding and borders
	contentWidth := appWidth - docStyle.GetHorizontalFrameSize() - windowStyle.GetHorizontalFrameSize()

	// Create renderer with proper width
	renderer, _ := glamour.NewTermRenderer(
		glamour.WithStylePath("dark"),
		glamour.WithWordWrap(contentWidth),
	)

	rendered, err := renderer.Render(string(mdContent))
	if err != nil {
		m.viewport.SetContent("Error rendering markdown: " + err.Error())
		return m
	}

	m.viewport.SetContent(rendered)
	m.viewport.GotoTop()
	return m
}

func tabGapBorder() lipgloss.Border {
	border := lipgloss.HiddenBorder()
	border.BottomLeft = "─"
	border.Bottom = "─"
	border.BottomRight = "┐"
	return border
}

var (
	docStyle       = lipgloss.NewStyle().Padding(1, 2, 1, 2)
	highlightColor = lipgloss.Color("#7D56F4")
	tabGapStyle    = lipgloss.NewStyle().Border(tabGapBorder(), true).BorderForeground(highlightColor)
	windowStyle    = lipgloss.NewStyle().BorderForeground(highlightColor).Padding(2, 2).Border(lipgloss.NormalBorder()).UnsetBorderTop()
)

func (m model) View() string {
	doc := strings.Builder{}

	// Render header with name and contact details
	nameStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(highlightColor).
		Align(lipgloss.Center)

	contactStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Align(lipgloss.Center)

	headerStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(highlightColor).
		Padding(1, 2).
		Width(appWidth - 8).
		Align(lipgloss.Center)

	header := lipgloss.JoinVertical(lipgloss.Center,
		nameStyle.Render(m.name),
		"",
		contactStyle.Render(m.email),
		contactStyle.Render(m.github+" • "+m.linkedin),
	)

	doc.WriteString(headerStyle.Render(header))
	doc.WriteString("\n\n")

	// Calculate width for content box
	contentBoxWidth := appWidth - docStyle.GetHorizontalFrameSize() - 2

	// Render tabs with position-specific borders
	var renderedTabs []string
	for i, t := range m.Tabs {
		isActive := i == m.activeTab
		// Define border based on position and active state
		border := lipgloss.RoundedBorder()

		isFirst := i == 0

		if isFirst {
			if isActive {
				border.BottomLeft = "│"
				border.Bottom = " "
				border.BottomRight = "└"
			} else {
				border.BottomLeft = "├"
				border.Bottom = "─"
				border.BottomRight = "┴"
			}
		} else {
			if isActive {
				border.BottomLeft = "┘"
				border.Bottom = " "
				border.BottomRight = "└"
			} else {
				border.BottomLeft = "┴"
				border.Bottom = "─"
				border.BottomRight = "┴"
			}
		}

		style := lipgloss.NewStyle().
			Border(border, true).
			BorderForeground(highlightColor).
			Padding(0, 1)

		if isActive {
			style = style.Bold(true)
		}

		renderedTabs = append(renderedTabs, style.Render(t))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	gap := contentBoxWidth - lipgloss.Width(row)
	if gap > 0 {
		row = lipgloss.JoinHorizontal(lipgloss.Bottom, row, tabGapStyle.Render(strings.Repeat(" ", gap)))
	}

	doc.WriteString(row)
	doc.WriteString("\n")

	// Render viewport content.
	doc.WriteString(windowStyle.Width(contentBoxWidth).Render(m.viewport.View()))

	// Render help bar with scroll percentage.
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	// Calculate scroll percentage: show how much of the content we've scrolled
	// through plus the visible portion on screen.
	var scrollPercent int
	if m.viewport.TotalLineCount() > 0 {
		// Calculate what percentage of content is above
		// the current view + visible content.
		visibleLines := m.viewport.Height
		totalLines := m.viewport.TotalLineCount()
		scrolledLines := m.viewport.YOffset

		if totalLines <= visibleLines {
			// All content fits on screen.
			scrollPercent = 100
		} else {
			// Show percentage up to the bottom of the visible area.
			scrollPercent = int(float64(scrolledLines+visibleLines) / float64(totalLines) * 100)
			if scrollPercent > 100 {
				scrollPercent = 100
			}
		}
	}

	leftHelp := "↑/↓: scroll • ←/→: switch tabs • q: quit"
	rightHelp := fmt.Sprintf("%d%%", scrollPercent)

	// Calculate spacing between left and right help text.
	spacing := contentBoxWidth - lipgloss.Width(leftHelp) - lipgloss.Width(rightHelp)
	if spacing < 0 {
		spacing = 0
	}

	helpLine := leftHelp + strings.Repeat(" ", spacing) + rightHelp

	doc.WriteString("\n")
	doc.WriteString(helpStyle.Render(helpLine))

	// Center the entire output.
	output := doc.String()
	centered := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, docStyle.Render(output))
	return centered
}

// NewHandler creates a new bubbletea handler with the specified configuration
func NewHandler(cfg *config.Config) func(ssh.Session) (tea.Model, []tea.ProgramOption) {
	return func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
		pty, _, _ := s.Pty()

		// Extract tab names and files from config
		tabs := make([]string, len(cfg.Tabs))
		tabFiles := make([]string, len(cfg.Tabs))
		for i, tab := range cfg.Tabs {
			tabs[i] = tab.Name
			tabFiles[i] = tab.File
		}

		m := model{
			Tabs:     tabs,
			TabFiles: tabFiles,
			width:    pty.Window.Width,
			height:   pty.Window.Height,
			name:     cfg.Profile.Name,
			email:    cfg.Profile.Email,
			github:   cfg.Profile.GitHub,
			linkedin: cfg.Profile.LinkedIn,
		}

		return m, []tea.ProgramOption{tea.WithAltScreen()}
	}
}
