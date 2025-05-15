package mapping

func MapSliceToDTOs[T, Y any](entites []T, mapFunc func(*T) *Y) []Y {
	dtos := make([]Y, 0, len(entites))
	for _, entity := range entites {
		dto := mapFunc(&entity)
		dtos = append(dtos, *dto)
	}
	return dtos
}
