# Model Context Protocol (MCP)

Loophole is a fully compliant MCP client. The Model Context Protocol is an open standard that allows AI applications to connect to external data sources and tools using a unified interface.

## Extending Loophole

By adding MCP servers to your configuration, you can give the AI access to:
- Google Search or other web search engines.
- Documentation for specific libraries.
- Database schemas and query tools.
- GitHub issues and pull requests.
- Custom internal scripts and tools.

## Configuration

MCP servers are configured in the `mcpServers` section of `.loophole.json`. Each server runs as a separate process that Loophole communicates with via JSON-RPC.

Example of adding a filesystem server:

```json
{
  "mcpServers": {
    "filesystem": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-filesystem", "/path/to/allowed/dir"]
    }
  }
}
```

## Tool Discovery

Once an MCP server is connected, Loophole automatically discovers the tools it provides and makes them available to the AI. You can see the list of active tools in the help menu.
