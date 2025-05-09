package car

import "time"

type Car struct {
	ID                 uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	ProductionYear     uint      `json:"production_year"`
	Mileage            uint      `json:"mileage"`
	NumberOfDoors      uint      `json:"number_of_door" gorm:"check:number_of_doors BETWEEN 1 AND 6"`
	NumberOfSeats      uint      `json:"number_of_seats" gorm:"check:number_of_seats BETWEEN 2 AND 100"`
	EnginePower        uint      `json:"engine_power" gorm:"check:engine_power <= 9999"`
	EngineCapacity     uint      `json:"engine_capacity" gorm:"check enginge_capacity <= 9000"`
	RegistrationDate   time.Time `json:"registration_date"`
	RegistrationNumber string    `json:"registration_number"`
}
