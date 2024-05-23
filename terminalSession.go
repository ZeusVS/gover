package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"golang.org/x/term"
)

const (
	// For now we will lock the framerate at 60fps
	framerate = 75
)

type terminalSession struct {
	mu     *sync.Mutex
	out    io.Writer
	buffer *bytes.Buffer
	ticker *time.Ticker
	done   chan struct{}

	originalState *term.State
	fdIn          int

	drawQueue    []drawInstruction
	cwd          string
	cwdFiles     []os.DirEntry
	previewLen   int
	selectionPos int
	width        int
	height       int

	mainOffset    int
	previewOffset int
}

type drawInstruction struct {
	x    int
	y    int
	line string
}

// Initialise the terminal screen
func StartTerminalSession() (terminalSession, error) {
	mu := &sync.Mutex{}
	out := os.Stdout
	buffer := new(bytes.Buffer)
	ticker := time.NewTicker(time.Millisecond * 1000 / framerate)
	done := make(chan struct{})

	fdIn := int(os.Stdin.Fd())
	// Put the terminal in raw mode and remember the original state
	originalState, err := term.MakeRaw(fdIn)
	if err != nil {
		return terminalSession{}, err
	}

	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return terminalSession{}, err
	}
	// Get the initial files in the current working directory
	cwdFiles, err := os.ReadDir(cwd)
	if err != nil {
		return terminalSession{}, err
	}

	ts := terminalSession{
		mu:     mu,
		out:    out,
		buffer: buffer,
		ticker: ticker,
		done:   done,

		originalState: originalState,
		fdIn:          fdIn,

		drawQueue:    []drawInstruction{},
		cwd:          cwd,
		selectionPos: 0,
		cwdFiles:     cwdFiles,

		mainOffset:    0,
		previewOffset: 0,
	}

	// Hide the cursor
	fmt.Fprint(ts.out, CSI+HideCursorSeq)
	// Enter the alt screen
	fmt.Fprint(ts.out, CSI+AltScreenSeq)

	// Set the initial state of the terminalSession
	ts.resize()

	return ts, nil
}

// Stop the session and return the terminal to it's initial state
func (ts *terminalSession) StopTerminalSession() {
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
