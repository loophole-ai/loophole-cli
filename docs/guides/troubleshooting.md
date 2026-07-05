# Troubleshooting Guide

Common issues and their solutions when using Loophole.

## Connection and API Issues

### Authentication Errors (401 Unauthorized)
- Verify your API keys are correctly exported as environment variables.
- Ensure there are no leading or trailing spaces in your keys.
- For GitHub Copilot, try re-authenticating with the GitHub CLI: `gh auth login`.

### Rate Limit Errors (429 Too Many Requests)
- Loophole has built-in exponential backoff, but you may need to wait if you have exceeded your provider's tier limits.
- Consider switching to a different model or provider if one is consistently blocked.

### Connection Timed Out
- Check your internet connection.
- If you are behind a corporate proxy, you may need to set the `HTTP_PROXY` and `HTTPS_PROXY` environment variables.

## Application Behavior

### Loophole won't start
- Run `loophole --debug` to see detailed startup logs.
- Verify that your `.loophole.json` is valid JSON. Use a linter if necessary.
- Ensure your Go version is 1.21 or higher.

### AI can't see specific files
- Verify that the files are not excluded by your `.gitignore`. Loophole respects git ignored files by default.
- Check file permissions to ensure the user running Loophole has read access.

### TUI Rendering Issues
- Some terminal emulators have issues with the Bubble Tea framework. Try a modern terminal like Alacritty, iTerm2, or Windows Terminal.
- Ensure your `$TERM` environment variable is set correctly (e.g., `xterm-256color`).

## Debugging

If you encounter an unhandled panic or persistent bug:
1. Enable debug mode: `export LOOPHOLE_DEBUG=true`.
2. Reproduce the issue.
3. Check the logs in `~/.loophole/logs/` for the exact error message and stack trace.
4. Open an issue on GitHub with these logs attached.
