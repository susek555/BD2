package auth

import (
	"errors"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/dto"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler { return &Handler{service: service} }

func (h *Handler) Register(ctx *gin.Context) {
	var request user.CreateUserDTO
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	access, refresh, err := h.service.Register(ctx, request)
	if err != nil {
		switch {
		case errors.Is(err, ErrEmailTaken):
			ctx.JSON(http.StatusConflict, gin.H{"error": "email taken"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		}
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"access_token":  access,
		"refresh_token": refresh,
	})
}

func (h *Handler) Login(c *gin.Context) {
	var req dto.LoginInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	access, refresh, err := h.service.Login(c, req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  access,
		"refresh_token": refresh,
	})
}

func (h *Handler) Refresh(c *gin.Context) {
	var req dto.RefreshInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	access, refresh, err := h.service.Refresh(c, req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"access_token":  access,
		"refresh_token": refresh,
	})
}

func (h *Handler) Logout(c *gin.Context) {
	var req dto.LogoutInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	userId, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing userID"})
		return
	}

	if err := h.service.Logout(c, userId.(uint), req.RefreshToken, req.AllDevices); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
