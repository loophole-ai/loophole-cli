package agent

import "testing"

func TestCleanTitle(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want string
	}{
		{
			name: "thinking preamble on its own line",
			raw:  "Here's a thinking process:\n\n1.  **Analyze User Input:** The user asked to fix the bug\n\n2.  **Identify:** the parser\n\nFix session parser",
			want: "Fix session parser",
		},
		{
			name: "preamble and steps on one line falls back",
			raw:  "Here's a thinking process: 1- analyze under input, 2- pick a fix, 3- write it",
			want: "",
		},
		{
			name: "thinking tag block",
			raw:  "<think>the user wants a title</think>\nAdd retry to upload",
			want: "Add retry to upload",
		},
		{
			name: "only steps remain after preamble",
			raw:  "Reasoning:\n1. first step\n2. second step",
			want: "",
		},
		{
			name: "markdown emphasis",
			raw:  "**Fix session parser**",
			want: "Fix session parser",
		},
		{
			name: "surrounding quotes",
			raw:  `"Add retry to upload"`,
			want: "Add retry to upload",
		},
		{
			name: "inline code",
			raw:  "Fix `run_command` timeout",
			want: "Fix run_command timeout",
		},
		{
			name: "heading markers stripped",
			raw:  "## Refactor the cache layer",
			want: "Refactor the cache layer",
		},
		{
			name: "plain answer untouched",
			raw:  "Add retry to upload",
			want: "Add retry to upload",
		},
		{
			name: "nothing usable",
			raw:  "<thinking>only scratchpad</thinking>",
			want: "",
		},
		{
			name: "empty",
			raw:  "",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cleanTitle(tt.raw); got != tt.want {
				t.Errorf("cleanTitle() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCleanTitleNeverEmptyForAnswer(t *testing.T) {
	raw := "Here's a thinking process:\n\n1.  **Analyze User Input:** look at the repo\n\n2.  **Decide:** add a flag\n\nThe repo needs a flag"
	got := cleanTitle(raw)
	if got == "" {
		t.Fatal("cleanTitle returned empty for a reply that contains an answer")
	}
	if got == "New Session" {
		t.Fatal("cleanTitle returned the placeholder")
	}
}
