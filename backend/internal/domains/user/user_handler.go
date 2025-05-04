package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/susek555/BD2/car-dealer-api/pkg/custom_errors"
	"gorm.io/gorm"
)

var ErrorMap = map[error]int{
	ErrInvalidSelector:     http.StatusBadRequest,
	ErrCreateCompany:       http.StatusBadRequest,
	ErrCreatePerson:        http.StatusBadRequest,
	ErrHashPassword:        http.StatusInternalServerError,
	strconv.ErrSyntax:      http.StatusBadRequest, // for parsing uint
	gorm.ErrRecordNotFound: http.StatusNotFound,
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
//	@Description	Creates user from DTO and inserts its data to database. Whenever you want to create user you have to specifiy subtype (selector, P or C), fullfil only respective fields, they are required.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			body	body		CreateUserDTO			true	"Creation form"
//	@Success		201		{object}	CreateUserDTO			"User created"
//	@Failure		400		{object}	custom_errors.HTTPError	"Invalid input dat(a, propably wrong selector (only "P" or "C" accepted)"
//	@Failure		500		{object}	custom_errors.HTTPError	"Internal server error"
//	@Router			/users [post]
func (h *Handler) CreateUser(c *gin.Context) {
	var userDTO CreateUserDTO
	if err := c.ShouldBindJSON(&userDTO); err != nil {
		custom_errors.HandleError(c, err, ErrorMap)
		return
	}
	if err := h.service.Create(userDTO); err != nil {
		custom_errors.HandleError(c, err, ErrorMap)
		return
	}
	c.JSON(http.StatusCreated, userDTO)
}

// GetAllUsers godoc
//
//	@Summary		Get all users
//	@Description	Get all users from database and return them as a list of DTOs. If user's subtype is person the company related fields will be ommitted and vice versa.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		RetrieveUserDTO			"List of users"
//	@Failure		500	{object}	custom_errors.HTTPError	"Internal server error"
//	@Router			/users [get]
func (h *Handler) GetAllUsers(c *gin.Context) {
	userDTOs, err := h.service.GetAll()
	if err != nil {
		custom_errors.HandleError(c, err, ErrorMap)
		return
	}
	c.JSON(http.StatusOK, userDTOs)
}

// GetUserById godoc
//
//	@Summary		Get user by id
//	@Description	Get user by id from database and return it as a DTO. If user's subtype is person the company related fields will be ommitted and vice versa.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int						true	"User ID"
//	@Success		200	{object}	RetrieveUserDTO			"User"
//	@Failure		400	{object}	custom_errors.HTTPError	"Id is not a number"
//	@Failure		404	{object}	custom_errors.HTTPError	"User not found"
//	@Failure		500	{object}	custom_errors.HTTPError	"Internal server error"
//	@Router			/users/id/{id} [get]
func (h *Handler) GetUserById(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		custom_errors.HandleError(c, err, ErrorMap)
		return
	}
	userDTO, err := h.service.GetById(uint(id))
	if err != nil {
		custom_errors.HandleError(c, err, ErrorMap)
		return
	}
	c.JSON(http.StatusOK, userDTO)
}

// GetUserByEmail godoc
//
//	@Summary		Get user by email
//	@Description	Get user by email from database and return it as a DTO. If user's subtype is person the company related fields will be ommitted and vice versa.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			email	path		string					true	"User email"
//	@Success		200		{object}	RetrieveUserDTO			"User found"
//	@Failure		404		{object}	custom_errors.HTTPError	"User not found"
//	@Failure		500		{object}	custom_errors.HTTPError	"Internal server error"
//	@Router			/users/email/{email} [get]
func (h *Handler) GetUserByEmail(c *gin.Context) {
	email := c.Param("email")
	userDTO, err := h.service.GetByEmail(email)
	if err != nil {
		custom_errors.HandleError(c, err, ErrorMap)
		return
	}
	c.JSON(http.StatusOK, userDTO)
}

// UpdateUser godoc
//
//	@Summary		Update user
//	@Description	Update user in database from DTO. Currently you can only change basic fields (email, username, password), not the subtype.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			body	body		UpdateUserDTO			true	"Update form"
//	@Success		200		{object}	UpdateUserDTO			"User updated"
//	@Failure		400		{object}	custom_errors.HTTPError	"Invalid input data"
//	@Failure		404		{object}	custom_errors.HTTPError	"User not found"
//	@Failure		500		{object}	custom_errors.HTTPError	"Internal server error"
//	@Router			/users [put]
func (h *Handler) UpdateUser(c *gin.Context) {
	var userDTO UpdateUserDTO
	if err := c.ShouldBindJSON(&userDTO); err != nil {
		custom_errors.HandleError(c, err, ErrorMap)
		return
	}
	if err := h.service.Update(userDTO); err != nil {
		custom_errors.HandleError(c, err, ErrorMap)
		return
	}
	c.JSON(http.StatusOK, userDTO)
}

// DeleteUser godoc
//
//	@Summary		Delete user
//	@Description	Delete user from database by id.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"User ID"
//	@Success		204	"User sucessfully deleted"
//	@Failure		400	{object}	custom_errors.HTTPError	"Id is not a number"
//	@Failure		404	{object}	custom_errors.HTTPError	"User not found"
//	@Failure		500	{object}	custom_errors.HTTPError	"Internal server error"
//	@Router			/users/{id} [delete]
func (h *Handler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		custom_errors.HandleError(c, err, ErrorMap)
		return
	}
	err = h.service.Delete(uint(id))
	if err != nil {
		custom_errors.HandleError(c, err, ErrorMap)
		return
	}
	c.Status(http.StatusNoContent)
}
