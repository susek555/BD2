package bid

type BidServiceInterface interface {
	BidRepositoryInterface
}

type BidService struct {
	Repo BidRepositoryInterface
}

func NewBidService(repo BidRepositoryInterface) BidServiceInterface {
	return &BidService{
		Repo: repo,
	}
}

func (service *BidService) Create(bid *Bid) error {
	return service.Repo.Create(bid)
}

func (service *BidService) GetAll() ([]Bid, error) {
	return service.Repo.GetAll()
}

func (service *BidService) GetById(id uint) (*Bid, error) {
	return service.Repo.GetById(id)
}

func (service *BidService) GetByBidderId(bidderId uint) ([]Bid, error) {
	return service.Repo.GetByBidderId(bidderId)
}

func (service *BidService) GetByAuctionId(auctionID uint) ([]Bid, error) {
	return service.Repo.GetByAuctionId(auctionID)
}

func (service *BidService) GetHighestBid(auctionId uint) (*Bid, error) {
	return service.Repo.GetHighestBid(auctionId)
}

func (service *BidService) GetHighestBidByUserId(auctionId, userId uint) (*Bid, error) {
	return service.Repo.GetHighestBidByUserId(auctionId, userId)
}
