package main

import (
	"bytes"
	"math"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
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
		previewFiles = ts.sortFunc(previewFiles)
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
			file.Name(),
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
	name string,
	contents string,
	offsetV int,
	offsetH int,
	col int,
	width int) {

	// First we apply the syntax highlighting
	lexer := lexers.Match(name)
	if lexer == nil {
		lexer = lexers.Fallback
	}
	// TODO: Make the syntax highlighting match the terminal colors
	style := styles.Get("tokyonight-night")

	formatter := formatters.Get("terminal256")
	iterator, err := lexer.Tokenise(nil, contents)
	if err != nil {
		// TODO
	}

	// Write the text with syntax highlighting to a buffer
	buffer := new(bytes.Buffer)
	err = formatter.Format(buffer, style, iterator)
	if err != nil {
		// TODO
	}

	lines := strings.Split(buffer.String(), "\n")
	ts.previewLen = len(lines)

	// Split up the buffer text per newline and add to the drawQueue
	for i, line := range lines {
		if i < offsetV || i > ts.height+offsetV-1-BottomRows {
			continue
		}

		line = strings.ReplaceAll(line, "\t", "    ")
		// Trim off the first "offsetH" characters
		line = removeFirstChars(line, offsetH)
		// Make every line the correct width
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

func removeFirstChars(line string, n int) string {

	runeLine := []rune(line)
	escapeCode := false
	returnLine := ""
	for _, rune := range runeLine {
		// Escape codes get started with escape and stop at m for the color codes
		// Skip adding n non escape characters to the return string
		if escapeCode {
			if rune == 'm' {
				escapeCode = false
			}
			returnLine += string(rune)
			continue
		} else if rune == inputMap["escape"] {
			escapeCode = true
			returnLine += string(rune)
		} else if n > 0 {
			n -= 1
		} else {
			returnLine += string(rune)
		}
	}

	return returnLine
}
