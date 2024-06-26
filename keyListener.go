package main

import "time"

// TODO:
// Command key ideas:
// ? show manual/all hotkeys

// <c-h> show/hide hidden files

// Will be hard to implement:
// u     undo
// <c-r> redo

// ?/? create new file/directory
// y copy
// d cut
// p paste

type command struct {
	callback   func()
	subCommand map[rune]command
}

func (ts *terminalSession) startKeyListener() {
	// Define the command tree
	ts.startCmd = command{
		subCommand: map[rune]command{
			// Actions
			// Quit gover
			'q': {callback: ts.quit},
			// Clear all actions
			inputMap["escape"]: {callback: ts.clearActions},
			// Show manual page
			'?': {callback: ts.showManual},
			// Custom terminal command
			':': {callback: ts.terminalCommand},
			// Open current selection
			inputMap["enter"]: {callback: ts.open},
			// Insert/create new file
			'i': {callback: ts.insertFile},
			// Insert/create new directory
			'I': {callback: ts.insertDir},
			// Copy the current selection
			'y': {callback: ts.copy},
			// Cut (move) the current selection
			'd': {callback: ts.cut},
			// Paste the current content of the copy/cut buffer
			'p': {callback: ts.paste},
			// Delete the current selection
			'D': {callback: ts.delete},
			// Rename current selection
			'R': {callback: ts.rename},
			// Search the current files
			'/': {callback: ts.search},
			// Search next
			'n': {callback: ts.searchN},
			// Search previous
			'N': {callback: ts.searchP},

			// Motions
			// Move to home directory
			'~': {callback: ts.goHome},
			// Go up 1 on main panel
			'k': {callback: func() { ts.moveUpSelection(1) }},
			// Go up 10 on main panel
			'K': {callback: func() { ts.moveUpSelection(10) }},
			// Go down 1 on main panel
			'j': {callback: func() { ts.moveDownSelection(1) }},
			// Go down 10 on main panel
			'J': {callback: func() { ts.moveDownSelection(10) }},
			// Go to bottom on main panel
			'G': {callback: func() { ts.moveDownSelection(len(ts.cwdFiles)) }},
			// Go up a directiry level
			'h': {callback: ts.moveUpDir},
			// Go down a directory level
			'l': {callback: ts.moveDownDir},
			// Scroll up preview panel half a page
			inputMap["<c-u>"]: {callback: func() { ts.moveUpPreview(ts.height / 2) }},
			// Scroll down preview panel half a page
			inputMap["<c-d>"]: {callback: func() { ts.moveDownPreview(ts.height / 2) }},
			// Scroll left preview panel half a page
			inputMap["<c-f>"]: {callback: func() { ts.moveLeftPreview(ts.width / 4) }},
			// Scroll right preview panel half a page
			inputMap["<c-k>"]: {callback: func() { ts.moveRightPreview(ts.width / 4) }},

			// Multi-char commands
			'g': {
				subCommand: map[rune]command{
					// Go to top on main panel
					'g': {callback: func() { ts.moveUpSelection(len(ts.cwdFiles)) }},
				},
			},

			// Sort commands
			's': {
				subCommand: map[rune]command{
					// Sort directories first
					'd': {callback: func() {
						ts.sortFunc = sortDirsFirst
						ts.refreshQueue()
					}},
					// Sort directories last
					'D': {callback: func() {
						ts.sortFunc = sortDirsLast
						ts.refreshQueue()
					}},
					// Sort alphabetically
					'a': {callback: func() {
						ts.sortFunc = sortAlpha
						ts.refreshQueue()
					}},
					// Sort alphabetically reversed
					'A': {callback: func() {
						ts.sortFunc = sortAlphaReverse
						ts.refreshQueue()
					}},
					// Sort by modification time, newest first
					't': {callback: func() {
						ts.sortFunc = sortDateNewest
						ts.refreshQueue()
					}},
					// Sort by modification time, oldest first
					'T': {callback: func() {
						ts.sortFunc = sortDateOldest
						ts.refreshQueue()
					}},
					// Sort by filesize, smallest first
					's': {callback: func() {
						ts.sortFunc = sortSizeSmallest
						ts.refreshQueue()
					}},
					// Sort by filesize, largest first
					'S': {callback: func() {
						ts.sortFunc = sortSizeLargest
						ts.refreshQueue()
					}},
				},
			},
		},
	}

	// Set the current command to the start command
	ts.curCmd = ts.startCmd

	// Make a channel that receives all runes read from Stdin
	go func() {
		for {
			ru, _, err := ts.in.ReadRune()
			if err != nil {
				continue
			}
			ts.inCh <- ru
		}
	}()

	// Loop where we check if ts.done is closed or a rune is read
	for {
		// Don't read from the channel if we are in inputMode
		if !ts.inputMode {
			select {
			case <-ts.done:
				return
			case ru := <-ts.inCh:
				ts.getCommand(ru)
			}
		}
	}
}

func (ts *terminalSession) getCommand(ru rune) {
	// If the input rune results in a nonexistant command we reset the command
	command, ok := ts.curCmd.subCommand[ru]
	if !ok {
		ts.curCmd = ts.startCmd
		ts.cmdStr = ""
		go func() {
			time.Sleep(time.Millisecond * 50)
			ts.clearCommand()
		}()
		return
	}

	// Add the string to the "commandline"
	cmdStr := string(ru)
	for keyStr, valRu := range inputMap {
		if valRu == ru {
			cmdStr = keyStr
			break
		}
	}

	ts.cmdStr += cmdStr
	ts.queueBottomBar()

	// If the command has a callback function we call it
	if command.callback != nil {
		callBackFunc := command.callback
		callBackFunc()
		ts.curCmd = ts.startCmd
		ts.cmdStr = ""
		// Make pressed commands display for 50ms before getting wiped
		// TODO: Input commands are blocking so these do not get wiped, fix
		go func() {
			time.Sleep(time.Millisecond * 50)
			ts.clearCommand()
		}()
		return
	}

	// Otherwise we will go to the remaining subcommands
	ts.curCmd = command
}
