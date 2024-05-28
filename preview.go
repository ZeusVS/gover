package main

import (
	"math"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

// TODO: Add preview types?
// Pdf?
// Images?
// Videos?

func (ts *terminalSession) queuePreview() {
	// The width of the preview pane is defined
	width := int(math.Ceil(float64(ts.width)/2.0) - 1)

	file, err := ts.cwdFiles[ts.selectionPos].Info()
	if err != nil {
		return
	}

	if file.Mode()&os.ModeSymlink != 0 {
		link, err := filepath.EvalSymlinks(filepath.Join(ts.cwd, file.Name()))
		// Only change the file if the link is found
		if err == nil {
			linkInfo, err := os.Stat(link)
			if err == nil {
				file = linkInfo
			}
		}
	}

	if file.IsDir() {
		// Get the files under the currently selected dir
		fileName := file.Name()
		previewDir := filepath.Join(ts.cwd, fileName)
		previewFiles, err := os.ReadDir(previewDir)
		ts.previewLen = len(previewFiles)
		if err != nil {
			return
		}

		ts.queueFiles(
			previewFiles,
			previewDir,
			ts.previewOffsetV,
			ts.width/2,
			width,
			false) // We do not want to get a selection in the preview panel
		return
	}

	// Check if file is valid utf8 and can be displayed as text
	// TODO: maybe change this to only read the first line of the file to check UTF8
	filePath := filepath.Join(ts.cwd, file.Name())
	b, _ := os.ReadFile(filePath)
	fileContent := string(b)
	if utf8.ValidString(fileContent) {
		ts.queueFileContents(
			fileContent,
			ts.previewOffsetV,
			ts.previewOffsetH,
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

func (ts *terminalSession) queueFileContents(
	contents string,
	offsetV int,
	offsetH int,
	col int,
	width int) {
	lines := strings.Split(contents, "\n")
	ts.previewLen = len(lines)

	for i, line := range lines {
		if i < offsetV || i > ts.height+offsetV-1-BottomRows {
			continue
		}
		// Replace all tabs with four spaces
		line = strings.ReplaceAll(line, "\t", "    ")
		// If line is longer than the offset
		if len(line) > offsetH {
			line = line[offsetH:]
			// Otherwise make line empty
		} else {
			line = ""
		}
		line = addPadding(line, " ", width)

		drawInstr := drawInstruction{
			x:    col,
			y:    i - offsetV,
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
