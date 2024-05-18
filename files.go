package main

import (
	"os"
)

func (ts *terminalSession) getFiles() (error) {
    var err error
    ts.cwdFiles, err = os.ReadDir(ts.cwd)
    return err
}

func (ts *terminalSession) addFilesToQueue() {
    for i, file := range ts.cwdFiles {
        // TODO: make this based on a view min and view max so that we can 'scroll'
        if i >= ts.height {
            break
        }
        name := file.Name()
        if file.IsDir() {
            name = StyleFgBlue + name + "/" + StyleReset
        // Check if file is executable
        } else if file.Type().Perm() & 0111 != 0 {
            name = StyleFgGreen + name + "*" + StyleReset
        }

        // drawQueue is 1 based and i is 0 based
        if (i + 1) == ts.selectionPos {
            name = StyleBgCyan + name + StyleReset
        }

        ts.drawQueue[i + 1] = name
    }
}
