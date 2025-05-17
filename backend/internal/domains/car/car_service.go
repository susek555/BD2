package car

import (
	"github.com/susek555/BD2/car-dealer-api/internal/domains/manufacturer"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/model"
	"github.com/susek555/BD2/car-dealer-api/pkg/mapping"
)

type CarServiceInterface interface {
	GetManufacturersModelsMap() (*ManufacturerModelMap, error)
}

type CarService struct {
	manufactuterRepo manufacturer.ManufacturerRepositoryInterface
	modelRepo        model.ModelRepositoryInterface
}

func NewCarService(manufacuterR manufacturer.ManufacturerRepositoryInterface, modelR model.ModelRepositoryInterface) CarServiceInterface {
	return &CarService{manufactuterRepo: manufacuterR, modelRepo: modelR}
}

func (s *CarService) GetManufacturersModelsMap() (*ManufacturerModelMap, error) {
	manufacturers, err := s.manufactuterRepo.GetAll()
	if err != nil {
		return nil, err
	}
	manufacturersNames := mapping.MapSliceToDTOs(manufacturers, (*manufacturer.Manufacturer).MapToName)
	var modelsNames [][]string
	for _, manufacturer := range manufacturers {
		models, err := s.modelRepo.GetByManufacturerID(manufacturer.ID)
		if err != nil {
			return nil, err
		}
		modelsNames = append(modelsNames, mapping.MapSliceToDTOs(models, (*model.Model).MapToName))
	}
	return &ManufacturerModelMap{Manufacturers: manufacturersNames, Models: modelsNames}, nil
}
