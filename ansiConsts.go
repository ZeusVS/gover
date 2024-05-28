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

// Use a var because we can't make a map a const
var (
	inputMap = map[string]rune{
		// Special input
		"enter":     0x0d,
		"backspace": 0x7f,
		"escape":    0x1b,
		// Ctrl + letter input
		"<c-a>": 0x01,
		"<c-b>": 0x02,
		"<c-c>": 0x03,
		"<c-d>": 0x04,
		"<c-e>": 0x05,
		"<c-f>": 0x06,
		"<c-g>": 0x07,
		"<c-h>": 0x08,
		"<c-i>": 0x09,
		"<c-j>": 0x0a,
		"<c-k>": 0x0b,
		"<c-l>": 0x0c,
		"<c-m>": 0x0d,
		"<c-n>": 0x0e,
		"<c-o>": 0x0f,
		"<c-p>": 0x10,
		"<c-q>": 0x11,
		"<c-r>": 0x12,
		"<c-s>": 0x13,
		"<c-t>": 0x14,
		"<c-u>": 0x15,
		"<c-v>": 0x16,
		"<c-w>": 0x17,
		"<c-x>": 0x18,
		"<c-y>": 0x19,
		"<c-z>": 0x1a,
	}
)
