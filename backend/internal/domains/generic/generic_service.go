package generic

type GenericService[T CRUDRepository[T]] struct {
	Repo CRUDRepository[T]
}

func (service GenericService[T]) Create(entity T) error {
	err := service.Repo.Create(entity)
	return err
}

func (service GenericService[T]) GetAll() ([]T, error) {
	return service.Repo.GetAll()
}

func (service GenericService[T]) GetByID(id uint) (T, error) {
	return service.Repo.GetByID(id)
}

func (service GenericService[T]) Update(entity T) error {
	return service.Repo.Update(entity)
}

func (service GenericService[T]) Delete(id uint) error {
	return service.Repo.Delete(id)
}
