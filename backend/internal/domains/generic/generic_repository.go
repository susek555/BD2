package generic

import "gorm.io/gorm"

type GormRepository[T any] struct {
	DB *gorm.DB
}

type CRUDRepository[T any] interface {
	Create(entity T) error
	GetAll() ([]T, error)
	GetById(id uint) (T, error)
	Update(entity T) error
	Delete(id uint) error
}

func GetGormRepository[T any](dbHandle *gorm.DB) *GormRepository[T] {
	return &GormRepository[T]{DB: dbHandle}
}

func (repo *GormRepository[T]) Create(entity T) error {
	err := repo.DB.Create(&entity).Error
	return err
}

func (repo *GormRepository[T]) GetAll() ([]T, error) {
	var entities []T
	err := repo.DB.Find(&entities).Error
	return entities, err
}

func (repo *GormRepository[T]) GetById(id uint) (T, error) {
	var entity T
	err := repo.DB.First(&entity, id).Error
	return entity, err
}

func (repo *GormRepository[T]) Update(entity T) error {
	err := repo.DB.Save(&entity).Error
	return err
}

func (repo *GormRepository[T]) Delete(id uint) error {
	var entity T
	err := repo.DB.Delete(&entity, id).Error
	return err
}
