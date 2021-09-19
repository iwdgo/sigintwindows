// Package is a sub-process waiting for a signal
package main

import (
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	// If parent process dies before emitting, no log is available on what happened
	f, err := os.Create("ctrlbreak" + ".log")
	if err != nil {
		log.Fatalf("Failed to create %v: %v", "sublog", err)
	}
	log.SetOutput(f)
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	log.Printf("sub-process %v started", os.Getpid())
	c := make(chan os.Signal, 10)
	signal.Notify(c)
	select {
	case s := <-c:
		if s != os.Interrupt {
			log.Fatalf("Wrong signal received: got %q, want %q\n", s, os.Interrupt)
		}
		log.Printf("graceful exit on %v", s)
	case <-time.After(10 * time.Second):
		// returns exit code 1 to parent process
		log.Fatalf("Timeout waiting for Ctrl+Break")
	}
}
