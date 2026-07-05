# AI Chat and Interaction

The chat interface is the primary way to interact with Loophole. It allows you to describe tasks, ask questions about your code, and review proposed changes.

## Message Handling

Loophole maintains a conversation history for each session. When you send a message, the AI analyzes your request and decides whether it needs to use any of its available tools (like reading a file or running a grep search) to provide a better answer.

## Model Selection

You can switch models in the middle of a conversation if needed. This is useful for:
- Using a cheaper, faster model for simple questions.
- Switching to a more powerful reasoning model for complex debugging or refactoring.

Use `Ctrl+O` to open the model selection dialog.

## Context Management

The AI can only see a certain amount of information at once. Loophole helps manage this by:
1. Allowing you to manually add files to context.
2. Automatically suggesting relevant files based on your query.
3. Compacting old messages to keep the conversation within the model's limits.

## Thinking Mode

Some models support a "Thinking" or "Reasoning" mode. When enabled, the model will output its internal chain of thought before providing a final answer. This is particularly helpful for complex logical tasks.
