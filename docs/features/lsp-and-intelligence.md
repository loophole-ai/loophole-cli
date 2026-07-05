# Language Server Protocol (LSP) Integration

Loophole features deep integration with the Language Server Protocol to provide "IDE-like" intelligence within a terminal environment.

## Capabilities

The LSP integration allows Loophole and the AI to:
- Detect syntax errors and potential bugs in real-time.
- View detailed diagnostic messages.
- Navigate to symbol definitions (Go to Definition).
- Find all occurrences of a variable or function (Find References).

## How it Works

When you open a project, Loophole attempts to detect the programming language and start the corresponding LSP server (e.g., `gopls` for Go, `pyright` for Python, `tsserver` for JavaScript/TypeScript).

The information from the LSP is fed back into the AI's context. For example, if the AI proposes a change that introduces a syntax error, Loophole will see the diagnostic from the LSP and can warn you or allow the AI to fix it immediately.

## Configuration

LSP servers are typically managed via the Model Context Protocol (MCP) using the `mcp-language-server`. You can configure which servers are active in your `.loophole.json` file.
