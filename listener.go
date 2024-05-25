package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"syscall"
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
// ? reverse search
// n search next
// N search previous

type command struct {
	callback   func()
	subCommand map[rune]command
}

func (ts *terminalSession) startListening() {
	go ts.startKeyListener()
	go ts.startResizeListener()
}

func (ts *terminalSession) startKeyListener() {
	// Define the command that contains all other commands
	startCommand := command{
		subCommand: map[rune]command{
			// Quit gover
			'q': {callback: ts.quit},

			// Go up 1 on main panel
			'k': {callback: func() { ts.moveUpSelection(1) }},
			// Go up 10 on main panel
			'K': {callback: func() { ts.moveUpSelection(10) }},
			// Go down 1 on main panel
			'j': {callback: func() { ts.moveDownSelection(1) }},
			// Go down 10 on main panel
			'J': {callback: func() { ts.moveDownSelection(1) }},
			// Go to bottom on main panel
			'G': {callback: func() { ts.moveDownSelection(len(ts.cwdFiles)) }},
			// Go up a directiry level
			'h': {callback: ts.moveUpDir},
			// Go down a directory level
			'l': {callback: ts.moveDownDir},
			// Scroll up selection panel half a page
			CtrlU: {callback: func() { ts.moveUpSelection(ts.height / 2) }},
			// Scroll down selection panel half a page
			CtrlD: {callback: func() { ts.moveDownSelection(ts.height / 2) }},

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

	r := bufio.NewReader(os.Stdin)
	for {
		ru, _, err := r.ReadRune()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: reading key from Stdin: %s\r\n", err)
		}

		// If the input rune results in a nonexistant command we reset the command
		command, ok := ts.command.subCommand[ru]
		if !ok {
			ts.command = startCommand
			continue
		}

		// If the command has a callback function we call it
		if command.callback != nil {
			callBackFunc := command.callback
			callBackFunc()
			continue
		}

		// Otherwise we will go to the remaining subcommands
		ts.command = command
	}
}

func (ts *terminalSession) startResizeListener() {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGWINCH)

	for {
		select {
		case <-ts.done:
			return
		case <-sigc:
			ts.resize()
		}
	}
}
