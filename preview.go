package main

// Here we will preview the currently selected file
func (ts *terminalSession) queuePreview() {
    var previewLines []string

    if ts.cwdFiles[ts.selectionPos].IsDir() {
        // previewLines := ...
    } else {
        // previewLines := ...
    }
    
    _ = previewLines

    /*
    for i, previewLine := range previewLines {
		drawInstr := drawInstruction{
			x:    ts.width / 2,
			y:    i,
			line: previewLine,
		}

        ts.drawQueue = append(ts.drawQueue, drawInstr)
    }
    */
}
