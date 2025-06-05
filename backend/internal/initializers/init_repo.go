package initializers

import (
	"github.com/susek555/BD2/car-dealer-api/internal/domains/auction"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/bid"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/image"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/liked_offer"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/manufacturer"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/model"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/notification"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/purchase"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/refresh_token"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/review"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/sale_offer"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/user"
)

var AuctionRepo auction.AuctionRepositoryInterface
var BidRepo bid.BidRepositoryInterface
var ClientNotificationRepo notification.ClientNotificationRepositoryInterface
var ImageRepo image.ImageRepositoryInterface
var LikedOfferRepo liked_offer.LikedOfferRepositoryInterface
var ManufacturerRepo manufacturer.ManufacturerRepositoryInterface
var ModelRepo model.ModelRepositoryInterface
var NotificationRepo notification.NotificationRepositoryInterface
var PurchaseRepo purchase.PurchaseRepositoryInterface
var RefreshTokenRepo refresh_token.RefreshTokenRepositoryInterface
var ReviewRepo review.ReviewRepositoryInterface
var SaleOfferRepo sale_offer.SaleOfferRepositoryInterface
var UserRepo user.UserRepositoryInterface

func InitializeRepos() {
	AuctionRepo = auction.NewAuctionRepository(DB)
	BidRepo = bid.NewBidRepository(DB)
	ClientNotificationRepo = notification.NewClientNotificationRepository(DB)
	ImageRepo = image.NewImageRepository(DB)
	LikedOfferRepo = liked_offer.NewLikedOfferRepository(DB)
	ManufacturerRepo = manufacturer.NewManufacturerRepository(DB)
	ModelRepo = model.NewModelRepository(DB)
	NotificationRepo = notification.NewNotificationRepository(DB)
	PurchaseRepo = purchase.NewPurchaseRepository(DB)
	RefreshTokenRepo = refresh_token.NewRefreshTokenRepository(DB)
	ReviewRepo = review.NewReviewRepository(DB)
	SaleOfferRepo = sale_offer.NewSaleOfferRepository(DB)
	UserRepo = user.NewUserRepository(DB)
}
