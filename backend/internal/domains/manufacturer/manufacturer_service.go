package manufacturer

import "github.com/susek555/BD2/car-dealer-api/internal/domains/models"

type ManufacturerServiceInterface interface {
	GetAll() ([]models.Manufacturer, error)
	GetAllAsNames() ([]string, error)
}

type ManufacturerService struct {
	repo ManufacturerRepositoryInterface
}

func NewManufacturerService(manufacturerRepository ManufacturerRepositoryInterface) ManufacturerServiceInterface {
	return &ManufacturerService{repo: manufacturerRepository}
}

func (s *ManufacturerService) GetAll() ([]models.Manufacturer, error) {
	return s.repo.GetAll()
}

func (s *ManufacturerService) GetAllAsNames() ([]string, error) {
	manufacturers, err := s.GetAll()
	if err != nil {
		return nil, err
	}
	result := make([]string, 0, len(manufacturers))
	for _, man := range manufacturers {
		result = append(result, man.Name)
	}
	return result, nil
}
