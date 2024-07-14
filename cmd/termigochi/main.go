package main

import (
	"flag"
	"fmt"
	"os"
	"termigochi/internal/config"
	"termigochi/internal/logger"
	"termigochi/internal/models"
	"termigochi/internal/termigochi"
	"time"
)

var (
	feed   = flag.String("feed", "", "Feed the Pet with specified food")
	play   = flag.String("play", "", "Play with the Pet using specified toy")
	status = flag.Bool("status", false, "Show the Pet's status")
)

func main() {
	if len(os.Args) > 1 {
		command := os.Args[1]
		switch command {
		case "start":
			termigochi.StartDaemon()
		case "stop":
			termigochi.StopDaemon()
		default:
			break
		}
	}

	flag.Parse()

	conf, err, created := config.LoadConfig(config.DefaultConfigPath)
	if err != nil {
		logger.ServiceLogger.Println("Error loading conf:", err)
		config.NewConfig(config.DefaultConfigPath)
	}

	if created {
		// We assume this is first run
		termigochi.StartOnboarding(conf)
	}

	// Get or Create active pet
	pet, err := models.LoadPetFromStateFile(conf.PetStateFilePath)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading state: %v\n", err)
		os.Exit(1)
	}

	// Check for hatch event
	if pet.HatchDate.Before(time.Now()) && !pet.Hatched {
		pet.Hatch()
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
}
