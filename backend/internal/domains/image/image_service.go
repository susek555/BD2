package image

import (
	"fmt"
	"mime/multipart"

	"github.com/susek555/BD2/car-dealer-api/internal/domains/sale_offer"
	"github.com/susek555/BD2/car-dealer-api/internal/enums"
	"github.com/susek555/BD2/car-dealer-api/internal/models"
)

type OfferAccessEvaluatorInterface interface {
	CanBeModifiedByUser(offer sale_offer.SaleOfferEntityInterface, userID *uint) error
}

type OfferRepositoryInterface interface {
	GetByID(offerID uint) (*models.SaleOffer, error)
	UpdateStatus(offer *models.SaleOffer, status enums.Status) error
}

type ImageServiceInterface interface {
	Store(offerID uint, image []*multipart.FileHeader, userID uint) error
	DeleteByURL(url string, userID uint) error
	DeleteByOfferID(offerID uint, userID uint) error
}

type ImageService struct {
	repo       ImageRepositoryInterface
	bucket     ImageBucketInterface
	offerRepo  OfferRepositoryInterface
	accessEval OfferAccessEvaluatorInterface
}

func NewImageService(
	r ImageRepositoryInterface,
	b ImageBucketInterface,
	offerRepo OfferRepositoryInterface,
	accessEval OfferAccessEvaluatorInterface,
) ImageServiceInterface {
	return &ImageService{repo: r, bucket: b, offerRepo: offerRepo, accessEval: accessEval}
}

func (s *ImageService) Store(offerID uint, images []*multipart.FileHeader, userID uint) error {
	offer, err := s.offerRepo.GetByID(offerID)
	if err != nil {
		return err
	}
	if err := s.accessEval.CanBeModifiedByUser(offer, &userID); err != nil {
		return err
	}
	storedImages, err := s.repo.GetByOfferID(offerID)
	if err != nil {
		return err
	}
	if s.wouldExceedImageLimit(storedImages, len(images), 10) {
		return ErrTooManyImages
	}
	if err := s.saveImagesToStorageAndDB(offerID, images); err != nil {
		return err
	}
	return s.setOfferStatus(offer)
}

func (s *ImageService) DeleteByURL(url string, userID uint) error {
	image, err := s.repo.GetByURL(url)
	if err != nil {
		return err
	}
	offer, err := s.offerRepo.GetByID(image.OfferID)
	if err != nil {
		return err
	}
	if err := s.accessEval.CanBeModifiedByUser(offer, &userID); err != nil {
		return err
	}
	if err := s.repo.Delete(image.ID); err != nil {
		return err
	}
	if err := s.bucket.Delete(image.PublicID); err != nil {
		if restoreErr := s.repo.Create(image); restoreErr != nil {
			return restoreErr
		}
		return err
	}
	return s.setOfferStatus(offer)
}

func (s *ImageService) DeleteByOfferID(offerID uint, userID uint) error {
	offer, err := s.offerRepo.GetByID(offerID)
	if err != nil {
		return err
	}
	if err := s.accessEval.CanBeModifiedByUser(offer, &userID); err != nil {
		return err
	}
	images, err := s.repo.GetByOfferID(offerID)
	if err != nil {
		return err
	}
	if !s.hasImages(images) {
		return ErrZeroImages
	}
	if err := s.repo.DeleteByOfferID(offerID); err != nil {
		return err
	}
	folder := fmt.Sprintf("sale-offer-%d", offerID)
	if err := s.bucket.DeleteByFolderName(folder); err != nil {
		if restoreErr := s.repo.BatchCreate(images); restoreErr != nil {
			return restoreErr
		}
		return err
	}
	return s.setOfferStatus(offer)
}

func (s *ImageService) wouldExceedImageLimit(images []models.Image, nImages int, maxImages int) bool {
	return len(images)+nImages > maxImages
}

func (s *ImageService) hasImages(images []models.Image) bool {
	return len(images) > 0
}

func (s *ImageService) saveImagesToStorageAndDB(offerID uint, images []*multipart.FileHeader) error {
	var (
		uploadedPublicIDs []string
		storedImages      []models.Image
	)
	folder := fmt.Sprintf("sale-offer-%d/", offerID)
	for _, image := range images {
		publicID, url, err := s.bucket.Upload(folder, image)
		if err != nil {
			s.partialCleanup(uploadedPublicIDs, storedImages)
			return nil
		}
		uploadedPublicIDs = append(uploadedPublicIDs, publicID)
		imageModel := &models.Image{OfferID: offerID, PublicID: publicID, Url: url}
		if err := s.repo.Create(imageModel); err != nil {
			s.partialCleanup(uploadedPublicIDs, storedImages)
			return nil
		}
		storedImages = append(storedImages, *imageModel)
	}
	return nil
}

func (s *ImageService) setOfferStatus(offer *models.SaleOffer) error {
	images, err := s.repo.GetByOfferID(offer.ID)
	if err != nil {
		return err
	}
	var status enums.Status
	switch {
	case len(images) < 3:
		status = enums.PENDING
	case len(images) >= 3:
		status = enums.READY
	}
	return s.offerRepo.UpdateStatus(offer, status)
}

func (s *ImageService) partialCleanup(publicIDs []string, images []models.Image) {
	for _, image := range images {
		_ = s.repo.Delete(image.ID)
	}
	for _, id := range publicIDs {
		_ = s.bucket.Delete(id)
	}
}
