// Package config holds runtime paths and defaults.
package config

import (
	"os"
	"path/filepath"
)

// DefaultDataPath returns ~/.tuitodo.json or a local fallback if home is unavailable.
func DefaultDataPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "todos.json"
	}
	return filepath.Join(home, ".tuitodo.json")
}
