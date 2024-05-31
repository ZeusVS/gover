package main

import (
	"os"
	"path/filepath"
)

func (ts *terminalSession) moveUpSelection(n int) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	if len(ts.cwdFiles) < 1 {
		return
	}

	ts.selectionPos -= n
	ts.previewOffsetV = 0
	ts.previewOffsetH = 0

	// Reset if we are before the beginning of the files
	if ts.selectionPos < 0 {
		ts.selectionPos = 0
	}

	// If the selection is outside of the range, adjust the offset
	if ts.selectionPos < ts.mainOffset {
		ts.mainOffset = ts.selectionPos
	}

	ts.refreshQueue()
}

func (ts *terminalSession) moveDownSelection(n int) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	if len(ts.cwdFiles) < 1 {
		return
	}

	ts.selectionPos += n
	ts.previewOffsetV = 0
	ts.previewOffsetH = 0

	// Reset if we are beyond the end of the files
	if ts.selectionPos > len(ts.cwdFiles)-1 {
		ts.selectionPos = len(ts.cwdFiles) - 1
	}

	// If the selection is outside of the range, adjust the offset
	if ts.selectionPos > ts.height+ts.mainOffset-1-BottomRows {
		ts.mainOffset = ts.selectionPos - ts.height + 1 + BottomRows
	}

	ts.refreshQueue()
}

func (ts *terminalSession) moveUpDir() {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	ts.mainOffset = 0
	ts.previewOffsetV = 0
	ts.previewOffsetH = 0

	_, fileName := filepath.Split(ts.cwd)
	ts.cwd = filepath.Dir(ts.cwd)

	cwdFiles, err := os.ReadDir(ts.cwd)
	if err != nil {
		return
	}
	ts.cwdFiles = ts.sortFunc(cwdFiles)

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

	// If the current directory has no files do nothing
	if len(ts.cwdFiles) == 0 {
		return
	}

	ts.mainOffset = 0
	ts.previewOffsetV = 0
	ts.previewOffsetH = 0

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

	ts.cwd = newDir
	ts.cwdFiles = newFiles
	ts.selectionPos = 0

	ts.refreshQueue()
}

func (ts *terminalSession) moveUpPreview(n int) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	ts.previewOffsetV -= n

	// Reset if we are before the beginning of the files
	if ts.previewOffsetV < 0 {
		ts.previewOffsetV = 0
	}

	ts.refreshQueue()
}

func (ts *terminalSession) moveDownPreview(n int) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	if ts.previewLen <= ts.height-BottomRows {
		return
	}

	ts.previewOffsetV += n

	// Reset if we are beyond the end of the files
	// TODO: it seems there is always an extra empty line at the end, check it out
	// Probably to do with splitting on \n and an empty last string
	if ts.previewOffsetV > ts.previewLen-(ts.height-BottomRows) {
		ts.previewOffsetV = ts.previewLen - (ts.height - BottomRows)
	}

	ts.refreshQueue()
}

func (ts *terminalSession) moveLeftPreview(n int) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	ts.previewOffsetH -= n

	// Reset if we are before the first char
	if ts.previewOffsetH < 0 {
		ts.previewOffsetH = 0
	}

	ts.refreshQueue()
}

func (ts *terminalSession) moveRightPreview(n int) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	ts.previewOffsetH += n

	// TODO: Reset if we are after the last char of the file
	// This might be a bit convoluted with the way the code is built atm

	ts.refreshQueue()
}

func (ts *terminalSession) goHome() {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	home, err := os.UserHomeDir()
	// If homedir is not found, just return
	if err != nil {
		return
	}

	ts.selectionPos = 0
	ts.mainOffset = 0
	ts.previewOffsetV = 0
	ts.previewOffsetH = 0
	ts.cwd = home
	// No error handling, needs to change?
	ts.cwdFiles, _ = os.ReadDir(home)

	ts.refreshQueue()
}
