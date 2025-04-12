package json

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// Storage provides a JSON-based storage implementation
type Storage struct {
	dataPath string
	mutex    sync.RWMutex
}

// NewStorage creates a new Storage instance
func NewStorage(dataPath string) *Storage {
	// Ensure the directory exists
	if _, err := os.Stat(dataPath); os.IsNotExist(err) {
		err := os.MkdirAll(dataPath, 0755)
		if err != nil {
			panic(fmt.Sprintf("Failed to create data directory: %v", err))
		}
	}
	return &Storage{
		dataPath: dataPath,
		mutex:    sync.RWMutex{},
	}
}

// Save stores data to a JSON file
func (s *Storage) Save(filename string, data interface{}) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	filePath := filepath.Join(s.dataPath, filename)
	
	// Marshal the data to JSON with indentation for readability
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}
	
	// Write the data to the file
	err = os.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	
	return nil
}

// Load reads data from a JSON file
func (s *Storage) Load(filename string, target interface{}) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	filePath := filepath.Join(s.dataPath, filename)
	
	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil // File doesn't exist, not an error
	}
	
	// Read the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	
	// Unmarshal the JSON data
	err = json.Unmarshal(data, target)
	if err != nil {
		return fmt.Errorf("failed to unmarshal data: %w", err)
	}
	
	return nil
}