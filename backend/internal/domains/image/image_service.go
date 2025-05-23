package image

import (
	"fmt"
	"mime/multipart"

	"github.com/susek555/BD2/car-dealer-api/internal/domains/models"
)

type ImageServiceInterface interface {
	StoreImages(offerID uint, image []multipart.FileHeader) error
}

type ImageService struct {
	repo   ImageRepositoryInterface
	bucket ImageBucketInterface
}

func NewImageService(r ImageRepositoryInterface, b ImageBucketInterface) ImageServiceInterface {
	return &ImageService{repo: r, bucket: b}
}

func (s *ImageService) StoreImages(offerID uint, images []multipart.FileHeader) error {
	prefix := fmt.Sprintf("sale-offer-%d/", offerID)
	for _, image := range images {
		url, err := s.bucket.UploadImage(prefix, &image)
		if err != nil {
			s.cleanup(offerID, prefix)
			return err
		}
		if err = s.repo.Create(&models.Image{OfferID: offerID, Url: url}); err != nil {
			s.cleanup(offerID, prefix)
			return err
		}
	}
	return nil
}

func (s *ImageService) cleanup(offerID uint, prefix string) {
	_ = s.repo.DeleteByOfferID(offerID)
	_ = s.bucket.DeleteImagesByPrefix(prefix)
}
