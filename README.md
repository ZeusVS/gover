# Gover
## Terminal file manager in go
The current state of the program is very alpha.
Please note, this file manager only works on linux for now!

### Etymology:
Word blend of
    - go: The programming language this file manager is written in
    - rover: The file manager this one is loosely based upon

### How to start:
Build it yourself: (make sure you have go 1.22.3 installed)
```bash
go build
```
Then run the executable to launch gover:
```bash
./gover
```

### Controls:
- 'q':  Quit the program

- 'h':  Go to parent directory
- 'l':  Go to selected directory
- 'j':  Move selection marker down
- 'k':  Move selection marker up

