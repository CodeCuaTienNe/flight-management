package flight

import (
	"fmt"
	"golang-airplane/internal/core/domain"
	"golang-airplane/internal/core/ports"
)

// ReservationService implements the ReservationService interface
type ReservationService struct {
	flightRepo      ports.FlightRepository
	reservationRepo ports.ReservationRepository
}

// NewReservationService creates a new ReservationService instance
func NewReservationService(flightRepo ports.FlightRepository, reservationRepo ports.ReservationRepository) *ReservationService {
	return &ReservationService{
		flightRepo:      flightRepo,
		reservationRepo: reservationRepo,
	}
}

// BookFlight creates a new reservation for a flight
func (s *ReservationService) BookFlight(name, address string, phoneNumber, identityCardNumber int64, flightNumber string) (*domain.Reservation, error) {
	// Verify that the flight exists
	flight, err := s.flightRepo.FindByID(flightNumber)
	if err != nil {
		return nil, fmt.Errorf("flight not found: %w", err)
	}
	
	// Check if there are available seats
	if flight.AvailableSeat <= 0 {
		return nil, fmt.Errorf("no available seats for flight %s", flightNumber)
	}
	
	// Create new reservation
	reservation := domain.NewReservation(name, address, phoneNumber, identityCardNumber, flightNumber)
	
	// Save the reservation
	err = s.reservationRepo.Save(reservation)
	if err != nil {
		return nil, fmt.Errorf("failed to save reservation: %w", err)
	}
	
	// Update flight available seats
	flight.AvailableSeat--
	err = s.flightRepo.Update(flight)
	if err != nil {
		return nil, fmt.Errorf("failed to update flight: %w", err)
	}
	
	return reservation, nil
}

// GetReservation retrieves a reservation by its ID
func (s *ReservationService) GetReservation(reservationID string) (*domain.Reservation, error) {
	return s.reservationRepo.FindByID(reservationID)
}

// CheckIn performs the check-in process for a reservation and assigns a seat
func (s *ReservationService) CheckIn(reservationID string, seatNumber string) error {
	// Get the reservation
	reservation, err := s.reservationRepo.FindByID(reservationID)
	if err != nil {
		return fmt.Errorf("reservation not found: %w", err)
	}
	
	// Don't allow check-in if already checked in
	if reservation.CheckedIn {
		return fmt.Errorf("reservation %s is already checked in", reservationID)
	}
	
	// Get the flight for this reservation
	flight, err := s.flightRepo.FindByID(reservation.ReservationFlightNumber)
	if err != nil {
		return fmt.Errorf("flight not found: %w", err)
	}
	
	// Check if the seat is available
	if available, exists := flight.SeatList[seatNumber]; !exists || !available {
		return fmt.Errorf("seat %s is not available", seatNumber)
	}
	
	// Mark the seat as occupied
	flight.SeatList[seatNumber] = false
	
	// Update the flight
	err = s.flightRepo.Update(flight)
	if err != nil {
		return fmt.Errorf("failed to update flight: %w", err)
	}
	
	// Assign the seat and mark as checked in
	reservation.SeatLocation = seatNumber
	reservation.CheckIn()
	
	// Update the reservation
	err = s.reservationRepo.Update(reservation)
	if err != nil {
		return fmt.Errorf("failed to update reservation: %w", err)
	}
	
	return nil
}

// GetReservationsForFlight retrieves all reservations for a specific flight
func (s *ReservationService) GetReservationsForFlight(flightNumber string) ([]*domain.Reservation, error) {
	// First check if the flight exists
	_, err := s.flightRepo.FindByID(flightNumber)
	if err != nil {
		return nil, fmt.Errorf("flight not found: %w", err)
	}
	
	// Find all reservations for the flight
	return s.reservationRepo.FindByFlightNumber(flightNumber)
}