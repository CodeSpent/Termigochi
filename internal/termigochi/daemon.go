package termigochi

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/sevlyar/go-daemon"
	"os"
	"os/signal"
	"strconv"
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

func CheckIfProcessIsRunning() bool {
	pidFile := config.DefaultProcessFilePath + "termigochi.pid"

	procFile, err := os.Open(pidFile)
	if err != nil {
		if os.IsNotExist(err) {
			logger.ServiceLogger.Println("PID file not found")
		}
		logger.ServiceLogger.Println("Error opening process:", err)
	}
	defer func(procFile *os.File) {
		err := procFile.Close()
		if err != nil {
			fmt.Printf("Error closing process file: %v", err)
		}
	}(procFile)

	r := bufio.NewReader(procFile)

	pidStr, err := r.ReadString('\n')
	if err != nil {
		logger.ServiceLogger.Println("Error reading PID file:", err)
	}

	pid, err := strconv.Atoi(pidStr)
	processRunning, err := PidExists(pid)

	return processRunning
}

func PidExists(pid int) (bool, error) {
	if pid <= 0 {
		return false, fmt.Errorf("invalid pid %v", pid)
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		return false, err
	}
	err = proc.Signal(syscall.Signal(0))
	if err == nil {
		return true, nil
	}
	if err.Error() == "os: process already finished" {
		return false, nil
	}
	var errno syscall.Errno
	ok := errors.As(err, &errno)
	if !ok {
		return false, err
	}
	switch {
	case errors.Is(errno, syscall.ESRCH):
		return false, nil
	case errors.Is(errno, syscall.EPERM):
		return true, nil
	}
	return false, err
}
