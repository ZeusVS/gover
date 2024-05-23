package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	BottomRows = 2
)

func (ts *terminalSession) queueFiles() {
	for i, dirEntry := range ts.cwdFiles {
		// TODO: make this based on a view min and view max so that we can 'scroll'
		if i > ts.height-1-BottomRows {
			break
		}

		file, err := dirEntry.Info()
		if err != nil {
			continue
		}

		var line string
		var link string

		// TODO: When there is a symlink I should check if the link points to
		// a directory or a file
		if file.Mode()&os.ModeSymlink != 0 {
			// Error handling???
			link, err = filepath.EvalSymlinks(filepath.Join(ts.cwd, file.Name()))
			line = ts.getLinkLine(i, file, link)

		} else if file.IsDir() {
			line = ts.getDirLine(i, file)

		} else if file.Mode()&0111 != 0 {
			line = ts.getExeLine(i, file)

		} else {
			line = ts.getFileLine(i, file)
		}

		drawInstr := drawInstruction{
			x:    0,
			y:    i,
			line: line,
		}

		ts.drawQueue = append(ts.drawQueue, drawInstr)
	}

    // We will write blank lines under the files to clear the files pane
    blanklines := ts.height - 1 - BottomRows -len(ts.cwdFiles)
    if blanklines > 0 {
        for i := len(ts.cwdFiles); i < ts.height - BottomRows; i++ {
            line := ts.addPadding("")

            drawInstr := drawInstruction{
                x:    0,
                y:    i,
                line: line,
            }
            ts.drawQueue = append(ts.drawQueue, drawInstr)
        }
    }
}

func (ts *terminalSession) getDirLine(i int, file os.FileInfo) string {
	line := DirectoryIcon + " " + file.Name()
	line = ts.addPadding(line)
	// Add amount of directories under this directory here

	if i == ts.selectionPos {
		line = StyleBgBlue + StyleFgBlack + line + StyleReset
	} else {
		line = StyleFgBlue + line + StyleReset
	}
	return line
}

func (ts *terminalSession) getExeLine(i int, file os.FileInfo) string {
	line := ExecutableIcon + " " + file.Name() + "*"
	line = ts.addPadding(line)
	// Add filesize here

	if i == ts.selectionPos {
		line = StyleBgRed + StyleFgBlack + line + StyleReset
	} else {
		line = StyleFgRed + line + StyleReset
	}
	return line
}

func (ts *terminalSession) getFileLine(i int, file os.FileInfo) string {
	line := FileIcon + " " + file.Name()
	line = ts.addPadding(line)
	// Add filesize here

	if i == ts.selectionPos {
		line = StyleBgWhite + StyleFgBlack + line + StyleReset
	}
	return line
}

func (ts *terminalSession) getLinkLine(i int, file os.FileInfo, link string) string {
	// TODO: Change icon based on link isdir
	line := LinkDirIcon + " " + file.Name() + " => " + link
	line = ts.addPadding(line)

	if i == ts.selectionPos {
		line = StyleBgCyan + StyleFgBlack + line + StyleReset
	} else {
		line = StyleFgCyan + line + StyleReset
	}
	return line
}

func (ts *terminalSession) addPadding(line string) string {
	// We use runes here because of the unicode character used
	// Make the selection box half the console's width wide
	// Minus one for the scrollbar
	addedSpaces := ts.width/2 - len([]rune(line))
	if addedSpaces > 0 {
		line = fmt.Sprintf("%s%s", line, strings.Repeat(" ", addedSpaces))
	}
	// line = string([]rune(line)[:ts.width/2])
	return line
}
