package teahandler

import (
	"os"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/tree"
	"github.com/charmbracelet/ssh"
	"github.com/muesli/termenv"
	"github.com/vinayakankugoyal/sshresume/pkg/config"
)

const (
	sidebarWidth    = 30  // Width of the left sidebar content
	maxContentWidth = 100 // Max width for content rendering
)

// treeItem represents a flattened tree item for display.
type treeItem struct {
	node     *config.TreeNode
	depth    int
	expanded bool
}

type model struct {
	tree         *config.TreeNode
	items        []treeItem
	expanded     map[string]bool // Track expanded state by path
	cursor       int
	selectedFile string

	viewport viewport.Model
	focused  int // 0: sidebar, 1: content

	width  int
	height int
	name   string
	ready  bool
}

func (m model) Init() tea.Cmd {
	return nil
}

// flattenTree converts the tree structure into a flat list for display.
func (m model) flattenTree() []treeItem {
	var items []treeItem
	var flatten func(node *config.TreeNode, depth int)

	flatten = func(node *config.TreeNode, depth int) {
		if node == nil {
			return
		}

		// Skip the root node itself, only show its children.
		if depth >= 0 {
			expanded := m.expanded[node.Path]
			items = append(items, treeItem{
				node:     node,
				depth:    depth,
				expanded: expanded,
			})

			// Only show children if this node is a directory and is expanded.
			if !node.IsDir || !expanded {
				return
			}
		}

		for _, child := range node.Children {
			flatten(child, depth+1)
		}
	}

	// Start from depth -1 so root's children appear at depth 0.
	flatten(m.tree, -1)
	return items
}

// selectFirstFile finds and selects the first markdown file in the tree.
func (m *model) selectFirstFile() {
	var findFirst func(node *config.TreeNode) string

	findFirst = func(node *config.TreeNode) string {
		if node == nil {
			return ""
		}

		if !node.IsDir {
			return node.Path
		}

		for _, child := range node.Children {
			if result := findFirst(child); result != "" {
				// Expand parent directories.
				m.expanded[node.Path] = true
				return result
			}
		}

		return ""
	}

	m.selectedFile = findFirst(m.tree)
	m.items = m.flattenTree()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	// Initialize on first update.
	if !m.ready {
		m.ready = true
		m.selectFirstFile()
		// Initialize viewport with defaults, will be resized immediately if WindowSizeMsg was cached or comes next
		m.viewport = viewport.New(0, 0)
		m.viewport.Style = contentStyle
		m = m.updateContent()
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			m.focused = (m.focused + 1) % 2
			return m, nil
		}

		if m.focused == 0 {
			// Sidebar Navigation
			switch keypress := msg.String(); keypress {
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.items)-1 {
					m.cursor++
				}
			case "enter", " ":
				if m.cursor < len(m.items) {
					item := m.items[m.cursor]
					if item.node.IsDir {
						m.expanded[item.node.Path] = !m.expanded[item.node.Path]
						m.items = m.flattenTree()
					} else {
						m.selectedFile = item.node.Path
						m = m.updateContent()
					}
				}
			}
		} else {
			// Content Navigation (Viewport)
			m.viewport, cmd = m.viewport.Update(msg)
			cmds = append(cmds, cmd)
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Calculate viewport dimensions
		// Sidebar width = sidebarWidth(30) + padding(2) + border(1) = 33
		const sidebarTotalWidth = sidebarWidth + 3
		const footerHeight = 2 // 1 line for help text + 1 line padding

		vpWidth := m.width - sidebarTotalWidth
		if vpWidth < 0 {
			vpWidth = 0
		}

		// Viewport Height = Screen Height - Top/Bottom Padding(2) - Footer
		vpHeight := m.height - 2 - footerHeight
		if vpHeight < 0 {
			vpHeight = 0
		}

		m.viewport.Width = vpWidth
		m.viewport.Height = vpHeight

		if m.ready {
			m = m.updateContent()
		}
	}

	return m, tea.Batch(cmds...)
}

