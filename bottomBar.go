package main

import (
	"os"
	"strconv"
	"strings"
)

func (ts *terminalSession) queueBottomBar() {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	// Get the index of the selected item
	selectionIndex := ts.selectionPos + 1
	totalFiles := len(ts.cwdFiles)
	// We want a space before and after the position
	position := " " + strconv.Itoa(selectionIndex) + "/" + strconv.Itoa(totalFiles) + " "
	// Make position atleast 7 long so that we prevent too much ui shifting
	if len(position) < 7 {
		position = strings.Repeat(" ", 7-len(position)) + position
	}

	// Get the path of the selected item
	// Start with a space
	line := ts.cwd + "/" + ts.cwdFiles[ts.selectionPos].Name()
	homeDir, err := os.UserHomeDir()

	var cut bool
	// Only cut prefix is a homeDir was found
	if err == nil {
		line, cut = strings.CutPrefix(line, homeDir)
		// Only add tilde if a prefix was cut
		if cut == true {
			line = "~" + line
		}
	}

	// Make the length fit
	if len(line)+len(position) <= ts.width {
		spacesToAdd := ts.width - len(line) - len(position)
		// TODO: add default background color
		// Will probably have to add termenv because it's not that easy
		line += strings.Repeat(" ", spacesToAdd)
	} else {
		// What to do if width is too long???
	}
	// Add color put the strings together
	position = StyleFgBlack + StyleBgBlue + position + StyleReset
	line = StyleBgWhite + StyleFgBlack + line + StyleReset + position

	drawInstr := drawInstruction{
		x:    0,
		y:    ts.height - 2,
		line: line,
	}

	ts.drawQueue = append(ts.drawQueue, drawInstr)

	// Now we get the second line, which for now only holds the file permissions
	fileInfo, err := ts.cwdFiles[ts.selectionPos].Info()
	if err != nil {
		// If we can't get the file info, just exit the function
		return
	}
	filePerms := fileInfo.Mode().String()

	drawInstr = drawInstruction{
		x:    0,
		y:    ts.height - 1,
		line: filePerms,
	}

	ts.drawQueue = append(ts.drawQueue, drawInstr)
}
