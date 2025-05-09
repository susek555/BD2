package model

func (m *Model) MapToDTO() RetrieveModelDTO {
	return RetrieveModelDTO{ID: m.ID, Name: m.Name}
}
