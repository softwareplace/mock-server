package config

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

var logFile *os.File // Global variable to keep the file open

func LogSetup(logDirPath string) {
	// Ensure the log directory exists
	if err := os.MkdirAll(logDirPath, 0755); err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}

	// Generate the log file name with the current date
	currentDate := time.Now().Format("2006-01-02") // yyyy-MM-dd format
	logFileName := "mock-server-" + currentDate + ".log"
	logFilePath := filepath.Join(logDirPath, logFileName)

	// Open or create the log file
	var err error
	logFile, err = os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	// Set the output of the logger to the file
	log.SetOutput(logFile)
}
