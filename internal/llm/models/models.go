package models

import "maps"

type (
	ModelID       string
	ModelProvider string
)

type Model struct {
	ID                  ModelID       `json:"id"`
	Name                string        `json:"name"`
	Provider            ModelProvider `json:"provider"`
	APIModel            string        `json:"api_model"`
	CostPer1MIn         float64       `json:"cost_per_1m_in"`
	CostPer1MOut        float64       `json:"cost_per_1m_out"`
	CostPer1MInCached   float64       `json:"cost_per_1m_in_cached"`
	CostPer1MOutCached  float64       `json:"cost_per_1m_out_cached"`
	ContextWindow       int64         `json:"context_window"`
	DefaultMaxTokens    int64         `json:"default_max_tokens"`
	CanReason           bool          `json:"can_reason"`
	SupportsAttachments bool          `json:"supports_attachments"`
	DisableTools        bool          `json:"disable_tools"`
}

const (
	// Bedrock
	BedrockClaude37Sonnet ModelID = "bedrock.claude-3.7-sonnet"
	BedrockClaude46Opus   ModelID = "bedrock.claude-4.6-opus"
	BedrockClaude5Sonnet  ModelID = "bedrock.claude-5-sonnet"

	ProviderBedrock ModelProvider = "bedrock"
	// ForTests
	ProviderMock ModelProvider = "__mock"
)

// Providers in order of popularity
var ProviderPopularity = map[ModelProvider]int{
	ProviderCopilot:    1,
	ProviderAnthropic:  2,
	ProviderOpenAI:     3,
	ProviderGemini:     4,
	ProviderDeepSeek:   5,
	ProviderMistral:    6,
	ProviderGROQ:       7,
	ProviderOpenRouter: 8,
	ProviderBedrock:    9,
	ProviderAzure:      10,
	ProviderXAI:        11,
	ProviderVertexAI:   12,
}

var SupportedModels = map[ModelID]Model{
	//
	// // GEMINI
	// GEMINI25: {
	// 	ID:                 GEMINI25,
	// 	Name:               "Gemini 2.5 Pro",
	// 	Provider:           ProviderGemini,
	// 	APIModel:           "gemini-2.5-pro-exp-03-25",
	// 	CostPer1MIn:        0,
	// 	CostPer1MInCached:  0,
	// 	CostPer1MOutCached: 0,
	// 	CostPer1MOut:       0,
	// },
	//
	// GRMINI20Flash: {
	// 	ID:                 GRMINI20Flash,
	// 	Name:               "Gemini 2.0 Flash",
	// 	Provider:           ProviderGemini,
	// 	APIModel:           "gemini-2.0-flash",
	// 	CostPer1MIn:        0.1,
	// 	CostPer1MInCached:  0,
	// 	CostPer1MOutCached: 0.025,
	// 	CostPer1MOut:       0.4,
	// },
	//
	// // Bedrock
	BedrockClaude37Sonnet: {
		ID:                 BedrockClaude37Sonnet,
		Name:               "Bedrock: Claude 3.7 Sonnet",
		Provider:           ProviderBedrock,
		APIModel:           "anthropic.claude-3-7-sonnet-20250219-v1:0",
		CostPer1MIn:        3.0,
		CostPer1MInCached:  3.75,
		CostPer1MOutCached: 0.30,
		CostPer1MOut:       15.0,
	},
	BedrockClaude46Opus: {
		ID:                 BedrockClaude46Opus,
		Name:               "Bedrock: Claude 4.6 Opus",
		Provider:           ProviderBedrock,
		APIModel:           "anthropic.claude-4-6-opus-20260205-v1:0",
		CostPer1MIn:        15.0,
		CostPer1MInCached:  18.75,
		CostPer1MOutCached: 1.50,
		CostPer1MOut:       75.0,
	},
	BedrockClaude5Sonnet: {
		ID:                 BedrockClaude5Sonnet,
		Name:               "Bedrock: Claude 5 Sonnet",
		Provider:           ProviderBedrock,
		APIModel:           "anthropic.claude-5-sonnet-20260530-v1:0", // Placeholder date
		CostPer1MIn:        2.0,
		CostPer1MInCached:  2.5,
		CostPer1MOutCached: 0.2,
		CostPer1MOut:       10.0,
	},
}

func init() {
	maps.Copy(SupportedModels, AnthropicModels)
	maps.Copy(SupportedModels, OpenAIModels)
	maps.Copy(SupportedModels, GeminiModels)
	maps.Copy(SupportedModels, GroqModels)
	maps.Copy(SupportedModels, AzureModels)
	maps.Copy(SupportedModels, OpenRouterModels)
	maps.Copy(SupportedModels, XAIModels)
	maps.Copy(SupportedModels, VertexAIGeminiModels)
	maps.Copy(SupportedModels, CopilotModels)
	maps.Copy(SupportedModels, DeepSeekModels)
	maps.Copy(SupportedModels, MistralModels)

	// Fetch the model catalogue in the background so it does not block startup.
	go func() {
		_ = FetchModels()
	}()
}
