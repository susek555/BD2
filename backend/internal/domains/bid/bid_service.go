package bid

type BidServiceInterface interface {
	Create(bid *Bid) error
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
