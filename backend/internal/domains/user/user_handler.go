package user

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service UserServiceInterface
}

func NewHandler(s UserServiceInterface) *Handler {
	return &Handler{service: s}
}

func (h *Handler) CreateUser(c *gin.Context) {
	var userDTO CreateUserDTO
	if err := c.ShouldBindJSON(&userDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	if err := h.service.Create(userDTO); err != nil {
		switch {
		case errors.Is(err, ErrInvalidSelector):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid selector (should be P or C)"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		}
		return
	}
	c.JSON(http.StatusCreated, gin.H{"user": userDTO})
}

func (h *Handler) UpdateUser(c *gin.Context) {
	var userDTO UpdateUserDTO
	if err := c.ShouldBindJSON(&userDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
	}
	if err := h.service.Update(userDTO); err != nil {
		switch {
		case errors.Is(err, ErrInvalidSelector):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid selector (should be P or C)"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": userDTO})
}

func (h *Handler) GetAllUsers(c *gin.Context) {
	users, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *Handler) GetUserById(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	user, err := h.service.GetById(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user with given id not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *Handler) GetUserByEmail(c *gin.Context) {
	email := c.Param("email")
	user, err := h.service.GetByEmail(email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user with given email not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *Handler) DeleteUser(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	err := h.service.Delete(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user deleted sucessfully"})
}
