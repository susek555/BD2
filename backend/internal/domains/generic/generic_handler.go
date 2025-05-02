package generic

type GenericHandler[T any, S CRUDService[T]] struct {
	Service CRUDService[T]
}

func (handler GenericHandler[T, S]) Create(entity T) error {
	err := handler.Service.Create(entity)
	return err
}

func (handler GenericHandler[T, S]) GetByID(id uint) (T, error) {
	return handler.Service.GetByID(id)
}

func (handler GenericHandler[T, S]) GetAll() ([]T, error) {
	return handler.Service.GetAll()
}

func (handler GenericHandler[T, S]) Update(entity T) error {
	return handler.Service.Update(entity)
}

func (handler GenericHandler[T, S]) Delete(id uint) error {
	return handler.Service.Delete(id)
}
