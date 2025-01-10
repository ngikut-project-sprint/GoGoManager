package services_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/ngikut-project-sprint/GoGoManager/internal/models"
	"github.com/ngikut-project-sprint/GoGoManager/internal/services"
	"github.com/ngikut-project-sprint/GoGoManager/internal/utils"
	mocksRepo "github.com/ngikut-project-sprint/GoGoManager/mocks/repository"
	mocksValidators "github.com/ngikut-project-sprint/GoGoManager/mocks/validators"
)

func TestManagerService_CreateManagers_Success(t *testing.T) {
	mockRepo := new(mocksRepo.ManagerRepository)
	mockEmailValidator := &mocksValidators.EmailValidator{}
	mockPwdValidator := &mocksValidators.PasswordValidator{}
	service := services.NewManagerService(mockRepo, mockEmailValidator.ValidateEmail, mockPwdValidator.ValidatePassword)

	id := 1
	email := "test@email.com"
	password := "securepassword123"

	mockEmailValidator.On("ValidateEmail", email).Return(nil)
	mockPwdValidator.On("ValidatePassword", password, 8, 52).Return(nil)
	mockRepo.On("Create", email, password).Return(id, nil)

	res, err := service.Create(email, password)

	utils.NoError(t, err)
	assert.Equal(t, id, res)

	mockRepo.AssertExpectations(t)
	mockEmailValidator.AssertExpectations(t)
	mockPwdValidator.AssertExpectations(t)
}

func TestManagerService_CreateManagers_InvalidEmail(t *testing.T) {
	mockRepo := new(mocksRepo.ManagerRepository)
	mockEmailValidator := &mocksValidators.EmailValidator{}
	mockPwdValidator := &mocksValidators.PasswordValidator{}
	service := services.NewManagerService(mockRepo, mockEmailValidator.ValidateEmail, mockPwdValidator.ValidatePassword)

	email := "test@email.c"
	password := "securepassword123"
	e := errors.New("Invalid email format")
	error := &utils.GoGoError{
		Type:    utils.InvalidEmailFormat,
		Message: "Invalid email format",
		Err:     e,
	}

	mockEmailValidator.On("ValidateEmail", email).Return(e)

	res, err := service.Create(email, password)

	assert.Equal(t, error, err)
	assert.Equal(t, 0, res)

	mockRepo.AssertExpectations(t)
	mockEmailValidator.AssertExpectations(t)
	mockPwdValidator.AssertExpectations(t)
}

func TestManagerService_CreateManagers_InvalidPassword(t *testing.T) {
	mockRepo := new(mocksRepo.ManagerRepository)
	mockEmailValidator := &mocksValidators.EmailValidator{}
	mockPwdValidator := &mocksValidators.PasswordValidator{}
	service := services.NewManagerService(mockRepo, mockEmailValidator.ValidateEmail, mockPwdValidator.ValidatePassword)

	email := "test@email.c"
	password := "securepassword123"
	e := errors.New("Invalid password length")
	error := &utils.GoGoError{
		Type:    utils.InvalidPasswordLength,
		Message: "Invalid password length",
		Err:     e,
	}

	mockEmailValidator.On("ValidateEmail", email).Return(nil)
	mockPwdValidator.On("ValidatePassword", password, 8, 52).Return(e)

	res, err := service.Create(email, password)

	assert.Equal(t, error, err)
	assert.Equal(t, 0, res)

	mockRepo.AssertExpectations(t)
	mockEmailValidator.AssertExpectations(t)
	mockPwdValidator.AssertExpectations(t)
}

func TestManagerService_GetAllManagers_Success(t *testing.T) {
	mockRepo := new(mocksRepo.ManagerRepository)
	mockEmailValidator := &mocksValidators.EmailValidator{}
	mockPwdValidator := &mocksValidators.PasswordValidator{}
	service := services.NewManagerService(mockRepo, mockEmailValidator.ValidateEmail, mockPwdValidator.ValidatePassword)

	mockManagers := []models.Manager{
		{
			ID:              1,
			Email:           "test1@example.com",
			Password:        "hashedpassword1",
			Name:            ptr("Manager One"),
			UserImageUri:    ptr("image1.png"),
			CompanyName:     ptr("Company A"),
			CompanyImageUri: ptr("company1.png"),
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
			DeletedAt:       nil,
		},
		{
			ID:              2,
			Email:           "test2@example.com",
			Password:        "hashedpassword2",
			Name:            ptr("Manager Two"),
			UserImageUri:    nil,
			CompanyName:     ptr("Company B"),
			CompanyImageUri: ptr("company2.png"),
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
			DeletedAt:       nil,
		},
	}

	mockRepo.On("GetAll").Return(mockManagers, nil)

	res, err := service.GetAll()

	utils.NoError(t, err)
	assert.Equal(t, mockManagers, res)

	mockRepo.AssertExpectations(t)
	mockEmailValidator.AssertExpectations(t)
	mockPwdValidator.AssertExpectations(t)
}

