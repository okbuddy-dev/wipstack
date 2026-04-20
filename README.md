# wipstack

A CLI for maintaining a stack of in-progress tasks. `wip` stores tasks in a JSON file and renders them with depth indention so the top of the stack is visually prominent.

## Install

```bash
go build -o wip .
```

## Usage

```
wip [flags] [command] [args...]

Commands:
  list              Show the stack (default)
  push <text>       Push text onto the stack
  pop               Remove and print the top item
  clear             Remove all items

Flags:
  --json, -j        Output in JSON format
```

### Examples

```bash
$ wip push fix the login bug
$ wip push write tests
$ wip push deploy to staging
$ wip list
fix the login bug
  write tests
    deploy to staging

$ wip pop
deploy to staging

$ wip --json list
["fix the login bug","write tests"]

$ wip --json clear
["fix the login bug","write tests"]

$ wip push start over
$ wip --json list
["start over"]
```

## Storage

State is stored at `$XDG_CONFIG_HOME/wipstack/stack.json` (defaulting to `~/.config/wipstack/stack.json`).

## License

Unlicense