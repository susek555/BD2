package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/susek555/BD2/car-dealer-api/pkg/custom_errors"
)

type Handler struct {
	service UserServiceInterface
}

func NewHandler(s UserServiceInterface) *Handler {
	return &Handler{service: s}
}

// GetAllUsers godoc
//
//	@Summary		Get all users
//	@Description	Return a list of all users stored in database. If user's subtype is person the company related fields will be omitted and vice versa.
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

// GetUserByID godoc
//
//	@Summary		Get user by id
//	@Description	Returns user who has provided id. If user's subtype is person the company related fields will be omitted and vice versa.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int						true	"User ID"
//	@Success		200	{object}	RetrieveUserDTO			"User"
//	@Failure		400	{object}	custom_errors.HTTPError	"ID is not a number"
//	@Failure		404	{object}	custom_errors.HTTPError	"User not found"
//	@Failure		500	{object}	custom_errors.HTTPError	"Internal server error"
//	@Router			/users/id/{id} [get]
func (h *Handler) GetUserByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	userDTO, err := h.service.GetByID(uint(id))
	if err != nil {
		custom_errors.HandleError(c, err, ErrorMap)
		return
	}
	c.JSON(http.StatusOK, userDTO)
}

// GetUserByEmail godoc
//
//	@Summary		Get user by email
//	@Description	Returns user who has provided email. If user's subtype is person the company related fields will be omitted and vice versa.
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
//	@Description	Updates user's data in database from given form. Currently, you can only change basic fields (email, username, password), not the subtype.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			body	body		UpdateUserDTO	true	"Update form"
//	@Success		200		{object}	UpdateUserDTO	"User updated"
//	@Failure		400		{object}	UpdateResponse	"Invalid input data - email, username or nip taken"
//	@Failure		401		{object}	UpdateResponse	"Unauthorized - user be logged in to update his data"
//	@Failure		403		{object}	UpdateResponse	"Forbidden - user can only update his own data"
//	@Failure		404		{object}	UpdateResponse	"User not found"
//	@Failure		500		{object}	UpdateResponse	"Internal server error"
//	@Router			/users [put]
//	@Security		Bearer
func (h *Handler) UpdateUser(c *gin.Context) {
	var userDTO UpdateUserDTO
	updateResponse := UpdateResponse{
		Errors: make(map[string][]string),
	}
	if err := c.ShouldBindJSON(&userDTO); err != nil {
		updateResponse.Errors["other"] = []string{err.Error()}
		c.JSON(http.StatusBadRequest, updateResponse)
		return
	}
	userID, ok := c.Get("userID")
	if !ok || userID != userDTO.ID {
		updateResponse.Errors["id"] = []string{ErrForbidden.Error()}
		c.JSON(http.StatusForbidden, updateResponse)
		return
	}
	err := h.service.Update(&userDTO)
	if len(err) > 0 {
		updateResponse.Errors = err
		c.JSON(http.StatusBadRequest, updateResponse)
		return
	}
	c.Status(http.StatusOK)
}

// DeleteUser godoc
//
//	@Summary		Delete user
//	@Description	Removes user who has provided id from database.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"User ID"
//	@Success		204	"User successfully deleted"
//	@Failure		400	{object}	custom_errors.HTTPError	"ID is not a number"
//	@Failure		401	{object}	custom_errors.HTTPError	"Unauthorized - user must be logged in to delete his account"
//	@Failure		403	{object}	custom_errors.HTTPError	"Forbidden - user can only delete his own account"
//	@Failure		404	{object}	custom_errors.HTTPError	"User not found"
//	@Failure		500	{object}	custom_errors.HTTPError	"Internal server error"
//	@Router			/users/{id} [delete]
//
//	@Security		Bearer
func (h *Handler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		custom_errors.HandleError(c, err, ErrorMap)
		return
	}
	userID, ok := c.Get("userID")
	if !ok || userID != uint(id) {
		custom_errors.HandleError(c, ErrInvalidUserID, ErrorMap)
		return
	}
	err = h.service.Delete(uint(id))
	if err != nil {
		custom_errors.HandleError(c, err, ErrorMap)
		return
	}
	c.Status(http.StatusNoContent)
}
