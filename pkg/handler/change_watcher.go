package handler

import (
	"github.com/fsnotify/fsnotify"
	"github.com/softwareplace/mock-server/pkg/env"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func watchAndReload() {
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
			log.Printf("Watching directory: %s", path)
			return watcher.Add(path)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Failed to watch directory: %v", err)
	}

	go func() {
		log.Println("Starting file watcher...")
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create || event.Op&fsnotify.Remove == fsnotify.Remove {
					log.Printf("File %s changed: %s", event.Name, event.Op)
					loadMockResponses()

					// If a new file is created, add it to the watcher
					if event.Op&fsnotify.Create == fsnotify.Create {
						info, err := os.Stat(event.Name)
						if err != nil {
							log.Printf("Failed to stat new file: %v", err)
							continue
						}
						if !info.IsDir() && isValidFileType(info) {
							err = watcher.Add(event.Name)
							if err != nil {
								log.Printf("Failed to add new file to watcher: %v", err)
							}
						}
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Printf("File watcher error: %v", err)
			}
		}
	}()

	defer func() {
		<-make(chan struct{})
		err := watcher.Close()
		if err != nil {
			log.Printf("Failed to close file watcher: %v", err)
		}
		log.Println("File watcher closed.")
	}()
}

func isValidFileType(info os.FileInfo) bool {
	return strings.HasSuffix(info.Name(), ".json") ||
		strings.HasSuffix(info.Name(), ".yaml") ||
		strings.HasSuffix(info.Name(), ".yml")
}
