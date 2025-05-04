package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	errors "github.com/susek555/BD2/car-dealer-api/pkg/errors"
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

// CreateUser godoc
//
//	@Summary		Create user
//	@Description	Creates user from DTO and inserts its data to database. Whenever you want to create user you have to specifiy subtype (selector, P or C), fullfil only respective fields.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			body	body		CreateUserDTO		true	"Creation form"
//	@Success		201		{object}	CreateUserDTO		"User created"
//	@Failure		400		{object}	errors.HTTPError	"Invalid input dat(a, propably wrong selector (only "P" or "C" accepted)"
//	@Failure		500		{object}	errors.HTTPError	"Internal server error"
//	@Router			/users [post]
func (h *Handler) CreateUser(c *gin.Context) {
	var userDTO CreateUserDTO
	if err := c.ShouldBindJSON(&userDTO); err != nil {
		errors.HandleError(c, err, ErrorMap)
		return
	}
	if err := h.service.Create(userDTO); err != nil {
		errors.HandleError(c, err, ErrorMap)
		return
	}
	c.JSON(http.StatusCreated, userDTO)
}

func (h *Handler) GetAllUsers(c *gin.Context) {
	userDTOs, err := h.service.GetAll()
	if err != nil {
		errors.HandleError(c, err, ErrorMap)
		return
	}
	c.JSON(http.StatusOK, userDTOs)
}

func (h *Handler) GetUserById(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	userDTO, err := h.service.GetById(uint(id))
	if err != nil {
		errors.HandleError(c, err, ErrorMap)
		return
	}
	c.JSON(http.StatusOK, userDTO)
}

func (h *Handler) GetUserByEmail(c *gin.Context) {
	email := c.Param("email")
	userDTO, err := h.service.GetByEmail(email)
	if err != nil {
		errors.HandleError(c, err, ErrorMap)
		return
	}
	c.JSON(http.StatusOK, userDTO)
}

func (h *Handler) UpdateUser(c *gin.Context) {
	var userDTO UpdateUserDTO
	if err := c.ShouldBindJSON(&userDTO); err != nil {
		errors.HandleError(c, err, ErrorMap)
		return
	}
	if err := h.service.Update(userDTO); err != nil {
		errors.HandleError(c, err, ErrorMap)
		return
	}
	c.JSON(http.StatusOK, userDTO)
}

func (h *Handler) DeleteUser(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	err := h.service.Delete(uint(id))
	if err != nil {
		errors.HandleError(c, err, ErrorMap)
		return
	}
	c.JSON(http.StatusOK, gin.H{"deleted id": id})
}
