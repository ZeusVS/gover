# Gover

<!-- TODO: Add a logo here -->

[![Go](https://img.shields.io/badge/go-blue?style=for-the-badge&logo=go&logoColor=white&logoSize=auto)](https://go.dev)
![Linux](https://img.shields.io/badge/Linux-FCC624?style=for-the-badge&logo=linux&logoColor=black)
![macOS](https://img.shields.io/badge/macOS-black?style=for-the-badge&logo=apple&logoColor=F0F0F0)

<!-- TODO: Add a gif here -->

<!-- TODO: Add a description here -->

## Table of Contents

* [Getting started](#Getting-started)
    * [Prerequisites](#Prerequisites)
    * [Installation](#Installation)
* [Usage](#Usage)
    * [All commands](#All-commands)
* [Contributing](#Contributing)

## Getting started

### Prerequisites

* Go 1.22.3+
<!-- TODO: Add go intallation instructions here -->
* Unix based operating system
    * Linux
    * macOS
* A [nerdfont](https://github.com/ryanoasis/nerd-fonts) installed if you want to properly display file/dir/link/exec icons

### Installation

```bash
go install github.com/ZeusVS/gover@latest
```

This will automatically build the code and add it to your $GOPATH/bin 

### Running

Run gover in the desired working directory

```bash
gover
```

    If you have problems running gover make sure your $GOPATH is set and $GOPATH/bin is added to your path

## Usage

### All commands

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

## Contributing

I would love your help! Contribute by forking the repo and opening pull requests or creating issues.

All pull requests should be submitted to the `main` branch.
