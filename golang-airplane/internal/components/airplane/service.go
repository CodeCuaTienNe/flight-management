package airplane

import (
	"errors"
	"golang-airplane/internal/core/domain"
	"golang-airplane/internal/core/ports"
)

// AirplaneService implements the service for airplane operations
type AirplaneService struct {
	repo ports.AirplaneRepository
}

// NewAirplaneService creates a new airplane service
func NewAirplaneService(repo ports.AirplaneRepository) *AirplaneService {
	return &AirplaneService{repo: repo}
}

// AddAirplane creates and stores a new airplane
func (s *AirplaneService) AddAirplane(id string, model string, capacity int) error {
	if model == "" || capacity <= 0 {
		return errors.New("invalid airplane data: model cannot be empty and capacity must be positive")
	}
	
	airplane := domain.Airplane{
		ID:       id,
		Model:    model,
		Capacity: capacity,
	}
	
	return s.repo.Save(airplane)
}

// GetAirplanes retrieves all airplanes
func (s *AirplaneService) GetAirplanes() ([]domain.Airplane, error) {
	return s.repo.FindAll()
}

// GetAirplaneByID retrieves an airplane by its ID
func (s *AirplaneService) GetAirplaneByID(id string) (domain.Airplane, error) {
	return s.repo.FindByID(id)
}