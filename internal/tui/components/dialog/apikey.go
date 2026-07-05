package dialog

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/loophole-ai/loophole-cli/internal/llm/models"
	"github.com/loophole-ai/loophole-cli/internal/tui/styles"
	"github.com/loophole-ai/loophole-cli/internal/tui/theme"
	"github.com/loophole-ai/loophole-cli/internal/tui/util"
)

type APIKeySelectedMsg struct {
	Provider models.ModelProvider
	APIKey   string
}

type CloseAPIKeyDialogMsg struct{}

type APIKeyDialogCmp struct {
	width, height int
	input         textinput.Model
	provider      models.ModelProvider
}

func NewAPIKeyDialogCmp(provider models.ModelProvider) APIKeyDialogCmp {
	t := theme.CurrentTheme()
	ti := textinput.New()
	ti.Placeholder = "Enter API Key..."
	ti.EchoMode = textinput.EchoPassword
	ti.EchoCharacter = '•'
	ti.Focus()
	ti.Width = 40
	ti.PromptStyle = ti.PromptStyle.Foreground(t.Primary()).Background(t.Background())
	ti.TextStyle = ti.TextStyle.Foreground(t.Text()).Background(t.Background())
	ti.PlaceholderStyle = ti.PlaceholderStyle.Foreground(t.TextMuted()).Background(t.Background())

	return APIKeyDialogCmp{
		input:    ti,
		provider: provider,
	}
}

func (m APIKeyDialogCmp) Init() tea.Cmd {
	return textinput.Blink
}

func (m APIKeyDialogCmp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.input.Value() != "" {
				return m, util.CmdHandler(APIKeySelectedMsg{
					Provider: m.provider,
					APIKey:   m.input.Value(),
				})
			}
		case "esc":
			return m, util.CmdHandler(CloseAPIKeyDialogMsg{})
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m APIKeyDialogCmp) View() string {
	t := theme.CurrentTheme()
	baseStyle := styles.BaseStyle()

	titleStr := styles.Bold().Foreground(t.Primary()).Render("Configure API Key")
	providerStr := styles.Regular().Foreground(t.TextMuted()).Render("Provider: " + string(m.provider))
	inputStr := m.input.View()

	lines := []string{
		titleStr,
		"",
		providerStr,
		"",
		inputStr,
	}

	// Calculate max width
	maxWidth := 0
	for _, l := range lines {
		w := lipgloss.Width(l)
		if w > maxWidth {
			maxWidth = w
		}
	}

	// Render each line with fixed width to ensure background fills
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

func (m *APIKeyDialogCmp) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m APIKeyDialogCmp) BindingKeys() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "save")),
		key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "cancel")),
	}
}
