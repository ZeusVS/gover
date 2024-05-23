package main

import (
	"os"
	"path/filepath"
)

func (ts *terminalSession) moveSelectionUp() {
    ts.mu.Lock()
    defer ts.mu.Unlock()

	ts.selectionPos = ts.selectionPos - 1
	if ts.selectionPos < 0 {
		ts.selectionPos = 0
        return
	}

	ts.refreshQueue()
}

func (ts *terminalSession) moveSelectionDown() {
    ts.mu.Lock()
    defer ts.mu.Unlock()

	ts.selectionPos = ts.selectionPos + 1
	if ts.selectionPos > len(ts.cwdFiles)-1 {
		ts.selectionPos = len(ts.cwdFiles) - 1
        return
	}

	ts.refreshQueue()
}

func (ts *terminalSession) moveUpDir() {
    ts.mu.Lock()
    defer ts.mu.Unlock()

	_, fileName := filepath.Split(ts.cwd)
	ts.cwd = filepath.Dir(ts.cwd)

    cwdFiles, err := os.ReadDir(ts.cwd)
    if err != nil {
        ts.mu.Unlock()
        return
    }
    ts.cwdFiles = cwdFiles

	// Get the index of the directory we just moved out of to select it
	for i, file := range ts.cwdFiles {
		if file.Name() == fileName {
			ts.selectionPos = i
		}
	}

	ts.refreshQueue()
}

func (ts *terminalSession) moveDownDir() {
    ts.mu.Lock()
    defer ts.mu.Unlock()

	// If the selection isn't a directory do nothing
	// Should I add symlink directories as well?
	if !ts.cwdFiles[ts.selectionPos].IsDir() {
		return
	}

	// Get the full path
	fileName := ts.cwdFiles[ts.selectionPos].Name()
	newDir := filepath.Join(ts.cwd, fileName)
	newFiles, err := os.ReadDir(newDir)
	if err != nil {
		// Better error handling?
		return
	}

	// If dir is empty, don't move directory down
	if len(newFiles) == 0 {
		return
	}

	ts.cwd = newDir
	ts.cwdFiles = newFiles
	ts.selectionPos = 0

	ts.refreshQueue()
}
