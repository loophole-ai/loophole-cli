package config

import (
	"testing"

	"github.com/loophole-ai/loophole-cli/internal/llm/models"
)

// A config that lists several providers used to send the title request to
// whichever provider key the defaults picked first, so a session named by the
// wrong provider came back as an authentication error and stayed on the
// "New Session" placeholder.
func TestApplyAgentDefaultsFollowsTheCoder(t *testing.T) {
	original := cfg
	t.Cleanup(func() { cfg = original })

	cfg = &Config{
		Agents: map[AgentName]Agent{
			AgentCoder:      {Model: "coralbricks/some-model", MaxTokens: 131072},
			AgentTitle:      {Model: models.Gemini25Flash},
			AgentSummarizer: {Model: models.Claude4Sonnet},
			AgentTask:       {Model: models.GPT41Mini},
		},
	}

	applyAgentDefaults(only(AgentCoder))

	title := cfg.Agents[AgentTitle]
	if title.Model != "coralbricks/some-model" {
		t.Errorf("title model = %q, want the coder's model", title.Model)
	}
	if title.MaxTokens <= 0 || title.MaxTokens > 1024 {
		t.Errorf("title max tokens = %d, want a small budget", title.MaxTokens)
	}
	if got := cfg.Agents[AgentSummarizer].Model; got != "coralbricks/some-model" {
		t.Errorf("summarizer model = %q, want the coder's model", got)
	}
	if got := cfg.Agents[AgentTask].Model; got != "coralbricks/some-model" {
		t.Errorf("task model = %q, want the coder's model", got)
	}
}

func TestApplyAgentDefaultsKeepsExplicitChoices(t *testing.T) {
	original := cfg
	t.Cleanup(func() { cfg = original })

	cfg = &Config{
		Agents: map[AgentName]Agent{
			AgentCoder: {Model: "coralbricks/some-model", MaxTokens: 131072},
			AgentTitle: {Model: models.Gemini25Flash, MaxTokens: 999},
		},
	}

	applyAgentDefaults(only(AgentCoder, AgentTitle))

	if got := cfg.Agents[AgentTitle].Model; got != models.Gemini25Flash {
		t.Errorf("title model = %q, want the configured model", got)
	}
	if got := cfg.Agents[AgentTitle].MaxTokens; got != titleMaxTokens(models.Gemini25Flash) {
		t.Errorf("title max tokens = %d, want %d", got, titleMaxTokens(models.Gemini25Flash))
	}
}

func TestApplyAgentDefaultsWithoutACoder(t *testing.T) {
	original := cfg
	t.Cleanup(func() { cfg = original })

	cfg = &Config{
		Agents: map[AgentName]Agent{
			AgentTitle: {Model: models.Gemini25Flash},
		},
	}

	// A seeded default is all there is to go on; dropping it would leave the
	// title agent with no model at all.
	applyAgentDefaults(only())

	if got := cfg.Agents[AgentTitle].Model; got != models.Gemini25Flash {
		t.Errorf("title model = %q, want the seeded model", got)
	}
}

func only(names ...AgentName) map[AgentName]bool {
	explicit := map[AgentName]bool{}
	for _, name := range names {
		explicit[name] = true
	}
	return explicit
}
