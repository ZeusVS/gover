package main

import "math"

func (ts *terminalSession) queueScrollbars() {
	// First screen scrollbar
	ts.queueScrollbar(ts.width/2-1, ts.height-BottomRows, ts.mainOffset, len(ts.cwdFiles))
	// Second screen scrollbar
	ts.queueScrollbar(ts.width-1, ts.height-BottomRows, ts.previewOffset, ts.previewLen)
}

func (ts *terminalSession) queueScrollbar(x int, height int, offset int, contentHeight int) {
	// To still show scrollbar if the content of the preview window is empty
	if contentHeight == 0 {
		contentHeight = 1
	}

	// Check amount of empty spaces to provide at the top
	offsetPercentage := float64(offset) / float64(contentHeight)
	offsetHeightFloat := offsetPercentage * float64(height)
	offsetHeight := int(math.Ceil(offsetHeightFloat))

	// Check amount of empty spaces to provide at the bottom
	offsetBottom := contentHeight - height - offset
	offsetPercentageBottom := (float64(offsetBottom)) / float64(contentHeight)
	offsetHeightFloatBottom := offsetPercentageBottom * float64(height)
	offsetHeightBottom := int(math.Ceil(offsetHeightFloatBottom))

	for i := range height {
		symbol := "┃"
		if i < offsetHeight || i >= height-offsetHeightBottom {
			symbol = " "
		}
		drawInstr := drawInstruction{
			x:    x,
			y:    i,
			line: StyleFgBlue + symbol + StyleReset,
		}
		ts.drawQueue = append(ts.drawQueue, drawInstr)
	}
}
