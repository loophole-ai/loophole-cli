package dialog

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/loophole-ai/loophole-cli/internal/session"
	"github.com/loophole-ai/loophole-cli/internal/tui/layout"
	"github.com/loophole-ai/loophole-cli/internal/tui/styles"
	"github.com/loophole-ai/loophole-cli/internal/tui/theme"
	"github.com/loophole-ai/loophole-cli/internal/tui/util"
)

// ConfirmDeleteSessionMsg asks whether the session should be removed.
type ConfirmDeleteSessionMsg struct {
	Session session.Session
}

// CloseDeleteSessionMsg is sent when the confirmation is dismissed.
type CloseDeleteSessionMsg struct{}

type DeleteSessionDialog interface {
	tea.Model
	layout.Bindings
	SetSession(session.Session)
}

type deleteSessionDialogCmp struct {
	target     session.Session
	selectedNo bool
}

// helpMapping is reused from the quit dialog so both confirmations answer the
// same keys the same way.
func NewDeleteSessionDialogCmp() DeleteSessionDialog {
	return &deleteSessionDialogCmp{
		// Defaulting to No means a stray Enter dismisses rather than destroys.
		selectedNo: true,
	}
}

func (d *deleteSessionDialogCmp) Init() tea.Cmd {
	return nil
}

func (d *deleteSessionDialogCmp) SetSession(s session.Session) {
	d.target = s
}

func (d *deleteSessionDialogCmp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, helpKeys.LeftRight) || key.Matches(msg, helpKeys.Tab):
			d.selectedNo = !d.selectedNo
			return d, nil
		case key.Matches(msg, helpKeys.EnterSpace):
			if d.selectedNo {
				return d, util.CmdHandler(CloseDeleteSessionMsg{})
			}
			return d, util.CmdHandler(ConfirmDeleteSessionMsg{Session: d.target})
		case key.Matches(msg, helpKeys.Yes):
			return d, util.CmdHandler(ConfirmDeleteSessionMsg{Session: d.target})
		case key.Matches(msg, helpKeys.No):
			return d, util.CmdHandler(CloseDeleteSessionMsg{})
		}
	}
	return d, nil
}

func (d *deleteSessionDialogCmp) View() string {
	t := theme.CurrentTheme()
	baseStyle := styles.BaseStyle()

	question := "Are you sure you want to delete this session?"
	if d.target.Title != "" {
		question += "\n" + clipTitle(d.target.Title, 44)
	}

	yesStyle := baseStyle
	noStyle := baseStyle
	spacerStyle := baseStyle.Background(t.Background())

	if d.selectedNo {
		noStyle = noStyle.Background(t.Primary()).Foreground(t.Background())
		yesStyle = yesStyle.Background(t.Background()).Foreground(t.Primary())
	} else {
		yesStyle = yesStyle.Background(t.Primary()).Foreground(t.Background())
		noStyle = noStyle.Background(t.Background()).Foreground(t.Primary())
	}

	yesButton := yesStyle.Padding(0, 1).Render("Yes")
	noButton := noStyle.Padding(0, 1).Render("No")

	buttons := lipgloss.JoinHorizontal(lipgloss.Left, yesButton, spacerStyle.Render("  "), noButton)

	width := lipgloss.Width(question)
	remainingWidth := width - lipgloss.Width(buttons)
	if remainingWidth > 0 {
		buttons = spacerStyle.Render(strings.Repeat(" ", remainingWidth)) + buttons
	}

	content := baseStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Center,
			question,
			"",
			buttons,
		),
	)

	return baseStyle.Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderBackground(t.Background()).
		BorderForeground(t.TextMuted()).
		Width(lipgloss.Width(content) + 4).
		Render(content)
}

// clipTitle keeps a long session title on one line inside the dialog.
func clipTitle(title string, limit int) string {
	if limit < 1 {
		limit = 1
	}
	runes := []rune(title)
	if len(runes) <= limit {
		return title
	}
	if limit <= 3 {
		return string(runes[:limit])
	}
	return string(runes[:limit-3]) + "..."
}

func (d *deleteSessionDialogCmp) BindingKeys() []key.Binding {
	return layout.KeyMapToSlice(helpKeys)
}
