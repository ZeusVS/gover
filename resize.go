package main

import "golang.org/x/term"

func (ts *terminalSession) resize() {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	err := ts.GetCurrentSize()
	if err != nil {
		return
	}

	ts.refreshQueue()
}

func (ts *terminalSession) GetCurrentSize() (err error) {
	width, height, err := term.GetSize(ts.fdIn)
	if err != nil {
		return err
	}
	ts.width, ts.height = width, height
	return nil
}
