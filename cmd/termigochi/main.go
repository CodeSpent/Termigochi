package main

import (
	"flag"
	"fmt"
	"os"
	"termigochi/internal/models"
	"termigochi/internal/termigochi"
)

const stateFile = "termigochi_state.json"

var (
	feed   = flag.String("feed", "", "Feed the Pet with specified food")
	play   = flag.String("play", "", "Play with the Pet using specified toy")
	status = flag.Bool("status", false, "Show the Pet's status")
)

func main() {
	flag.Parse()

	pet, err := models.LoadState(stateFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading state: %v\n", err)
		os.Exit(1)
	}

	if *feed != "" {
		termigochi.FeedPet(pet, *feed)
		termigochi.ReportState(pet)
	}

	if *play != "" {
		termigochi.PlayWithPet(pet, *play)
		termigochi.ReportState(pet)
	}

	if *status {
		termigochi.PrintStatus(pet)
		termigochi.ReportState(pet)
	}

	if len(os.Args) < 2 {
		termigochi.ReportState(pet)
		return
	}
	command := os.Args[1]
	switch command {
	case "start":
		termigochi.StartDaemon()
	case "stop":
		termigochi.StopDaemon()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}
