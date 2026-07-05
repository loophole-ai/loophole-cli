# Introduction to Loophole

Loophole is a high-performance, terminal-based AI coding assistant. It integrates directly into your command-line environment, providing a Terminal User Interface (TUI) for interacting with large language models that have direct access to your codebase.

## Purpose

The main goal of Loophole is to bridge the gap between high-level AI reasoning and the low-level reality of a local development environment. Instead of copying and pasting code snippets between a web browser and an editor, Loophole lives where your code lives.

Unlike IDE plugins that can be heavy or gated behind proprietary ecosystems, Loophole is a lightweight, standalone tool designed for the unix-philosophy enthusiast. It works everywhere a terminal can run and respects your local environment configurations.

## Core Philosophies

1. Safety First: Any destructive action or external command execution requires explicit user authorization. The AI lives in a "sandboxed" session until you verify and apply its suggestions.
2. Context Awareness: By using advanced search and indexing tools, Loophole provides the AI with relevant file context without overwhelming the token limit. It uses RAG-like techniques (Retrieval Augmented Generation) locally to find the right files for your query.
3. Extensibility: Through the Model Context Protocol (MCP), Loophole can be extended to interact with any external service or data source.
4. Privacy: All conversation history and session data are stored locally in a SQLite database. Your code is only sent to the AI providers you explicitly configure.
5. Latency Optimized: Built in Go, Loophole is designed to be snappy and responsive, even when handling large projects with thousands of files.
