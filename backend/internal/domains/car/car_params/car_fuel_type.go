package car_params

type FuelType string

const (
	DIESEL   FuelType = "Diesel"
	PETROL   FuelType = "Petrol"
	ELECTRIC FuelType = "Electric"
	ETHANOL  FuelType = "Ethanol"
	LPG      FuelType = "LPG"
	BIOFUEL  FuelType = "Biofuel"
	HYBRID   FuelType = "Hybird"
	HYDROGEN FuelType = "Hydrogen"
)

var Types = []FuelType{
	DIESEL, PETROL, ELECTRIC, ETHANOL,
	LPG, BIOFUEL, HYBRID, HYDROGEN}
