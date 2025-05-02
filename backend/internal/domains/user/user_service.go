package user

import "github.com/susek555/BD2/car-dealer-api/internal/domains/generic"

type UserServiceInterface interface {
	generic.CRUDService[UserDTO]
	GetByEmail(email string) (UserDTO, error)
}

type UserService struct {
	repo UserRepositoryInterface
}

func NewService(userRepository UserRepositoryInterface) UserServiceInterface {
	return &UserService{repo: userRepository}
}

func (s *UserService) Create(in UserDTO) error {
	user, err := in.MapToUser()
	if err != nil {
		return err
	}
	return s.repo.Create(user)
}

func (s *UserService) Update(in UserDTO) error {
	user, err := in.MapToUser()
	if err != nil {
		return nil
	}
	return s.repo.Update(user)
}

func (s *UserService) GetAll() ([]UserDTO, error) {
	users, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	userDTOs := make([]UserDTO, len(users))
	for _, user := range users {
		dto, _ := user.MapToDTO()
		userDTOs = append(userDTOs, dto)
	}
	return userDTOs, nil
}

func (s *UserService) GetById(id uint) (UserDTO, error) {
	user, err := s.repo.GetById(id)
	if err != nil {
		return UserDTO{}, err
	}
	userDTO, _ := user.MapToDTO()
	return userDTO, nil
}

func (s *UserService) GetByEmail(email string) (UserDTO, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return UserDTO{}, err
	}
	userDTO, _ := user.MapToDTO()
	return userDTO, nil
}

func (s *UserService) Delete(id uint) error {
	return s.repo.Delete(id)
}
