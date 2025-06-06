package liked_offer

import "github.com/susek555/BD2/car-dealer-api/internal/models"

type SaleOfferRetrieverInterface interface {
	GetByID(id uint) (*models.SaleOffer, error)
}

type LikedOfferServiceInterface interface {
	LikeOffer(offerID, userID uint) error
	DislikeOffer(offerID, userID uint) error
}

type LikedOfferService struct {
	likedOfferRepo     LikedOfferRepositoryInterface
	saleOfferRetriever SaleOfferRetrieverInterface
}

func NewLikedOfferService(likedOfferRepo LikedOfferRepositoryInterface, saleOfferRetriever SaleOfferRetrieverInterface) LikedOfferServiceInterface {
	return &LikedOfferService{likedOfferRepo: likedOfferRepo, saleOfferRetriever: saleOfferRetriever}
}

func (s *LikedOfferService) LikeOffer(offerID, userID uint) error {
	offer, err := s.saleOfferRetriever.GetByID(offerID)
	if err != nil {
		return err
	}
	if s.likedOfferRepo.IsOfferLikedByUser(offerID, userID) {
		return ErrLikeAlreadyLikedOffer
	}
	if offer.BelongsToUser(userID) {
		return ErrLikeOwnOffer
	}
	return s.likedOfferRepo.Create(&models.LikedOffer{OfferID: offerID, UserID: userID})
}

func (s *LikedOfferService) DislikeOffer(offerID, userID uint) error {
	if _, err := s.saleOfferRetriever.GetByID(offerID); err != nil {
		return err
	}
	if !s.likedOfferRepo.IsOfferLikedByUser(offerID, userID) {
		return ErrDislikeNotLikedOffer
	}
	return s.likedOfferRepo.Delete(offerID, userID)
}
