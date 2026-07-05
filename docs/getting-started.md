# Getting Started with Loophole

This guide will help you install, configure, and start using Loophole for the first time.

## Prerequisites

- Go 1.21 or higher (if building from source)
- Node.js (for npm installation)
- An API key from a supported provider (Anthropic, OpenAI, or Google Gemini)

## Installation

### Via npm

This is the recommended method for most users.

```bash
npm install -g @loophole-ai/loophole-cli
```

### Via Raw Script

For environments without npm:

```bash
curl -fsSL https://raw.githubusercontent.com/loophole-ai/loophole-cli/main/install | bash
```

## First Run

1. Open your terminal in a project directory.
2. Set your API key as an environment variable:
   ```bash
   export ANTHROPIC_API_KEY="your-api-key"
   ```
3. Launch the application:
   ```bash
   loophole
   ```

## Initial Setup

When you first launch Loophole, it will look for a `.loophole.json` configuration file in your home directory or the current working directory. If none is found, it will use default settings. It is recommended to create a configuration file to specify your preferred models and settings.

### Interactive Commands

Loophole supports slash commands for quick actions:
- `/help` - Show all commands and keybindings
- `/new` - Start a fresh conversation session
- `/sessions` - Open the session switcher
- `/clear` - Clear the current chat view
- `/about` - Show version and author info

### Custom Commands

You can create your own commands by adding Markdown files to the `.loophole/commands/` directory in your project root. These files will be parsed, and you can use them via `/your-command-name`.

Example `.loophole/commands/test.md`:
```markdown
---
description: Run tests for a specific file
---
I need you to run tests for $FILENAME and report the results.
```

When you type `/test`, Loophole will prompt you for the `$FILENAME` argument.
