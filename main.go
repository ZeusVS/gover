package main

import (
	"fmt"
	"sync"
)

func main() {
	// Start the terminal session
	ts, err := StartTerminalSession()
	if err != nil {
		fmt.Printf("Error starting the session: %s\n", err)
		return
	}

	// Add a waitgroup so that we don't quit untill we are done with all goroutines
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		ts.startKeyListener()
	}()
	go func() {
		defer wg.Done()
		ts.startResizeListener()
	}()
	go func() {
		defer wg.Done()
		ts.startRendering()
		// Stop the session when we are done rendering
		ts.StopTerminalSession()
	}()

	// Wait for all goroutines to finish before closing
	wg.Wait()
}
