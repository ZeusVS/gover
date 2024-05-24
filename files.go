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

func (ts *terminalSession) queueMainFiles() {
	// The width of the main file pane is defined
	width := ts.width/2 - 1
	ts.queueFiles(ts.cwdFiles, ts.cwd, ts.mainOffset, 0, width)
}

func (ts *terminalSession) queueFiles(dirEntries []os.DirEntry, dir string, offset int, col int, width int) {
	for i, dirEntry := range dirEntries {
		if i < offset || i > ts.height+offset-1-BottomRows {
			continue
		}

		file, err := dirEntry.Info()
		if err != nil {
			continue
		}

		var line string
		var link string

		// a directory or a file
		if file.Mode()&os.ModeSymlink != 0 {
			link, err = filepath.EvalSymlinks(filepath.Join(dir, file.Name()))
            if err != nil {
                link = "Error: Link not found"
            }

            symbol := LinkDirIcon
            linkInfo, err := os.Stat(link)
            if err != nil {
                symbol = "?"
            }
            if !linkInfo.IsDir() {
                symbol = LinkFileIcon
            }
			line = ts.getLinkLine(width, i, file, link, symbol)

		} else if file.IsDir() {
			line = ts.getDirLine(width, i, file)

		} else if file.Mode()&0111 != 0 {
			line = ts.getExeLine(width, i, file)

		} else {
			line = ts.getFileLine(width, i, file)
		}

		drawInstr := drawInstruction{
			x:    col,
			y:    i - offset,
			line: line,
		}

		ts.drawQueue = append(ts.drawQueue, drawInstr)
	}

	// We will write blank lines under the files to clear the files pane
	blanklines := ts.height - BottomRows - len(dirEntries)
	for i := range blanklines {
		line := addPadding("", " ", width)

		drawInstr := drawInstruction{
			x:    col,
			y:    len(dirEntries) + i,
			line: line,
		}
		ts.drawQueue = append(ts.drawQueue, drawInstr)
	}
}

func (ts *terminalSession) getDirLine(width int, i int, file os.FileInfo) string {
	line := " " + DirectoryIcon + " " + file.Name()
	line = addPadding(line, " ", width)
	// Add amount of directories under this directory here

	if i == ts.selectionPos {
		line = StyleBgBlue + StyleFgBlack + line + StyleReset
	} else {
		line = StyleFgBlue + line + StyleReset
	}
	return line
}

func (ts *terminalSession) getExeLine(width int, i int, file os.FileInfo) string {
	line := " " + ExecutableIcon + " " + file.Name() + "*"
	line = addPadding(line, " ", width)
	// Add filesize here

	if i == ts.selectionPos {
		line = StyleBgRed + StyleFgBlack + line + StyleReset
	} else {
		line = StyleFgRed + line + StyleReset
	}
	return line
}

func (ts *terminalSession) getFileLine(width int, i int, file os.FileInfo) string {
	line := " " + FileIcon + " " + file.Name()
	line = addPadding(line, " ", width)
	// Add filesize here

	if i == ts.selectionPos {
		line = StyleBgWhite + StyleFgBlack + line + StyleReset
	}
	return line
}

func (ts *terminalSession) getLinkLine(width int, i int, file os.FileInfo, link string, icon string) string {
    
	// TODO: Change icon based on link isdir
	line := " " + icon + " " + file.Name() + " => " + link
	line = addPadding(line, " ", width)

	if i == ts.selectionPos {
		line = StyleBgCyan + StyleFgBlack + line + StyleReset
	} else {
		line = StyleFgCyan + line + StyleReset
	}
	return line
}

func addPadding(line string, padChar string, padWidth int) string {
	// Add spaces to make it fill the file pane slot
	// Rounded down width
	addedSpaces := padWidth - len([]rune(line))
	if addedSpaces > 0 {
		line = fmt.Sprintf("%s%s", line, strings.Repeat(padChar, addedSpaces))
	}
	// line = string([]rune(line)[:ts.width/2])
	return line
}
