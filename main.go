package main

import (
	"fmt"
)

func main() {
	// Start the terminal session
	ts, err := StartTerminalSession()
	if err != nil {
		fmt.Printf("Error starting the session: %s\n", err)
		return
	}
	defer ts.StopTerminalSession()

	go ts.startListening()
	ts.startRendering()
}
