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
			ts.render()
		}
	}
}

func (ts *terminalSession) render() {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	// draw everything waiting in the queue to the screen
	for pos, line := range ts.drawQueue {
		// The line with index 0 is drawn on position 1
		ts.moveCursorTo(1, pos+1)
		ts.eraseLine()
		ts.drawLine(line)

		delete(ts.drawQueue, pos)
	}
	// This gets called every frame, hugely unnececary
	// Has to change
	ts.drawScrollbars()
}

func (ts *terminalSession) emptyDrawQueue() {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	for i := range ts.drawQueue {
		delete(ts.drawQueue, i)
	}
}

func (ts *terminalSession) drawLine(line string) {
	fmt.Fprint(ts.out, line)
}

func (ts *terminalSession) clearScreen() {
	fmt.Fprint(ts.out, CSI+ClearScreenSeq)
}

// Unused
func (ts *terminalSession) eraseLine() {
	fmt.Fprint(ts.out, CSI+EraseLineSeq)
}

func (ts *terminalSession) moveCursorTo(x, y int) {
	fmt.Fprintf(ts.out, CSI+MoseCursorToSeq, y, x)
}

func (ts *terminalSession) GetCurrentSize() (err error) {
	width, height, err := term.GetSize(ts.fdIn)
	if err != nil {
		return err
	}
	ts.width, ts.height = width, height
	return nil
}
