package services

import (
	"errors"

	"github.com/ngikut-project-sprint/GoGoManager/internal/models"
	"github.com/ngikut-project-sprint/GoGoManager/internal/repositories"
	"github.com/ngikut-project-sprint/GoGoManager/internal/utils"
)

type ManagerService interface {
	Create(email string, password string) (int, *utils.GoGoError)
	GetAll() ([]models.Manager, *utils.GoGoError)
	GetByID(id int) (*models.Manager, *utils.GoGoError)
	GetByEmail(email string) (*models.Manager, *utils.GoGoError)
	Update(manager *models.Manager) *utils.GoGoError
}

type managerService struct {
	managerRepo repositories.ManagerRepository
}

func NewManagerService(managerRepo repositories.ManagerRepository) ManagerService {
	return &managerService{managerRepo: managerRepo}
}

func (s *managerService) Create(email string, password string) (int, *utils.GoGoError) {
	err := utils.ValidateEmail(email)
	if err != nil {
		return 0, utils.WrapError(err, utils.InvalidEmailFormat, "Invalid email format")
	}

	passwordLength := len(password)
	if passwordLength < 8 {
		return 0, utils.WrapError(err, utils.InvalidPasswordLength, "Password too short (> 8 characters)")
	}

	if passwordLength > 52 {
		return 0, utils.WrapError(err, utils.InvalidPasswordLength, "Password too long (< 52 characters)")
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
	err := utils.ValidateEmail(email)
	if err != nil {
		return nil, utils.WrapError(err, utils.InvalidEmailFormat, "Invalid email format")
	}

	return s.managerRepo.GetByEmail(email)
}

func (s *managerService) Update(manager *models.Manager) *utils.GoGoError {
	return s.managerRepo.Update(manager)
}