func (m model) updateContent() model {
	if m.selectedFile == "" {
		m.viewport.SetContent("No file selected")
		return m
	}

	mdContent, err := os.ReadFile(m.selectedFile)
	if err != nil {
		m.viewport.SetContent("Error reading file: " + err.Error())
		return m
	}

	// Calculate content width for word wrapping
	// Viewport width already accounts for sidebar.
	// We need to account for viewport padding (2 horizontal).
	contentWidth := m.viewport.Width - 4
	if contentWidth > maxContentWidth {
		contentWidth = maxContentWidth
	}
	if contentWidth < 10 {
		contentWidth = 10
	}

	renderer, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithColorProfile(termenv.ANSI256),
		glamour.WithWordWrap(contentWidth),
	)

	rendered, err := renderer.Render(string(mdContent))
	if err != nil {
		m.viewport.SetContent("Error rendering markdown: " + err.Error())
		return m
	}

	m.viewport.SetContent(rendered)

	return m
}

// --- Styling ---
var (
	cursorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true)
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))
	sidebarStyle  = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, true, false, false).
			Padding(1)
	contentStyle = lipgloss.NewStyle().Padding(1, 2)
	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).PaddingTop(1)
)

func (m model) View() string {
	if !m.ready {
		return "Initializing..."
	}

	// Build lipgloss tree recursively
	idx := -1
	t := m.buildLipglossTree(m.tree, &idx).
		Enumerator(tree.RoundedEnumerator)

	const footerHeight = 2

	// Sidebar styling
	// Force height to fill screen (minus padding and footer)
	sidebarHeight := m.height - 2 - footerHeight
	if sidebarHeight < 0 {
		sidebarHeight = 0
	}

	currentSidebarStyle := sidebarStyle.
		Width(sidebarWidth).
		Height(sidebarHeight)

	// Dynamic border color based on focus
	if m.focused == 0 {
		currentSidebarStyle = currentSidebarStyle.BorderForeground(lipgloss.Color("212"))
	} else {
		currentSidebarStyle = currentSidebarStyle.BorderForeground(lipgloss.Color("240"))
	}

	sidebar := currentSidebarStyle.Render(t.String())

	// Viewport render
	content := m.viewport.View()

	mainView := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, content)

	// Help footer
	helpText := "Tab: Switch Focus • j/k: Navigate • Enter: Open • q: Quit"
	footer := helpStyle.Render(helpText)

	return lipgloss.JoinVertical(lipgloss.Top, mainView, footer)
}

// buildLipglossTree recursively builds a lipgloss tree from the config tree.
func (m model) buildLipglossTree(node *config.TreeNode, currentIdx *int) *tree.Tree {
	if node == nil {
		return tree.Root(".")
	}

	var children []any

	// Process children if this is the root or an expanded directory.
	for _, child := range node.Children {
		*currentIdx++
		idx := *currentIdx

		// Create label with styling.
		label := child.Name

		// Apply styling based on state.
		if idx == m.cursor {
			label = cursorStyle.Render("▸ " + label)
		} else if !child.IsDir && child.Path == m.selectedFile {
			label = selectedStyle.Render("• " + label)
		} else {
			label = "  " + label
		}

		if child.IsDir {
			if m.expanded[child.Path] {
				// Directory is expanded, recursively build its children.
				childTree := m.buildLipglossTree(child, currentIdx)
				children = append(children, childTree.Root(label))
			} else {
				// Directory is collapsed, don't show children.
				children = append(children, tree.Root(label))
			}
		} else {
			// File node (leaf).
			children = append(children, tree.Root(label))
		}
	}

	return tree.Root("").Child(children...)
}

// NewHandler creates a new bubbletea handler with the specified configuration.
func NewHandler(tree *config.TreeNode) func(ssh.Session) (tea.Model, []tea.ProgramOption) {
	return func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
		pty, _, _ := s.Pty()

		m := model{
			tree:     tree,
			expanded: make(map[string]bool),
			width:    pty.Window.Width,
			height:   pty.Window.Height,
			focused:  0, // Start focused on sidebar
		}

		return m, []tea.ProgramOption{tea.WithAltScreen()}
	}
}
