package model

type ModelServiceInterace interface {
	GetByManufactuerID(id uint) ([]Model, error)
	GetByManufacuterName(name string) ([]Model, error)
}

type ModelService struct {
	repo ModelRepositoryInterface
}

func NewModelService(modelRepository ModelRepositoryInterface) ModelServiceInterace {
	return &ModelService{repo: modelRepository}
}

func (s *ModelService) GetByManufactuerID(id uint) ([]Model, error) {
	return s.repo.GetByManufactuerID(id)
}

func (s *ModelService) GetByManufacuterName(name string) ([]Model, error) {
	return s.repo.GetByManufacuterName(name)
}
