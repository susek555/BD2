package review

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/auth"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/generic"
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
//	@Success		200	{array}		RetrieveReviewDTO		"OK – list of reviews"
//	@Failure		400	{object}	custom_errors.HTTPError	"Bad Request – query failed"
//	@Router			/review [get]
func (h *Handler) GetAllReviews(c *gin.Context) {
	reviews, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	generic.HandleListResponse(c, reviews)
}

// GetReviewByID godoc
//
//	@ID				getReviewByID
//	@Summary		Get review by ID
//	@Description	Returns review that match given ID as an DTO.
//	@Tags			reviews
//	@Produce		json
//	@Param			ID	path		int						true	"Review ID"
//	@Success		200	{object}	RetrieveReviewDTO		"OK – review with given ID"
//	@Failure		400	{object}	custom_errors.HTTPError	"Bad Request – query failed"
//	@Router			/review/{id} [get]
func (h *Handler) GetReviewByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	review, err := h.service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, custom_errors.NewHTTPError(err.Error()))
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
//	@Param			body	body		CreateReviewDTO			true	"Review payload"
//	@Success		201		{object}	RetrieveReviewDTO		"Created – review stored"
//	@Failure		400		{object}	custom_errors.HTTPError	"Bad Request – valIDation or persistence error"
//	@Router			/review [post]
func (h *Handler) CreateReview(c *gin.Context) {
	var reviewInput CreateReviewDTO
	reviewerID, err := auth.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	if err := c.ShouldBindJSON(&reviewInput); err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	reviewOutput, err := h.service.Create(uint(reviewerID), &reviewInput)
	if err != nil {
		custom_errors.HandleError(c, err, ErrorMap)
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
//	@Param			body	body		UpdateReviewDTO			true	"Review payload"
//	@Success		200		{object}	RetrieveReviewDTO		"OK – review updated"
//	@Failure		400		{object}	custom_errors.HTTPError	"Bad Request – valIDation or update error"
//	@Router			/review [put]
func (h *Handler) UpdateReview(c *gin.Context) {
	var reviewInput UpdateReviewDTO
	reviewerID, err := auth.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	if err := c.ShouldBindJSON(&reviewInput); err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	reviewOutput, err := h.service.Update(uint(reviewerID), &reviewInput)
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
//	@Description	Deletes the review IDentified by its ID.
//	@Tags			reviews
//	@Param			id	path		int						true	"Review ID"
//	@Success		204	{string}	string					"No Content – review deleted"
//	@Failure		400	{object}	custom_errors.HTTPError	"Bad Request – invalID ID format or delete failed"
//	@Router			/review/{id} [delete]
func (h *Handler) DeleteReview(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	reviewerID, err := auth.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	if err := h.service.Delete(reviewerID, uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	c.Status(http.StatusNoContent)
}

// GetReviewsByReviewerID godoc
//
//	@ID				getReviewsByReviewerID
//	@Summary		List reviews written by a reviewer
//	@Description	Returns all reviews authored by the reviewer specified by ID.
//	@Tags			reviews
//	@Produce		json
//	@Param			id	path		int						true	"Reviewer ID"
//	@Success		200	{array}		RetrieveReviewDTO		"OK – list of reviews"
//	@Failure		400	{object}	custom_errors.HTTPError	"Bad Request – invalID ID format or query failed"
//	@Router			/review/reviewer/{id} [post]
func (h *Handler) GetReviewsByReviewerID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	userID := uint(id)
	filter := NewReviewFilter()

	if err := c.ShouldBindJSON(filter); err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	filter.ReviewerID = &userID

	reviews, err := h.service.GetFiltered(filter)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, reviews)
}

// GetReviewsByRevieweeID godoc
//
//	@ID				getReviewsByRevieweeID
//	@Summary		List reviews about a reviewee
//	@Description	Returns all reviews where the given user is the reviewee.
//	@Tags			reviews
//	@Produce		json
//	@Param			id	path		int						true	"Reviewee ID"
//	@Success		200	{array}		RetrieveReviewDTO		"OK – list of reviews"
//	@Failure		400	{object}	custom_errors.HTTPError	"Bad Request – invalID ID format or query failed"
//	@Router			/review/reviewee/{id} [post]
func (h *Handler) GetReviewsByRevieweeID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	userID := uint(id)
	filter := NewReviewFilter()

	if err := c.ShouldBindJSON(filter); err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	filter.RevieweeID = &userID

	reviews, err := h.service.GetFiltered(filter)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, reviews)
}

