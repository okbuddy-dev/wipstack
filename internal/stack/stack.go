package stack

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Stack struct {
	items []string
}

func New() Stack {
	return Stack{items: []string{}}
}

func (s *Stack) Push(text string) {
	s.items = append(s.items, text)
}

func (s *Stack) Pop() (string, bool) {
	if len(s.items) == 0 {
		return "", false
	}
	idx := len(s.items) - 1
	item := s.items[idx]
	s.items = s.items[:idx]
	return item, true
}

func (s *Stack) Items() []string {
	result := make([]string, len(s.items))
	copy(result, s.items)
	return result
}

func (s *Stack) Clear() []string {
	cleared := s.Items()
	s.items = s.items[:0]
	return cleared
}

func (s *Stack) IsEmpty() bool {
	return len(s.items) == 0
}

func (s *Stack) Save(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	data, err := json.Marshal(s.items)
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func Load(path string) (Stack, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return New(), nil
		}
		return New(), err
	}
	var items []string
	if err := json.Unmarshal(data, &items); err != nil {
		return New(), err
	}
	return Stack{items: items}, nil
}