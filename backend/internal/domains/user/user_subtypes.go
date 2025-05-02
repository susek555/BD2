package user

import "gorm.io/gorm"

type Company struct {
	UserID uint   `json:"id" gorm:"primaryKey"`
	NIP    string `json:"nip"`
	Name   string `json:"name"`
	User   User   `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
}

type Person struct {
	UserID  uint   `json:"id" gorm:"primaryKey"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
	User    User   `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
}

type UserSubtype interface {
	SetUserID(id uint)
	SaveSubtype(db *gorm.DB) error
}

func (c *Company) SetUserID(id uint) {
	c.UserID = id
}

func (c *Company) SaveSubtype(tx *gorm.DB) error {
	return tx.Save(c).Error
}

func (p *Person) SetUserID(id uint) {
	p.UserID = id
}

func (p *Person) SaveSubtype(tx *gorm.DB) error {
	return tx.Save(p).Error
}
