package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
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
			return
		case "stop":
			termigochi.StopDaemon()
			return
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

	// Verify background service is running
	// before pet interactions
	processRunning := termigochi.CheckIfProcessIsRunning()

	if !processRunning {
		daemonConfirmation := confirm("Would you like to start the service?", 2)

		if daemonConfirmation {
			termigochi.StartDaemon()
		} else {
			os.Exit(0)
		}
	}

	// Get or Create active pet
	pet, err := models.LoadPetFromStateFile(conf.PetStateFilePath)

	if err != nil {
		fmt.Printf("Error loading state: %v\n\n", err)
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

func confirm(s string, tries int) bool {
	r := bufio.NewReader(os.Stdin)

	for ; tries > 0; tries-- {
		fmt.Printf("%s [y/n]: ", s)

		res, err := r.ReadString('\n')
		if err != nil {
			log.Fatalln(err)
		}

		if len(res) < 2 {
			continue
		}
		return strings.ToLower(strings.TrimSpace(res))[0] == 'y'
	}
	return false
}
