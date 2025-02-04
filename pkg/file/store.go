package file

import (
	"os"
	"path/filepath"
)

// SaveToFile takes a string `data` and a file path `filePath`,
// creates the directory if it does not exist, and writes the data to the file.
func SaveToFile(data []byte, filePath string) error {
	// Get the directory from the filePath.
	dir := filepath.Dir(filePath)

	// Create the directory if it does not exist.
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	// Write the data to the file.
	if err := os.WriteFile(filePath, data, os.ModePerm); err != nil {
		return err
	}

	return nil
}
