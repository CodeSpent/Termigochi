package logger

import (
	"log"
	"os"
)

var (
	ServiceLogger *log.Logger
	EventLogger   *log.Logger
)

func init() {
	serviceLogFile, err := os.OpenFile("termigochi_service.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening service log file: %v\n", err)
	}

	eventLogFile, err := os.OpenFile("termigochi_events.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening events log file: %v\n", err)
	}

	ServiceLogger = log.New(serviceLogFile, "SERVICE: ", log.Ldate|log.Ltime|log.Lshortfile)
	EventLogger = log.New(eventLogFile, "EVENT: ", log.Ldate|log.Ltime|log.Lshortfile)
}
