# Architecture Overview

Loophole is written in Go and leverages several powerful libraries to provide a responsive and reliable terminal experience.

## Components

### Core (`/internal`)
- **llm**: Handles communication with AI providers (Anthropic, OpenAI, etc.).
- **tui**: Built with the Bubble Tea framework. Manages the terminal rendering and user input.
- **lsp**: Implements the Language Server Protocol client.
- **db**: Manages the SQLite database for persistence.
- **message**: Defines the data structures for user and AI messages.

### CLI (`/cmd`)
- Uses the Cobra library to handle command-line arguments and subcommands.

## Data Flow

1. **Input**: User types a message in the TUI.
2. **Processing**: The message is sent to the configured `llm` provider.
3. **Reasoning**: The AI decides to use a tool.
4. **Execution**: Loophole executes the tool locally (e.g., reads a file) and sends the result back to the AI.
5. **Modification**: If the AI proposes a code change, the TUI displays a diff.
6. **Commit**: Upon user approval, the changes are written to the disk using the `fileutil` and `diff` packages.

## Frameworks and Libraries

- **Bubble Tea**: The core TUI framework for state management and rendering.
- **Lip Gloss**: Used for styling the TUI (colors, borders, layouts).
- **Cobra**: For building the CLI interface.
- **Viper**: For configuration management.
- **SQLite**: For local data persistence.
