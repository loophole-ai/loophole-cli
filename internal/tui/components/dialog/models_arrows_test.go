package dialog

import (
	"strings"
	"testing"

	"github.com/loophole-ai/loophole-cli/internal/llm/models"
)

func scrollIndicatorRunes(m *modelDialogCmp) string {
	var sb strings.Builder
	for _, r := range m.getScrollIndicators(40) {
		if r > 126 {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

func TestScrollIndicatorArrows(t *testing.T) {
	providers := []models.ModelProvider{"a", "b"}

	cases := []struct {
		name string
		m    *modelDialogCmp
		want string
	}{
		{
			name: "middle of a long list with a provider to the right",
			m: &modelDialogCmp{
				mode:               modeModels,
				models:             manyModels(40),
				availableProviders: providers,
				scrollOffset:       5,
				hScrollOffset:      0,
				hScrollPossible:    true,
			},
			want: "↑↓→",
		},
		{
			name: "top of the list hides the up arrow",
			m: &modelDialogCmp{
				mode:               modeModels,
				models:             manyModels(40),
				availableProviders: providers,
				scrollOffset:       0,
				hScrollOffset:      0,
				hScrollPossible:    true,
			},
			want: "↓→",
		},
		{
			name: "bottom of the list and last provider",
			m: &modelDialogCmp{
				mode:               modeModels,
				models:             manyModels(40),
				availableProviders: providers,
				scrollOffset:       30,
				hScrollOffset:      1,
				hScrollPossible:    true,
			},
			want: "←↑",
		},
		{
			name: "provider with models on both sides",
			m: &modelDialogCmp{
				mode:               modeModels,
				models:             manyModels(40),
				availableProviders: providers,
				scrollOffset:       5,
				hScrollOffset:      1,
				hScrollPossible:    true,
			},
			want: "←↑↓",
		},
		{
			name: "nothing to scroll shows nothing",
			m: &modelDialogCmp{
				mode:               modeModels,
				models:             manyModels(3),
				availableProviders: providers,
				hScrollPossible:    false,
			},
			want: "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := scrollIndicatorRunes(tc.m); got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestScrollIndicatorArrowsAreRealRunes(t *testing.T) {
	m := &modelDialogCmp{
		mode:               modeModels,
		models:             manyModels(40),
		availableProviders: []models.ModelProvider{"a", "b"},
		scrollOffset:       5,
		hScrollOffset:      1,
		hScrollPossible:    true,
	}

	got := scrollIndicatorRunes(m)
	for _, r := range got {
		switch r {
		case '←', '↑', '→', '↓':
		default:
			t.Errorf("unexpected rune %q (U+%04X) in scroll indicator %q", r, r, got)
		}
	}
}

func manyModels(n int) []models.Model {
	out := make([]models.Model, 0, n)
	for i := 0; i < n; i++ {
		out = append(out, models.Model{ID: models.ModelID(string(rune('a' + i%26))), Name: "m"})
	}
	return out
}
