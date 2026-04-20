package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"wipstack/internal/stack"
)

func configPath() string {
	dir := os.Getenv("XDG_CONFIG_HOME")
	if dir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, "cannot determine home directory")
			os.Exit(1)
		}
		dir = filepath.Join(home, ".config")
	}
	return filepath.Join(dir, "wipstack", "stack.json")
}

func main() {
	jsonFlag := false
	fs := flagSet()
	fs.BoolVar(&jsonFlag, "json", false, "output in JSON format")
	fs.BoolVar(&jsonFlag, "j", false, "output in JSON format (shorthand)")
	fs.Parse(os.Args[1:])

	args := fs.Args()
	cmd := "list"
	if len(args) > 0 {
		cmd = args[0]
		args = args[1:]
	}

	path := configPath()

	switch cmd {
	case "list":
		cmdList(path, jsonFlag)
	case "push":
		cmdPush(path, jsonFlag, args)
	case "pop":
		cmdPop(path, jsonFlag)
	case "clear":
		cmdClear(path, jsonFlag)
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", cmd)
		os.Exit(1)
	}
}

func flagSet() *flag.FlagSet {
	return flag.NewFlagSet("wip", flag.ExitOnError)
}

func cmdList(path string, jsonOut bool) {
	s, err := stack.Load(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	items := s.Items()
	if jsonOut {
		data, _ := json.Marshal(items)
		fmt.Println(string(data))
		return
	}
	for i, item := range items {
		indent := strings.Repeat("  ", i)
		fmt.Printf("%s%s\n", indent, item)
	}
}

func cmdPush(path string, jsonOut bool, args []string) {
	text := strings.Join(args, " ")
	if text == "" {
		fmt.Fprintln(os.Stderr, "push requires text")
		os.Exit(1)
	}
	s, err := stack.Load(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	s.Push(text)
	if err := s.Save(path); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if jsonOut {
		data, _ := json.Marshal(s.Items())
		fmt.Println(string(data))
	}
}

func cmdPop(path string, jsonOut bool) {
	s, err := stack.Load(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	item, ok := s.Pop()
	if !ok {
		if jsonOut {
			fmt.Println("null")
		}
		os.Exit(1)
	}
	if err := s.Save(path); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if jsonOut {
		data, _ := json.Marshal(item)
		fmt.Println(string(data))
		return
	}
	fmt.Println(item)
}

func cmdClear(path string, jsonOut bool) {
	s, err := stack.Load(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	cleared := s.Clear()
	if err := s.Save(path); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if jsonOut {
		data, _ := json.Marshal(cleared)
		fmt.Println(string(data))
	}
}