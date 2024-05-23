package main

const (
	CSI = "\033["

	ShowCursorSeq    = "?25h"
	HideCursorSeq    = "?25l"
	AltScreenSeq     = "?1049h"
	ExitAltScreenSeq = "?1049l"
	ClearScreenSeq   = "2J"
	EraseLineSeq     = "2K"
	MoseCursorToSeq  = "%d;%dH"

	StyleReset = CSI + "0m"
	// Regular foreground colors
	StyleFgBlack  = CSI + "30m"
	StyleFgRed    = CSI + "31m"
	StyleFgGreen  = CSI + "32m"
	StyleFgYellow = CSI + "33m"
	StyleFgBlue   = CSI + "34m"
	StyleFgPurple = CSI + "35m"
	StyleFgCyan   = CSI + "36m"
	StyleFgWhite  = CSI + "37m"

	// Bright foreground colors
	StyleFgBlackBright = CSI + "90m"

	// Regular background colors
	StyleBgBlack  = CSI + "40m"
	StyleBgRed    = CSI + "41m"
	StyleBgGreen  = CSI + "42m"
	StyleBgYellow = CSI + "43m"
	StyleBgBlue   = CSI + "44m"
	StyleBgPurple = CSI + "45m"
	StyleBgCyan   = CSI + "46m"
	StyleBgWhite  = CSI + "47m"

	// Icons
	ExecutableIcon = "\uf489 "
	DirectoryIcon  = "\uf115 "
	FileIcon       = "\uf016 "
	LinkDirIcon    = "\uf482 "
	LinkFileIcon   = "\uf481 "
)
