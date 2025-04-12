package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// DataManager handles loading and saving application data
type DataManager struct {
	dataDir string
}

// NewDataManager creates a new DataManager with specified data directory
func NewDataManager(dataDir string) *DataManager {
	// Ensure data directory exists
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		err := os.MkdirAll(dataDir, 0755)
		if err != nil {
			fmt.Printf("Error creating data directory: %v\n", err)
			return nil
		}
	}
	
	return &DataManager{
		dataDir: dataDir,
	}
}

// LoadData loads data from a JSON file into the provided destination
func (dm *DataManager) LoadData(filename string, dest interface{}) error {
	filePath := filepath.Join(dm.dataDir, filename)
	
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Printf("Data file %s does not exist, will create when saving.\n", filename)
		return nil
	}
	
	// Read file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file %s: %w", filename, err)
	}
	
	// Unmarshal JSON data
	err = json.Unmarshal(data, dest)
	if err != nil {
		return fmt.Errorf("error parsing JSON data from %s: %w", filename, err)
	}
	
	fmt.Printf("Data loaded successfully from %s\n", filename)
	return nil
}

// SaveData saves data to a JSON file
func (dm *DataManager) SaveData(filename string, data interface{}) error {
	filePath := filepath.Join(dm.dataDir, filename)
	
	// Marshal data to JSON with indentation for readability
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("error encoding JSON data for %s: %w", filename, err)
	}
	
	// Write to file
	err = os.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("error writing to file %s: %w", filename, err)
	}
	
	fmt.Printf("Data saved successfully to %s\n", filename)
	return nil
}