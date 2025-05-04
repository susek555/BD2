package user

import (
	"errors"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

var (
	ErrInvalidSelector error = errors.New("user selector has to be P (person) or C (company)")
	ErrHashPassword    error = errors.New("error occured while hashing password")
	ErrCreateUser      error = errors.New("not all required fields (username, password, email) provided ")
	ErrCreateCompany   error = errors.New("company_name and company_nip must be provided")
	ErrCreatePerson    error = errors.New("person_name and person_surname must be provided")
)

// Map to return error code based on catched error
var ErrorMap = map[error]int{
	ErrInvalidSelector:     http.StatusBadRequest,
	ErrCreateCompany:       http.StatusBadRequest,
	ErrCreatePerson:        http.StatusBadRequest,
	ErrHashPassword:        http.StatusInternalServerError,
	strconv.ErrSyntax:      http.StatusBadRequest,
	gorm.ErrRecordNotFound: http.StatusNotFound,
}
