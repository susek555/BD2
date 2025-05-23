package auth

import (
	"github.com/susek555/BD2/car-dealer-api/internal/domains/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/user"
	"github.com/susek555/BD2/car-dealer-api/pkg/custom_errors"
)

type Handler struct {
	Service AuthServiceInterface
}

var validate = validator.New()

func NewHandler(service AuthServiceInterface) *Handler { return &Handler{Service: service} }

// Register godoc
//
//	@Summary		Register new user
//	@Description	Set up account and return the status
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		user.CreateUserDTO	true	"Registration form"
//	@Success		201		{object}	RegisterResponse	"Created - returns tokens"
//	@Failure		400		{object}	RegisterResponse	"Invalid input data"
//	@Failure		409		{object}	RegisterResponse	"Login taken"
//	@Failure		500		{object}	RegisterResponse	"Internal server error"
//	@Router			/auth/register [post]
func (h *Handler) Register(ctx *gin.Context) {
	var request user.CreateUserDTO
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, RegisterResponse{
			Errors: map[string][]string{"other": {ErrInvalidBody.Error()}},
		})
		return
	}

	if err := validate.Struct(request); err != nil {
		errs := make(map[string][]string)
		for _, fe := range err.(validator.ValidationErrors) {
			errs[fe.Field()] = []string{fe.Error()}
		}
		ctx.JSON(http.StatusBadRequest, RegisterResponse{Errors: errs})
		return
	}

	err := h.Service.Register(ctx, request)
	if len(err) > 0 {
		ctx.JSON(http.StatusConflict, RegisterResponse{Errors: err})
		return
	}
	ctx.Status(http.StatusCreated)
}

// Login godoc
//
//	@Summary		Login user
//	@Description	Authenticate user and return a pair of tokens and user data
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		LoginInput		true	"Login form"
//	@Success		200		{object}	LoginResponse	"OK - returns tokens and data of the user"
//	@Failure		400		{object}	LoginResponse	"Invalid input data"
//	@Failure		401		{object}	LoginResponse	"Unauthorized"
//	@Failure		500		{object}	LoginResponse	"Internal server error"
//	@Router			/auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req LoginInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, LoginResponse{Errors: map[string][]string{"server": {ErrInvalidBody.Error()}}})
		return
	}

	access, refresh, user_, err := h.Service.Login(c, req)
	loginResponse := prepareLoginResponse(access, refresh, *user_)
	if err != nil {
		c.JSON(http.StatusUnauthorized, LoginResponse{Errors: map[string][]string{"credentials": {ErrInvalidCredentials.Error()}}})
		return
	}
	c.JSON(http.StatusOK, loginResponse)
}

func prepareLoginResponse(access, refresh string, user models.User) *LoginResponse {
	loginResponse := LoginResponse{
		RefreshToken: refresh,
		AccessToken:  access,
		Selector:     user.Selector,
		Username:     user.Username,
		UserID:       user.ID,
		Email:        user.Email,
	}
	if user.Selector == "C" {
		loginResponse.CompanyName = user.Company.Name
		loginResponse.CompanyNip = user.Company.NIP
	} else if user.Selector == "P" {
		loginResponse.PersonName = user.Person.Name
		loginResponse.PersonSurname = user.Person.Surname
	}
	return &loginResponse
}

// Refresh godoc
//
//	@Summary		Refresh tokens
//	@Description	Refresh access tokens
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		RefreshInput			true	"Refresh token form"
//	@Success		200		{object}	LoginResponse			"OK - returns new access token"
//	@Failure		400		{object}	custom_errors.HTTPError	"Invalid input data"
//	@Failure		401		{object}	custom_errors.HTTPError	"Unauthorized"
//	@Failure		500		{object}	custom_errors.HTTPError	"Internal server error"
//	@Router			/auth/refresh [post]
func (h *Handler) Refresh(c *gin.Context) {
	var req RefreshInput
	if err := c.ShouldBindJSON(&req); err != nil {
		custom_errors.HandleError(c, ErrInvalidBody, ErrorMap)
		return
	}
	access, err := h.Service.Refresh(c, req.RefreshToken)
	if err != nil {
		custom_errors.HandleError(c, err, ErrorMap)
		return
	}
	response := LoginResponse{AccessToken: access}
	c.JSON(http.StatusOK, response)
}

// Logout godoc
//
//	@Summary		Logout user
//	@Description	Logout user and invalidate refresh token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		LogoutInput				true	"Logout form"
//	@Success		204		{object}	custom_errors.HTTPError	"No content"
//	@Failure		400		{object}	custom_errors.HTTPError	"Invalid input data"
//	@Failure		401		{object}	custom_errors.HTTPError	"Unauthorized"
//	@Failure		500		{object}	custom_errors.HTTPError	"Internal server error"
//	@Router			/logout [post]
func (h *Handler) Logout(c *gin.Context) {
	var req LogoutInput
	if err := c.ShouldBindJSON(&req); err != nil {
		custom_errors.HandleError(c, ErrInvalidBody, ErrorMap)
		return
	}

	userId, ok := c.Get("userID")
	if !ok {
		custom_errors.HandleError(c, ErrUnauthorized, ErrorMap)
		return
	}

	if err := h.Service.Logout(c, userId.(uint), req.RefreshToken, req.AllDevices); err != nil {
		custom_errors.HandleError(c, err, ErrorMap)
		return
	}
	c.Status(http.StatusNoContent)
}
