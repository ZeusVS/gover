package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"golang.org/x/term"
)

const (
	refreshrate = 75
)

type terminalSession struct {
	mu     *sync.Mutex
	in     *bufio.Reader
	out    io.Writer
	buffer *bytes.Buffer
	ticker *time.Ticker
	done   chan struct{}
	inCh   chan rune

	originalState *term.State
	fdIn          int

	startCmd  command
	curCmd    command
	cmdStr    string
	searchStr string
	inputMode bool

	drawQueue    []drawInstruction
	cwd          string
	cwdFiles     []os.DirEntry
	previewLen   int
	selectionPos int
	width        int
	height       int

	mainOffset     int
	previewOffsetV int
	previewOffsetH int
}

type drawInstruction struct {
	x    int
	y    int
	line string
}

// Initialise the terminal screen
func StartTerminalSession() (terminalSession, error) {
	mu := &sync.Mutex{}
	in := bufio.NewReader(os.Stdin)
	out := os.Stdout
	buffer := new(bytes.Buffer)
	ticker := time.NewTicker(time.Millisecond * 1000 / refreshrate)
	done := make(chan struct{})
	inCh := make(chan rune)
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
		in:     in,
		out:    out,
		buffer: buffer,
		ticker: ticker,
		done:   done,
		inCh:   inCh,

		originalState: originalState,
		fdIn:          fdIn,

		// Initialise these variables to the default state
		cmdStr:    "",
		searchStr: "",
		inputMode: false,

		drawQueue:    []drawInstruction{},
		cwd:          cwd,
		selectionPos: 0,
		cwdFiles:     cwdFiles,

		mainOffset:     0,
		previewOffsetV: 0,
		previewOffsetH: 0,
	}

	// Hide the cursor
	fmt.Fprint(ts.out, CSI+HideCursorSeq)
	// Enter the alt screen
	fmt.Fprint(ts.out, CSI+AltScreenSeq)

	// Set the initial size of the terminalSession
	// This function also adds the initial state to the drawQueue
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
		fmt.Fprintf(os.Stderr, "Error restoring terminal's initial state: %s", err)
	}
}
