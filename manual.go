package main

import "strings"

func (ts *terminalSession) showManual() {
	// TODO: Add scrolling when the cmd amount gets much higher

	// start inputmode so that we block all input untill we leave the manual page
	ts.inputMode = true
	defer func() { ts.inputMode = false }()

	// Clear entire screen first (only leave very last row)
	for i := range ts.height - 1 {
		line := addPadding("", " ", ts.width)
		drawInstr := drawInstruction{
			x:    0,
			y:    i,
			line: line,
		}

		ts.drawQueue = append(ts.drawQueue, drawInstr)
	}

	// Draw the manual
	lines := `All available commands:
Actions:
'q'      Quit Gover
'?'      Show manual page
':'      Enter console command from the current directory
'escape' Clear all actions

'i'     Insert/create new file in the current directory
'I'     Insert/create new directory in the current directory

'd'     Mark the currently selected file for cutting/moving
'y'     Mark the currently selected file for copying
'p'     Cut/Copy the marked file to the current directory

'D'     (Recursively) delete the current selection - will ask for confirmation
'R'     Rename the current selection

'/'     Search the main panel for specific text
'n'     Jump to next occurrence of the searchstring
'N'     Jump to previous occurrence of the searchstring

Motions:
'~'     Go to your home directory
'h'     Go to parent directory
'l'     Go to selected directory
'j'     Move selection marker down
'J'     Move selection marker down by 10
'k'     Move selection marker up
'K'     Move selection marker up by 10
'gg':   Move selection marker to the top
'G':    Move selection marker to the bottom

'<c-u>' Scroll the preview panel up
'<c-d>' Scroll the preview panel down
'<c-f>' Scroll the preview panel left
'<c-k>' Scroll the preview panel right

Sorting commands:
'sd'    Sort directories first
'sD'    Sort directories last
'sa'    Sort alphabetically
'sA'    Sort alphabetically reversed
'st'    Sort by modification time, newest first
'sT'    Sort by modification time, oldest first
'ss'    Sort by filesize, smallest first
'sS'    Sort by filesize, largest first`

	for i, line := range strings.Split(lines, "\n") {
		drawInstr := drawInstruction{
			x:    0,
			y:    i,
			line: line,
		}

		ts.drawQueue = append(ts.drawQueue, drawInstr)
	}

	// Leave the manual screen when we press escape
	for {
		ts.queueInputLine("Press the escape button to leave the manual page")
		ru := <-ts.inCh
		if ru == inputMap["escape"] {
			break
		}
	}

	// Add the original screen back to the drawQueue
	ts.refreshQueue()
}
