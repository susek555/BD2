package model

type ModelServiceInterace interface {
	GetByManufactuerID(id uint) ([]RetrieveModelDTO, error)
	GetByManufacuterName(name string) ([]RetrieveModelDTO, error)
}

type ModelService struct {
	repo ModelRepositoryInterface
}

func NewModelService(modelRepository ModelRepositoryInterface) ModelServiceInterace {
	return &ModelService{repo: modelRepository}
}

func (s *ModelService) GetByManufactuerID(id uint) ([]RetrieveModelDTO, error) {
	models, err := s.repo.GetByManufactuerID(id)
	if err != nil {
		return nil, err
	}
	return MapModelListToDTO(models), nil
}

func (s *ModelService) GetByManufacuterName(name string) ([]RetrieveModelDTO, error) {
	models, err := s.repo.GetByManufacuterName(name)
	if err != nil {
		return nil, err
	}
	return MapModelListToDTO(models), nil
}

func MapModelListToDTO(models []Model) []RetrieveModelDTO {
	modelDTOs := make([]RetrieveModelDTO, 0, len(models))
	for _, model := range models {
		modelDTOs = append(modelDTOs, model.MapToDTO())
	}
	return modelDTOs
}
