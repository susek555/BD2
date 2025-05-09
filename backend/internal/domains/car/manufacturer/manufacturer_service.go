package manufacturer

type ManufacturerServiceInterface interface {
	GetAll() ([]Manufacturer, error)
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
