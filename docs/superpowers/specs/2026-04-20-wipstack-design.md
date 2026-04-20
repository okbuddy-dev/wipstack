# Wipstack Design Spec

A CLI tool called `wip` for maintaining a stack of in-progress tasks. Simple, composable, stdlib-only Go.

## CLI Interface

```
wip [flags] [command] [args...]

Commands:
  list              Show the stack (default when no command given)
  push <text>       Push text onto the stack
  pop               Remove and print the top item
  clear             Remove all items

Flags:
  --json, -j        Output in JSON format
```

- Running `wip` with no args is equivalent to `wip list`.
- `push` text is all remaining args joined by spaces — no quoting required. `wip push fix the login bug` pushes "fix the login bug".

## Data Model & Storage

**Stack model**: Ordered list of strings. Index 0 = bottom, last index = top.

**File location**: `$XDG_CONFIG_HOME/wipstack/stack.json`, falling back to `~/.config/wipstack/stack.json` when `$XDG_CONFIG_HOME` is unset.

**File format**: A JSON array of strings, bottom to top:
```json
["bottom item", "middle item", "top item"]
```

**File behavior**:
- Directory created on first write if it doesn't exist.
- File created on first `push` if it doesn't exist.
- `list` with no file present: treat as empty stack, no error.
- File is read entirely on each command, modified in memory, written back atomically (write to temp file, then rename).

## Output Behavior

### Text output (default)

- `list`: One line per task, indented by 2 spaces per position from bottom. Bottom at column 0, top most indented.
  ```
  first task
    second task
      third task
  ```
- `push`: Silent, exit 0.
- `pop`: Print the popped item as a single line, exit 0. Empty stack: nothing printed, exit 1.
- `clear`: Silent, exit 0. Empty stack: silent, exit 0.

### JSON output (`--json` / `-j`)

- `list`: JSON array, bottom to top. `["first task","second task","third task"]`. Empty: `[]`.
- `push`: JSON array — the full stack after push. `["first task","pushed task"]`.
- `pop`: JSON string — the popped item. `"third task"`. Empty stack: `null`, exit 1.
- `clear`: JSON array — the items that were cleared. `["first task","second task"]`. Empty: `[]`.

## Architecture

```
wipstack/
├── main.go                  # CLI parsing, dispatch, output formatting
├── internal/
│   └── stack/
│       └── stack.go         # Stack type: Push, Pop, Items, Clear, Load, Save
├── go.mod
├── go.sum
├── .gitignore
├── LICENSE
└── README.md
```

**`internal/stack`** owns:
- `Stack` type (slice of strings)
- `Push(s string)`, `Pop() (string, bool)`, `Items() []string`, `Clear() []string`
- `Load(path string) (Stack, error)` — reads JSON file, returns empty stack if file missing
- `Save(path string) error` — writes JSON atomically, creates dir if needed

**`main.go`** owns:
- Arg parsing (stdlib `os`, `flag`, `strings` — no third-party packages)
- Dispatching to stack operations
- Formatting output for text vs JSON
- Resolving config path (XDG logic)
- Exit codes

## Error Handling

- **No file on disk**: Treated as empty stack. `list` returns nothing, `pop` exits 1. File created on first `push`.
- **Corrupt JSON file**: Print error to stderr, exit 1.
- **Permission error on write**: Print error to stderr, exit 1.
- **`push` with no text**: Print error to stderr, exit 1.
- **`pop` on empty stack**: Nothing printed, exit 1. JSON mode: print `null`, exit 1.

## Constraints

- Go standard library only — no third-party packages or modules.
- No subtasks or nesting — each entry is a flat string.
- Exit code 0 on success, 1 on error/empty-pop.