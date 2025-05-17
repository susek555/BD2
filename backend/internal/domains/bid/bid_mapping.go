package bid

import "time"

func (b *Bid) MapToDTO() *RetrieveBidDTO {
	return &RetrieveBidDTO{
		AuctionID: b.AuctionID,
		BidderID:  b.BidderID,
		Amount:    b.Amount,
	}
}

func (cb *CreateBidDTO) MapToBid(userID uint) *Bid {
	return &Bid{
		AuctionID: cb.AuctionID,
		BidderID:  userID,
		Amount:    cb.Amount,
		CreatedAt: time.Now(),
	}
}
