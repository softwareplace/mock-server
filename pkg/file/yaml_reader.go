package file

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

func FromYaml[T any](filePath string, outStruct T) (*T, error) {
	// Read the YAML file
	data, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	defer func(data *os.File) {
		err := data.Close()
		if err != nil {
			fmt.Println("Failed to close file: ", err)
		}
	}(data)

	// Ensure the output struct is not nil
	if &outStruct == nil {
		return nil, fmt.Errorf("the output struct must not be nil")
	}

	// Parse the YAML data into the struct
	decoder := yaml.NewDecoder(data)
	err = decoder.Decode(&outStruct)
	if err != nil {
		return nil, fmt.Errorf("failed to decode YAML data: %w", err)
	}

	return &outStruct, nil
}
