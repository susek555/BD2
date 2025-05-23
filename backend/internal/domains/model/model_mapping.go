package model

import "github.com/susek555/BD2/car-dealer-api/internal/domains/models"

func MapToDTO(m *models.Model) *RetrieveModelDTO {
	return &RetrieveModelDTO{ID: m.ID, Name: m.Name}
}

func MapToName(m *models.Model) *string {
	return &m.Name
}
