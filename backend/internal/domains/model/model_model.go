package model

import "github.com/susek555/BD2/car-dealer-api/internal/domains/manufacturer"

type Model struct {
	ID             uint                      `json:"id" gorm:"primaryKey"`
	Name           string                    `json:"name"`
	ManufacturerID uint                      `json:"manufacturer_id"`
	Manufacturer   manufacturer.Manufacturer `gorm:"foreignKey:ManufacturerID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
