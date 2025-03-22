package config

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	logFile     *os.File // Global variable to keep the file open
	logDirPath  string   // Global variable to store the log directory path
	logFilePath string   // Global variable to store the full log file path
	multiWriter io.Writer
)

func LogSetup(dirPath string) {
	logDirPath = dirPath
	createLogFile()

	// Start a goroutine to monitor the log file
	go monitorLogFile()
}

// createLogFile creates or reopens the log file
func createLogFile() {
	if err := os.MkdirAll(logDirPath, 0755); err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}

	// Generate the log file name with the current date
	currentDate := time.Now().Format("2006-01-02") // yyyy-MM-dd format
	logFileName := "mock-server-" + currentDate + ".log"
	logFilePath = filepath.Join(logDirPath, logFileName)

	// Open or create the log file
	var err error
	logFile, err = os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	// Create a multi-writer that writes to both the file and the terminal
	multiWriter = io.MultiWriter(os.Stdout, logFile)

	// Set the output of the logger to the multi-writer
	log.SetOutput(&timestampWriter{
		writer: multiWriter,
		format: "2006-01-02 15:04:05", // Custom timestamp format
	})
}

// monitorLogFile periodically checks if the log file exists and recreates it if necessary
func monitorLogFile() {
	for {
		time.Sleep(5 * time.Second) // Check every 5 seconds

		// Check if the log file still exists
		if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
			log.Println("Log file was deleted. Recreating...")
			createLogFile() // Recreate the log file
		}
	}
}

// CloseLogFile can be called to properly close the log file when the program exits
func CloseLogFile() {
	if logFile != nil {
		err := logFile.Close()
		if err != nil {
			log.Printf("Failed to close log file: %v", err)
		}
	}
}

// timestampWriter is a custom writer that adds a timestamp to each log message
type timestampWriter struct {
	writer io.Writer
	format string
}

func (tw *timestampWriter) Write(p []byte) (n int, err error) {
	timestamp := time.Now().Format(tw.format)
	message := timestamp + " " + string(p)
	return tw.writer.Write([]byte(message))
}
