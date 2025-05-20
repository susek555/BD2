package liked_offer

type LikeOfferServiceInterface interface {
	Create(offer *LikedOffer) error
	Delete(offerID, userID uint) error
}

type LikedOfferService struct {
	repo LikedOfferRepositoryInterface
}

func NewLikedOfferService(likedOfferRepo LikedOfferRepositoryInterface) LikeOfferServiceInterface {
	return &LikedOfferService{repo: likedOfferRepo}
}

func (s *LikedOfferService) Create(offer *LikedOffer) error {
	return s.repo.Create(offer)
}

func (s *LikedOfferService) Delete(offerID, userID uint) error {
	return s.repo.Delete(offerID, userID)
}
