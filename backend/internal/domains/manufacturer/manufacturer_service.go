package manufacturer

type ManufacturerServiceInterface interface {
	GetAll() ([]Manufacturer, error)
	GetAllAsNames() ([]string, error)
}

type ManufacturerService struct {
	repo ManufacturerRepositoryInterface
}

func NewManufacturerService(manufacturerRepository ManufacturerRepositoryInterface) ManufacturerServiceInterface {
	return &ManufacturerService{repo: manufacturerRepository}
}

func (s *ManufacturerService) GetAll() ([]Manufacturer, error) {
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
