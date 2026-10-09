package dialog

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/loophole-ai/loophole-cli/internal/llm/models"
)

func TestModelRowsStayOneLine(t *testing.T) {
	names := []string{
		"Qwen3 Coder 480B A35B Instruct (Vertex AI (OpenAI-compatible))",
		"Gemini Pro Latest (Gemini 3.1 Pro Preview, Vertex AI)",
		"North Mini Code (free)",
		"",
		strings.Repeat("x", 200),
	}

	for _, name := range names {
		row := ansi.Truncate(name, maxDialogWidth, "…")
		if got := lipgloss.Width(row); got > maxDialogWidth {
			t.Errorf("row for %q measures %d, wider than the %d column dialog", name, got, maxDialogWidth)
		}
		if strings.ContainsAny(row, "\n\r") {
			t.Errorf("row for %q contains a line break: %q", name, row)
		}
	}
}

func TestModelDialogRendersLongNamesOnOneLine(t *testing.T) {
	long := "Qwen3 Coder 480B A35B Instruct (Vertex AI (OpenAI-compatible))"
	m := &modelDialogCmp{
		mode:               modeModels,
		provider:           "google",
		availableProviders: []models.ModelProvider{"google"},
		models:             []models.Model{{Name: long}},
		selectedIdx:        0,
	}

	rendered := m.View()

	if strings.Contains(rendered, "\n"+long[:20]) {
		t.Errorf("long model name wrapped onto its own line:\n%s", rendered)
	}
	if !strings.Contains(rendered, "…") {
		t.Errorf("long model name was not truncated:\n%s", rendered)
	}

	for _, line := range strings.Split(rendered, "\n") {
		if w := lipgloss.Width(line); w > maxDialogWidth+6 {
			t.Errorf("rendered line is %d wide, dialog cannot hold it: %q", w, line)
		}
	}
}
