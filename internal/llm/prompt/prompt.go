package prompt

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/loophole-ai/loophole-cli/internal/config"
	"github.com/loophole-ai/loophole-cli/internal/llm/models"
	"github.com/loophole-ai/loophole-cli/internal/logging"
)

func GetAgentPrompt(agentName config.AgentName, provider models.ModelProvider) string {
	basePrompt := ""
	switch agentName {
	case config.AgentCoder:
		basePrompt = CoderPrompt(provider)
	case config.AgentTitle:
		basePrompt = TitlePrompt(provider)
	case config.AgentTask:
		basePrompt = TaskPrompt(provider)
	case config.AgentSummarizer:
		basePrompt = SummarizerPrompt(provider)
	default:
		basePrompt = "You are a helpful assistant"
	}

	if agentName == config.AgentCoder || agentName == config.AgentTask {
		// Add context from project-specific instruction files if they exist
		contextContent := getContextFromPaths()
		logging.Debug("Context content", "Context", contextContent)
		if contextContent != "" {
			return fmt.Sprintf("%s\n\n# Project-Specific Context\n Make sure to follow the instructions in the context below\n%s", basePrompt, contextContent)
		}
	}
	return basePrompt
}

var (
	onceContext    sync.Once
	contextContent string
)

func getContextFromPaths() string {
	onceContext.Do(func() {
		var (
			cfg          = config.Get()
			workDir      = cfg.WorkingDir
			contextPaths = cfg.ContextPaths
		)

		contextContent = processContextPaths(workDir, contextPaths)
	})

	return contextContent
}

// processContextPaths reads the configured context paths into one block of text.
//
// The paths are read concurrently but kept in the order they were configured.
// Collecting results as they finished made the block reorder itself between
// runs, so the same prompt produced a different context each time.
func processContextPaths(workDir string, paths []string) string {
	var (
		wg      sync.WaitGroup
		results = make([]string, len(paths))
	)

	// Track processed files to avoid duplicates
	processedFiles := make(map[string]bool)
	var processedMutex sync.Mutex

	for i, p := range paths {
		wg.Add(1)
		go func(i int, p string) {
			defer wg.Done()

			var found []string
			claim := func(path string) bool {
				processedMutex.Lock()
				defer processedMutex.Unlock()
				lowerPath := strings.ToLower(path)
				if processedFiles[lowerPath] {
					return false
				}
				processedFiles[lowerPath] = true
				return true
			}

			if strings.HasSuffix(p, "/") {
				filepath.WalkDir(filepath.Join(workDir, p), func(path string, d os.DirEntry, err error) error {
					if err != nil {
						return err
					}
					if !d.IsDir() && claim(path) {
						if result := processFile(path); result != "" {
							found = append(found, result)
						}
					}
					return nil
				})
			} else {
				fullPath := filepath.Join(workDir, p)
				if claim(fullPath) {
					if result := processFile(fullPath); result != "" {
						found = append(found, result)
					}
				}
			}

			results[i] = strings.Join(found, "\n")
		}(i, p)
	}

	wg.Wait()

	kept := make([]string, 0, len(results))
	for _, result := range results {
		if result != "" {
			kept = append(kept, result)
		}
	}

	return strings.Join(kept, "\n")
}

func processFile(filePath string) string {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return ""
	}
	return "# From:" + filePath + "\n" + string(content)
}
