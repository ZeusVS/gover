package main

import (
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
}

func (ts *terminalSession) moveSelectionDown() {
    ts.selectionPos = ts.selectionPos + 1
    if ts.selectionPos > len(ts.cwdFiles) - 1 {
        ts.selectionPos = len(ts.cwdFiles) - 1
    }

    // TODO: make the program only redraw the 2 changed lines instead of entire screen
    ts.getFiles()
    ts.addFilesToQueue()
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

    ts.clearScreen()
    ts.addFilesToQueue()
}

func (ts *terminalSession) moveDownDir() {
    // If the selection isn't a directory do nothing
    if !ts.cwdFiles[ts.selectionPos].IsDir() {
        return
    }

    // TODO: prevent going down if the directory is empty

    // Get the full path
    fileName := ts.cwdFiles[ts.selectionPos].Name()
    ts.cwd = filepath.Join(ts.cwd, fileName)

    ts.selectionPos = 0
    ts.clearScreen()
    ts.getFiles()
    ts.addFilesToQueue()
}

