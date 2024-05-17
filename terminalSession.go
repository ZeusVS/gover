package main

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"golang.org/x/term"
)

const (
	// For now we will lock the framerate at 60fps
	framerate = 60
)

type terminalSession struct {
	mu     *sync.Mutex
	out    io.Writer
	ticker *time.Ticker
	done   chan bool

	originalState *term.State
	fdIn          int

	drawQueue map[int]string
}

// Initialise the terminal screen
func StartTerminalSession() (terminalSession, error) {
	// Get the terminal input file descriptor
	fdIn := int(os.Stdin.Fd())

    // Create a new mutex
    mu := &sync.Mutex{}

	// Put the terminal in raw mode
	originalState, err := term.MakeRaw(fdIn)
	if err != nil {
		return terminalSession{}, err
	}

	ticker := time.NewTicker(time.Millisecond * 1000 / framerate)
	done := make(chan bool)

	ts := terminalSession{
		mu:     mu,
		out:    os.Stdout,
		ticker: ticker,
		done:   done,

		originalState: originalState,
		fdIn:          fdIn,

        drawQueue: make(map[int]string),
	}

	// Hide the cursor
	fmt.Fprint(ts.out, CSI+HideCursorSeq)
	// Enter the alt screen
	fmt.Fprint(ts.out, CSI+AltScreenSeq)

	return ts, nil
}

// Stop the session and return the terminal to it's initial state
func (ts terminalSession) StopTerminalSession() {
	// Exit the alt screen
	fmt.Fprint(ts.out, CSI+ExitAltScreenSeq)
	// Show the cursor again
	fmt.Fprint(ts.out, CSI+ShowCursorSeq)

	// Restore original state
	err := term.Restore(ts.fdIn, ts.originalState)
	if err != nil {
		fmt.Printf("Error restoring terminal's initial state: %s", err)
	}
}
