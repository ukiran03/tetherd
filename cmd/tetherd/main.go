package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ukiran.com/tetherd/internal/procnet"
)

const stateFileName = "/var/lib/tetherd/state.bin"

func main() {
	state, err := NewState(stateFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer state.Close()

	if err := state.Init(); err != nil {
		log.Fatal(err)
	}

	// channel to listen for OS signals (Ctrl+C, SIGTERM)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-sigChan:
			fmt.Println("Shutting down gracefully...")
			finalStats, _ := procnet.ReadTotalDataUsage()
			state.Sync(finalStats)
			return

		case <-ticker.C:
			newStats, err := procnet.ReadTotalDataUsage()
			if err != nil {
				fmt.Fprintf(os.Stderr, "reading usage error: %v\n", err)
				continue
			}

			if err := state.Sync(newStats); err != nil {
				fmt.Fprintf(os.Stderr, "error syncing: %v\n", err)
			}
		}
	}
}
