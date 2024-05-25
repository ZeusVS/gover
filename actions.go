package main

import (
	"os"
	"path/filepath"
)

func (ts *terminalSession) quit() {
	close(ts.done)
	// This does not close the for loop of the keylistener
	// But it does stop ts.startRendering() which in turn stops main()
}

func (ts *terminalSession) up() {
	ts.moveUpSelection(1)
}

func (ts *terminalSession) top() {
	ts.moveUpSelection(len(ts.cwdFiles))
}

func (ts *terminalSession) moveUpSelection(n int) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	ts.selectionPos -= n
	ts.previewOffset = 0

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

func (ts *terminalSession) down() {
	ts.moveDownSelection(1)
}

func (ts *terminalSession) bottom() {
	ts.moveDownSelection(len(ts.cwdFiles))
}

func (ts *terminalSession) moveDownSelection(n int) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	ts.selectionPos += n
	ts.previewOffset = 0

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

	ts.mainOffset = 0
	ts.previewOffset = 0

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

	ts.mainOffset = 0
	ts.previewOffset = 0

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

func (ts *terminalSession) scrollUpPreview() {
	ts.moveUpPreview(ts.height / 2)
}

func (ts *terminalSession) moveUpPreview(n int) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	ts.previewOffset -= n

	// Go back if we are before the beginning of the files
	if ts.previewOffset < 0 {
		ts.previewOffset = 0
	}

	ts.refreshQueue()
}

func (ts *terminalSession) scrollDownPreview() {
	ts.moveDownPreview(ts.height / 2)
}

func (ts *terminalSession) moveDownPreview(n int) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	ts.previewOffset += n

	// Go back if we are before the beginning of the files
	// TODO: it seems there is always an extra empty line at the end, check it out
	if ts.previewOffset > ts.previewLen-(ts.height-BottomRows) {
		ts.previewOffset = ts.previewLen - (ts.height - BottomRows)
	}

	ts.refreshQueue()
}
