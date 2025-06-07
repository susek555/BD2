package notification

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/auth"
	"github.com/susek555/BD2/car-dealer-api/pkg/custom_errors"
)

type Handler struct {
	service NotificationServiceInterface
}

func NewHandler(service NotificationServiceInterface) *Handler {
	return &Handler{
		service: service,
	}
}

// GetFilteredNotifications godoc
// @Summary Retrieve filtered notifications for authenticated user
// @Description Gets a list of notifications based on provided filter criteria for the current user
// @Tags notification
// @Accept json
// @Produce json
// @Param body body NotificationFilter true "Filter criteria for notifications"
// @Security BearerAuth
// @Success 200 {array} RetrieveNotificationsWithPagination "List of notifications"
// @Failure 400 {object} custom_errors.HTTPError "Invalid body or bad request"
// @Failure 401 {object} custom_errors.HTTPError "Unauthorized"
// @Router /notification/filter [post]
func (h *Handler) GetFilteredNotifications(c *gin.Context) {
	filter := NewNotificationFilter()
	receiverID, err := auth.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, custom_errors.NewHTTPError(err.Error()))
		return
	}
	if err := c.ShouldBindJSON(filter); err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	filter.ReceiverID = &receiverID

	notifications, err := h.service.GetFilteredNotifications(filter)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, notifications)
}

// MarkAsSeen godoc
// @Summary Mark notification as seen
// @Description Updates the seen status of a specific notification to true for the authenticated user
// @Tags notification
// @Accept json
// @Produce json
// @Param id path int true "Notification ID"
// @Success 200 "Status OK"
// @Failure 400 {object} custom_errors.HTTPError "Bad Request - Invalid notification ID or update failed"
// @Failure 401 {object} custom_errors.HTTPError "Unauthorized - User authentication required"
// @Router /notification/seen/{id} [put]
// @Security BearerAuth
func (h *Handler) MarkAsSeen(c *gin.Context) {
	notificationID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}

	userID, err := auth.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, custom_errors.NewHTTPError(err.Error()))
		return
	}

	err = h.service.UpdateSeenStatus(uint(notificationID), userID, true)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}

	c.Status(http.StatusOK)
}

// MarkAsUnseen godoc
// @Summary Mark notification as unseen
// @Description Updates the seen status of a notification to unseen for the authenticated user
// @Tags notification
// @Accept json
// @Produce json
// @Param id path int true "Notification ID"
// @Success 200 "OK"
// @Failure 400 {object} custom_errors.HTTPError "Bad request"
// @Failure 401 {object} custom_errors.HTTPError "Unauthorized"
// @Router /notification/unseen/{id} [put]
// @Security BearerAuth
func (h *Handler) MarkAsUnseen(c *gin.Context) {
	notificationID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}

	userID, err := auth.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, custom_errors.NewHTTPError(err.Error()))
		return
	}

	err = h.service.UpdateSeenStatus(uint(notificationID), userID, false)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}

	c.Status(http.StatusOK)
}

// MarkAllAsSeen godoc
// @Summary Mark all notifications as seen
// @Description Updates the seen status of all notifications for the authenticated user to true
// @Tags notification
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 "Successfully marked all notifications as seen"
// @Failure 400 {object} custom_errors.HTTPError "Bad request"
// @Failure 401 {object} custom_errors.HTTPError "Unauthorized"
// @Router /notification/seen [put]
func (h *Handler) MarkAllAsSeen(c *gin.Context) {
	userID, err := auth.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, custom_errors.NewHTTPError(err.Error()))
		return
	}

	err = h.service.UpdateSeenStatusForAll(userID, true)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}

	c.Status(http.StatusOK)
}

// MarkAllAsUnseen godoc
// @Summary Mark all notifications as unseen
// @Description Updates the seen status of all notifications for the authenticated user to false
// @Tags notification
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 "Successfully marked all notifications as unseen"
// @Failure 400 {object} custom_errors.HTTPError "Bad request"
// @Failure 401 {object} custom_errors.HTTPError "Unauthorized"
// @Router /notification/unseen [put]
func (h *Handler) MarkAllAsUnseen(c *gin.Context) {
	userID, err := auth.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, custom_errors.NewHTTPError(err.Error()))
		return
	}

	err = h.service.UpdateSeenStatusForAll(userID, false)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}

	c.Status(http.StatusOK)
}
