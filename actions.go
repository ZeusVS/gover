package main

import (
	"os"
	"path/filepath"
)

func (ts *terminalSession) moveSelectionUp(n int) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	ts.command = ""
	ts.selectionPos -= n

	// Go back if we are before the beginning of the files
	if ts.selectionPos < 0 {
		ts.selectionPos = 0
	}

	// If the selection is outside of the range, set the offset
	if ts.selectionPos < ts.mainOffset {
		ts.mainOffset = ts.selectionPos
	}

	ts.refreshQueue()
}

func (ts *terminalSession) moveSelectionDown(n int) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	ts.command = ""
	ts.selectionPos += n

	// Go back if we are beyond the end of the files
	if ts.selectionPos > len(ts.cwdFiles)-1 {
		ts.selectionPos = len(ts.cwdFiles) - 1
	}

	// If the selection is outside of the range, set the offset
	if ts.selectionPos > ts.height+ts.mainOffset-1-BottomRows {
		ts.mainOffset = ts.selectionPos - ts.height + 1 + BottomRows
	}

	ts.refreshQueue()
}

func (ts *terminalSession) moveUpDir() {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	ts.command = ""
	ts.mainOffset = 0

	_, fileName := filepath.Split(ts.cwd)
	ts.cwd = filepath.Dir(ts.cwd)

	cwdFiles, err := os.ReadDir(ts.cwd)
	if err != nil {
		return
	}
	ts.cwdFiles = cwdFiles

	// Get the index of the directory we just moved out of to select it
	for i, file := range ts.cwdFiles {
		if file.Name() == fileName {
			ts.selectionPos = i
		}
	}

	if ts.selectionPos > ts.height+ts.mainOffset-1-BottomRows {
		ts.mainOffset += ts.selectionPos - (ts.height - 1 - BottomRows)
	}

	ts.refreshQueue()
}

func (ts *terminalSession) moveDownDir() {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	ts.command = ""
	ts.mainOffset = 0

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
