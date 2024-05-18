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

        // Check if file is a directory
        if file.IsDir() {
            name = DirectoryIcon + " " + name + "/"
            if i == ts.selectionPos {
                name = StyleBgBlue + StyleFgBlack + name + StyleReset
            } else {
                name = StyleFgBlue + name + StyleReset
            }

        // Check if file is executable
        } else if file.Type().Perm() & 0111 != 0 {
            name = ExecutableIcon + " " + name + "*"
            if i == ts.selectionPos {
                name = StyleBgGreen + StyleFgBlack + name + StyleReset
            } else {
                name = StyleFgGreen + name + StyleReset
            }

        // Regular files
        } else {
            name = FileIcon + " " + name
            if i == ts.selectionPos {
                name = StyleBgWhite + StyleFgBlack + name + StyleReset
            }
        }

        ts.drawQueue[i] = name
    }
}
