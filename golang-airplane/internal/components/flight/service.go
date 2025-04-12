package flight

import (
	"fmt"
	"golang-airplane/internal/core/domain"
	"golang-airplane/internal/core/ports"
	"sort"
	"time"
)

// Service implements the FlightService interface
type Service struct {
	flightRepo      ports.FlightRepository
	reservationRepo ports.ReservationRepository
}

// NewService creates a new flight service instance
func NewService(flightRepo ports.FlightRepository, reservationRepo ports.ReservationRepository) *Service {
	return &Service{
		flightRepo:      flightRepo,
		reservationRepo: reservationRepo,
	}
}

// AddFlight adds a new flight
func (s *Service) AddFlight(flightNumber, departureCity, destinationCity string, 
	departureTime, arrivalTime time.Time, availableSeat int) (*domain.Flight, error) {
	
	// Check if flight with the same number already exists
	existingFlight, err := s.flightRepo.FindByID(flightNumber)
	if err == nil && existingFlight != nil {
		return nil, fmt.Errorf("flight with number %s already exists", flightNumber)
	}
	
	// Create new flight
	flight := domain.NewFlight(flightNumber, departureCity, destinationCity, 
		departureTime, arrivalTime, availableSeat)
	
	// Store the flight
	err = s.flightRepo.Save(flight)
	if err != nil {
		return nil, fmt.Errorf("failed to save flight: %w", err)
	}
	
	return flight, nil
}

// GetFlight retrieves a flight by its flight number
func (s *Service) GetFlight(flightNumber string) (*domain.Flight, error) {
	return s.flightRepo.FindByID(flightNumber)
}

// SearchFlights searches for flights by location and date
func (s *Service) SearchFlights(location string, date time.Time) ([]*domain.Flight, error) {
	dateStr := date.Format("02/01/2006")
	return s.flightRepo.SearchFlights(location, dateStr)
}

// AssignCrew assigns crew members to a flight
func (s *Service) AssignCrew(flightNumber string, crewMembers []domain.Crew) error {
	// Get the flight
	flight, err := s.flightRepo.FindByID(flightNumber)
	if err != nil {
		return err
	}
	
	// Verify that the flight doesn't already have a crew assigned
	if len(flight.CrewMembers) > 0 {
		return fmt.Errorf("flight %s already has crew assigned", flightNumber)
	}
	
	// Assign crew members
	flight.AssignCrew(crewMembers)
	
	// Update flight
	return s.flightRepo.Update(flight)
}

// ListAllFlights retrieves all flights sorted by departure time (descending)
func (s *Service) ListAllFlights() ([]*domain.Flight, error) {
	flights, err := s.flightRepo.FindAll()
	if err != nil {
		return nil, err
	}
	
	// Sort flights by departure time in descending order
	sort.Slice(flights, func(i, j int) bool {
		return flights[i].DepartureTime.After(flights[j].DepartureTime)
	})
	
	return flights, nil
}