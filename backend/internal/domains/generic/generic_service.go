package generic

type GenericService[T any, R CRUDService[T]] struct {
	Repo R
}

type CRUDService[T any] interface {
	CRUDRepository[T]
}

func (service GenericService[T, R]) Create(entity *T) error {
	err := service.Repo.Create(entity)
	return err
}

func (service GenericService[T, R]) GetAll() ([]T, error) {
	return service.Repo.GetAll()
}

func (service GenericService[T, R]) GetById(id uint) (T, error) {
	return service.Repo.GetById(id)
}

func (service GenericService[T, R]) Update(entity *T) error {
	return service.Repo.Update(entity)
}

func (service GenericService[T, R]) Delete(id uint) error {
	return service.Repo.Delete(id)
}
