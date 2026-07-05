# Best Practices for Using Loophole

To get the most out of your AI coding assistant, follow these guidelines for effective interaction and project management.

## Crafting Prompts

1. Be Specific: Instead of "Fix the bug," try "Fix the null pointer exception when the user submits an empty search string in the search handler."
2. Provide Context: Mention relevant files or architectural patterns. Even though Loophole can search, giving a starting point saves time and tokens.
3. Iterative Development: Break down large features into smaller, manageable sub-tasks. Ask the AI to implement the data model first, then the logic, then the interface.

## Managing Context

1. Use Glob and Grep: Use the built-in search tools to narrow down where a feature might be implemented. 
2. Explicit Files: If you know exactly which files need to change, tell the AI: "Read internal/app/handler.go and update the Save method."
3. Session hygiene: Start a new session (`Ctrl+N`) when switching to a completely different task to avoid polluting the AI's memory with irrelevant code.

## Code Review and Safety

1. Review Every Diff: The AI can sometimes make mistakes or introduce subtle bugs. Always review the diff in the sidebar before accepting a change.
2. Run Tests Immediately: After applying a change, run your project's test suite to verify that existing functionality remains intact.
3. Compiler/LSP Feedback: Pay attention to the LSP diagnostics. If the AI introduces a syntax error, the LSP will flag it immediately. Ask the AI to fix the diagnostic if you are unsure how to resolve it.

## Advanced Usage

1. Tool Chaining: The AI can perform multiple actions in a single turn. It might read a file, run a grep search, and then propose an edit all at once.
2. MCP Optimization: Only enable the MCP servers you actually need for your current project to keep the interface clean and responsive.
3. External Editors: Use `Ctrl+E` to write long or complex prompts in your system's default editor (like Vim or VS Code) instead of the terminal buffer.
