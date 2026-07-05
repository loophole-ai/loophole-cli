# Configuration Reference

Loophole is configured via a JSON file named `.loophole.json`. The application looks for this file in the following locations (in order of priority):

1. Path specified by the `LOOPHOLE_CONFIG` environment variable.
2. The current working directory.
3. Your home directory (`~/.loophole.json` or `%USERPROFILE%\.loophole.json`).

## Example Configuration

```json
{
  "agents": {
    "coder": {
      "model": "claude-3.7-sonnet",
      "maxTokens": 8192,
      "temperature": 0.7
    }
  },
  "autoCompact": true,
  "mcpServers": {
    "memory": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-memory"]
    }
  }
}
```

## Top-Level Fields

### agents
A map of agent configurations. Currently, the most important agent is `coder`.
- **model**: The ID of the model to use.
- **maxTokens**: The maximum number of tokens the model can generate.
- **temperature**: Controls randomness in the output (0.0 to 1.0).

### Provider Specific Options

#### Anthropic
- **disableCache**: Set to `true` to disable Anthropic prompt caching (default: `false`).

#### OpenAI
- **baseURL**: Override the default API URL. Useful for local LLMs (like Ollama) or custom proxies.
- **reasoningEffort**: For `o1` and `o3` models, set to `low`, `medium`, or `high` to control reasoning time.
- **extraHeaders**: A map of custom HTTP headers to send with every request.

#### GitHub Copilot
- **bearerToken**: Manually provide a GitHub Copilot bearer token if automatic discovery fails.

### autoCompact
Boolean. When set to true, Loophole will automatically compress older conversation history when the context window is near its limit. This preserves recent context while discarding older, less relevant messages to stay within token limits.

### mcpServers
A map of Model Context Protocol server configurations.
- **command**: The executable to run (e.g., `npx`, `python3`, `node`).
- **args**: An array of arguments to pass to the command.
- **env**: (Optional) Map of environment variables specific to this server process.

## Environment Variables

- `ANTHROPIC_API_KEY`: API key for Anthropic models.
- `OPENAI_API_KEY`: API key for OpenAI models.
- `GEMINI_API_KEY`: API key for Google Gemini models.
- `GITHUB_TOKEN`: Used for authenticating with GitHub Copilot if NOT using the GitHub CLI.
- `LOOPHOLE_DEBUG`: Set to `true` to enable verbose logging to `~/.loophole/logs/`.
- `LOOPHOLE_CONFIG`: Explicitly set the path to your `.loophole.json` configuration file.
