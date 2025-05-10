package review

import (
	"github.com/gin-gonic/gin"
	"github.com/susek555/BD2/car-dealer-api/pkg/custom_errors"
	"net/http"
	"strconv"
)

type Handler struct {
	service ReviewServiceInterface
}

func NewHandler(service ReviewServiceInterface) *Handler {
	return &Handler{service: service}
}

// GetAllReviews godoc
//
// @ID           getAllReviews
// @Summary      List all reviews
// @Description  Returns every review in the system as an array of DTOs.
// @Tags         reviews
// @Produce      json
// @Success      200  {array}   ReviewDTO      "OK – list of reviews"
// @Failure      400  {object}  custom_errors.HTTPError  "Bad Request – query failed"
// @Router       /review [get]
func (h *Handler) GetAllReviews(c *gin.Context) {
	reviews, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	var reviewsDTO []ReviewDTO
	for _, review := range reviews {
		reviewsDTO = append(reviewsDTO, review.MapToDTO())
	}
	c.JSON(http.StatusOK, reviewsDTO)
}

// GetReviewById godoc
//
// @ID           getReviewById
// @Summary      Get review by id
// @Description  Returns review that match given id as an DTO.
// @Tags         reviews
// @Produce      json
// @Param 		 id  path	int	true "Review id"
// @Success      200  {object}   ReviewDTO      "OK – review with given id"
// @Failure      400  {object}  custom_errors.HTTPError  "Bad Request – query failed"
// @Router       /review/{id} [get]
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
	reviewDTO := review.MapToDTO()
	c.JSON(http.StatusOK, reviewDTO)
}

func (h *Handler) CreateReview(c *gin.Context) {
	var review Review
	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	if err := h.service.Create(&review); err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	reviewDTO := review.MapToDTO()
	c.JSON(http.StatusCreated, reviewDTO)
}

func (h *Handler) UpdateReview(c *gin.Context) {
	var review Review
	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
	}
	if err := h.service.Update(&review); err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
	}
	reviewDTO := review.MapToDTO()
	c.JSON(http.StatusOK, reviewDTO)
}

func (h *Handler) DeleteReview(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
	}
	if err := h.service.Delete(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
	}
	c.JSON(http.StatusNoContent, nil)
}

func (h *Handler) GetReviewsByReviewerId(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
	}
	reviews, err := h.service.GetByReviewerId(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
	}
	var reviewsDTO []ReviewDTO
	for _, review := range reviews {
		reviewDTO := review.MapToDTO()
		reviewsDTO = append(reviewsDTO, reviewDTO)
	}
	c.JSON(http.StatusOK, reviewsDTO)
}

func (h *Handler) GetReviewsByRevieweeId(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
	}
	reviews, err := h.service.GetByRevieweeId(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
	}
	var reviewsDTO []ReviewDTO
	for _, review := range reviews {
		reviewDTO := review.MapToDTO()
		reviewsDTO = append(reviewsDTO, reviewDTO)
	}
	c.JSON(http.StatusOK, reviewsDTO)
}

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