func TestManagerService_GetAllManagers_Error(t *testing.T) {
	mockRepo := new(mocksRepo.ManagerRepository)
	mockEmailValidator := &mocksValidators.EmailValidator{}
	mockPwdValidator := &mocksValidators.PasswordValidator{}
	service := services.NewManagerService(mockRepo, mockEmailValidator.ValidateEmail, mockPwdValidator.ValidatePassword)

	error := &utils.GoGoError{Err: errors.New("Database error")}

	mockRepo.On("GetAll").Return(nil, error)

	res, err := service.GetAll()

	assert.Equal(t, error, err)
	assert.Nil(t, res)

	mockRepo.AssertExpectations(t)
}

func TestManagerService_GetByIDManager_Success(t *testing.T) {
	mockRepo := new(mocksRepo.ManagerRepository)
	mockEmailValidator := &mocksValidators.EmailValidator{}
	mockPwdValidator := &mocksValidators.PasswordValidator{}
	service := services.NewManagerService(mockRepo, mockEmailValidator.ValidateEmail, mockPwdValidator.ValidatePassword)

	id := 1

	manager := &models.Manager{
		ID:              1,
		Email:           "test1@example.com",
		Password:        "hashedpassword1",
		Name:            ptr("Manager One"),
		UserImageUri:    ptr("image1.png"),
		CompanyName:     ptr("Company A"),
		CompanyImageUri: ptr("company1.png"),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		DeletedAt:       nil,
	}

	mockRepo.On("GetByID", id).Return(manager, nil)

	res, err := service.GetByID(id)

	utils.NoError(t, err)
	assert.Equal(t, manager, res)

	mockRepo.AssertExpectations(t)
	mockEmailValidator.AssertExpectations(t)
	mockPwdValidator.AssertExpectations(t)
}

func TestManagerService_GetByIDManager_InvalidID(t *testing.T) {
	mockRepo := new(mocksRepo.ManagerRepository)
	mockEmailValidator := &mocksValidators.EmailValidator{}
	mockPwdValidator := &mocksValidators.PasswordValidator{}
	service := services.NewManagerService(mockRepo, mockEmailValidator.ValidateEmail, mockPwdValidator.ValidatePassword)

	id := 0
	e := errors.New("Invalid sql id")
	error := &utils.GoGoError{
		Type:    utils.InvalidUserId,
		Message: "Invalid manager id",
		Err:     e,
	}

	res, err := service.GetByID(id)

	assert.Equal(t, error, err)
	assert.Nil(t, res)

	mockRepo.AssertExpectations(t)
	mockEmailValidator.AssertExpectations(t)
	mockPwdValidator.AssertExpectations(t)
}

func TestManagerService_GetByIDManager_Error(t *testing.T) {
	mockRepo := new(mocksRepo.ManagerRepository)
	mockEmailValidator := &mocksValidators.EmailValidator{}
	mockPwdValidator := &mocksValidators.PasswordValidator{}
	service := services.NewManagerService(mockRepo, mockEmailValidator.ValidateEmail, mockPwdValidator.ValidatePassword)

	id := 1
	error := &utils.GoGoError{Err: errors.New("Database error")}

	mockRepo.On("GetByID", id).Return(nil, error)

	res, err := service.GetByID(id)

	assert.Equal(t, error, err)
	assert.Nil(t, res)

	mockRepo.AssertExpectations(t)
	mockEmailValidator.AssertExpectations(t)
	mockPwdValidator.AssertExpectations(t)
}

func TestManagerService_GetByEmailManager_Success(t *testing.T) {
	mockRepo := new(mocksRepo.ManagerRepository)
	mockEmailValidator := &mocksValidators.EmailValidator{}
	mockPwdValidator := &mocksValidators.PasswordValidator{}
	service := services.NewManagerService(mockRepo, mockEmailValidator.ValidateEmail, mockPwdValidator.ValidatePassword)

	email := "test1@example.com"

	manager := &models.Manager{
		ID:              1,
		Email:           "test1@example.com",
		Password:        "hashedpassword1",
		Name:            ptr("Manager One"),
		UserImageUri:    ptr("image1.png"),
		CompanyName:     ptr("Company A"),
		CompanyImageUri: ptr("company1.png"),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		DeletedAt:       nil,
	}

	mockEmailValidator.On("ValidateEmail", email).Return(nil)
	mockRepo.On("GetByEmail", email).Return(manager, nil)

	res, err := service.GetByEmail(email)

	utils.NoError(t, err)
	assert.Equal(t, manager, res)

	mockRepo.AssertExpectations(t)
	mockEmailValidator.AssertExpectations(t)
	mockPwdValidator.AssertExpectations(t)
}

