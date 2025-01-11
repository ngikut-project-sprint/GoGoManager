package services

import (
	"errors"

	"github.com/ngikut-project-sprint/GoGoManager/internal/models"
	"github.com/ngikut-project-sprint/GoGoManager/internal/repository"
	"github.com/ngikut-project-sprint/GoGoManager/internal/utils"
)

type ManagerService interface {
	Create(email string, password string) (int, *utils.GoGoError)
	GetAll() ([]models.Manager, *utils.GoGoError)
	GetByID(id int) (*models.Manager, *utils.GoGoError)
	GetByEmail(email string) (*models.Manager, *utils.GoGoError)
	Update(manager *utils.ManagerRequest) *utils.GoGoError
}

type ValidEmailFunc func(email string) error

type ValidPasswordFunc func(email string, min int, max int) error

type managerService struct {
	managerRepo      repository.ManagerRepository
	validateEmail    ValidEmailFunc
	validatePassword ValidPasswordFunc
}

func NewManagerService(
	managerRepo repository.ManagerRepository,
	validateEmail ValidEmailFunc,
	validatePassword ValidPasswordFunc,
) ManagerService {
	return &managerService{managerRepo: managerRepo, validateEmail: validateEmail, validatePassword: validatePassword}
}

func (s *managerService) Create(email string, password string) (int, *utils.GoGoError) {
	emailErr := s.validateEmail(email)
	if emailErr != nil {
		return 0, utils.WrapError(emailErr, utils.InvalidEmailFormat, "Invalid email format")
	}

	pwdErr := s.validatePassword(password, 8, 52)
	if pwdErr != nil {
		return 0, utils.WrapError(pwdErr, utils.InvalidPasswordLength, "Invalid password length")
	}

	return s.managerRepo.Create(email, password)
}

func (s *managerService) GetAll() ([]models.Manager, *utils.GoGoError) {
	return s.managerRepo.GetAll()
}

func (s *managerService) GetByID(id int) (*models.Manager, *utils.GoGoError) {
	if id <= 0 {
		err := errors.New("Invalid sql id")
		return nil, utils.WrapError(err, utils.InvalidUserId, "Invalid manager id")
	}

	return s.managerRepo.GetByID(id)
}

func (s *managerService) GetByEmail(email string) (*models.Manager, *utils.GoGoError) {
	err := s.validateEmail(email)
	if err != nil {
		return nil, utils.WrapError(err, utils.InvalidEmailFormat, "Invalid email format")
	}

	return s.managerRepo.GetByEmail(email)
}

func (s *managerService) Update(manager *utils.ManagerRequest) *utils.GoGoError {
	return s.managerRepo.Update(manager)
}
