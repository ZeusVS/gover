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
		"backspace": 0x7F,
		// Ctrl + letter input
		"ctrl-a": 0x01,
		"ctrl-b": 0x02,
		"ctrl-c": 0x03,
		"ctrl-d": 0x04,
		"ctrl-e": 0x05,
		"ctrl-f": 0x06,
		"ctrl-g": 0x07,
		"ctrl-h": 0x08,
		"ctrl-i": 0x09,
		"ctrl-j": 0x0a,
		"ctrl-k": 0x0b,
		"ctrl-l": 0x0c,
		"ctrl-m": 0x0d,
		"ctrl-n": 0x0e,
		"ctrl-o": 0x0f,
		"ctrl-p": 0x10,
		"ctrl-q": 0x11,
		"ctrl-r": 0x12,
		"ctrl-s": 0x13,
		"ctrl-t": 0x14,
		"ctrl-u": 0x15,
		"ctrl-v": 0x16,
		"ctrl-w": 0x17,
		"ctrl-x": 0x18,
		"ctrl-y": 0x19,
		"ctrl-z": 0x1a,
	}
)
