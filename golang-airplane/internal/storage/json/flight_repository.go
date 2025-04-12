package json

import (
	"fmt"
	"golang-airplane/internal/core/domain"
	"golang-airplane/internal/core/ports"
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
	return &FlightRepositoryJSON{
		storage: storage,
	}
}

// FindAll returns all flights in the repository
func (r *FlightRepositoryJSON) FindAll() ([]*domain.Flight, error) {
	flightsMap := make(map[string]*domain.Flight)
	err := r.storage.Load("flights.json", &flightsMap)
	if err != nil {
		return nil, err
	}

	// Convert map to slice
	flights := make([]*domain.Flight, 0, len(flightsMap))
	for _, flight := range flightsMap {
		flights = append(flights, flight)
	}
	
	return flights, nil
}

// FindByID finds a flight by its flight number
func (r *FlightRepositoryJSON) FindByID(flightNumber string) (*domain.Flight, error) {
	flightsMap := make(map[string]*domain.Flight)
	err := r.storage.Load("flights.json", &flightsMap)
	if err != nil {
		return nil, err
	}
	
	flight, exists := flightsMap[flightNumber]
	if !exists {
		return nil, fmt.Errorf("flight with number %s not found", flightNumber)
	}
	
	return flight, nil
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
	flightsMap := make(map[string]*domain.Flight)
	
	// Load existing flights
	err := r.storage.Load("flights.json", &flightsMap)
	if err != nil {
		return err
	}
	
	// Add or update flight
	flightsMap[flight.FlightNumber] = flight
	
	// Save updated flights map
	return r.storage.Save("flights.json", flightsMap)
}

// Update updates an existing flight in the repository
func (r *FlightRepositoryJSON) Update(flight *domain.Flight) error {
	// Update and save are the same operation in this implementation
	return r.Save(flight)
}

// SortFlightsByDepartureTimeDesc sorts flights by departure time in descending order
func SortFlightsByDepartureTimeDesc(flights []*domain.Flight) {
	sort.Slice(flights, func(i, j int) bool {
		return flights[i].DepartureTime.After(flights[j].DepartureTime)
	})
}