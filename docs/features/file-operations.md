# File Operations and Code Modification

Loophole provides the AI with advanced tools to safely interact with your local filesystem.

## Reading Files

The AI can explore your project using several tools:
- **ls**: List contents of directories to understand project structure.
- **view**: Read the full content of a specific file.
- **grep**: Search for specific strings or patterns across multiple files.
- **glob**: Find files matching complex wildcard patterns.

## Modifying Code

Loophole uses atomic operations to modify your code, ensuring that files are either fully updated or left untouched in case of an error.

### Edit Tool
The AI sends specific replacement chunks. It identifies a block of code and provides the new version. Loophole then performs a precise string replacement.

### Patch Tool
For larger or more complex refactors, the AI can generate unified diffs (patches) which are then applied to the file.

## User Authorization

Before any file is modified on your disk, Loophole will display a diff in the sidebar or an overlay. You must explicitly accept the change by pressing the confirmation key (usually 'a' or 'Enter'). This prevents the AI from making unwanted changes to your codebase.
