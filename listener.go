package main

import (
	"bufio"
	"os"
)

// Ideas:
// ? show manual

// s + ...
// d/D sort dirs first/last (default)
// a/A sort alphabetically
// n/N sort files by last open? date
// s/S sort files by size
// ???

// S + ...
// h show/hide hidden files

// y copy
// d cut
// p paste
// r rename
// i new ???

// / search
// n search next
// N search previous

// ? show all commands

type command struct {
	callback   func()
	subCommand map[rune]command
}

func (ts *terminalSession) startKeyListener() {
	// Define the command tree
	startCommand := command{
		subCommand: map[rune]command{
			// Actions
			// Quit gover
			'q': {callback: ts.quit},
			// Open current selection
			inputMap["enter"]: {callback: ts.open},

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
			inputMap["ctrl-u"]: {callback: func() { ts.moveUpPreview(ts.height / 2) }},
			// Scroll down preview panel half a page
			inputMap["ctrl-d"]: {callback: func() { ts.moveDownPreview(ts.height / 2) }},
			// Scroll left preview panel half a page
			inputMap["ctrl-f"]: {callback: func() { ts.moveLeftPreview(ts.width / 4) }},
			// Scroll right preview panel half a page
			inputMap["ctrl-k"]: {callback: func() { ts.moveRightPreview(ts.width / 4) }},

			// Multiline command starting with 'g'
			'g': {
				subCommand: map[rune]command{
					// Go to top on main panel
					'g': {callback: func() { ts.moveUpSelection(len(ts.cwdFiles)) }},
				},
			},
		},
	}

	// Set the current command to the start command
	ts.command = startCommand

	// Make a channel that recieves all runes read from Stdin
	r := bufio.NewReader(os.Stdin)
	readCh := make(chan (rune))
	readCh = runeReader(r)

	// Loop where we check if ts.done is closed or a rune is read
	for {
		select {
		case <-ts.done:
			return
		case ru := <-readCh:
			ts.getCommand(ru, startCommand)
		}
	}
}

func (ts *terminalSession) getCommand(ru rune, startCommand command) {
	// If the input rune results in a nonexistant command we reset the command
	command, ok := ts.command.subCommand[ru]
	if !ok {
		ts.command = startCommand
		return
	}

	// If the command has a callback function we call it
	if command.callback != nil {
		callBackFunc := command.callback
		callBackFunc()
		return
	}

	// Otherwise we will go to the remaining subcommands
	ts.command = command
}

// Function that makes a runereader goroutine
func runeReader(r *bufio.Reader) chan rune {
	ch := make(chan (rune))
	go func() {
		for {
			for {
				ru, _, err := r.ReadRune()
				if err != nil {
					continue
				}
				ch <- ru
			}
		}
	}()

	return ch
}
