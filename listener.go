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
	test = 'f'

	up      = 'k'
	down    = 'j'
	dirup   = 'h'
	dirdown = 'l'
)

func (ts *terminalSession) startListening() {
    go ts.startKeyListener()
    go ts.startResizeListener()
}

// TODO: Add multi char inputs
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

		case ru == up:
			ts.moveSelectionUp()
		case ru == down:
			ts.moveSelectionDown()
		case ru == dirup:
			ts.moveUpDir()
		case ru == dirdown:
			ts.moveDownDir()

		case ru == test:
            // Nothing to test atm
		}
	}
}

func (ts *terminalSession) startResizeListener() {
    sigc := make(chan os.Signal, 1)
    signal.Notify(sigc, syscall.SIGWINCH)

    for {
        select {
        case <- ts.done:
            return
        case <- sigc:
            err := ts.GetCurrentSize()
            if err != nil {
                continue
            }
            err = ts.getFiles()
            if err != nil {
                // Better error handling needed probably
                continue
            }
            ts.addFilesToQueue()
        }
    }
    
}
