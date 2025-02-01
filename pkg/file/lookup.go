package file

import (
	"os"
	"path/filepath"
	"strings"
)

// FindFilesWithExtension searches a directory and its subdirectories
// for files that match the given extension.
func FindFilesWithExtension(dirPath string, extension string) ([]string, error) {
	var files []string

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if it's a file and matches the extension
		if !info.IsDir() && strings.HasSuffix(info.Name(), extension) {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}
