package main

import (
	"bufio"
	"fmt"
	"math/rand/v2"
	"os"
)

const (
    exit = 'q'
    test = 'f'
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
        case ru == test:
            for i := 0; i < 50; i++ {
                ts.drawQueue[i] = fmt.Sprintf("%d", rand.IntN(10000000000))
            }
    }
    }
}
