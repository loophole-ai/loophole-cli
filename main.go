package main

import (
	"github.com/loophole-ai/loophole-cli/cmd"
	"github.com/loophole-ai/loophole-cli/internal/logging"
)

func main() {
	defer logging.RecoverPanic("main", func() {
		logging.ErrorPersist("Application terminated due to unhandled panic")
	})

	cmd.Execute()
}
