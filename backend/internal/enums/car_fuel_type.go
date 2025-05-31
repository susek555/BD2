package enums

import (
	"database/sql/driver"
)

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
	var sValue string
	switch v := value.(type) {
	case string:
		sValue = v
	case []byte:
		sValue = string(v)
	default:
		return ErrStringConversion
	}
	*f = FuelType(convertDBFormatToAppFormat(sValue, sValue == string(LPG)))
	return nil
}

func (f FuelType) Value() (driver.Value, error) {
	return convertAppFormatToDBFormat(string(f)), nil
}

var Types = []FuelType{
	DIESEL, PETROL, ELECTRIC, ETHANOL,
	LPG, BIOFUEL, HYBRID, HYDROGEN}
