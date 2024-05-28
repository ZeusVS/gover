package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

func (ts *terminalSession) quit() {
	close(ts.done)
}

func (ts *terminalSession) moveUpSelection(n int) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

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

	// If dir is empty, don't move directory down
	if len(newFiles) == 0 {
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

	ts.previewOffsetV += n

	// Reset if we are beyond the end of the files
	// TODO: it seems there is always an extra empty line at the end, check it out
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

func (ts *terminalSession) open() {
	// This will open up a terminal session in a new terminal window
	// This is not ideal, but it's impossible to change the cwd of a parent process
	ts.mu.Lock()
	defer ts.mu.Unlock()

	selectionName := ts.cwdFiles[ts.selectionPos].Name()
	filePath := filepath.Join(ts.cwd, selectionName)

	// Get the default terminal
	terminal := os.Getenv("TERM")
	if terminal == "" {
		return
	}

	// If selection is a directory open a new terminal window in that directory
	if ts.cwdFiles[ts.selectionPos].IsDir() {
		os.Chdir(filePath)
		cmd := exec.Command(terminal)
		cmd.Run()
		return
	}

	// If selection is a valid utf8 encoded file open in default editor
	b, _ := os.ReadFile(filePath)
	fileContent := string(b)
	if utf8.ValidString(fileContent) {
		// Get the default editor
		editor := os.Getenv("EDITOR")
		if editor == "" {
			return
		}

		cmd := exec.Command(terminal, "-e", editor, filePath)
		cmd.Run()
		return
	}

	// Executables and links remain

	// Code for executables:
	// cmd := exec.Command(terminal, "-e", path)
	// cmd.Run()
}

func (ts *terminalSession) rename() {
	// TODO: Move the selectionPos to the correct location after rename
	ts.inputMode = true
	defer func() { ts.inputMode = false }()
	runeSlice := []rune{}
	for {
		ts.queueInputLine("Rename: " + string(runeSlice))
		ru := <-ts.inCh
		if ru == inputMap["escape"] {
			// Redraw the original bottomBar
			ts.queueBottomBar()
			return
		}
		if ru == inputMap["enter"] {
			break
		}
		if ru == inputMap["backspace"] {
			if len(runeSlice) == 0 {
				continue
			}
			runeSlice = runeSlice[:len(runeSlice)-1]
		} else {
			runeSlice = append(runeSlice, ru)
		}
	}
	name := string(runeSlice)

	selectionName := ts.cwdFiles[ts.selectionPos].Name()
	oldPath := filepath.Join(ts.cwd, selectionName)

	newPath := filepath.Join(ts.cwd, name)
	os.Rename(oldPath, newPath)

	// Refresh files
	cwdFiles, err := os.ReadDir(ts.cwd)
	if err != nil {
		return
	}
	ts.cwdFiles = cwdFiles
	ts.refreshQueue()
}

func (ts *terminalSession) search() {
	// TODO: Highlight all found search entries
	// TODO: Add regex to the search methods
	ts.inputMode = true
	defer func() { ts.inputMode = false }()
	runeSlice := []rune{}
	for {
		ts.queueInputLine("Search: " + string(runeSlice))
		ru := <-ts.inCh
		if ru == inputMap["escape"] {
			// Redraw the original bottomBar
			ts.queueBottomBar()
			// Reset the highlights
			ts.queueMainFiles()
			return
		}
		if ru == inputMap["enter"] {
			break
		}
		if ru == inputMap["backspace"] {
			if len(runeSlice) == 0 {
				continue
			}
			runeSlice = runeSlice[:len(runeSlice)-1]
		} else {
			runeSlice = append(runeSlice, ru)
		}

		// Reset the highlights
		ts.queueMainFiles()
		// Add highlights to all matching strings
		ts.searchStr = string(runeSlice)
		for i, file := range ts.cwdFiles {
			fileStr := strings.ToLower(file.Name())
			searchStr := strings.ToLower(ts.searchStr)
			if strI := strings.Index(fileStr, searchStr); strI >= 0 {
				line := file.Name()[strI : strI+len(searchStr)]
				line = StyleBgYellow + StyleFgBlack + line + StyleReset
				drawInstr := drawInstruction{
					// +4 because of spaces + icons
					x:    strI + 4,
					y:    i - ts.mainOffset,
					line: line,
				}
				ts.drawQueue = append(ts.drawQueue, drawInstr)
			}
		}
	}

	ts.searchN()
}

func (ts *terminalSession) searchN() {
	if ts.searchStr == "" {
		return
	}
	for i := ts.selectionPos + 1; i < len(ts.cwdFiles); i++ {
		// For now we ignore char casing
		fileStr := strings.ToLower(ts.cwdFiles[i].Name())
		searchStr := strings.ToLower(ts.searchStr)
		if strings.Contains(fileStr, searchStr) {
			ts.selectionPos = i
			break
		}
	}

	ts.refreshQueue()
}

func (ts *terminalSession) searchP() {
	if ts.searchStr == "" {
		return
	}
	for i := ts.selectionPos - 1; i >= 0; i-- {
		// For now we ignore char casing
		fileStr := strings.ToLower(ts.cwdFiles[i].Name())
		searchStr := strings.ToLower(ts.searchStr)
		if strings.Contains(fileStr, searchStr) {
			ts.selectionPos = i
			break
		}
	}

	ts.refreshQueue()
}
