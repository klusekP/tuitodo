package app

import (
	"errors"
	"testing"

	"github.com/klusekP/tuitodo/internal/todo"
)

type fakeRepo struct {
	items   []todo.Item
	saveErr error
}

func (f *fakeRepo) Load() ([]todo.Item, error) {
	return append([]todo.Item(nil), f.items...), nil
}

func (f *fakeRepo) Save(items []todo.Item) error {
	if f.saveErr != nil {
		return f.saveErr
	}
	f.items = append([]todo.Item(nil), items...)
	return nil
}

func TestServiceAddAndToggle(t *testing.T) {
	r := &fakeRepo{items: []todo.Item{}}
	s := NewService(r)
	if _, err := s.Add("buy milk"); err != nil {
		t.Fatal(err)
	}
	if s.Len() != 1 {
		t.Fatalf("len = %d", s.Len())
	}
	done, err := s.Toggle(0)
	if err != nil || !done {
		t.Fatalf("toggle: done=%v err=%v", done, err)
	}
	if len(r.items) != 1 || !r.items[0].Done {
		t.Fatalf("repo: %+v", r.items)
	}
}

func TestServiceSaveErrorRollsBackAdd(t *testing.T) {
	r := &fakeRepo{items: []todo.Item{}, saveErr: errors.New("disk full")}
	s := NewService(r)
	if _, err := s.Add("x"); err == nil {
		t.Fatal("expected error")
	}
	if s.Len() != 0 {
		t.Fatalf("expected rollback, len=%d", s.Len())
	}
}

func TestServiceClearDone(t *testing.T) {
	r := &fakeRepo{items: []todo.Item{
		{ID: 1, Title: "a", Done: true},
		{ID: 2, Title: "b", Done: false},
	}}
	s := NewService(r)
	n, err := s.ClearDone()
	if err != nil || n != 1 {
		t.Fatalf("n=%d err=%v", n, err)
	}
	if s.Len() != 1 || s.Items()[0].Title != "b" {
		t.Fatal("unexpected items")
	}
}
