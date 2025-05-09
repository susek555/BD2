package car_params

type Model struct {
	ID             uint         `json:"id" gorm:"primaryKey"`
	Name           string       `json:"name"`
	ManufacturerID uint         `json:"manufacturer_id"`
	Manufacturer   Manufacturer `gorm:"foreignKey:ManufacturerID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
