package auth

import (
	"github.com/susek555/BD2/car-dealer-api/internal/domains/dto"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/service/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc auth.Service
}

func NewHandler(svc auth.Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Register(c *gin.Context) {
	var req dto.RegisterInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	token, err := h.svc.Register(c, req)
	if err != nil {
		switch err {
		case auth.ErrEmailTaken:
			c.JSON(http.StatusConflict, gin.H{"error": "email taken"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"token": token})
}

func (h *Handler) Login(c *gin.Context) {
	var req dto.LoginInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	token, err := h.svc.Login(c, req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
