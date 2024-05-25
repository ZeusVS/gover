package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

const (
	CtrlU = 0x15
	CtrlD = 0x04
)

// Ideas:
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
			'q':   {callback: ts.quit},
			'k':   {callback: ts.up},
			'j':   {callback: ts.down},
			'G':   {callback: ts.bottom},
			'h':   {callback: ts.moveUpDir},
			'l':   {callback: ts.moveDownDir},
			CtrlU: {callback: ts.scrollUpPreview},
			CtrlD: {callback: ts.scrollDownPreview},
			'g': {
				subCommand: map[rune]command{
					'g': {callback: ts.top},
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
			command.callback()
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
