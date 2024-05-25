package main

import (
	"math"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
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

		ts.queueFiles(
			previewFiles,
			previewDir,
			ts.previewOffset,
			ts.width/2,
			width,
			false) // We do not want to get a selection in the preview panel
		return
	}

	// Check if file is valid utf8 and can be displayed as text
	// TODO: maybe change this to only read the first line of the file to check UTF8
	filePath := filepath.Join(ts.cwd, ts.cwdFiles[ts.selectionPos].Name())
	b, _ := os.ReadFile(filePath)
	fileContent := string(b)
	if utf8.ValidString(fileContent) {
		ts.queueFileContents(
			fileContent,
			ts.previewOffset,
			ts.width/2,
			width)
		return
	}

	// Otherwise we just display a hatch
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

func (ts *terminalSession) queueFileContents(contents string, offset int, col int, width int) {
	lines := strings.Split(contents, "\n")
	ts.previewLen = len(lines)

	for i, line := range lines {
		if i < offset || i > ts.height+offset-1-BottomRows {
			continue
		}
		// Replace all tabs with four spaces
		line = strings.ReplaceAll(line, "\t", "    ")
		line = addPadding(line, " ", width)

		drawInstr := drawInstruction{
			x:    col,
			y:    i - offset,
			line: line,
		}

		ts.drawQueue = append(ts.drawQueue, drawInstr)
	}

	// We will write blank lines under the files to clear the files pane
	blanklines := ts.height - BottomRows - len(lines)
	for i := range blanklines {
		line := addPadding("", " ", width)

		drawInstr := drawInstruction{
			x:    col,
			y:    len(lines) + i,
			line: line,
		}
		ts.drawQueue = append(ts.drawQueue, drawInstr)
	}
}
