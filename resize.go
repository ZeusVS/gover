package main

import "golang.org/x/term"

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
