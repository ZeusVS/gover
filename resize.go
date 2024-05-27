package main

import (
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/term"
)

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

func (ts *terminalSession) resize() {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	width, height, err := term.GetSize(ts.fdIn)
	if err != nil {
		return
	}

	ts.width = width
	ts.height = height

	ts.refreshQueue()
}
