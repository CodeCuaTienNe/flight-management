package ports

import (
	"time"

	"golang-airplane/internal/core/domain"
)

type AirplaneService interface {
	// AddAirplane creates and stores a new airplane
	AddAirplane(id string, model string, capacity int) error
	
	// GetAirplanes retrieves all airplanes
	GetAirplanes() ([]domain.Airplane, error)
	
	// GetAirplaneByID retrieves an airplane by its ID
	GetAirplaneByID(id string) (domain.Airplane, error)
}

type FlightService interface {
	// AddFlight adds a new flight
	AddFlight(flightNumber, departureCity, destinationCity string, departureTime, arrivalTime time.Time, availableSeat int) (*domain.Flight, error)
	
	// GetFlight retrieves a flight by its flight number
	GetFlight(flightNumber string) (*domain.Flight, error)
	
	// SearchFlights searches for flights by location and date
	SearchFlights(location string, date time.Time) ([]*domain.Flight, error)
	
	// AssignCrew assigns crew members to a flight
	AssignCrew(flightNumber string, crewMembers []domain.Crew) error
	
	// ListAllFlights retrieves all flights sorted by departure time (descending)
	ListAllFlights() ([]*domain.Flight, error)
}

type ReservationService interface {
	// BookFlight creates a new reservation for a flight
	BookFlight(name, address string, phoneNumber, identityCardNumber int64, flightNumber string) (*domain.Reservation, error)
	
	// GetReservation retrieves a reservation by its ID
	GetReservation(reservationID string) (*domain.Reservation, error)
	
	// CheckIn performs the check-in process for a reservation and assigns a seat
	CheckIn(reservationID string, seatNumber string) error
	
	// GetReservationsForFlight retrieves all reservations for a specific flight
	GetReservationsForFlight(flightNumber string) ([]*domain.Reservation, error)
}

type ValidationService interface {
	// ValidateFlightNumber checks if a flight number matches the required format
	ValidateFlightNumber(flightNumber string) bool
	
	// ValidateDates checks if the departure and arrival dates are valid
	ValidateDates(departureTime, arrivalTime time.Time) bool
	
	// CheckYesOrNo prompts the user with a yes/no question and returns the result
	CheckYesOrNo(prompt string) bool
}