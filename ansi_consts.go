package main

const (
	CSI = "\x1b["

	ShowCursorSeq    = "?25h"
	HideCursorSeq    = "?25l"
	AltScreenSeq     = "?1049h"
	ExitAltScreenSeq = "?1049l"
	ClearScreenSeq   = "2J"
	EraseLineSeq     = "2K"
	MoseCursorToSeq  = "%d;%dH"
)
