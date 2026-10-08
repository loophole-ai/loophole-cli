package models

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// DefaultAPI is the models.dev catalogue endpoint. It is the same source the
// agent engine reads, so the CLI and the engine agree on model ids and pricing.
const DefaultAPI = "https://models.dev/api.json"

// EnvAPI overrides the endpoint, mostly so the catalogue can be served locally.
const EnvAPI = "MODELS_DEV_URL"

var (
	remoteModels    map[ModelID]Model
	remoteProviders map[ModelProvider]bool
	remoteOnce      *sync.Once = &sync.Once{}
	remoteMu        sync.RWMutex
	remoteErr       error
)

// modelsDevPayload is the shape of https://models.dev/api.json.
type modelsDevPayload map[string]modelsDevProvider

type modelsDevProvider struct {
	ID     string                    `json:"id"`
	Name   string                    `json:"name"`
	Env    []string                  `json:"env"`
	Models map[string]modelsDevModel `json:"models"`
}

type modelsDevModel struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Attachment bool   `json:"attachment"`
	Reasoning  bool   `json:"reasoning"`
	ToolCall   bool   `json:"tool_call"`
	Limit      struct {
		Context int64 `json:"context"`
		Output  int64 `json:"output"`
	} `json:"limit"`
	Cost struct {
		Input      float64 `json:"input"`
		Output     float64 `json:"output"`
		CacheRead  float64 `json:"cache_read"`
		CacheWrite float64 `json:"cache_write"`
	} `json:"cost"`
}

// providerAliases maps a configured provider id onto its models.dev name. The
// two catalogues do not agree on spelling, and a provider that cannot be
// resolved simply contributes no models.
var providerAliases = map[ModelProvider]string{
	"fireworks":  "fireworks-ai",
	"gemini":     "google",
	"google":     "google",
	"copilot":    "github-copilot",
	"vertex":     "google-vertex",
	"vertexai":   "google-vertex",
	"bedrock":    "amazon-bedrock",
	"openai":     "openai",
	"anthropic":  "anthropic",
	"mistral":    "mistral",
	"openrouter": "openrouter",
	"groq":       "groq",
	"xai":        "xai",
	"deepseek":   "deepseek",
}

// ResolveProvider maps a configured provider onto its models.dev id.
func ResolveProvider(provider ModelProvider) string {
	if alias, ok := providerAliases[provider]; ok {
		return alias
	}
	// Most ids are identical on both sides, so try the name as given before
	// giving up.
	return string(provider)
}

// FetchModels loads the models.dev catalogue. It runs once and caches the
// result, so a failure leaves the static catalogue in place rather than an
// empty picker.
func FetchModels() error {
	remoteOnce.Do(func() {
		endpoint := DefaultAPI
		if override := os.Getenv(EnvAPI); override != "" {
			endpoint = override
		}

		slog.Info("Fetching models from models.dev", "url", endpoint)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
		if err != nil {
			remoteErr = fmt.Errorf("build models.dev request: %w", err)
			slog.Warn("Failed to reach models.dev, using static models only", "error", err)
			return
		}
		req.Header.Set("Accept", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			remoteErr = fmt.Errorf("fetch models.dev: %w", err)
			slog.Warn("Failed to reach models.dev, using static models only", "error", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			remoteErr = fmt.Errorf("models.dev returned %s", resp.Status)
			slog.Warn("Failed to reach models.dev, using static models only", "error", remoteErr)
			return
		}

		var payload modelsDevPayload
		if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
			remoteErr = fmt.Errorf("decode models.dev: %w", err)
			slog.Warn("Failed to read models.dev, using static models only", "error", err)
			return
		}

		fetched, providers := convert(payload)

		remoteMu.Lock()
		remoteModels = fetched
		remoteProviders = providers
		remoteMu.Unlock()

		slog.Info("Fetched models from models.dev", "providers", len(providers), "models", len(fetched))
	})

	return remoteErr
}

// convert flattens the catalogue into model ids of the form provider/model,
// which is the format the engine and the configuration both use.
func convert(payload modelsDevPayload) (map[ModelID]Model, map[ModelProvider]bool) {
	out := make(map[ModelID]Model, len(payload))
	providers := make(map[ModelProvider]bool, len(payload))

	for key, provider := range payload {
		id := provider.ID
		if id == "" {
			id = key
		}
		local := ModelProvider(strings.ToLower(id))
		providers[local] = true

		for modelKey, remote := range provider.Models {
			modelID := remote.ID
			if modelID == "" {
				modelID = modelKey
			}

			out[ModelID(id+"/"+modelID)] = Model{
				ID:                  ModelID(id + "/" + modelID),
				Name:                remote.Name,
				Provider:            local,
				APIModel:            modelID,
				CostPer1MIn:         remote.Cost.Input,
				CostPer1MOut:        remote.Cost.Output,
				CostPer1MInCached:   remote.Cost.CacheRead,
				CostPer1MOutCached:  remote.Cost.CacheWrite,
				ContextWindow:       remote.Limit.Context,
				DefaultMaxTokens:    remote.Limit.Output,
				CanReason:           remote.Reasoning,
				SupportsAttachments: remote.Attachment,
				// models.dev reports whether a model can call tools, so the
				// inverse is what the rest of the CLI reads.
				DisableTools: !remote.ToolCall,
			}
		}
	}

	return out, providers
}

// GetAllModels returns the static catalogue merged with models.dev. Static
// models are the fallback, so a failed fetch still leaves a usable list.
func GetAllModels() map[ModelID]Model {
	_ = FetchModels()

	all := make(map[ModelID]Model)
	for id, model := range SupportedModels {
		all[id] = model
	}

	remoteMu.RLock()
	defer remoteMu.RUnlock()

	for id, model := range remoteModels {
		if static, exists := all[id]; exists {
			// Keep the static flags where the catalogue says nothing.
			if static.SupportsAttachments && !model.SupportsAttachments {
				model.SupportsAttachments = true
			}
			if static.DisableTools && !model.DisableTools {
				model.DisableTools = true
			}
		}
		all[id] = model
	}

	return all
}

// GetRemoteModels returns only the models.dev entries.
func GetRemoteModels() map[ModelID]Model {
	_ = FetchModels()

	remoteMu.RLock()
	defer remoteMu.RUnlock()

	out := make(map[ModelID]Model, len(remoteModels))
	for id, model := range remoteModels {
		out[id] = model
	}
	return out
}

// IsRemoteProvider reports whether the catalogue knows a provider. The name is
// resolved first, so a configured provider such as fireworks or gemini is
// matched against its models.dev id rather than reported as unknown.
func IsRemoteProvider(provider ModelProvider) bool {
	_ = FetchModels()

	resolved := ResolveProvider(provider)

	remoteMu.RLock()
	defer remoteMu.RUnlock()

	if remoteProviders[ModelProvider(resolved)] {
		return true
	}
	return remoteProviders[provider]
}

// RefreshModels forces a refetch by discarding the cached catalogue.
func RefreshModels() error {
	remoteMu.Lock()
	remoteOnce = &sync.Once{}
	remoteMu.Unlock()
	return FetchModels()
}
