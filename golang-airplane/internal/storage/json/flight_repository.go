package json

import (
	"encoding/json"
	"fmt"
	"golang-airplane/internal/core/domain"
	"golang-airplane/internal/core/ports"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// FlightRepositoryJSON implements the FlightRepository interface using JSON files
type FlightRepositoryJSON struct {
	storage *Storage
}

// NewFlightRepository creates a new FlightRepositoryJSON instance
func NewFlightRepository(storage *Storage) ports.FlightRepository {
	repo := &FlightRepositoryJSON{
		storage: storage,
	}
	// Migrate data if needed
	repo.migrateDataIfNeeded()
	return repo
}

// migrateDataIfNeeded checks if the data is in the old format (map) and migrates it to the new format (slice)
func (r *FlightRepositoryJSON) migrateDataIfNeeded() error {
	filePath := filepath.Join(r.storage.dataPath, "flights.json")
	
	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil // File doesn't exist, no need to migrate
	}
	
	// Read the file content
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file for migration: %w", err)
	}
	
	// Try to unmarshal as a slice first
	var flights []*domain.Flight
	err = json.Unmarshal(data, &flights)
	if err == nil {
		// Data is already in the new format, no migration needed
		return nil
	}
	
	// Try to unmarshal as a map (old format)
	flightsMap := make(map[string]*domain.Flight)
	err = json.Unmarshal(data, &flightsMap)
	if err != nil {
		return fmt.Errorf("failed to unmarshal data in either format: %w", err)
	}
	
	// Convert map to slice
	migratedFlights := make([]*domain.Flight, 0, len(flightsMap))
	for _, flight := range flightsMap {
		migratedFlights = append(migratedFlights, flight)
	}
	
	// Save migrated data
	return r.storage.Save("flights.json", migratedFlights)
}

// FindAll returns all flights in the repository
func (r *FlightRepositoryJSON) FindAll() ([]*domain.Flight, error) {
	var flights []*domain.Flight
	
	// First try to load as a slice
	err := r.storage.Load("flights.json", &flights)
	if err == nil {
		return flights, nil
	}
	
	// If that fails, try to load as a map and convert
	flightsMap := make(map[string]*domain.Flight)
	err = r.storage.Load("flights.json", &flightsMap)
	if err != nil {
		return nil, err
	}
	
	// Convert map to slice
	flights = make([]*domain.Flight, 0, len(flightsMap))
	for _, flight := range flightsMap {
		flights = append(flights, flight)
	}
	
	// Save in the new format for future use
	r.storage.Save("flights.json", flights)
	
	return flights, nil
}

// FindByID finds a flight by its flight number
func (r *FlightRepositoryJSON) FindByID(flightNumber string) (*domain.Flight, error) {
	flights, err := r.FindAll()
	if err != nil {
		return nil, err
	}
	
	// Search for the flight with the given flight number
	for _, flight := range flights {
		if flight.FlightNumber == flightNumber {
			return flight, nil
		}
	}
	
	return nil, fmt.Errorf("flight with number %s not found", flightNumber)
}

// SearchFlights searches for flights by location (departure or destination) and date
func (r *FlightRepositoryJSON) SearchFlights(location string, dateStr string) ([]*domain.Flight, error) {
	// Parse date string to time.Time format to validate it's a correct date
	_, err := time.Parse("02/01/2006", dateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %w", err)
	}
	
	// Get all flights
	flights, err := r.FindAll()
	if err != nil {
		return nil, err
	}
	
	// Filter flights by location and date
	var matchedFlights []*domain.Flight
	location = strings.ToLower(location)
	
	for _, flight := range flights {
		flightDepartureCity := strings.ToLower(flight.DepartureCity)
		flightDestinationCity := strings.ToLower(flight.DestinationCity)
		flightDepartureDate := flight.DepartureTime.Format("02/01/2006")
		
		if (strings.Contains(flightDepartureCity, location) || strings.Contains(flightDestinationCity, location)) && 
		   flightDepartureDate == dateStr {
			matchedFlights = append(matchedFlights, flight)
		}
	}
	
	return matchedFlights, nil
}

// Save stores a flight in the repository
func (r *FlightRepositoryJSON) Save(flight *domain.Flight) error {
	flights, err := r.FindAll()
	if err != nil {
		return err
	}
	
	// Check if the flight already exists
	found := false
	for i, existingFlight := range flights {
		if existingFlight.FlightNumber == flight.FlightNumber {
			// Update existing flight
			flights[i] = flight
			found = true
			break
		}
	}
	
	// Add new flight if not found
	if !found {
		flights = append(flights, flight)
	}
	
	// Save updated flights list
	return r.storage.Save("flights.json", flights)
}

// Update updates an existing flight in the repository
func (r *FlightRepositoryJSON) Update(flight *domain.Flight) error {
	flights, err := r.FindAll()
	if err != nil {
		return err
	}
	
	// Find and update the flight
	found := false
	for i, existingFlight := range flights {
		if existingFlight.FlightNumber == flight.FlightNumber {
			flights[i] = flight
			found = true
			break
		}
	}
	
	if !found {
		return fmt.Errorf("flight with number %s not found", flight.FlightNumber)
	}
	
	// Save updated flights list
	return r.storage.Save("flights.json", flights)
}

// SortFlightsByDepartureTimeDesc sorts flights by departure time in descending order
func SortFlightsByDepartureTimeDesc(flights []*domain.Flight) {
	sort.Slice(flights, func(i, j int) bool {
		return flights[i].DepartureTime.After(flights[j].DepartureTime)
	})
}