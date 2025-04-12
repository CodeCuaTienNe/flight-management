package json

import (
	"encoding/json"
	"fmt"
	"golang-airplane/internal/core/domain"
	"golang-airplane/internal/core/ports"
	"os"
	"path/filepath"
)

// ReservationRepositoryJSON implements the ReservationRepository interface using JSON files
type ReservationRepositoryJSON struct {
	storage *Storage
}

// NewReservationRepository creates a new ReservationRepositoryJSON instance
func NewReservationRepository(storage *Storage) ports.ReservationRepository {
	repo := &ReservationRepositoryJSON{
		storage: storage,
	}
	// Migrate data if needed
	repo.migrateDataIfNeeded()
	return repo
}

// migrateDataIfNeeded checks if the data is in the old format (map) and migrates it to the new format (slice)
func (r *ReservationRepositoryJSON) migrateDataIfNeeded() error {
	filePath := filepath.Join(r.storage.dataPath, "reservations.json")
	
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
	var reservations []*domain.Reservation
	err = json.Unmarshal(data, &reservations)
	if err == nil {
		// Data is already in the new format, no migration needed
		return nil
	}
	
	// Try to unmarshal as a map (old format)
	reservationsMap := make(map[string]*domain.Reservation)
	err = json.Unmarshal(data, &reservationsMap)
	if err != nil {
		return fmt.Errorf("failed to unmarshal data in either format: %w", err)
	}
	
	// Convert map to slice
	migratedReservations := make([]*domain.Reservation, 0, len(reservationsMap))
	for _, reservation := range reservationsMap {
		migratedReservations = append(migratedReservations, reservation)
	}
	
	// Save migrated data
	return r.storage.Save("reservations.json", migratedReservations)
}

// FindAll returns all reservations in the repository
func (r *ReservationRepositoryJSON) FindAll() ([]*domain.Reservation, error) {
	var reservations []*domain.Reservation
	
	// First try to load as a slice
	err := r.storage.Load("reservations.json", &reservations)
	if err == nil {
		return reservations, nil
	}
	
	// If that fails, try to load as a map and convert
	reservationsMap := make(map[string]*domain.Reservation)
	err = r.storage.Load("reservations.json", &reservationsMap)
	if err != nil {
		return nil, err
	}
	
	// Convert map to slice
	reservations = make([]*domain.Reservation, 0, len(reservationsMap))
	for _, reservation := range reservationsMap {
		reservations = append(reservations, reservation)
	}
	
	// Save in the new format for future use
	r.storage.Save("reservations.json", reservations)
	
	return reservations, nil
}

// FindByID finds a reservation by its ID
func (r *ReservationRepositoryJSON) FindByID(reservationID string) (*domain.Reservation, error) {
	reservations, err := r.FindAll()
	if err != nil {
		return nil, err
	}
	
	// Search for the reservation with the given ID
	for _, reservation := range reservations {
		if reservation.ReservationID == reservationID {
			return reservation, nil
		}
	}
	
	return nil, fmt.Errorf("reservation with ID %s not found", reservationID)
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
	reservations, err := r.FindAll()
	if err != nil {
		return err
	}
	
	// Check if the reservation already exists
	found := false
	for i, existingReservation := range reservations {
		if existingReservation.ReservationID == reservation.ReservationID {
			// Update existing reservation
			reservations[i] = reservation
			found = true
			break
		}
	}
	
	// Add new reservation if not found
	if !found {
		reservations = append(reservations, reservation)
	}
	
	// Save updated reservations list
	return r.storage.Save("reservations.json", reservations)
}

// Update updates an existing reservation in the repository
func (r *ReservationRepositoryJSON) Update(reservation *domain.Reservation) error {
	reservations, err := r.FindAll()
	if err != nil {
		return err
	}
	
	// Find and update the reservation
	found := false
	for i, existingReservation := range reservations {
		if existingReservation.ReservationID == reservation.ReservationID {
			reservations[i] = reservation
			found = true
			break
		}
	}
	
	if !found {
		return fmt.Errorf("reservation with ID %s not found", reservation.ReservationID)
	}
	
	// Save updated reservations list
	return r.storage.Save("reservations.json", reservations)
}