// Package jsonfile implements [todo.Repo] with a JSON file on disk.
package jsonfile

import (
	"encoding/json"
	"os"

	"github.com/klusekP/tuitodo/internal/todo"
)

// Repo reads and writes the task list as formatted JSON.
type Repo struct {
	path string
}

// New returns a [Repo] that uses path as the JSON file.
func New(path string) *Repo {
	return &Repo{path: path}
}

// Path returns the backing file path.
func (r *Repo) Path() string { return r.path }

// Load implements [todo.Repo]. Missing or empty file yields an empty slice, no error.
func (r *Repo) Load() ([]todo.Item, error) {
	data, err := os.ReadFile(r.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	if len(data) == 0 {
		return nil, nil
	}
	var items []todo.Item
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, err
	}
	return items, nil
}

// Save implements [todo.Repo].
func (r *Repo) Save(items []todo.Item) error {
	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(r.path, data, 0o644)
}
