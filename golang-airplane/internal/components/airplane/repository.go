package airplane

import (
	"errors"
	"golang-airplane/internal/core/domain"
	"golang-airplane/internal/core/ports"
)

// AirplaneRepositoryImpl implements the ports.AirplaneRepository interface
type AirplaneRepositoryImpl struct {
	storage interface {
		Load(filename string, target interface{}) error
		Save(filename string, data interface{}) error
	}
}

// NewAirplaneRepository creates a new repository instance
func NewAirplaneRepository(storage interface {
	Load(filename string, target interface{}) error
	Save(filename string, data interface{}) error
}) ports.AirplaneRepository {
	return &AirplaneRepositoryImpl{
		storage: storage,
	}
}

// Save stores an airplane in the repository
func (r *AirplaneRepositoryImpl) Save(airplane domain.Airplane) error {
	airplanesMap := make(map[string]*domain.Airplane)
	
	// Load existing airplanes
	err := r.storage.Load("airplanes.json", &airplanesMap)
	if err != nil {
		return err
	}
	
	// Add or update airplane
	airplanesMap[airplane.ID] = &airplane
	
	// Save updated airplanes map
	return r.storage.Save("airplanes.json", airplanesMap)
}

// FindByID finds an airplane by its ID
func (r *AirplaneRepositoryImpl) FindByID(id string) (domain.Airplane, error) {
	airplanesMap := make(map[string]*domain.Airplane)
	err := r.storage.Load("airplanes.json", &airplanesMap)
	if err != nil {
		return domain.Airplane{}, err
	}
	
	airplane, exists := airplanesMap[id]
	if !exists {
		return domain.Airplane{}, errors.New("airplane not found")
	}
	
	return *airplane, nil
}

// FindAll returns all airplanes in the repository
func (r *AirplaneRepositoryImpl) FindAll() ([]domain.Airplane, error) {
	airplanesMap := make(map[string]*domain.Airplane)
	err := r.storage.Load("airplanes.json", &airplanesMap)
	if err != nil {
		return nil, err
	}

	// Convert map to slice
	airplanes := make([]domain.Airplane, 0, len(airplanesMap))
	for _, airplane := range airplanesMap {
		airplanes = append(airplanes, *airplane)
	}
	
	return airplanes, nil
}