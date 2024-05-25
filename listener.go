package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

const (
	exit = 'q'

	up      = 'k'
	down    = 'j'
	dirUp   = 'h'
	dirDown = 'l'

	goTo     = 'g'
	goBottom = 'G'

	// Scroll preview up
	CtrlU = 0x15
	// Scroll preview down
	CtrlD = 0x04
	// TODO: add consts

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
)

func (ts *terminalSession) startListening() {
	go ts.startKeyListener()
	go ts.startResizeListener()
}

func (ts *terminalSession) startKeyListener() {
	r := bufio.NewReader(os.Stdin)
	for {
		ru, _, err := r.ReadRune()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: reading key from Stdin: %s\r\n", err)
		}

		switch {
		case ru == exit:
			close(ts.done)
			return

		// Main panel move selection up
		case ru == up:
			ts.moveSelectionUp(1)
		// Main panel move selection down
		case ru == down:
			ts.moveSelectionDown(1)
		// Main panel go dir level higher
		case ru == dirUp:
			ts.moveUpDir()
		// Main panel go dir level lower
		case ru == dirDown:
			ts.moveDownDir()
		case ru == goTo:
			// Main panel go to top
			if ts.command == "g" {
				ts.moveSelectionUp(len(ts.cwdFiles))
			} else if ts.command == "" {
				ts.command = "g"
			} else {
				ts.command = ""
			}
		// Main panel scroll to bottom
		case ru == goBottom:
			ts.moveSelectionDown(len(ts.cwdFiles))
		case ru == CtrlD:
			// Scroll half a page
			n := ts.height / 2
			ts.movePreviewDown(n)
		case ru == CtrlU:
			// Scroll half a page
			n := ts.height / 2
			ts.movePreviewUp(n)
		}
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
