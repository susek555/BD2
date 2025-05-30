package enums

import "database/sql/driver"

type FuelType string

const (
	DIESEL   FuelType = "Diesel"
	PETROL   FuelType = "Petrol"
	ELECTRIC FuelType = "Electric"
	ETHANOL  FuelType = "Ethanol"
	LPG      FuelType = "LPG"
	BIOFUEL  FuelType = "Biofuel"
	HYBRID   FuelType = "Hybrid"
	HYDROGEN FuelType = "Hydrogen"
)

func (f *FuelType) Scan(value any) error {
	*f = FuelType(value.([]byte))
	return nil
}

func (f FuelType) Value() (driver.Value, error) {
	return string(f), nil
}

var Types = []FuelType{
	DIESEL, PETROL, ELECTRIC, ETHANOL,
	LPG, BIOFUEL, HYBRID, HYDROGEN}
