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
		// Moves the cursor to pos x+1:y+1
		// Terminal is 1 based and our instr's are 0 based
		fmt.Fprintf(ts.out, CSI+MoseCursorToSeq, drawInstr.y+1, drawInstr.x+1)
		// Draw the line
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
