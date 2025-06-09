package sale_offer

import (
	"github.com/susek555/BD2/car-dealer-api/internal/enums"
	"github.com/susek555/BD2/car-dealer-api/pkg/pagination"
)

type CreateSaleOfferDTO struct {
	UserID             uint               `json:"-"`
	Description        string             `json:"description" validate:"required"`
	Price              uint               `json:"price" validate:"required"`
	Margin             enums.MarginValue  `json:"margin" validate:"required"`
	Vin                string             `json:"vin" validate:"required"`
	ProductionYear     uint               `json:"production_year" validate:"required"`
	Mileage            uint               `json:"mileage" validate:"required"`
	NumberOfDoors      uint               `json:"number_of_doors" validate:"required"`
	NumberOfSeats      uint               `json:"number_of_seats" validate:"required"`
	EnginePower        uint               `json:"engine_power" validate:"required"`
	EngineCapacity     uint               `json:"engine_capacity" validate:"required"`
	RegistrationNumber string             `json:"registration_number" validate:"required"`
	RegistrationDate   string             `json:"registration_date" validate:"required"`
	Color              enums.Color        `json:"color" validate:"required"`
	FuelType           enums.FuelType     `json:"fuel_type" validate:"required"`
	Transmission       enums.Transmission `json:"transmission" validate:"required"`
	NumberOfGears      uint               `json:"number_of_gears" validate:"required"`
	Drive              enums.Drive        `json:"drive" validate:"required"`
	ManufacturerName   string             `json:"manufacturer" validate:"required"`
	ModelName          string             `json:"model" validate:"required"`
}

type UpdateSaleOfferDTO struct {
	ID                 uint                `json:"id"`
	Description        *string             `json:"description"`
	Price              *uint               `json:"price"`
	Margin             *enums.MarginValue  `json:"margin"`
	Vin                *string             `json:"vin"`
	ProductionYear     *uint               `json:"production_year"`
	Mileage            *uint               `json:"mileage"`
	NumberOfDoors      *uint               `json:"number_of_doors"`
	NumberOfSeats      *uint               `json:"number_of_seats"`
	EnginePower        *uint               `json:"engine_power"`
	EngineCapacity     *uint               `json:"engine_capacity"`
	RegistrationNumber *string             `json:"registration_number"`
	RegistrationDate   *string             `json:"registration_date"`
	Color              *enums.Color        `json:"color"`
	FuelType           *enums.FuelType     `json:"fuel_type"`
	Transmission       *enums.Transmission `json:"transmission"`
	NumberOfGears      *uint               `json:"number_of_gears"`
	Drive              *enums.Drive        `json:"drive"`
	ManufacturerName   *string             `json:"manufacturer"`
	ModelName          *string             `json:"model"`
}

type UserContext struct {
	IsLiked   bool `json:"is_liked"`
	CanModify bool `json:"can_modify"`
}

type RetrieveSaleOfferDTO struct {
	ID             uint         `json:"id"`
	UserID         uint         `json:"user_id"`
	Username       string       `json:"username"`
	Name           string       `json:"name"`
	Price          uint         `json:"price"`
	Mileage        uint         `json:"mileage"`
	ProductionYear uint         `json:"production_year"`
	Color          enums.Color  `json:"color"`
	MainURL        string       `json:"main_url"`
	IsAuction      bool         `json:"is_auction"`
	Status         enums.Status `json:"status"`
	IssueDate      *string      `json:"issue_date,omitempty"`
	UserContext
}

type RetrieveDetailedSaleOfferDTO struct {
	ID                 uint               `json:"id"`
	UserID             uint               `json:"user_id"`
	Username           string             `json:"username"`
	Description        string             `json:"description"`
	Price              uint               `json:"price"`
	DateOfIssue        string             `json:"date_of_issue"`
	Margin             enums.MarginValue  `json:"margin"`
	Status             enums.Status       `json:"status,omitempty"`
	Vin                string             `json:"vin"`
	ProductionYear     uint               `json:"production_year"`
	Mileage            uint               `json:"mileage"`
	NumberOfDoors      uint               `json:"number_of_doors"`
	NumberOfSeats      uint               `json:"number_of_seats"`
	EnginePower        uint               `json:"engine_power"`
	EngineCapacity     uint               `json:"engine_capacity"`
	RegistrationNumber string             `json:"registration_number"`
	RegistrationDate   string             `json:"registration_date"`
	Color              enums.Color        `json:"color"`
	FuelType           enums.FuelType     `json:"fuel_type"`
	Transmission       enums.Transmission `json:"transmission"`
	NumberOfGears      uint               `json:"number_of_gears"`
	Drive              enums.Drive        `json:"drive"`
	Brand              string             `json:"brand"`
	Model              string             `json:"model"`
	ImagesUrls         []string           `json:"images_urls"`
	IsAuction          bool               `json:"is_auction"`
	DateEnd            *string            `json:"date_end,omitempty"`
	BuyNowPrice        *uint              `json:"buy_now_price,omitempty"`
	IssueDate          *string            `json:"issue_date,omitempty"`
	UserContext
}

type RetrieveOffersWithPagination struct {
	PaginationResponse pagination.PaginationResponse `json:"pagination"`
	Offers             []RetrieveSaleOfferDTO        `json:"offers"`
}
