package user

import (
	"github.com/susek555/BD2/car-dealer-api/pkg/mapping"
)

type UserServiceInterface interface {
	Create(CreateUserDTO) error
	GetAll() ([]RetrieveUserDTO, error)
	GetById(id uint) (*RetrieveUserDTO, error)
	GetByEmail(email string) (*RetrieveUserDTO, error)
	GetByCompanyNip(email string) (*RetrieveUserDTO, error)
	GetByUsername(username string) (*RetrieveUserDTO, error)
	Update(*UpdateUserDTO) map[string][]string
	Delete(id uint) error
}

type UserService struct {
	repo UserRepositoryInterface
}

func NewUserService(userRepository UserRepositoryInterface) UserServiceInterface {
	return &UserService{repo: userRepository}
}

func (s *UserService) Create(in CreateUserDTO) error {
	user, err := in.MapToUser()
	if err != nil {
		return err
	}
	return s.repo.Create(user)
}

func (s *UserService) GetAll() ([]RetrieveUserDTO, error) {
	users, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	userDTOs := mapping.MapSliceToDTOs(users, MapToDTO)
	return userDTOs, nil
}

func (s *UserService) GetById(id uint) (*RetrieveUserDTO, error) {
	user, err := s.repo.GetById(id)
	if err != nil {
		return nil, err
	}
	userDTO := MapToDTO(user)
	return userDTO, nil
}

func (s *UserService) GetByEmail(email string) (*RetrieveUserDTO, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	userDTO := MapToDTO(&user)
	return userDTO, nil
}

func (s *UserService) GetByUsername(username string) (*RetrieveUserDTO, error) {
	user, err := s.repo.GetByUsername(username)
	if err != nil {
		return nil, err
	}
	userDTO := MapToDTO(&user)
	return userDTO, nil
}

func (s *UserService) GetByCompanyNip(nip string) (*RetrieveUserDTO, error) {
	user, err := s.repo.GetByCompanyNip(nip)
	if err != nil {
		return nil, err
	}
	userDTO := MapToDTO(&user)
	return userDTO, nil
}

func (s *UserService) Update(in *UpdateUserDTO) map[string][]string {
	user, err := s.repo.GetById(in.ID)
	var errs = make(map[string][]string)
	if err != nil {
		errs["id"] = []string{err.Error()}
		return errs
	}
	updatedUser, err := in.UpdateUserFromDTO(user)
	if err != nil {
		errs["other"] = []string{err.Error()}
		return errs
	}
	u, noUsername := s.repo.GetByUsername(updatedUser.Username)
	if noUsername == nil && u.ID != updatedUser.ID {
		errs["username"] = []string{ErrUsernameTaken.Error()}
	}
	u, noEmail := s.repo.GetByEmail(updatedUser.Email)
	if noEmail == nil && u.ID != updatedUser.ID {
		errs["email"] = []string{ErrEmailTaken.Error()}
	}
	if updatedUser.Selector == "C" {
		u, noNip := s.repo.GetByCompanyNip(updatedUser.Company.Nip)
		if noNip == nil && u.ID != updatedUser.ID {
			errs["company_nip"] = []string{ErrNipAlreadyTaken.Error()}
		}
	}
	if len(errs) > 0 {
		return errs
	}
	if err := s.repo.Update(updatedUser); err != nil {
		errs["other"] = []string{err.Error()}
	}
	return errs
}

func (s *UserService) Delete(id uint) error {
	return s.repo.Delete(id)
}
