package ports

import (
	"golang-airplane/internal/core/domain"
)

type AirplaneRepository interface {
	Save(airplane domain.Airplane) error
	FindByID(id string) (domain.Airplane, error)
	FindAll() ([]domain.Airplane, error)
}

// FlightRepository defines the interface for flight data operations
type FlightRepository interface {
	// FindAll returns all flights in the repository
	FindAll() ([]*domain.Flight, error)
	
	// FindByID finds a flight by its flight number
	FindByID(flightNumber string) (*domain.Flight, error)
	
	// SearchFlights searches for flights by location (departure or destination) and date
	SearchFlights(location string, dateStr string) ([]*domain.Flight, error)
	
	// Save stores a flight in the repository
	Save(flight *domain.Flight) error
	
	// Update updates an existing flight in the repository
	Update(flight *domain.Flight) error
}

// ReservationRepository defines the interface for reservation data operations
type ReservationRepository interface {
	// FindAll returns all reservations in the repository
	FindAll() ([]*domain.Reservation, error)
	
	// FindByID finds a reservation by its ID
	FindByID(reservationID string) (*domain.Reservation, error)
	
	// FindByFlightNumber finds all reservations for a specific flight
	FindByFlightNumber(flightNumber string) ([]*domain.Reservation, error)
	
	// Save stores a reservation in the repository
	Save(reservation *domain.Reservation) error
	
	// Update updates an existing reservation in the repository
	Update(reservation *domain.Reservation) error
}