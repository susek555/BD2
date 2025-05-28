package image

import (
	"fmt"
	"mime/multipart"

	"github.com/susek555/BD2/car-dealer-api/internal/domains/models"
)

type ImageServiceInterface interface {
	StoreImages(offerID uint, image []*multipart.FileHeader) error
}

type ImageService struct {
	repo   ImageRepositoryInterface
	bucket ImageBucketInterface
}

func NewImageService(r ImageRepositoryInterface, b ImageBucketInterface) ImageServiceInterface {
	return &ImageService{repo: r, bucket: b}
}

func (s *ImageService) StoreImages(offerID uint, images []*multipart.FileHeader) error {
	if err := s.validateImageLimit(offerID, len(images), 10); err != nil {
		return nil
	}
	return s.saveImagesToStorageAndDB(offerID, images)
}

func (s *ImageService) validateImageLimit(offerID uint, nImages int, maxImages int) error {
	storedImages, err := s.repo.GetImagesByOfferID(offerID)
	if err != nil {
		return err
	}
	if nImages+len(storedImages) > maxImages {
		return ErrTooManyImages
	}
	return nil
}

func (s *ImageService) saveImagesToStorageAndDB(offerID uint, images []*multipart.FileHeader) error {
	var (
		uploadedPublicIDs []string
		storedImages      []models.Image
	)
	folder := fmt.Sprintf("sale-offer=%d/", offerID)
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

func (s *ImageService) partialCleanup(publicIDs []string, images []models.Image) {
	for _, image := range images {
		_ = s.repo.Delete(image.ID)
	}
	for _, id := range publicIDs {
		_ = s.bucket.Delete(id)
	}
}

func (s *ImageService) totalCleanup(offerID uint) error {
	folder := fmt.Sprintf("sale-offer-%d/", offerID)
	if err := s.repo.DeleteByOfferID(offerID); err != nil {
		return err
	}
	if err := s.bucket.DeleteByFolderName(folder); err != nil {
		return err
	}
	return nil
}
