package generic

type GenericService[T any, R CRUDRepository[T]] struct {
	Repo R
}

func (service GenericService[T, R]) Create(entity T) error {
	err := service.Repo.Create(entity)
	return err
}

func (service GenericService[T, R]) GetAll() ([]T, error) {
	return service.Repo.GetAll()
}

func (service GenericService[T, R]) GetByID(id uint) (T, error) {
	return service.Repo.GetByID(id)
}

func (service GenericService[T, R]) Update(entity T) error {
	return service.Repo.Update(entity)
}

func (service GenericService[T, R]) Delete(id uint) error {
	return service.Repo.Delete(id)
}