func TestManagerService_GetByEmailManager_InvalidEmail(t *testing.T) {
	mockRepo := new(mocksRepo.ManagerRepository)
	mockEmailValidator := &mocksValidators.EmailValidator{}
	mockPwdValidator := &mocksValidators.PasswordValidator{}
	service := services.NewManagerService(mockRepo, mockEmailValidator.ValidateEmail, mockPwdValidator.ValidatePassword)

	email := "test1@example.com"
	e := errors.New("Invalid email format")
	error := &utils.GoGoError{
		Type:    utils.InvalidEmailFormat,
		Message: "Invalid email format",
		Err:     e,
	}

	mockEmailValidator.On("ValidateEmail", email).Return(e)

	res, err := service.GetByEmail(email)

	assert.Equal(t, error, err)
	assert.Nil(t, res)

	mockRepo.AssertExpectations(t)
	mockEmailValidator.AssertExpectations(t)
	mockPwdValidator.AssertExpectations(t)
}

func TestManagerService_GetByEmailManager_Error(t *testing.T) {
	mockRepo := new(mocksRepo.ManagerRepository)
	mockEmailValidator := &mocksValidators.EmailValidator{}
	mockPwdValidator := &mocksValidators.PasswordValidator{}
	service := services.NewManagerService(mockRepo, mockEmailValidator.ValidateEmail, mockPwdValidator.ValidatePassword)

	email := "test1@example.com"
	error := &utils.GoGoError{Err: errors.New("Database error")}

	mockEmailValidator.On("ValidateEmail", email).Return(nil)
	mockRepo.On("GetByEmail", email).Return(nil, error)

	res, err := service.GetByEmail(email)

	assert.Equal(t, error, err)
	assert.Nil(t, res)

	mockRepo.AssertExpectations(t)
	mockEmailValidator.AssertExpectations(t)
	mockPwdValidator.AssertExpectations(t)
}

func TestManagerService_UpdateManager_Success(t *testing.T) {
	mockRepo := new(mocksRepo.ManagerRepository)
	mockEmailValidator := &mocksValidators.EmailValidator{}
	mockPwdValidator := &mocksValidators.PasswordValidator{}
	service := services.NewManagerService(mockRepo, mockEmailValidator.ValidateEmail, mockPwdValidator.ValidatePassword)

	manager := &utils.ManagerRequest{
		ID:              1,
		Email:           ptr("test1@example.com"),
		Password:        ptr("hashedpassword1"),
		Name:            ptr("Manager One"),
		UserImageUri:    ptr("image1.png"),
		CompanyName:     ptr("Company A"),
		CompanyImageUri: ptr("company1.png"),
	}

	mockRepo.On("Update", manager).Return(nil)

	err := service.Update(manager)

	utils.NoError(t, err)

	mockRepo.AssertExpectations(t)
	mockEmailValidator.AssertExpectations(t)
	mockPwdValidator.AssertExpectations(t)
}

func TestManagerService_UpdateManager_Error(t *testing.T) {
	mockRepo := new(mocksRepo.ManagerRepository)
	mockEmailValidator := &mocksValidators.EmailValidator{}
	mockPwdValidator := &mocksValidators.PasswordValidator{}
	service := services.NewManagerService(mockRepo, mockEmailValidator.ValidateEmail, mockPwdValidator.ValidatePassword)

	error := &utils.GoGoError{Err: errors.New("Database error")}

	manager := &utils.ManagerRequest{
		ID:              1,
		Email:           ptr("test1@example.com"),
		Password:        ptr("hashedpassword1"),
		Name:            ptr("Manager One"),
		UserImageUri:    ptr("image1.png"),
		CompanyName:     ptr("Company A"),
		CompanyImageUri: ptr("company1.png"),
	}

	mockRepo.On("Update", manager).Return(error)

	err := service.Update(manager)

	utils.Error(t, err)

	mockRepo.AssertExpectations(t)
	mockEmailValidator.AssertExpectations(t)
	mockPwdValidator.AssertExpectations(t)
}

// Helper

func ptr(s string) *string {
	return &s
}
