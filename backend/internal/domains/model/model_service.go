package model

import (
	"github.com/susek555/BD2/car-dealer-api/pkg/mapping"
)

type ModelServiceInterface interface {
	GetByManufacturerID(id uint) ([]RetrieveModelDTO, error)
	GetByManufacturerName(name string) ([]RetrieveModelDTO, error)
}

type ModelService struct {
	repo ModelRepositoryInterface
}

func NewModelService(modelRepository ModelRepositoryInterface) ModelServiceInterface {
	return &ModelService{repo: modelRepository}
}

func (s *ModelService) GetByManufacturerID(id uint) ([]RetrieveModelDTO, error) {
	models, err := s.repo.GetByManufacturerID(id)
	if err != nil {
		return nil, err
	}
	return mapping.MapSliceToDTOs(models, MapToDTO), nil
}

func (s *ModelService) GetByManufacturerName(name string) ([]RetrieveModelDTO, error) {
	models, err := s.repo.GetByManufacturerName(name)
	if err != nil {
		return nil, err
	}
	return mapping.MapSliceToDTOs(models, MapToDTO), nil
}
