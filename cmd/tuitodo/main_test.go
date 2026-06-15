package main

import (
	"errors"
	"testing"

	"github.com/klusekP/tuitodo/internal/tui"
)

func TestRun_Success(t *testing.T) {
	orig := runProgram
	runProgram = func(_ tui.Model) error { return nil }
	t.Cleanup(func() { runProgram = orig })

	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	if err := run(); err != nil {
		t.Fatalf("run() = %v, want nil", err)
	}
}

func TestRun_ProgramError(t *testing.T) {
	want := errors.New("bubbletea error")
	orig := runProgram
	runProgram = func(_ tui.Model) error { return want }
	t.Cleanup(func() { runProgram = orig })

	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	if err := run(); err != want {
		t.Fatalf("run() = %v, want %v", err, want)
	}
}
