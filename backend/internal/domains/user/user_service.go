package user

type UserServiceInterface interface {
	Create(CreateUserDTO) error
	GetAll() ([]RetrieveUserDTO, error)
	GetById(id uint) (RetrieveUserDTO, error)
	GetByEmail(email string) (RetrieveUserDTO, error)
	GetByCompanyNip(email string) (RetrieveUserDTO, error)
	GetByUsername(username string) (RetrieveUserDTO, error)
	Update(UpdateUserDTO) error
	Delete(id uint) error
}

type UserService struct {
	repo UserRepositoryInterface
}

func NewService(userRepository UserRepositoryInterface) UserServiceInterface {
	return &UserService{repo: userRepository}
}

func (s *UserService) Create(in CreateUserDTO) error {
	user, err := in.MapToUser()
	if err != nil {
		return err
	}
	return s.repo.Create(&user)
}

func (s *UserService) GetAll() ([]RetrieveUserDTO, error) {
	users, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	userDTOs := make([]RetrieveUserDTO, 0, len(users))
	for _, user := range users {
		dto, _ := user.MapToDTO()
		userDTOs = append(userDTOs, dto)
	}
	return userDTOs, nil
}

func (s *UserService) GetById(id uint) (RetrieveUserDTO, error) {
	user, err := s.repo.GetById(id)
	if err != nil {
		return RetrieveUserDTO{}, err
	}
	userDTO, _ := user.MapToDTO()
	return userDTO, nil
}

func (s *UserService) GetByEmail(email string) (RetrieveUserDTO, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return RetrieveUserDTO{}, err
	}
	userDTO, _ := user.MapToDTO()
	return userDTO, nil
}

func (s *UserService) GetByUsername(username string) (RetrieveUserDTO, error) {
	user, err := s.repo.GetByUsername(username)
	if err != nil {
		return RetrieveUserDTO{}, err
	}
	userDTO, _ := user.MapToDTO()
	return userDTO, nil
}

func (s *UserService) GetByCompanyNip(nip string) (RetrieveUserDTO, error) {
	user, err := s.repo.GetByCompanyNip(nip)
	if err != nil {
		return RetrieveUserDTO{}, err
	}
	userDTO, _ := user.MapToDTO()
	return userDTO, nil
}

func (s *UserService) Update(in UpdateUserDTO) error {
	user, err := s.repo.GetById(in.ID)
	if err != nil {
		return err
	}
	updatedUser, err := in.UpdateUserFromDTO(&user)
	if err != nil {
		return nil
	}
	return s.repo.Update(updatedUser)
}

func (s *UserService) Delete(id uint) error {
	return s.repo.Delete(id)
}
