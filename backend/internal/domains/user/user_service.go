package user

import "github.com/susek555/BD2/car-dealer-api/internal/domains/generic"

type UserService struct {
	generic.GenericService[User, UserRepositoryInterface]
}

func GetUserService(repo UserRepositoryInterface) *UserService {
	return &UserService{
		GenericService: generic.GenericService[User, UserRepositoryInterface]{Repo: repo},
	}
}

func (service *UserService) GetUserByEmail(email string) (User, error) {
	return service.Repo.GetUserByEmail(email)
}
