package termigochi

import (
	"fmt"
	"github.com/sevlyar/go-daemon"
	"os"
	"os/signal"
	"syscall"
	"termigochi/internal/config"
	"termigochi/internal/logger"
	"termigochi/internal/models"
	"time"
)

func StartDaemon() {
	cntxt := &daemon.Context{
		PidFileName: "termigochi.pid",
		PidFilePerm: 0644,
		LogFileName: "termigochi_service.log",
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
		Args:        []string{"termigochi", "start"},
	}

	d, err := cntxt.Reborn()
	if err != nil {
		fmt.Println("Unable to start:", err)
		os.Exit(1)
	}
	if d != nil {
		fmt.Println("Termigochi service started")
		logger.ServiceLogger.Println("Termigochi service started")
		return
	}
	defer cntxt.Release()

	go petBackgroundService()

	// Handle graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	logger.ServiceLogger.Println("Termigochi service stopped")
}

func StopDaemon() {
	cntxt := &daemon.Context{
		PidFileName: "termigochi.pid",
	}

	d, err := cntxt.Search()
	if err != nil {
		fmt.Println("Unable to find daemon:", err)
		return
	}

	err = d.Kill()
	if err != nil {
		fmt.Println("Unable to stop daemon:", err)
		return
	}

	fmt.Println("Termigochi service stopped")
	logger.ServiceLogger.Println("Termigochi service stopped")
}

func petBackgroundService() {
	pet, err := models.LoadPetFromStateFile(config.DefaultPetStateFilePath)
	if err != nil {
		logger.ServiceLogger.Printf("Error loading state: %v\n", err)
		os.Exit(1)
	}

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		prevHunger := pet.Hunger
		prevHappiness := pet.Happiness

		pet.TickState()
		pet.SaveState(config.DefaultPetStateFilePath)
		logger.EventLogger.Printf("Updated Hunger: %d -> %d, Happiness: %d -> %d", prevHunger, pet.Hunger,
			prevHappiness, pet.Happiness)
	}
}
