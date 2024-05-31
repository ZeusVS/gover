package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

func (ts *terminalSession) quit() {
	close(ts.done)
}

func (ts *terminalSession) open() {
	// This will open up a terminal session in a new terminal window
	// This is not ideal, but it's impossible to change the cwd of a parent process
	ts.mu.Lock()
	defer ts.mu.Unlock()

	selection, err := ts.cwdFiles[ts.selectionPos].Info()
	if err != nil {
		return
	}
	filePath := filepath.Join(ts.cwd, selection.Name())

	// Get the default terminal
	terminal := os.Getenv("TERM")
	if terminal == "" {
		return
	}

	if selection.Mode()&os.ModeSymlink != 0 {
		link, err := filepath.EvalSymlinks(filePath)
		if err == nil {
			linkInfo, err := os.Stat(link)
			if err == nil {
				selection = linkInfo
				filePath = link
			}
		}
	}

	// If selection is a directory open a new terminal window in that directory
	if selection.IsDir() {
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

	// Code for executables:
	cmd := exec.Command(terminal, "-e", filePath)
	cmd.Run()
}

func (ts *terminalSession) delete() {
	// TODO: Move the selectionPos to the correct location after rename
	ts.inputMode = true
	defer func() { ts.inputMode = false }()
	runeSlice := []rune{}

	selectionName := ts.cwdFiles[ts.selectionPos].Name()

	confirmText := "Confirm deletion of \""
	if ts.cwdFiles[ts.selectionPos].IsDir() {
		confirmText = "Confirm recursive deletion of "
	}

	keys := "\" [y(es), n(o)]: "

	for {
		for {
			ts.queueInputLine(confirmText + selectionName + keys + string(runeSlice))
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
		cmd := string(runeSlice)

		// Stay in the loop untill we say y/yes/n/no
		if strings.ToLower(cmd) == "y" || strings.ToLower(cmd) == "yes" {
			break
		} else if strings.ToLower(cmd) == "n" || strings.ToLower(cmd) == "no" {
			ts.queueBottomBar()
			return
		}
		runeSlice = nil
	}

	selectionPath := filepath.Join(ts.cwd, selectionName)

	err := os.RemoveAll(selectionPath)
	if err != nil {
		return
	}

	// Refresh files
	cwdFiles, err := os.ReadDir(ts.cwd)
	if err != nil {
		return
	}
	ts.cwdFiles = ts.sortFunc(cwdFiles)
	ts.refreshQueue()
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
	ts.cwdFiles = ts.sortFunc(cwdFiles)
	ts.refreshQueue()
}

func (ts *terminalSession) search() {
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
				// Get the part of the string from the original filename to make cases match
				line := file.Name()[strI : strI+len(searchStr)]
				// Style the highlight
				// TODO: Change the yellow background?
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
	found := false
	// first search after the current selection
	for i := ts.selectionPos + 1; i < len(ts.cwdFiles); i++ {
		// For now we ignore char casing
		fileStr := strings.ToLower(ts.cwdFiles[i].Name())
		searchStr := strings.ToLower(ts.searchStr)
		if strings.Contains(fileStr, searchStr) {
			ts.selectionPos = i
			found = true
			break
		}
	}

	if found == false {
		// then search from beginning up to current selection
		for i := 0; i <= ts.selectionPos; i++ {
			// For now we ignore char casing
			fileStr := strings.ToLower(ts.cwdFiles[i].Name())
			searchStr := strings.ToLower(ts.searchStr)
			if strings.Contains(fileStr, searchStr) {
				ts.selectionPos = i
				found = true
				break
			}
		}
	}

	if found == false {
		line := StyleFgRed + "Pattern not found: " + ts.searchStr + StyleReset
		ts.queueInputLine(line)
		return
	}

	ts.refreshQueue()
	ts.queueInputLine("Search: " + ts.searchStr)
}

func (ts *terminalSession) searchP() {
	if ts.searchStr == "" {
		return
	}

	found := false

	for i := ts.selectionPos - 1; i >= 0; i-- {
		// For now we ignore char casing
		fileStr := strings.ToLower(ts.cwdFiles[i].Name())
		searchStr := strings.ToLower(ts.searchStr)
		if strings.Contains(fileStr, searchStr) {
			ts.selectionPos = i
			found = true
			break
		}
	}

	if found == false {
		for i := len(ts.cwdFiles) - 1; i >= ts.selectionPos; i-- {
			// For now we ignore char casing
			fileStr := strings.ToLower(ts.cwdFiles[i].Name())
			searchStr := strings.ToLower(ts.searchStr)
			if strings.Contains(fileStr, searchStr) {
				ts.selectionPos = i
				found = true
				break
			}
		}
	}

	if found == false {
		line := StyleFgRed + "Pattern not found: " + ts.searchStr + StyleReset
		ts.queueInputLine(line)
		return
	}

	ts.refreshQueue()
	ts.queueInputLine("Search: " + ts.searchStr)
}

func (ts *terminalSession) terminalCommand() {
	ts.inputMode = true
	defer func() { ts.inputMode = false }()

	// Get the default terminal
	terminal := os.Getenv("TERM")
	if terminal == "" {
		return
	}
	os.Chdir(ts.cwd)

	cwd := ts.cwd
	homeDir, err := os.UserHomeDir()
	var cut bool
	// Only try to cut prefix if a homeDir was found
	if err == nil {
		cwd, cut = strings.CutPrefix(cwd, homeDir)
		// Only add tilde if a prefix was cut
		if cut == true {
			cwd = "~" + cwd
		}
	}

	consoleString := StyleFgBlue + cwd + StyleFgGreen + " ❱ " + StyleReset
	runeSlice := []rune{}
	for {
		ts.queueInputLine(consoleString + string(runeSlice))
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
	}

	splitCmd := strings.Split(string(runeSlice), " ")
	name := splitCmd[0]
	args := []string{}
	if len(splitCmd) > 1 {
		args = splitCmd[1:]
	}

	cmd := exec.Command(name, args...)
	err = cmd.Run()
	if err != nil {
		// If command fails only redraw bottom bar
		ts.queueBottomBar()
		fmt.Fprintf(os.Stderr, "Error executing command: %v", err)
		return
	}

	// Refresh files
	cwdFiles, err := os.ReadDir(ts.cwd)
	if err != nil {
		return
	}
	ts.cwdFiles = ts.sortFunc(cwdFiles)

	ts.refreshQueue()
}
