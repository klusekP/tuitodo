// Package todo defines the task entity and related factories.
package todo

import "time"

// Item is a single todo entry (persistence-agnostic).
type Item struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
}

// NewItem creates a new task with a fresh ID and creation timestamp.
func NewItem(title string) Item {
	return Item{
		ID:        time.Now().UnixNano(),
		Title:     title,
		CreatedAt: time.Now().UTC(),
	}
}
