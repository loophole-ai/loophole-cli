package models

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// fixture is a trimmed copy of the models.dev shape, taken from a real
// response, so the parser is tested against the fields it actually carries.
const fixture = `{
  "anthropic": {
    "id": "anthropic",
    "name": "Anthropic",
    "env": ["ANTHROPIC_API_KEY"],
    "models": {
      "claude-haiku-4-5": {
        "id": "claude-haiku-4-5",
        "name": "Claude Haiku 4.5 (latest)",
        "attachment": true,
        "reasoning": true,
        "tool_call": true,
        "limit": { "context": 200000, "output": 64000 },
        "cost": { "input": 1, "output": 5, "cache_read": 0.1, "cache_write": 1.25 }
      },
      "no-tools": {
        "id": "no-tools",
        "name": "No Tools",
        "tool_call": false
      }
    }
  },
  "fireworks-ai": {
    "id": "fireworks-ai",
    "name": "Fireworks",
    "models": {
      "kimi-k2": { "id": "kimi-k2", "name": "Kimi K2", "tool_call": true }
    }
  }
}`

func TestConvertProducesEngineStyleIDs(t *testing.T) {
	var payload modelsDevPayload
	if err := json.Unmarshal([]byte(fixture), &payload); err != nil {
		t.Fatalf("decode fixture: %v", err)
	}

	models, providers := convert(payload)

	// The engine and the configuration both use provider/model, so that is the
	// key. A bare model id here would be rejected by the engine.
	key := ModelID("anthropic/claude-haiku-4-5")
	model, ok := models[key]
	if !ok {
		t.Fatalf("model %q missing, got %v", key, len(models))
	}

	if model.Provider != "anthropic" {
		t.Errorf("provider = %q", model.Provider)
	}
	// APIModel is what gets sent to the provider, so it must not carry the
	// provider prefix.
	if model.APIModel != "claude-haiku-4-5" {
		t.Errorf("api model = %q, want the bare id", model.APIModel)
	}
	if model.Name != "Claude Haiku 4.5 (latest)" {
		t.Errorf("name = %q", model.Name)
	}
	if model.ContextWindow != 200000 {
		t.Errorf("context = %d", model.ContextWindow)
	}
	if model.DefaultMaxTokens != 64000 {
		t.Errorf("max tokens = %d", model.DefaultMaxTokens)
	}
	if model.CostPer1MIn != 1 || model.CostPer1MOut != 5 {
		t.Errorf("cost = %v/%v", model.CostPer1MIn, model.CostPer1MOut)
	}
	if model.CostPer1MInCached != 0.1 {
		t.Errorf("cached cost = %v", model.CostPer1MInCached)
	}
	if !model.SupportsAttachments {
		t.Error("attachment flag lost")
	}
	if !model.CanReason {
		t.Error("reasoning flag lost")
	}
	if !providers["anthropic"] || !providers["fireworks-ai"] {
		t.Errorf("providers = %v", providers)
	}
}

func TestConvertInvertsToolCallFlag(t *testing.T) {
	var payload modelsDevPayload
	if err := json.Unmarshal([]byte(fixture), &payload); err != nil {
		t.Fatalf("decode fixture: %v", err)
	}
	models, _ := convert(payload)

	// The rest of the CLI reads DisableTools, so it is the inverse of the
	// catalogue's tool_call.
	withTools := models[ModelID("anthropic/claude-haiku-4-5")]
	if withTools.DisableTools {
		t.Error("a model that can call tools was marked as unable to")
	}
	without := models[ModelID("anthropic/no-tools")]
	if !without.DisableTools {
		t.Error("a model that cannot call tools was marked as able to")
	}
}

func TestConvertSurvivesMissingFields(t *testing.T) {
	// A model with nothing but an id must not panic or disappear.
	var payload modelsDevPayload
	if err := json.Unmarshal([]byte(`{"x":{"id":"x","models":{"bare":{}}}}`), &payload); err != nil {
		t.Fatalf("decode: %v", err)
	}
	models, providers := convert(payload)

	if _, ok := models[ModelID("x/bare")]; !ok {
		t.Fatalf("bare model missing, got %v", models)
	}
	if !providers["x"] {
		t.Error("provider not recorded")
	}
}

// TestResolveProvider covers the names that differ between the CLI's
// configuration and models.dev. A provider that fails to resolve contributes no
// models at all, which is exactly how a configured provider ends up looking
// broken.
func TestResolveProvider(t *testing.T) {
	cases := map[ModelProvider]string{
		"fireworks":   "fireworks-ai",
		"gemini":      "google",
		"vertexai":    "google-vertex",
		"copilot":     "github-copilot",
		"bedrock":     "amazon-bedrock",
		"anthropic":   "anthropic",
		"mistral":     "mistral",
		"openrouter":  "openrouter",
		"coralbricks": "coralbricks",
	}
	for in, want := range cases {
		if got := ResolveProvider(in); got != want {
			t.Errorf("ResolveProvider(%q) = %q, want %q", in, got, want)
		}
	}
}

// TestFetchModelsOverHTTP exercises the whole fetch path against a local server,
// including the fallback when the catalogue is unreachable.
func TestFetchModelsOverHTTP(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(fixture))
	}))
	defer server.Close()

	t.Setenv(EnvAPI, server.URL)
	if err := RefreshModels(); err != nil {
		t.Fatalf("RefreshModels: %v", err)
	}

	remote := GetRemoteModels()
	if len(remote) != 3 {
		t.Fatalf("got %d remote models, want 3: %v", len(remote), remote)
	}
	all := GetAllModels()
	if _, ok := all[ModelID("anthropic/claude-haiku-4-5")]; !ok {
		t.Error("remote model missing from GetAllModels")
	}
	// Static models must survive the merge.
	if _, ok := all[MistralSmall3_1]; !ok {
		t.Error("a static model was dropped from GetAllModels")
	}
}

// TestFetchModelsFallsBackWhenUnreachable proves a bad endpoint leaves the
// static catalogue usable rather than an empty picker.
func TestFetchModelsFallsBackWhenUnreachable(t *testing.T) {
	t.Setenv(EnvAPI, "http://127.0.0.1:1/nope")
	if err := RefreshModels(); err == nil {
		t.Error("expected an error for an unreachable endpoint")
	}

	all := GetAllModels()
	if len(all) == 0 {
		t.Fatal("catalogue is empty after a failed fetch")
	}
	if _, ok := all[MistralSmall3_1]; !ok {
		t.Error("static models should still be present")
	}
}
