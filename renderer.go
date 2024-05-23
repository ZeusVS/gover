package main

import (
	"fmt"
	"io"
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
		// The line with index 0 is drawn on position 1 in both x and y direction
		moveCursorTo(ts.out, drawInstr.x+1, drawInstr.y+1)
		fmt.Fprint(ts.out, drawInstr.line)
	}

	// Empty the queue after we are done drawing
	ts.emptyDrawQueue()
}

func (ts *terminalSession) refreshQueue() {
	ts.emptyDrawQueue()
	ts.queueMainFiles()
	ts.queuePreview()
	ts.queueScrollbars()
	ts.queueBottomBar()
}

func (ts *terminalSession) emptyDrawQueue() {
	ts.drawQueue = ts.drawQueue[:0]
}

func moveCursorTo(out io.Writer, x, y int) {
	fmt.Fprintf(out, CSI+MoseCursorToSeq, y, x)
}
