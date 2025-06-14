package manufacturer

import "github.com/susek555/BD2/car-dealer-api/internal/models"

func MapToName(m *models.Manufacturer) *string {
	return &m.Name
}

func MapToDTO(m *models.Manufacturer) *RetrieveManufacturerDTO {
	return &RetrieveManufacturerDTO{ID: m.ID, Name: m.Name}
}
