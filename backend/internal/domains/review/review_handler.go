package review

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/auth"
	"github.com/susek555/BD2/car-dealer-api/pkg/custom_errors"
)

type Handler struct {
	service ReviewServiceInterface
}

func NewHandler(service ReviewServiceInterface) *Handler {
	return &Handler{service: service}
}

// GetAllReviews godoc
//
//	@ID				getAllReviews
//	@Summary		List all reviews
//	@Description	Returns every review in the system as an array of DTOs.
//	@Tags			reviews
//	@Produce		json
//	@Success		200	{array}		ReviewOutput			"OK – list of reviews"
//	@Failure		400	{object}	custom_errors.HTTPError	"Bad Request – query failed"
//	@Router			/review [get]
func (h *Handler) GetAllReviews(c *gin.Context) {
	reviews, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, reviews)
}

// GetReviewById godoc
//
//	@ID				getReviewById
//	@Summary		Get review by id
//	@Description	Returns review that match given id as an DTO.
//	@Tags			reviews
//	@Produce		json
//	@Param			id	path		int						true	"Review id"
//	@Success		200	{object}	ReviewOutput			"OK – review with given id"
//	@Failure		400	{object}	custom_errors.HTTPError	"Bad Request – query failed"
//	@Router			/review/{id} [get]
func (h *Handler) GetReviewById(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	review, err := h.service.GetById(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, review)
}

// CreateReview godoc
//
//	@ID				createReview
//	@Summary		Create a new review
//	@Description	Persists a new review entity and returns the created review.
//	@Tags			reviews
//	@Accept			json
//	@Produce		json
//	@Param			body	body		ReviewInput				true	"Review payload"
//	@Success		201		{object}	ReviewOutput			"Created – review stored"
//	@Failure		400		{object}	custom_errors.HTTPError	"Bad Request – validation or persistence error"
//	@Router			/review [post]
func (h *Handler) CreateReview(c *gin.Context) {
	var reviewInput CreateReviewDTO
	reviewerId, err := auth.GetUserId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	if err := c.ShouldBindJSON(&reviewInput); err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	reviewOutput, err := h.service.Create(uint(reviewerId), &reviewInput)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, reviewOutput)
}

// UpdateReview godoc
//
//	@ID				updateReview
//	@Summary		Update an existing review
//	@Description	Updates a review and returns the updated entity.
//	@Tags			reviews
//	@Accept			json
//	@Produce		json
//	@Param			body	body		ReviewInput				true	"Review payload"
//	@Success		200		{object}	ReviewOutput			"OK – review updated"
//	@Failure		400		{object}	custom_errors.HTTPError	"Bad Request – validation or update error"
//	@Router			/review [put]
func (h *Handler) UpdateReview(c *gin.Context) {
	var reviewInput UpdateReviewDTO
	reviewerId, err := auth.GetUserId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	if err := c.ShouldBindJSON(&reviewInput); err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	reviewOutput, err := h.service.Update(uint(reviewerId), &reviewInput)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, reviewOutput)
}

// DeleteReview godoc
//
//	@ID				deleteReview
//	@Summary		Delete a review
//	@Description	Deletes the review identified by its ID.
//	@Tags			reviews
//	@Param			id	path		int						true	"Review ID"
//	@Success		204	{string}	string					"No Content – review deleted"
//	@Failure		400	{object}	custom_errors.HTTPError	"Bad Request – invalid ID format or delete failed"
//	@Router			/review/{id} [delete]
func (h *Handler) DeleteReview(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	reviewerId, err := auth.GetUserId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
	}
	if err := h.service.Delete(reviewerId, uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
	}
	c.JSON(http.StatusNoContent, nil)
}

// GetReviewsByReviewerId godoc
//
//	@ID				getReviewsByReviewerId
//	@Summary		List reviews written by a reviewer
//	@Description	Returns all reviews authored by the reviewer specified by ID.
//	@Tags			reviews
//	@Produce		json
//	@Param			id	path		int						true	"Reviewer ID"
//	@Success		200	{array}		ReviewOutput			"OK – list of reviews"
//	@Failure		400	{object}	custom_errors.HTTPError	"Bad Request – invalid ID format or query failed"
//	@Router			/review/reviewer/{id} [get]
func (h *Handler) GetReviewsByReviewerId(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
	}
	reviews, err := h.service.GetByReviewerId(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
	}
	var reviewsDTO []RetrieveReviewDTO
	for _, review := range reviews {
		reviewDTO := review.MapToDTO()
		reviewsDTO = append(reviewsDTO, reviewDTO)
	}
	c.JSON(http.StatusOK, reviewsDTO)
}

// GetReviewsByRevieweeId godoc
//
//	@ID				getReviewsByRevieweeId
//	@Summary		List reviews about a reviewee
//	@Description	Returns all reviews where the given user is the reviewee.
//	@Tags			reviews
//	@Produce		json
//	@Param			id	path		int						true	"Reviewee ID"
//	@Success		200	{array}		ReviewOutput			"OK – list of reviews"
//	@Failure		400	{object}	custom_errors.HTTPError	"Bad Request – invalid ID format or query failed"
//	@Router			/review/reviewee/{id} [get]
func (h *Handler) GetReviewsByRevieweeId(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
	}
	reviews, err := h.service.GetByRevieweeId(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
	}
	var reviewsDTO []RetrieveReviewDTO
	for _, review := range reviews {
		reviewDTO := review.MapToDTO()
		reviewsDTO = append(reviewsDTO, reviewDTO)
	}
	c.JSON(http.StatusOK, reviewsDTO)
}

// GetReviewsByReviewerIdAndRevieweeId godoc
//
//	@ID				getReviewByReviewerAndReviewee
//	@Summary		Get review written by one user about another
//	@Description	Returns the review where <reviewerId> is the author and <revieweeId> is the subject.
//	@Tags			reviews
//	@Produce		json
//	@Param			reviewerId	path		int						true	"Reviewer ID"
//	@Param			revieweeId	path		int						true	"Reviewee ID"
//	@Success		200			{object}	ReviewOutput			"OK – review found"
//	@Failure		400			{object}	custom_errors.HTTPError	"Bad Request – invalid ID format or query failed"
//	@Router			/review/reviewer/{reviewerId}/reviewee/{revieweeId} [get]
func (h *Handler) GetReviewsByReviewerIdAndRevieweeId(c *gin.Context) {
	reviewerId, err := strconv.ParseUint(c.Param("reviewerId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
	}
	revieweeId, err := strconv.ParseUint(c.Param("revieweeId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
	}
	review, err := h.service.GetByReviewerIdAndRevieweeId(uint(reviewerId), uint(revieweeId))
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
	}
	reviewDTO := review.MapToDTO()
	c.JSON(http.StatusOK, reviewDTO)
}
