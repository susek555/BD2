package manufacturer

import "github.com/susek555/BD2/car-dealer-api/internal/domains/models"

func MapToName(m *models.Manufacturer) *string {
	return &m.Name
}
