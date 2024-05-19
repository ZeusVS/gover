package main

import (
	"fmt"

	"golang.org/x/term"
)

func (ts *terminalSession) startRendering() {
	for {
		select {
		case <-ts.done:
			return
		case <-ts.ticker.C:
			ts.redraw()
		}
	}
}

func (ts *terminalSession) redraw() {
	// draw everything waiting in the queue to the screen
	ts.mu.Lock()
	defer ts.mu.Unlock()

	for pos, line := range ts.drawQueue {
		// The line with index 0 is drawn on position 1
		ts.moveCursorTo(1, pos+1)
		ts.eraseLine()
		ts.drawLine(line)

		delete(ts.drawQueue, pos)
	}
}

func (ts *terminalSession) drawLine(line string) {
	fmt.Fprint(ts.out, line)
}

// Unused for now
func (ts *terminalSession) clearScreen() {
	fmt.Fprint(ts.out, CSI+ClearScreenSeq)
}

func (ts *terminalSession) eraseLine() {
	fmt.Fprint(ts.out, CSI+EraseLineSeq)
}

func (ts *terminalSession) moveCursorTo(x, y int) {
	fmt.Fprintf(ts.out, CSI+MoseCursorToSeq, y, x)
}

func (ts *terminalSession) GetCurrentSize() (err error) {
	ts.width, ts.height, err = term.GetSize(ts.fdIn)
	return err
}