// GetReviewsByReviewerIDAndRevieweeID godoc
//
//	@ID				getReviewByReviewerAndReviewee
//	@Summary		Get review written by one user about another
//	@Description	Returns the review where <reviewerID> is the author and <revieweeID> is the subject.
//	@Tags			reviews
//	@Produce		json
//	@Param			reviewerID	path		int						true	"Reviewer ID"
//	@Param			revieweeID	path		int						true	"Reviewee ID"
//	@Success		200			{object}	RetrieveReviewDTO		"OK – review found"
//	@Failure		400			{object}	custom_errors.HTTPError	"Bad Request – invalID ID format or query failed"
//	@Router			/review/reviewer/{reviewerID}/reviewee/{revieweeID} [get]
func (h *Handler) GetReviewsByReviewerIDAndRevieweeID(c *gin.Context) {
	reviewerID, err := strconv.ParseUint(c.Param("reviewerID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	revieweeID, err := strconv.ParseUint(c.Param("revieweeID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	review, err := h.service.GetByReviewerIDAndRevieweeID(uint(reviewerID), uint(revieweeID))
	if err != nil {
		c.JSON(http.StatusOK, custom_errors.NewHTTPError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, review)
}

// GetFilteredReviews godoc
//
//	@ID				getFilteredReviews
//	@Summary		Filter reviews with pagination
//	@Description	Returns reviews matching podanych kryteriów filtrowania wraz z paginacją.
//	@Tags			reviews
//	@Accept			json
//	@Produce		json
//	@Param			body	body		review.ReviewFilter						true	"Filter payload"
//	@Success		200		{object}	review.RetrieveReviewsWithPagination	"OK "
//	@Failure		400		{object}	custom_errors.HTTPError					"Bad Request – invalID filter or query failed"
//	@Router			/review/filter [post]
func (h *Handler) GetFilteredReviews(c *gin.Context) {
	filter := NewReviewFilter()
	if err := c.ShouldBindJSON(filter); err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	reviews, err := h.service.GetFiltered(filter)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, reviews)
}

// GetAverageRatingByRevieweeID godoc
//
//	@ID				getAverageRatingByRevieweeID
//	@Summary		Get average rating for a reviewee
//	@Description	Returns the average rating value calculated over all reviews for the given reviewee.
//	@Tags			reviews
//	@Produce		json
//	@Param			id	path		int						true	"Reviewee ID"
//	@Success		200	{number}	float64					"OK – average rating (rounded to two decimals)"
//	@Failure		400	{object}	custom_errors.HTTPError	"Bad Request – invalID ID format or query failed"
//	@Router			/review/average-rating/{id} [get]
func (h *Handler) GetAverageRatingByRevieweeID(c *gin.Context) {
	ID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	revieweeID := uint(ID)
	averageRating, err := h.service.GetAverageRatingByRevieweeID(revieweeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, averageRating)
}

// GetFrequencyOfRatingByRevieweeID godoc
//
//	@ID				getFrequencyOfRatingByRevieweeID
//	@Summary		Get distribution of ratings for a reviewee
//	@Description	Returns a map from rating value (1–5) to percentage frequency among all reviews for the given reviewee.
//	@Tags			reviews
//	@Produce		json
//	@Param			id	path		int						true	"Reviewee ID"
//	@Success		200	{object}	map[int]int				"OK – percentage frequencies for ratings 1 through 5"
//	@Failure		400	{object}	custom_errors.HTTPError	"Bad Request – invalID ID format or query failed"
//	@Router			/review/frequency/{id} [get]
func (h *Handler) GetFrequencyOfRatingByRevieweeID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	revieweeID := uint(id)
	frequency, err := h.service.GetFrequencyOfRatingByRevieweeID(revieweeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, frequency)
}
