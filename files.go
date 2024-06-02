package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	BottomRows = 2
)

func (ts *terminalSession) queueMainFiles() {
	// The width of the main file pane is defined
	width := ts.width/2 - 1
	// Sort the files
	ts.cwdFiles = ts.sortFunc(ts.cwdFiles)
	ts.queueFiles(
		ts.cwdFiles,
		ts.cwd,
		ts.mainOffset,
		0,
		width,
		true) // We want to show the current selection in this panel
}

func (ts *terminalSession) queueFiles(
	dirEntries []os.DirEntry,
	dir string,
	offset int,
	col int,
	width int,
	selection bool,
) {
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
		if file.IsDir() {
			line = ts.getDirLine(width, i, file, dir, selection)
		} else if file.Mode()&os.ModeSymlink != 0 {
			link, err = filepath.EvalSymlinks(filepath.Join(dir, file.Name()))
			if err != nil {
				link = "Error: Link not found"
			}

			symbol := LinkDirIcon
			linkInfo, err := os.Stat(link)
			if err != nil {
				symbol = '?'
			} else if !linkInfo.IsDir() {
				symbol = LinkFileIcon
			}
			line = ts.getLinkLine(width, i, file, link, symbol, selection)

		} else if file.Mode()&0o111 != 0 {
			line = ts.getExeLine(width, i, file, selection)
		} else {
			line = ts.getFileLine(width, i, file, selection)
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

func (ts *terminalSession) getDirLine(width int, i int, file os.FileInfo, dir string, selection bool) string {
	line := " " + string(DirectoryIcon) + "  " + file.Name()
	dirPath, err := os.ReadDir(filepath.Join(dir, file.Name()))
	// Get the number of files underneath the directory
	var dirAmt string
	if err != nil {
		dirAmt = " ? "
	} else {
		dirAmt = " " + strconv.Itoa(len(dirPath)) + " "
	}

	// Do not display dirAmt if there is no space for it
	if len(dirAmt) > width {
		dirAmt = ""
	}

	line = addPadding(line, " ", width-len(dirAmt))
	line += dirAmt

	if i == ts.selectionPos && selection {
		line = StyleBgBlue + StyleFgBlack + line + StyleReset
	} else {
		line = StyleFgBlue + line + StyleReset
	}
	return line
}

func (ts *terminalSession) getExeLine(width int, i int, file os.FileInfo, selection bool) string {
	line := " " + string(ExecutableIcon) + "  " + file.Name() + "*"
	fileSize := getFileSize(file.Size())

	// Do not display filesize if there is no space for it
	if len(fileSize) > width {
		fileSize = ""
	}

	line = addPadding(line, " ", width-len(fileSize))

	line += fileSize

	if i == ts.selectionPos && selection {
		line = StyleBgRed + StyleFgBlack + line + StyleReset
	} else {
		line = StyleFgRed + line + StyleReset
	}
	return line
}

func (ts *terminalSession) getFileLine(width int, i int, file os.FileInfo, selection bool) string {
	line := " " + string(FileIcon) + "  " + file.Name()
	fileSize := getFileSize(file.Size())

	// Do not display filesize if there is no space for it
	if len(fileSize) > width {
		fileSize = ""
	}

	line = addPadding(line, " ", width-len(fileSize))
	line += fileSize

	if i == ts.selectionPos && selection {
		line = StyleBgWhite + StyleFgBlack + line + StyleReset
	}
	return line
}

func (ts *terminalSession) getLinkLine(width int, i int, file os.FileInfo, link string, icon rune, selection bool) string {
	line := " " + string(icon) + "  " + file.Name() + " => " + link
	line = addPadding(line, " ", width)

	if i == ts.selectionPos && selection {
		line = StyleBgCyan + StyleFgBlack + line + StyleReset
	} else {
		line = StyleFgCyan + line + StyleReset
	}
	return line
}

func addPadding(line string, padChar string, padWidth int) string {
	// Chop off part of the string if it's too large
	var trimmedLine string
	chars := 0
	escapeCode := false
	for _, rune := range line {
		// Escape codes get started with escape and stop at m for the color codes
		// Only add 1 to char for non escape code runes
		// Stop adding these chars when past the 'padWidth'
		if escapeCode {
			if rune == 'm' {
				escapeCode = false
			}
			trimmedLine += string(rune)
			continue
		} else if rune == inputMap["escape"] {
			escapeCode = true
			trimmedLine += string(rune)
		} else if chars >= padWidth {
			continue
		} else {
			trimmedLine += string(rune)
			chars += 1
		}
	}

	// Add spaces to make it fill the file panel slot
	addedSpaces := padWidth - chars
	if addedSpaces > 0 {
		trimmedLine = fmt.Sprintf("%s%s", trimmedLine, strings.Repeat(padChar, addedSpaces))
	}

	return trimmedLine
}

func getFileSize(size int64) string {
	var sizeStr string
	var unit string
	// bytes
	if size <= 1024 {
		sizeStr = strconv.Itoa(int(size))
		unit = "B"
	} else {
		sizeFloat := float64(size) / 1024.0
		// kilobytes
		if sizeFloat <= 1024 {
			sizeStr = fmt.Sprintf("%.2f", sizeFloat)
			unit = "K"
		} else {
			sizeFloat /= 1024.0
			// megabytes
			if sizeFloat <= 1024 {
				sizeStr = fmt.Sprintf("%.2f", sizeFloat)
				unit = "M"
			} else {
				sizeFloat /= 1024.0
				// gigabytes
				if sizeFloat <= 1024 {
					sizeStr = fmt.Sprintf("%.2f", sizeFloat)
					unit = "G"
				} else {
					sizeFloat /= 1024.0
					// terabytes
					if sizeFloat <= 1024 {
						sizeStr = fmt.Sprintf("%.2f", sizeFloat)
						unit = "T"
						// I think we can stop here, no?
					}
				}
			}
		}
	}

	// Pad the sizestring with spaces for legibility
	sizeStr = " " + sizeStr + " " + unit + " "
	return sizeStr
}
