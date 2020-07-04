package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	// When using go run, os.Stdout is held and output of sub process does not display
	// TODO Use a logger
	logger := "ctrlbreak.log"
	os.Remove(logger)
	f, err := os.Create(logger)
	if err != nil {
		log.Fatalf("Failed to create %v: %v", "sublog", err)
	}
	defer f.Close()

	fmt.Fprintln(f, "sub-process started")
	c := make(chan os.Signal, 10)
	signal.Notify(c)
	select {
	case s := <-c:
		if s != os.Interrupt {
			log.Fatalf("Wrong signal received: got %q, want %q\n", s, os.Interrupt)
		}
		fmt.Fprintln(f, "exit on ", s)
	case <-time.After(10 * time.Second):
		log.Fatalf("Timeout waiting for Ctrl+Break\n")
	}
	fmt.Fprintln(f, "graceful exit")
}

