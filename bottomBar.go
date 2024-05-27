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

	// If there is not enough space to display the position, just remove it
	if len(position) > ts.width {
		position = ""
	}

	// Get the path of the selected item
	// Prevent "//" at root
	cwd := ts.cwd
	if cwd == "/" {
		cwd = ""
	}
	lineTop := cwd + "/" + ts.cwdFiles[ts.selectionPos].Name()
	homeDir, err := os.UserHomeDir()

	var cut bool
	// Only try to cut prefix if a homeDir was found
	if err == nil {
		lineTop, cut = strings.CutPrefix(lineTop, homeDir)
		// Only add tilde if a prefix was cut
		if cut == true {
			lineTop = "~" + lineTop
		}
	}
	// Add a single space padding before and after the title
	lineTop = " " + lineTop + " "

	// Make the length short enough
	for len(lineTop)+len(position) > ts.width {
		lineTop, _ = strings.CutPrefix(lineTop, " </")
		splitLine := strings.SplitN(lineTop, "/", 2)
		if len(splitLine) > 1 {
			lineTop = " </" + strings.SplitN(lineTop, "/", 2)[1]
		} else {
			// Now we need to remove letter per letter
			lineTop = " </" + lineTop[1:]

			// If the line is " </ " we need to break out of this loop
			if len(lineTop) == 4 {
				break
			}
		}
	}

	spacesToAdd := ts.width - len(lineTop) - len(position)
	if spacesToAdd > 0 {
		lineTop += strings.Repeat(" ", spacesToAdd)
	} else {
		// Just trim if the width of the terminal is this small
		lineTop = lineTop[:ts.width-len(position)]
	}

	// TODO: add terminal's default background color
	// Will probably have to add termenv because it's not that easy
	// Add color put the strings together
	position = StyleFgBlack + StyleBgBlue + position + StyleReset
	lineTop = StyleBgWhite + StyleFgBlack + lineTop + StyleReset + position

	drawInstrTop := drawInstruction{
		x:    0,
		y:    ts.height - 2,
		line: lineTop,
	}

	ts.drawQueue = append(ts.drawQueue, drawInstrTop)

	// We add the file permissions on the second line
	// In input mode this line will be overwritten in a separate function
	fileInfo, err := ts.cwdFiles[ts.selectionPos].Info()
	if err != nil {
		// If we can't get the file info, just exit the function
		return
	}
	lineBottom := fileInfo.Mode().String()
	lineBottom = addPadding(lineBottom, " ", ts.width)

	drawInstrBottom := drawInstruction{
		x:    0,
		y:    ts.height - 1,
		line: lineBottom,
	}

	ts.drawQueue = append(ts.drawQueue, drawInstrBottom)
}

func (ts *terminalSession) queueInputLine(input string) {
	input = addPadding(input, " ", ts.width)

	drawInstrBottom := drawInstruction{
		x:    0,
		y:    ts.height - 1,
		line: input,
	}

	ts.drawQueue = append(ts.drawQueue, drawInstrBottom)
}
