package sale_offer

import (
	"slices"
	"time"

	"github.com/susek555/BD2/car-dealer-api/pkg/formats"
)

func ParseDate(date string) (*time.Time, error) {
	t, err := time.Parse(formats.DateLayout, date)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func IsParamValid[T comparable](param T, validParams []T) bool {
	return slices.Contains(validParams, param)
}

func AreParamsValid[T comparable](params *[]T, validParams *[]T) bool {
	for _, param := range *params {
		if !IsParamValid(param, *validParams) {
			return false
		}
	}
	return true
}
