package handler

import (
	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"github.com/softwareplace/mock-server/pkg/env"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type OnFileChangDetected func(restartServer bool)

func watchAndReload(onFileChangeDetected OnFileChangDetected) {
	mockJsonFilesBasePath := env.GetAppEnv().MockPath

	// Set up file watcher to reload mock responses and redirect rules on file changes
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Failed to create file watcher: %v", err)
	}

	// Watch the data directory and its subdirectories
	err = filepath.Walk(mockJsonFilesBasePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			log.Infof("Watching directory: %s", path)
			return watcher.Add(path)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Failed to watch directory: %v", err)
	}

	// Debouncing mechanism
	var (
		debounceDuration = 250 * time.Millisecond // Set the debounce duration to 1 second
		lastEventTime    time.Time
		timer            *time.Timer
	)

	go func() {
		log.Infof("Starting file watcher...")
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create || event.Op&fsnotify.Remove == fsnotify.Remove {

					// Update the last event time
					lastEventTime = time.Now()

					// If a timer is already running, stop it
					if timer != nil {
						timer.Stop()
					}

					// Start a new timer
					timer = time.AfterFunc(debounceDuration, func() {
						// Check if the last event was more than debounceDuration ago
						if time.Since(lastEventTime) >= debounceDuration {
							log.Infof("File %s has changed. Reloading the server...", event.Name)
							loadMockResponses()
							onFileChangeDetected(true)
						}
					})

					// If a new file is created, add it to the watcher
					if event.Op&fsnotify.Create == fsnotify.Create {
						info, err := os.Stat(event.Name)
						if err != nil {
							log.Infof("Failed to stat new file: %v", err)
							continue
						}
						if !info.IsDir() && isValidFileType(info) {
							err = watcher.Add(event.Name)
							if err != nil {
								log.Infof("Failed to add new file to watcher: %v", err)
							}
						}
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Infof("File watcher error: %v", err)
			}
		}
	}()

	defer func() {
		<-make(chan struct{})
		err := watcher.Close()
		if err != nil {
			log.Infof("Failed to close file watcher: %v", err)
		}
		log.Infof("File watcher closed.")
	}()
}

func isValidFileType(info os.FileInfo) bool {
	return strings.HasSuffix(info.Name(), ".json") ||
		strings.HasSuffix(info.Name(), ".yaml") ||
		strings.HasSuffix(info.Name(), ".yml")
}
