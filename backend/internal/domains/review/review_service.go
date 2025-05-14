package review

import (
	"errors"
)

type ReviewServiceInterface interface {
	Create(userId uint, review *CreateReviewDTO) (*RetrieveReviewDTO, error)
	GetAll() ([]RetrieveReviewDTO, error)
	GetById(id uint) (*RetrieveReviewDTO, error)
	Update(userId uint, review *UpdateReviewDTO) (*RetrieveReviewDTO, error)
	Delete(userId, id uint) error
	GetByReviewerId(reviewerId uint) ([]Review, error)
	GetByRevieweeId(reviewedId uint) ([]Review, error)
	GetByReviewerIdAndRevieweeId(reviewerId uint, revieweeId uint) (Review, error)
}

type ReviewService struct {
	Repo ReviewRepositoryInterface
}

func NewReviewService(repo ReviewRepositoryInterface) ReviewServiceInterface {
	return &ReviewService{
		Repo: repo,
	}
}

func (service *ReviewService) Create(userId uint, review *CreateReviewDTO) (*RetrieveReviewDTO, error) {
	reviewObj := review.MapToObject(userId)
	err := service.Repo.Create(&reviewObj)
	if err != nil {
		return nil, err
	}
	reviewDTO := reviewObj.MapToDTO()
	return &reviewDTO, nil
}

func (service *ReviewService) GetAll() ([]RetrieveReviewDTO, error) {
	reviews, err := service.Repo.GetAll()
	if err != nil {
		return nil, err
	}
	var reviewsDTO []RetrieveReviewDTO
	for _, review := range reviews {
		reviewsDTO = append(reviewsDTO, review.MapToDTO())
	}
	return reviewsDTO, nil
}

func (service *ReviewService) GetById(id uint) (*RetrieveReviewDTO, error) {
	review, err := service.Repo.GetById(id)
	if err != nil {
		return nil, err
	}
	reviewDTO := review.MapToDTO()
	return &reviewDTO, nil
}

func (service *ReviewService) Update(reviewerId uint, review *UpdateReviewDTO) (*RetrieveReviewDTO, error) {
	revieweeId, err := service.getRevieweeId(review.ID)
	if err != nil {
		return nil, err
	}
	reviewObj := review.MapToObject(reviewerId, revieweeId)
	err = service.Repo.Update(&reviewObj)
	if err != nil {
		return nil, err
	}
	reviewDTO := reviewObj.MapToDTO()
	return &reviewDTO, nil
}

func (service *ReviewService) Delete(userId, id uint) error {
	review, err := service.Repo.GetById(id)
	if err != nil {
		return err
	}
	if review.ReviewerID != userId {
		return errors.New("you are not the reviewer of this review")
	}
	return service.Repo.Delete(id)
}

func (service *ReviewService) GetByReviewerId(reviewerId uint) ([]Review, error) {
	return service.Repo.GetByReviewerId(reviewerId)
}

func (service *ReviewService) GetByRevieweeId(reviewedId uint) ([]Review, error) {
	return service.Repo.GetByRevieweeId(reviewedId)
}

func (service *ReviewService) GetByReviewerIdAndRevieweeId(reviewerId uint, reviewedId uint) (Review, error) {
	return service.Repo.GetByReviewerIdAndRevieweeId(reviewerId, reviewedId)
}

func (service *ReviewService) getRevieweeId(reviewId uint) (uint, error) {
	review, err := service.Repo.GetById(reviewId)
	if err != nil {
		return 0, err
	}
	return review.RevieweeId, nil
}
