package auction

import "errors"

type AuctionServiceInterface interface {
	Create(auction *CreateAuctionDTO) (*RetrieveAuctionDTO, error)
	GetAll() ([]RetrieveAuctionDTO, error)
	GetById(id uint) (*RetrieveAuctionDTO, error)
	Update(auction *UpdateAuctionDTO) (*RetrieveAuctionDTO, error)
	Delete(id, userId uint) error
}

type AuctionService struct {
	repo AuctionRepositoryInterface
}

func NewAuctionService(repo AuctionRepositoryInterface) AuctionServiceInterface {
	return &AuctionService{
		repo: repo,
	}
}

func (s *AuctionService) Create(auction *CreateAuctionDTO) (*RetrieveAuctionDTO, error) {
	auctionEntity, err := auction.MapToAuction()
	if err != nil {
		return nil, err
	}
	err = s.repo.Create(auctionEntity)
	if err != nil {
		return nil, err
	}
	dto := MapToDTO(auctionEntity)
	return dto, nil
}

func (s *AuctionService) GetAll() ([]RetrieveAuctionDTO, error) {
	auctions, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	var auctionsDTO []RetrieveAuctionDTO
	for _, auction := range auctions {
		dto := MapToDTO(&auction)
		auctionsDTO = append(auctionsDTO, *dto)
	}
	return auctionsDTO, nil
}

func (s *AuctionService) GetById(id uint) (*RetrieveAuctionDTO, error) {
	auction, err := s.repo.GetById(id)
	if err != nil {
		return nil, err
	}
	dto := MapToDTO(auction)
	return dto, nil
}

func (s *AuctionService) Update(auction *UpdateAuctionDTO) (*RetrieveAuctionDTO, error) {
	auctionEntity, err := auction.MapToAuction()
	if err != nil {
		return nil, err
	}
	err = s.repo.Update(auctionEntity)
	if err != nil {
		return nil, err
	}
	dto := MapToDTO(auctionEntity)
	return dto, nil
}

func (s *AuctionService) Delete(id, userId uint) error {
	auction, err := s.repo.GetById(id)
	if err != nil {
		return err
	}
	if auction.Offer.UserID != userId {
		return errors.New("you are not the owner of this auction")
	}
	err = s.repo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
