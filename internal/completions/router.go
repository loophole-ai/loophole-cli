package completions

import (
	"github.com/loophole-ai/loophole-cli/internal/tui/components/dialog"
)

type completionRouter struct {
	fileProvider    dialog.CompletionProvider
	commandProvider dialog.CompletionProvider
}

func (r *completionRouter) GetId() string {
	return "router"
}

func (r *completionRouter) GetEntry() dialog.CompletionItemI {
	return nil
}

func (r *completionRouter) GetChildEntries(trigger string, query string) ([]dialog.CompletionItemI, error) {
	if trigger == "/" {
		if r.commandProvider != nil {
			return r.commandProvider.GetChildEntries(trigger, query)
		}
	} else if trigger == "@" {
		if r.fileProvider != nil {
			return r.fileProvider.GetChildEntries(trigger, query)
		}
	}
	
	// Default to file provider if trigger is unknown or empty
	if r.fileProvider != nil {
		return r.fileProvider.GetChildEntries(trigger, query)
	}
	
	return nil, nil
}

func NewCompletionRouter(fileProvider, commandProvider dialog.CompletionProvider) dialog.CompletionProvider {
	return &completionRouter{
		fileProvider:    fileProvider,
		commandProvider: commandProvider,
	}
}
