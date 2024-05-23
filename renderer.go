package main

import (
	"fmt"
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

// Draw everything waiting in the queue
func (ts *terminalSession) render() {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	for _, drawInstr := range ts.drawQueue {
		// The line with index 0 is drawn on position 1
		ts.moveCursorTo(drawInstr.x+1, drawInstr.y+1)
		ts.drawLine(drawInstr.line)
	}

	// Empty the queue after we are done drawing
	ts.drawQueue = ts.drawQueue[:0]
}

func (ts *terminalSession) refreshQueue() {
	ts.emptyDrawQueue()
	ts.queueFiles()
	ts.queueScrollbars()
	ts.queueBottomBar()
}

func (ts *terminalSession) emptyDrawQueue() {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	ts.drawQueue = ts.drawQueue[:0]
}

func (ts *terminalSession) drawLine(line string) {
	fmt.Fprint(ts.out, line)
}

func (ts *terminalSession) moveCursorTo(x, y int) {
	fmt.Fprintf(ts.out, CSI+MoseCursorToSeq, y, x)
}
