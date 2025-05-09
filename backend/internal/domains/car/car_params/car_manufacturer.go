package car_params

type Manufacturer struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name" `
}
