package sale_offer

import (
	"github.com/susek555/BD2/car-dealer-api/internal/models"
	"github.com/susek555/BD2/car-dealer-api/internal/views"
)

type LikedOfferCheckerInterface interface {
	IsOfferLikedByUser(offerID, userID uint) bool
}

type BidRetrieverInterface interface {
	GetByAuctionId(auctionID uint) ([]models.Bid, error)
}

type OfferAccessEvaluatorInterface interface {
	CanBeModifiedByUser(*views.SaleOfferView, *uint) (bool, error)
	IsOfferLikedByUser(*views.SaleOfferView, *uint) bool
}

type OfferAccessEvaluator struct {
	bidRetriever BidRetrieverInterface
	likedChecker LikedOfferCheckerInterface
}

func NewAccessEvaluator(bidRetriever BidRetrieverInterface, likedChecker LikedOfferCheckerInterface) OfferAccessEvaluatorInterface {
	return &OfferAccessEvaluator{bidRetriever: bidRetriever, likedChecker: likedChecker}
}

func (e *OfferAccessEvaluator) CanBeModifiedByUser(offerView *views.SaleOfferView, userID *uint) (bool, error) {
	if userID == nil {
		return false, nil
	}
	if !offerView.BelongsToUser(*userID) {
		return false, nil
	}
	if !offerView.IsAuction {
		return true, nil
	}
	hasBids, err := e.hasBids(offerView)
	if err != nil {
		return false, err
	}
	return !hasBids, nil
}

func (e *OfferAccessEvaluator) IsOfferLikedByUser(offer *views.SaleOfferView, userID *uint) bool {
	if userID == nil {
		return false
	}
	return e.likedChecker.IsOfferLikedByUser(offer.ID, *userID)
}

func (e *OfferAccessEvaluator) hasBids(offer *views.SaleOfferView) (bool, error) {
	bids, err := e.bidRetriever.GetByAuctionId(offer.ID)
	if err != nil {
		return false, err
	}
	return len(bids) > 0, nil
}
