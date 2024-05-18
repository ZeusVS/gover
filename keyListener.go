package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
    exit = 'q'
    test = 'f'

    up =      'k'
    down =    'j'
    dirup =   'h'
    dirdown = 'l'
)
func (ts terminalSession) startListening() {
    r := bufio.NewReader(os.Stdin)
    for {
        ru, _, err := r.ReadRune()

        if err != nil {
            fmt.Fprintf(os.Stderr, "Error: reading key from Stdin: %s\r\n", err)
        }

        switch {
        case ru == exit:
                ts.done <- true

        case ru == up:
            ts.moveSelectionUp()
        case ru == down:
            ts.moveSelectionDown()
        case ru == dirup:
            ts.moveUpDir()
        case ru == dirdown:
            ts.moveDownDir()

        //case ru == test:
        }
    }
}
