package json

import (
	"fmt"
	"golang-airplane/internal/core/domain"
	"golang-airplane/internal/core/ports"
)

// ReservationRepositoryJSON implements the ReservationRepository interface using JSON files
type ReservationRepositoryJSON struct {
	storage *Storage
}

// NewReservationRepository creates a new ReservationRepositoryJSON instance
func NewReservationRepository(storage *Storage) ports.ReservationRepository {
	return &ReservationRepositoryJSON{
		storage: storage,
	}
}

// FindAll returns all reservations in the repository
func (r *ReservationRepositoryJSON) FindAll() ([]*domain.Reservation, error) {
	reservationsMap := make(map[string]*domain.Reservation)
	err := r.storage.Load("reservations.json", &reservationsMap)
	if err != nil {
		return nil, err
	}

	// Convert map to slice
	reservations := make([]*domain.Reservation, 0, len(reservationsMap))
	for _, res := range reservationsMap {
		reservations = append(reservations, res)
	}
	
	return reservations, nil
}

// FindByID finds a reservation by its ID
func (r *ReservationRepositoryJSON) FindByID(reservationID string) (*domain.Reservation, error) {
	reservationsMap := make(map[string]*domain.Reservation)
	err := r.storage.Load("reservations.json", &reservationsMap)
	if err != nil {
		return nil, err
	}
	
	reservation, exists := reservationsMap[reservationID]
	if !exists {
		return nil, fmt.Errorf("reservation with ID %s not found", reservationID)
	}
	
	return reservation, nil
}

// FindByFlightNumber finds all reservations for a specific flight
func (r *ReservationRepositoryJSON) FindByFlightNumber(flightNumber string) ([]*domain.Reservation, error) {
	// Get all reservations
	allReservations, err := r.FindAll()
	if err != nil {
		return nil, err
	}
	
	// Filter reservations by flight number
	var flightReservations []*domain.Reservation
	for _, res := range allReservations {
		if res.ReservationFlightNumber == flightNumber {
			flightReservations = append(flightReservations, res)
		}
	}
	
	return flightReservations, nil
}

// Save stores a reservation in the repository
func (r *ReservationRepositoryJSON) Save(reservation *domain.Reservation) error {
	reservationsMap := make(map[string]*domain.Reservation)
	
	// Load existing reservations
	err := r.storage.Load("reservations.json", &reservationsMap)
	if err != nil {
		return err
	}
	
	// Add or update reservation
	reservationsMap[reservation.ReservationID] = reservation
	
	// Save updated reservations map
	return r.storage.Save("reservations.json", reservationsMap)
}

// Update updates an existing reservation in the repository
func (r *ReservationRepositoryJSON) Update(reservation *domain.Reservation) error {
	// Update is the same as Save in this implementation
	return r.Save(reservation)
}