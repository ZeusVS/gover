package main

func (ts *terminalSession) queueScrollbars() {
	// First screen scrollbar
	ts.queueScrollbar(ts.width/2-1, ts.height-BottomRows)
	// Second screen scrollbar
	ts.queueScrollbar(ts.width-1, ts.height-BottomRows)
}

func (ts *terminalSession) queueScrollbar(x int, height int) {
	for i := range height {
		drawInstr := drawInstruction{
			x:    x,
			y:    i,
			line: StyleFgBlue + "┃" + StyleReset,
		}
		ts.drawQueue = append(ts.drawQueue, drawInstr)
	}
}
