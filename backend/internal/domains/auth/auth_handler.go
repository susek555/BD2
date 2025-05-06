package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/user"
	"github.com/susek555/BD2/car-dealer-api/pkg/custom_errors"
	"net/http"
)

type Handler struct {
	service Service
}

var validate = validator.New()

func NewHandler(service Service) *Handler { return &Handler{service: service} }

// Register godoc
//
//	@Summary		Register new user
//	@Description	Set up account and return a pair of tokens
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		user.CreateUserDTO		true	"Registration form"
//	@Success		201		{object}	TokenResponse			"Created - returns tokens"
//	@Failure		400		{object}	custom_errors.HTTPError	"Invalid input data"
//	@Failure		409		{object}	custom_errors.HTTPError	"Login taken"
//	@Failure		500		{object}	custom_errors.HTTPError	"Internal server error"
//	@Router			/auth/register [post]
func (h *Handler) Register(ctx *gin.Context) {
	var request user.CreateUserDTO
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, RegisterResponse{
			Errors: map[string][]string{"other": {"invalid body"}},
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

	err := h.service.Register(ctx, request)
	if len(err) > 0 {
		ctx.JSON(http.StatusConflict, RegisterResponse{Errors: err})
		return
	}
	ctx.JSON(http.StatusCreated, RegisterResponse{})
}

// Login godoc
//
//	@Summary		Login user
//	@Description	Authenticate user and return a pair of tokens
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		LoginInput				true	"Login form"
//	@Success		200		{object}	TokenResponse			"OK - returns tokens"
//	@Failure		400		{object}	custom_errors.HTTPError	"Invalid input data"
//	@Failure		401		{object}	custom_errors.HTTPError	"Unauthorized"
//	@Failure		500		{object}	custom_errors.HTTPError	"Internal server error"
//	@Router			/auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req LoginInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, LoginResponse{Errors: map[string][]string{"server": {"invalid body"}}})
		return
	}

	access, refresh, err, user := h.service.Login(c, req)
	loginResponse := prepareLoginResponse(access, refresh, *user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, LoginResponse{Errors: map[string][]string{"credentials": {"invalid credentials"}}})
		return
	}
	c.JSON(http.StatusOK, loginResponse)
}

func prepareLoginResponse(access, refresh string, user user.User) *LoginResponse {
	loginResponse := LoginResponse{
		RefreshToken: refresh,
		AccessToken:  access,
		Selector:     user.Selector,
		Username:     user.Username,
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
//	@Description	Refresh access and refresh tokens
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		RefreshInput			true	"Refresh token form"
//	@Success		200		{object}	TokenResponse			"OK - returns tokens"
//	@Failure		400		{object}	custom_errors.HTTPError	"Invalid input data"
//	@Failure		401		{object}	custom_errors.HTTPError	"Unauthorized"
//	@Failure		500		{object}	custom_errors.HTTPError	"Internal server error"
//	@Router			/auth/refresh [post]
func (h *Handler) Refresh(c *gin.Context) {
	var req RefreshInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError("invalid body"))
		return
	}
	access, refresh, err := h.service.Refresh(c, req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, custom_errors.NewHTTPError("unauthorized"))
		return
	}
	response := LoginResponse{AccessToken: access, RefreshToken: refresh}
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
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError("invalid body"))
		return
	}

	userId, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, custom_errors.NewHTTPError("unauthorized"))
		return
	}

	if err := h.service.Logout(c, userId.(uint), req.RefreshToken, req.AllDevices); err != nil {
		c.JSON(http.StatusInternalServerError, custom_errors.NewHTTPError(err.Error()))
		return
	}
	c.Status(http.StatusNoContent)
}
