package completions

import (


	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/loophole-ai/loophole-cli/internal/tui/components/dialog"
)

type commandCompletionProvider struct {
	commands []dialog.Command
}

func (p *commandCompletionProvider) GetId() string {
	return "command"
}

func (p *commandCompletionProvider) GetEntry() dialog.CompletionItemI {
	return dialog.NewCompletionItem(dialog.CompletionItem{
		Title: "Commands",
		Value: "command",
	})
}

func (p *commandCompletionProvider) GetChildEntries(trigger string, query string) ([]dialog.CompletionItemI, error) {
	var titles []string
	titleToCmd := make(map[string]dialog.Command)
	for _, cmd := range p.commands {
		titles = append(titles, cmd.Title)
		titleToCmd[cmd.Title] = cmd
	}

	matches := fuzzy.Find(query, titles)
	items := make([]dialog.CompletionItemI, 0, len(matches))
	for _, title := range matches {
		cmd := titleToCmd[title]
		items = append(items, dialog.NewCompletionItem(dialog.CompletionItem{
			Title: title,
			Value: cmd.Title, // We use the title as the value to be inserted
		}))
	}

	return items, nil
}

func NewCommandCompletionProvider(commands []dialog.Command) dialog.CompletionProvider {
	return &commandCompletionProvider{
		commands: commands,
	}
}
