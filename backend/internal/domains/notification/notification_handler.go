package notification

import (
	"net/http"

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

func (h *Handler) GetFilteredNotifications(c *gin.Context) {
	filter := NewNotificationFilter()
	receiverID, err := auth.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, custom_errors.NewHTTPError(err.Error()))
		return
	}
	filter.ReceiverID = &receiverID
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(400, gin.H{"error": "Invalid query parameters"})
		return
	}

	notifications, err := h.service.GetFilteredNotifications(filter)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, notifications)
}
