package main

import (
	"math"
	"os"
	"path/filepath"
)

// Issues:
// ts.selectionPos also colors selected field in the preview pane
// Links don't show up properly in preview pane (only see => not the actual destination)

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
		fileName := ts.cwdFiles[ts.selectionPos].Name()
		previewDir := filepath.Join(ts.cwd, fileName)
		previewFiles, err := os.ReadDir(previewDir)
		ts.previewLen = len(previewFiles)
		if err != nil {
			return
		}

		ts.queueFiles(previewFiles, previewDir, ts.previewOffset, ts.width/2, width)
	} else {
		ts.previewLen = ts.height - BottomRows
		for i := range ts.previewLen {
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
