package bid

import (
	"time"

	"github.com/susek555/BD2/car-dealer-api/internal/domains/models"
)

func MapToProcessingDTO(b *models.Bid) *ProcessingBidDTO {
	return &ProcessingBidDTO{
		AuctionID: b.AuctionID,
		BidderID:  b.BidderID,
		Amount:    b.Amount,
		Auction:   b.Auction,
	}
}

func MapToDTO(b *models.Bid) *RetrieveBidDTO {
	return &RetrieveBidDTO{
		AuctionID: b.AuctionID,
		BidderID:  b.BidderID,
		Amount:    b.Amount,
	}
}

func ProcessingToRetrieve(b *ProcessingBidDTO) *RetrieveBidDTO {
	return &RetrieveBidDTO{
		AuctionID: b.AuctionID,
		BidderID:  b.BidderID,
		Amount:    b.Amount,
	}
}

func (cb *CreateBidDTO) MapToBid(userID uint) *models.Bid {
	return &models.Bid{
		AuctionID: cb.AuctionID,
		BidderID:  userID,
		Amount:    cb.Amount,
		CreatedAt: time.Now(),
	}
}
