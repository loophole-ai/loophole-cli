# Session Management and Persistence

Loophole is designed to handle multiple tasks simultaneously through an organized session management system.

## Data Persistence

All application data is stored in a local SQLite database, typically located in `~/.loophole/loophole.db`. This includes:
- Conversation history (messages, tool calls, and results).
- Metadata about sessions (titles, creation dates).
- User permissions and preferences.

## Managing Sessions

You can maintain multiple isolated conversations. This is useful for tracking different features or bugs without mixing context.

### Commands
- **New Session**: Start a fresh conversation with no previous context (`Ctrl+N`).
- **Switch Session**: Open a dialog to browse and switch between existing sessions (`Ctrl+S`).
- **Persistence**: Sessions are saved automatically. You can close Loophole and resume your conversation later.

## Security

Because the database is stored locally, your conversation history remains on your machine. Loophole does not upload your history to any central server, except for the parts of the code you explicitly send to your configured AI provider during a chat.
