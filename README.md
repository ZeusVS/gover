# Gover
## Terminal file manager for unix systems in go
The current state of the program is very alpha.\
Please note, this file manager is only for unix based systems!

### Etymology:
Word blend of:\
    - go: The programming language this file manager is written in\
    - rover: The file manager this one is loosely based upon

### How to start:
- Make sure you have go 1.22.3 installed
- Install using the following:
```bash
go install github.com/ZeusVS/gover@latest
```
- This will automatically build the code and add it to $GOPATH/bin for you
- Now you should be able to execute gover in your terminal
    - If you have problems make sure your $GOPATH is set and $GOPATH/bin is added to your $PATH

### Commands:
```
- 'q': Quit the program

- ':': Enter console command from the current directory

- 'enter': Open the current selection in a new window:
               Directory: default terminal
               Text file: default editor

- 'R': Rename the current selection
- '/': Search the main panel for specific text
- 'n': Jump to next occurrence of the searchstring
- 'N': Jump to previous occurrence of the searchstring

- '~':  Go to your home directory
- 'h':  Go to parent directory
- 'l':  Go to selected directory
- 'j':  Move selection marker down
- 'J':  Move selection marker down by 10
- 'k':  Move selection marker up
- 'K':  Move selection marker up by 10
- 'gg': Move selection marker to the top
- 'G':  Move selection marker to the bottom

- '<c-u>':  Scroll the preview panel up
- '<c-d>':  Scroll the preview panel down
- '<c-f>':  Scroll the preview panel left
- '<c-k>':  Scroll the preview panel right

// Sorting commands
- 'sd': Sort directories first
- 'sD': Sort directories last
- 'sa': Sort alphabetically
- 'sA': Sort alphabetically reversed
- 'st': Sort by modification time, newest first
- 'sT': Sort by modification time, oldest first
- 'ss': Sort by filesize, smallest first
- 'sS': Sort by filesize, largest first
```

### Important information:
To open files and directories the environment variables $TERM and $EDITOR are used,\
to make this functionality work as intended make sure these are set.
