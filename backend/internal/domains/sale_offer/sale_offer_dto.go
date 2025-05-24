package sale_offer

import (
	"time"

	"github.com/susek555/BD2/car-dealer-api/internal/domains/models"

	"github.com/susek555/BD2/car-dealer-api/internal/domains/car/car_params"
	"github.com/susek555/BD2/car-dealer-api/pkg/pagination"
)

type CreateSaleOfferDTO struct {
	UserID             uint                    `json:"-"`
	Description        string                  `json:"description" validate:"required"`
	Price              uint                    `json:"price" validate:"required"`
	Margin             models.MarginValue      `json:"margin" validate:"required"`
	Vin                string                  `json:"vin" validate:"required"`
	ProductionYear     uint                    `json:"production_year" validate:"required"`
	Mileage            uint                    `json:"mileage" validate:"required"`
	NumberOfDoors      uint                    `json:"number_of_doors" validate:"required"`
	NumberOfSeats      uint                    `json:"number_of_seats" validate:"required"`
	EnginePower        uint                    `json:"engine_power" validate:"required"`
	EngineCapacity     uint                    `json:"engine_capacity" validate:"required"`
	RegistrationNumber string                  `json:"registration_number" validate:"required"`
	RegistrationDate   string                  `json:"registration_date" validate:"required"`
	Color              car_params.Color        `json:"color" validate:"required"`
	FuelType           car_params.FuelType     `json:"fuel_type" validate:"required"`
	Transmission       car_params.Transmission `json:"transmission" validate:"required"`
	NumberOfGears      uint                    `json:"number_of_gears" validate:"required"`
	Drive              car_params.Drive        `json:"drive" validate:"required"`
	Manufacturer       string                  `json:"manufacturer" validate:"required"`
	Model              string                  `json:"model" validate:"required"`
}

type UpdateSaleOfferDTO struct {
	ID                 uint                     `json:"id"`
	Description        *string                  `json:"description"`
	Price              *uint                    `json:"price"`
	Margin             *models.MarginValue      `json:"margin"`
	Vin                *string                  `json:"vin"`
	ProductionYear     *uint                    `json:"production_year"`
	Mileage            *uint                    `json:"mileage"`
	NumberOfDoors      *uint                    `json:"number_of_doors"`
	NumberOfSeats      *uint                    `json:"number_of_seats"`
	EnginePower        *uint                    `json:"engine_power"`
	EngineCapacity     *uint                    `json:"engine_capacity"`
	RegistrationNumber *string                  `json:"registration_number"`
	RegistrationDate   *string                  `json:"registration_date"`
	Color              *car_params.Color        `json:"color"`
	FuelType           *car_params.FuelType     `json:"fuel_type"`
	Transmission       *car_params.Transmission `json:"transmission"`
	NumberOfGears      *uint                    `json:"number_of_gears"`
	Drive              *car_params.Drive        `json:"drive"`
	Manufacturer       *string                  `json:"manufacturer"`
	Model              *string                  `json:"model"`
	Status             *models.Status           `json:"-"`
}

type UserContext struct {
	IsLiked   bool `json:"is_liked"`
	CanModify bool `json:"can_modify"`
}

type RetrieveSaleOfferDTO struct {
	ID             uint             `json:"id"`
	Username       string           `json:"username"`
	Name           string           `json:"name"`
	Price          uint             `json:"price"`
	Mileage        uint             `json:"mileage"`
	ProductionYear uint             `json:"production_year"`
	Color          car_params.Color `json:"color"`
	IsAuction      bool             `json:"is_auction"`
	UserContext
}

type RetrieveDetailedSaleOfferDTO struct {
	ID                 uint                    `json:"id"`
	UserID             uint                    `json:"user_id"`
	Username           string                  `json:"username"`
	Description        string                  `json:"description"`
	Price              uint                    `json:"price"`
	Margin             models.MarginValue      `json:"margin"`
	DateOfIssue        string                  `json:"date_of_issue"`
	Vin                string                  `json:"vin"`
	ProductionYear     uint                    `json:"production_year"`
	Mileage            uint                    `json:"mileage"`
	NumberOfDoors      uint                    `json:"number_of_doors"`
	NumberOfSeats      uint                    `json:"number_of_seats"`
	EnginePower        uint                    `json:"engine_power"`
	EngineCapacity     uint                    `json:"engine_capacity"`
	RegistrationNumber string                  `json:"registration_number"`
	RegistrationDate   string                  `json:"registration_date"`
	Color              car_params.Color        `json:"color"`
	FuelType           car_params.FuelType     `json:"fuel_type"`
	Transmission       car_params.Transmission `json:"transmission"`
	NumberOfGears      uint                    `json:"number_of_gears"`
	Drive              car_params.Drive        `json:"drive"`
	Brand              string                  `json:"brand"`
	Model              string                  `json:"model"`
	Status             models.Status           `json:"status"`
	ImagesUrls         []string                `json:"images_urls"`
	DateEnd            *time.Time              `json:"date_end,omitempty"`
	BuyNowPrice        *uint                   `json:"buy_now_price,omitempty"`
	IsAuction          bool                    `json:"is_auction"`
	UserContext
}

type RetrieveOffersWithPagination struct {
	PaginationResponse *pagination.PaginationResponse `json:"pagination"`
	Offers             []RetrieveSaleOfferDTO         `json:"offers"`
}
