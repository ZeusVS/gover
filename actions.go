package main

import (
	"os"
	"path/filepath"
)

func (ts *terminalSession) moveSelectionUp() {
	ts.selectionPos = ts.selectionPos - 1
	if ts.selectionPos < 0 {
		ts.selectionPos = 0
	}

	// TODO: make the program only redraw the 2 changed lines instead of entire screen
	ts.getFiles()
	ts.addFilesToQueue()
    ts.addBottomBarToQueue()
}

func (ts *terminalSession) moveSelectionDown() {
	ts.selectionPos = ts.selectionPos + 1
	if ts.selectionPos > len(ts.cwdFiles)-1 {
		ts.selectionPos = len(ts.cwdFiles) - 1
	}

	// TODO: make the program only redraw the 2 changed lines instead of entire screen
	ts.getFiles()
	ts.addFilesToQueue()
    ts.addBottomBarToQueue()
}

func (ts *terminalSession) moveUpDir() {
	_, fileName := filepath.Split(ts.cwd)
	ts.cwd = filepath.Dir(ts.cwd)
	ts.getFiles()
	// Get the index of the directory we just moved out of to select it
	for i, file := range ts.cwdFiles {
		if file.Name() == fileName {
			ts.selectionPos = i
		}
	}

	// If we mode up and down the dir faster than the refresh rate, we will clear
	// the screen and add files to the queue before the previous files were drawn
	// this will cause a weird merging of these two draws, should look into solving
	ts.clearScreen()
	ts.addFilesToQueue()
    ts.addBottomBarToQueue()
}

func (ts *terminalSession) moveDownDir() {
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
	ts.clearScreen()
	ts.addFilesToQueue()
    ts.addBottomBarToQueue()
}
