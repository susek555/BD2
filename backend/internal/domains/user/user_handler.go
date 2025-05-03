package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var ErrorMap = map[error]int{
	ErrInvalidSelector: http.StatusBadRequest,
	ErrCreateCompany:   http.StatusBadRequest,
	ErrCreatePerson:    http.StatusBadRequest,
}

type Handler struct {
	service UserServiceInterface
}

func NewHandler(s UserServiceInterface) *Handler {
	return &Handler{service: s}
}

func (h *Handler) CreateUser(c *gin.Context) {
	var userDTO CreateUserDTO
	if err := c.ShouldBindJSON(&userDTO); err != nil {
		h.handleError(c, err)
		return
	}
	if err := h.service.Create(userDTO); err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, userDTO)
}

func (h *Handler) GetAllUsers(c *gin.Context) {
	userDTOs, err := h.service.GetAll()
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, userDTOs)
}

func (h *Handler) GetUserById(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	userDTO, err := h.service.GetById(uint(id))
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, userDTO)
}

func (h *Handler) GetUserByEmail(c *gin.Context) {
	email := c.Param("email")
	userDTO, err := h.service.GetByEmail(email)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, userDTO)
}

func (h *Handler) UpdateUser(c *gin.Context) {
	var userDTO UpdateUserDTO
	if err := c.ShouldBindJSON(&userDTO); err != nil {
		h.handleError(c, err)
		return
	}
	if err := h.service.Update(userDTO); err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, userDTO)
}

func (h *Handler) DeleteUser(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	err := h.service.Delete(uint(id))
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"deleted id": id})
}

func (h *Handler) getStatusCode(err error) int {
	if statusCode, ok := ErrorMap[err]; ok {
		return statusCode
	}
	return http.StatusInternalServerError
}

func (h *Handler) handleError(c *gin.Context, err error) {
	code := h.getStatusCode(err)
	c.JSON(code, err.Error())
}
