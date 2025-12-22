package main

import (
	"bytes"
	"fmt"
	"os"
	"time"
)

func promptRawXML(buf *bytes.Buffer) {
	if buf.Len() == 0 {
		return
	}
	fmt.Fprintf(os.Stderr, "\nShow raw XML? (y/N) [5s]: ")

	// Create a channel to signal input
	input := make(chan string, 1)

	go func() {
		var s string
		fmt.Scanln(&s)
		input <- s
	}()

	select {
	case res := <-input:
		if len(res) > 0 && (res[0] == 'y' || res[0] == 'Y') {
			fmt.Fprintln(os.Stderr, "\n--- Raw XML Log ---")
			fmt.Fprint(os.Stderr, buf.String())
			fmt.Fprintln(os.Stderr, "-------------------")
		}
	case <-time.After(5 * time.Second):
		fmt.Fprintln(os.Stderr, "") // Newline
	}
}
