// Package app contains application services: use cases over [todo.Item] and [todo.Repo].
package app

import (
	"errors"
	"fmt"

	"github.com/klusekP/tuitodo/internal/todo"
)

// Service coordinates task operations and persistence (Single Responsibility: task lifecycle + save).
type Service struct {
	repo  todo.Repo
	items []todo.Item
	// initErr is set if initial Load failed; caller may show it once.
	initErr error
}

// NewService loads tasks from repo. A load error is stored on the service and does
// not prevent using an empty list (mirrors lenient shell UX).
func NewService(r todo.Repo) *Service {
	s := &Service{repo: r}
	s.items, s.initErr = r.Load()
	if s.items == nil {
		s.items = []todo.Item{}
	}
	s.sortPendingFirst()
	return s
}

// InitError is the error from the initial [todo.Repo.Load], if any.
func (s *Service) InitError() error { return s.initErr }

// Items is the current list (same backing slice; do not mutate — use [Service] methods).
func (s *Service) Items() []todo.Item { return s.items }

// Len returns the number of tasks.
func (s *Service) Len() int { return len(s.items) }

// persist writes the current list to the repo.
func (s *Service) persist() error {
	if err := s.repo.Save(s.items); err != nil {
		return fmt.Errorf("save: %w", err)
	}
	return nil
}

// ErrInvalidIndex is returned when an index is out of range.
var ErrInvalidIndex = errors.New("invalid task index")

// Add appends a new task, moves it to the pending group (invariant: pending before done), and saves.
func (s *Service) Add(title string) (newIndex int, err error) {
	before := append([]todo.Item(nil), s.items...)
	newItem := todo.NewItem(title)
	s.items = append(s.items, newItem)
	s.sortPendingFirst()
	if err := s.persist(); err != nil {
		s.items = before
		return 0, err
	}
	return s.indexByID(newItem.ID), nil
}

// SetTitle updates the task title at i and saves.
func (s *Service) SetTitle(i int, title string) error {
	if i < 0 || i >= len(s.items) {
		return ErrInvalidIndex
	}
	prev := s.items[i].Title
	s.items[i].Title = title
	if err := s.persist(); err != nil {
		s.items[i].Title = prev
		return err
	}
	return nil
}

// Toggle flips the done flag at i, reorders so all pending items stay above completed ones, saves,
// and returns the task's new index (cursor should follow the same item).
func (s *Service) Toggle(i int) (done bool, newIndex int, err error) {
	if i < 0 || i >= len(s.items) {
		return false, 0, ErrInvalidIndex
	}
	before := append([]todo.Item(nil), s.items...)
	id := s.items[i].ID
	s.items[i].Done = !s.items[i].Done
	done = s.items[i].Done
	s.sortPendingFirst()
	if err := s.persist(); err != nil {
		s.items = before
		return false, 0, err
	}
	return done, s.indexByID(id), nil
}

// Remove deletes the task at i and saves.
func (s *Service) Remove(i int) (removedTitle string, err error) {
	if i < 0 || i >= len(s.items) {
		return "", ErrInvalidIndex
	}
	removedTitle = s.items[i].Title
	next := make([]todo.Item, 0, len(s.items)-1)
	next = append(next, s.items[:i]...)
	next = append(next, s.items[i+1:]...)
	if err := s.repo.Save(next); err != nil {
		return "", err
	}
	s.items = next
	return removedTitle, nil
}

// ClearDone removes all completed tasks and returns how many were removed.
func (s *Service) ClearDone() (cleared int, err error) {
	kept := make([]todo.Item, 0, len(s.items))
	for _, t := range s.items {
		if !t.Done {
			kept = append(kept, t)
		}
	}
	cleared = len(s.items) - len(kept)
	if cleared == 0 {
		return 0, nil
	}
	if err := s.repo.Save(kept); err != nil {
		return 0, err
	}
	s.items = kept
	return cleared, nil
}

// Counts returns (done, total) for the status line.
func (s *Service) Counts() (done, total int) {
	total = len(s.items)
	for _, t := range s.items {
		if t.Done {
			done++
		}
	}
	return done, total
}

// sortPendingFirst keeps relative order within each group: all pending first, then all completed.
func (s *Service) sortPendingFirst() {
	var pending, done []todo.Item
	for _, t := range s.items {
		if !t.Done {
			pending = append(pending, t)
		}
	}
	for _, t := range s.items {
		if t.Done {
			done = append(done, t)
		}
	}
	s.items = append(pending, done...)
}

func (s *Service) indexByID(id int64) int {
	for i, t := range s.items {
		if t.ID == id {
			return i
		}
	}
	return -1
}
