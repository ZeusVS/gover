<div align="center">

# Gover
#### Terminal file manager in go
[![Go](https://img.shields.io/badge/golang-blue?style=for-the-badge&logo=go&logoColor=white&logoSize=auto)](https://go.dev)
![macOS](https://img.shields.io/badge/macOS-black?style=for-the-badge&logo=apple&logoColor=F0F0F0)
![Linux](https://img.shields.io/badge/Linux-FCC624?style=for-the-badge&logo=linux&logoColor=black)

</div>

## - TOC
* [Etymology](#Etymology)
* [Installation](#Installation)
* [Commands](#Commands)
* [Important information](#Important-information)

## Etymology
Word blend of\
* go: The programming language this file manager is written in\
* rover: The file manager this one is loosely based upon

## Installation
* Go 1.22.3+ required
* Install using the following command
```bash
go install github.com/ZeusVS/gover@latest
```
* This will automatically build the code and add it to $GOPATH/bin for you
* Now you should be able to execute gover in your terminal
    * If you have problems make sure your $GOPATH is set and $GOPATH/bin is added to your $PATH

## Commands
```
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

'enter' Open the current selection in a new window
            Directory: default terminal
            Text file: default editor
            Executable: run in a new window

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
'sS'    Sort by filesize, largest first
```

## Important information
To open files and directories the environment variables $TERM and $EDITOR are used,\
to make this functionality work as intended make sure these are set.
