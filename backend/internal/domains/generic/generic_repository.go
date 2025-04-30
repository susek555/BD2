package generic

import "gorm.io/gorm"

type GormRepository[T any] struct {
	db *gorm.DB
}

func GetGormRepository[T any](dbHandle *gorm.DB) *GormRepository[T] {
	return &GormRepository[T]{db: dbHandle}
}

func (repo *GormRepository[T]) Create(entity T) error {
	err := repo.db.Create(&entity).Error
	return err
}

func (repo *GormRepository[T]) GetAll() ([]T, error) {
	var entities []T
	err := repo.db.Find(&entities).Error
	return entities, err
}

func (repo *GormRepository[T]) GetById(id uint) (T, error) {
	var entity T
	err := repo.db.First(&entity, id).Error
	return entity, err
}

func (repo *GormRepository[T]) Update(entity T) error {
	err := repo.db.Save(&entity).Error
	return err
}

func (repo *GormRepository[T]) Delete(id uint) error {
	var entity T
	err := repo.db.Delete(&entity, id).Error
	return err
}
