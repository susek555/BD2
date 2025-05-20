package user

import (
	"errors"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

var (
	ErrInvalidSelector = errors.New("user selector has to be P (person) or C (company)")
	ErrHashPassword    = errors.New("error occurred while hashing password")
	ErrCreateUser      = errors.New("not all required fields (username, password, email) provided ")
	ErrCreateCompany   = errors.New("company_name and company_nip must be provided")
	ErrCreatePerson    = errors.New("person_name and person_surname must be provided")
	ErrUpdatePerson    = errors.New("you cannot update person fields if user is not person")
	ErrUpdateCompany   = errors.New("you cannot update company fields if user is not company")
	ErrForbidden       = errors.New("provided id does not match the id of the logged in user")
)

var ErrorMap = map[error]int{
	ErrInvalidSelector:     http.StatusBadRequest,
	ErrCreateCompany:       http.StatusBadRequest,
	ErrCreatePerson:        http.StatusBadRequest,
	ErrHashPassword:        http.StatusInternalServerError,
	ErrForbidden:           http.StatusForbidden,
	ErrUpdateCompany:       http.StatusBadRequest,
	ErrUpdatePerson:        http.StatusBadRequest,
	strconv.ErrSyntax:      http.StatusBadRequest,
	gorm.ErrRecordNotFound: http.StatusNotFound,
}
