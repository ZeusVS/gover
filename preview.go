package main

import (
	"math"
	"os"
	"path/filepath"
)

// TODO: Very rough draft
// Issues:
// ts.selectionPos also colors selected field in the preview pane

// Preview types:
// Directories DONE
// Text files TODO
// Pdf?
// Images?
// Videos?
// Everything else will be filled with ╱╱╱╱ for now

// Here we will preview the currently selected file
func (ts *terminalSession) queuePreview() {
	// The width of the preview pane is defined
	width := int(math.Ceil(float64(ts.width)/2.0) - 1)

	if ts.cwdFiles[ts.selectionPos].IsDir() {
		// Get the files under the currently selected dir
		// TODO: We could add this to our terminalSession struct so that we do not have to do this again in the goDownDir action?
		fileName := ts.cwdFiles[ts.selectionPos].Name()
		newDir := filepath.Join(ts.cwd, fileName)
		newFiles, err := os.ReadDir(newDir)
		if err != nil {
			return
		}

		ts.queueFiles(newFiles, ts.width/2, width)
	} else {
		for i := range ts.height - BottomRows {
			line := StyleFgBlackBright + addPadding("", "╱", width) + StyleReset
			drawInstr := drawInstruction{
				x:    ts.width / 2,
				y:    i,
				line: line,
			}

			ts.drawQueue = append(ts.drawQueue, drawInstr)
		}
	}
}
