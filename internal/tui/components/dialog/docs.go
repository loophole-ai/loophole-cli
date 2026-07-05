package dialog

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/loophole-ai/loophole-cli/internal/tui/layout"
	"github.com/loophole-ai/loophole-cli/internal/tui/styles"
	"github.com/loophole-ai/loophole-cli/internal/tui/theme"
	"github.com/loophole-ai/loophole-cli/internal/tui/util"
)

type DocsDialogCmp interface {
	tea.Model
	layout.Bindings
}

type CloseDocsDialogMsg struct{}

type docsDialogCmp struct{}

func (d *docsDialogCmp) Init() tea.Cmd {
	return nil
}

func (d *docsDialogCmp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("esc", "q", "ctrl+d"))):
			return d, util.CmdHandler(CloseDocsDialogMsg{})
		}
	}
	return d, nil
}

func (d *docsDialogCmp) View() string {
	t := theme.CurrentTheme()
	baseStyle := styles.BaseStyle()

	// Initial render of content to determine max width
	header := styles.Bold().Foreground(t.Primary()).Render("Documentation")
	repoLabel := "GitHub Repository:"
	repoLink := styles.Regular().Foreground(t.MarkdownLink()).Render("https://github.com/loophole-ai/loophole-cli")
	localLabel := "Local Documentation Folder:"
	localLink := styles.Regular().Foreground(t.MarkdownLink()).Render("./docs/")
	filesLabel := "Key Documentation Files:"
	file1 := " - Introduction: " + styles.Regular().Foreground(t.MarkdownLink()).Render("docs/introduction.md")
	file2 := " - Getting Started: " + styles.Regular().Foreground(t.MarkdownLink()).Render("docs/getting-started.md")
	file3 := " - Configuration: " + styles.Regular().Foreground(t.MarkdownLink()).Render("docs/configuration.md")
	file4 := " - Key Bindings: " + styles.Regular().Foreground(t.MarkdownLink()).Render("docs/key-bindings.md")
	file5 := " - AI Chat: " + styles.Regular().Foreground(t.MarkdownLink()).Render("docs/features/ai-chat.md")
	file6 := " - File Operations: " + styles.Regular().Foreground(t.MarkdownLink()).Render("docs/features/file-operations.md")
	file7 := " - LSP & Intelligence: " + styles.Regular().Foreground(t.MarkdownLink()).Render("docs/features/lsp-and-intelligence.md")
	file8 := " - Best Practices: " + styles.Regular().Foreground(t.MarkdownLink()).Render("docs/guides/best-practices.md")
	footer := styles.Regular().Foreground(t.TextMuted()).Render("Press esc or q to close")

	lines := []string{
		header,
		"",
		repoLabel,
		repoLink,
		"",
		localLabel,
		localLink,
		"",
		filesLabel,
		file1,
		file2,
		file3,
		file4,
		file5,
		file6,
		file7,
		file8,
		"",
		footer,
	}

	// Calculate max width
	maxWidth := 0
	for _, l := range lines {
		w := lipgloss.Width(l)
		if w > maxWidth {
			maxWidth = w
		}
	}

	// Render each line with the base style and fixed width to ensure solid background
	var styledLines []string
	itemStyle := baseStyle.Copy().Width(maxWidth)
	for _, l := range lines {
		styledLines = append(styledLines, itemStyle.Render(l))
	}

	content := lipgloss.JoinVertical(lipgloss.Left, styledLines...)

	return baseStyle.
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderBackground(t.Background()).
		BorderForeground(t.TextMuted()).
		Render(content)
}

func (d *docsDialogCmp) BindingKeys() []key.Binding {
	return []key.Binding{
		key.NewBinding(
			key.WithKeys("esc", "q"),
			key.WithHelp("esc/q", "close"),
		),
	}
}

func NewDocsDialogCmp() DocsDialogCmp {
	return &docsDialogCmp{}
}
