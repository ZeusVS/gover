package main

import (
	"os"
	"strconv"
	"strings"
)

func (ts *terminalSession) queueBottomBar() {
	// Get the index of the selected item
	selectionIndex := ts.selectionPos + 1
	totalFiles := len(ts.cwdFiles)
	// We want a space before and after the position
	position := " " + strconv.Itoa(selectionIndex) + "/" + strconv.Itoa(totalFiles) + " "

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
	// Add a single space padding before and after the title
	line = " " + line + " "

	// Make the length short enough
	for len(line)+len(position) > ts.width {
		line, _ = strings.CutPrefix(line, " </")
		splitLine := strings.SplitN(line, "/", 2)
		if len(splitLine) > 1 {
			line = " </" + strings.SplitN(line, "/", 2)[1]
		} else {
			// Now we need to remove letter per letter
			line = " </" + line[1:]

			// If the line is " </ " we need to break out of this loop
			if len(line) == 4 {
				break
			}
		}
	}

	spacesToAdd := ts.width - len(line) - len(position)
	if spacesToAdd >= 0 {
		line += strings.Repeat(" ", spacesToAdd)
	} else {
		// Just trim if the width of the terminal is this small
		line = line[:ts.width-len(position)]
	}

	// TODO: add terminal's default background color
	// Will probably have to add termenv because it's not that easy
	// Add color put the strings together
	position = StyleFgBlack + StyleBgBlue + position + StyleReset
	line = StyleBgWhite + StyleFgBlack + line + StyleReset + position

	drawInstr := drawInstruction{
		x:    0,
		y:    ts.height - 2,
		line: line,
	}

	ts.drawQueue = append(ts.drawQueue, drawInstr)

	// Do we need to pad this line with spaces to clear the screen here?
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
