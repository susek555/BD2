package sale_offer

import (
	"github.com/susek555/BD2/car-dealer-api/internal/enums"
	"github.com/susek555/BD2/car-dealer-api/internal/models"
)

type LikedOfferCheckerInterface interface {
	IsOfferLikedByUser(offerID, userID uint) error
}

type BidRetrieverInterface interface {
	GetByAuctionID(auctionID uint) ([]models.Bid, error)
}

type OfferAccessEvaluatorInterface interface {
	CanBeModifiedByUser(SaleOfferEntityInterface, *uint) error
	IsOfferLikedByUser(SaleOfferEntityInterface, *uint) error
}

type OfferAccessEvaluator struct {
	bidRetriever BidRetrieverInterface
	likedChecker LikedOfferCheckerInterface
}

func NewAccessEvaluator(bidRetriever BidRetrieverInterface, likedChecker LikedOfferCheckerInterface) OfferAccessEvaluatorInterface {
	return &OfferAccessEvaluator{bidRetriever: bidRetriever, likedChecker: likedChecker}
}

func (e *OfferAccessEvaluator) CanBeModifiedByUser(offer SaleOfferEntityInterface, userID *uint) error {
	if userID == nil || !offer.BelongsToUser(*userID) {
		return ErrOfferNotOwned
	}
	if offer.GetStatus() == enums.SOLD || offer.GetStatus() == enums.EXPIRED {
		return ErrOfferAlreadySold
	}
	if !offer.IsAuctionOffer() {
		return nil
	}
	hasBids, err := e.hasBids(offer)
	if err != nil {
		return err
	}
	if hasBids {
		return ErrOfferHasBids
	}
	return nil
}

func (e *OfferAccessEvaluator) IsOfferLikedByUser(offer SaleOfferEntityInterface, userID *uint) error {
	if userID == nil {
		return ErrOfferNotOwned
	}
	return e.likedChecker.IsOfferLikedByUser(offer.GetID(), *userID)
}

func (e *OfferAccessEvaluator) hasBids(offer SaleOfferEntityInterface) (bool, error) {
	bids, err := e.bidRetriever.GetByAuctionID(offer.GetID())
	if err != nil {
		return false, err
	}
	return len(bids) > 0, nil
}
