# Gover

<!-- TODO: Add a logo here -->

[![Go](https://img.shields.io/badge/go-blue?style=for-the-badge&logo=go&logoColor=white&logoSize=auto)](https://go.dev)
![Linux](https://img.shields.io/badge/Linux-FCC624?style=for-the-badge&logo=linux&logoColor=black)
![macOS](https://img.shields.io/badge/macOS-black?style=for-the-badge&logo=apple&logoColor=F0F0F0)

<!-- TODO: Add a gif here -->

Gover is a minimalistic terminal-based file manager with vim based key bindings.

## Table of Contents

* [Motivation](#Motivation)
    * [Goal](#Goal)
* [Key features](#Key-features)
* [Getting started](#Getting-started)
    * [Prerequisites](#Prerequisites)
    * [Installation](#Installation)
    * [Running](#Running)
* [Usage](#Usage)
    * [All commands](#All-commands)
* [Contributing](#Contributing)

## Motivation

The currently available file managers are either:

* made in a slower language like python
* bloated with too many unneeded features and not simple/minimalistic anymore
* lacking some core features that are must-haves for a complete file manager experience.

### Goal

The primary goal of Gover is to be a blazingly fast and uncompromising file manager.\
The aim is to add as much functionality as possible while keeping the program still simple to use.\
Each feature will be thoughtfully considered to maintain an intuitive and efficient experience.

## Key features

* Vim based motions and commands
    * Moving around
    * Scrolling
    * Searching
* Basic file manager functionality
    * Opening directories and text based files
    * Creating new files and directories
    * Cutting/copying/pasting files and directories
    * Deleting files and directories
    * Renaming files and directories
* Custom file manager features
    * Custom file/directory sorting
    * Text file previews with syntax highlighting
    * Run console/terminal commands from within Gover
    * Built in manual page showing all available commands

## Getting started

### Prerequisites

* Go 1.22.3+ installed
    * If you do not yet have a working Go environment of at least V1.22.3, please check out [this page](https://go.dev/doc/install)
* Unix based operating system
    * Linux
    * macOS
* A nerdfont installed to properly display Gover's icons
    * If you do not have a nerdfont installed, please check out [this page](https://github.com/ryanoasis/nerd-fonts)

### Installation

```bash
go install github.com/ZeusVS/gover@latest
```

### Running

To start Gover in the current terminal directory:

```bash
gover
```

If you have problems installing and running Gover make sure:

* your $GOPATH is set
* $GOPATH/bin is added to your path

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
