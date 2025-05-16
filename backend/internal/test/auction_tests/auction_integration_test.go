package auction_test

import (
	"github.com/gin-gonic/gin"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/auction"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/manufacturer"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/model"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/sale_offer"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// ------
// Setup
// ------

func setupDB(manufacturers []manufacturer.Manufacturer, models []model.Model, cars []sale_offer.Car, saleOffers []sale_offer.SaleOffer, auctions []sale_offer.Auction) (auction.AuctionServiceInterface, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(
		&manufacturer.Manufacturer{},
		&model.Model{},
		&sale_offer.Car{},
		&sale_offer.SaleOffer{},
		&sale_offer.Auction{},
	)
	if err != nil {
		return nil, err
	}
	err = db.Create(&manufacturers).Error
	if err != nil {
		return nil, err
	}
	err = db.Create(&models).Error
	if err != nil {
		return nil, err
	}
	err = db.Create(&cars).Error
	if err != nil {
		return nil, err
	}
	err = db.Create(&saleOffers).Error
	if err != nil {
		return nil, err
	}
	err = db.Create(&auctions).Error
	if err != nil {
		return nil, err
	}
	repo := auction.NewAuctionRepository(db)
	service := auction.NewAuctionService(repo)
	return service, nil
}

func newTestServer(seedManufacturers []manufacturer.Manufacturer, seedModels []model.Model, seedCars []sale_offer.Car, seedSaleOffers []sale_offer.SaleOffer, seedAuctions []sale_offer.Auction) (*gin.Engine, auction.AuctionServiceInterface, error) {
	service, err := setupDB(seedManufacturers, seedModels, seedCars, seedSaleOffers, seedAuctions)
	if err != nil {
		return nil, nil, err
	}
	r := gin.Default()
	auctionHandler := auction.NewHandler(service)
	auctionRoutes := r.Group("/auction")
	auctionRoutes.GET("/", auctionHandler.GetAllAuctions)
	auctionRoutes.GET("/:id", auctionHandler.GetAuctionById)
	auctionRoutes.POST("/", auctionHandler.CreateAuction)
	auctionRoutes.PUT("/", auctionHandler.UpdateAuction)
	auctionRoutes.DELETE("/:id", auctionHandler.DeleteAuctionById)
	return r, service, nil
}


