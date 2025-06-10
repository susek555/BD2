package manufacturer

import (
	"github.com/susek555/BD2/car-dealer-api/pkg/mapping"
)

type ManufacturerServiceInterface interface {
	GetAll() ([]RetrieveManufacturerDTO, error)
	GetAllAsNames() ([]string, error)
}

type ManufacturerService struct {
	repo ManufacturerRepositoryInterface
}

func NewManufacturerService(manufacturerRepository ManufacturerRepositoryInterface) ManufacturerServiceInterface {
	return &ManufacturerService{repo: manufacturerRepository}
}

func (s *ManufacturerService) GetAll() ([]RetrieveManufacturerDTO, error) {
	manufacturers, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	return mapping.MapSliceToDTOs(manufacturers, MapToDTO), nil
}

func (s *ManufacturerService) GetAllAsNames() ([]string, error) {
	manufacturers, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	result := make([]string, 0, len(manufacturers))
	for _, man := range manufacturers {
		result = append(result, man.Name)
	}
	return result, nil
}
