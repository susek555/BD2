package car_params

type FuelType string

const (
	DIESEL   FuelType = "Diesel"
	PETROL   FuelType = "Petrol"
	ELECTIRC FuelType = "Electric"
	ETHANOL  FuelType = "Ethanol"
	LPG      FuelType = "LPG"
	BIOFUEL  FuelType = "Biofuel"
	HYBRID   FuelType = "Hybird"
	HYDROGEN FuelType = "Hydrogen"
)

var Types = []FuelType{
	DIESEL, PETROL, ELECTIRC, ETHANOL,
	LPG, BIOFUEL, HYBRID, HYDROGEN}
