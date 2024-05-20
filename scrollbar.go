package main

func (ts terminalSession) drawScrollbars() {
	// First screen scrollbar
	ts.drawScrollbar(ts.width/2, ts.height-2)
	// Second screen scrollbar
	ts.drawScrollbar(ts.width, ts.height-2)
}

func (ts terminalSession) drawScrollbar(x int, height int) {
	for i := range height {
		ts.moveCursorTo(x, i+1)
		ts.drawLine(StyleFgBlue + "┃" + StyleReset)
	}
}
