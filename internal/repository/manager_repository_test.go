package repository_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"

	"github.com/ngikut-project-sprint/GoGoManager/internal/models"
	"github.com/ngikut-project-sprint/GoGoManager/internal/repository"
	"github.com/ngikut-project-sprint/GoGoManager/internal/utils"
	mocksDatabase "github.com/ngikut-project-sprint/GoGoManager/mocks/database"
	mocksUtils "github.com/ngikut-project-sprint/GoGoManager/mocks/utils"
)

func TestManagerRepository_Create_Success(t *testing.T) {
	mockDB := &mocksDatabase.DB{}
	mockRow := &mocksDatabase.Row{}
	mockEncrypt := &mocksUtils.Encryption{}

	repo := repository.NewManagerRepository(mockDB, mockEncrypt.GenerateFromPassword)

	email := "test@email.com"
	password := "securepassword123"
	hashedPassword := []byte("$2a$10$hashedpasswordexample")
	query := `
  INSERT INTO managers (email, password)
  VALUES ($1, $2)
  RETURNING id`

	mockEncrypt.On("GenerateFromPassword", []byte(password), bcrypt.DefaultCost).Return(hashedPassword, nil)

	mockRow.On("Scan", mock.AnythingOfType("*int")).Run(func(args mock.Arguments) {
		*(args[0].(*int)) = 1
	}).Return(nil)

	mockDB.On(
		"QueryRow",
		query,
		email,
		string(hashedPassword),
	).Return(mockRow)

	id, error := repo.Create(email, password)

	utils.NoError(t, error)
	assert.Equal(t, 1, id)

	mockEncrypt.AssertExpectations(t)
	mockRow.AssertExpectations(t)
	mockDB.AssertExpectations(t)
}

func TestManagerRepository_Create_FailedHashPassword(t *testing.T) {
	mockDB := &mocksDatabase.DB{}
	mockRow := &mocksDatabase.Row{}
	mockEncrypt := &mocksUtils.Encryption{}

	repo := repository.NewManagerRepository(mockDB, mockEncrypt.GenerateFromPassword)

	email := "test@email.com"
	password := "securepassword123"

	mockEncrypt.On("GenerateFromPassword", []byte(password), bcrypt.DefaultCost).Return(nil, errors.New("Failed hash password"))

	id, error := repo.Create(email, password)

	utils.Error(t, error)
	assert.Equal(t, 0, id)

	mockEncrypt.AssertExpectations(t)
	mockRow.AssertExpectations(t)
	mockDB.AssertExpectations(t)
}

func TestManagerRepository_Create_EmailAlreadyRegistered(t *testing.T) {
	mockDB := &mocksDatabase.DB{}
	mockRow := &mocksDatabase.Row{}
	mockEncrypt := &mocksUtils.Encryption{}

	repo := repository.NewManagerRepository(mockDB, mockEncrypt.GenerateFromPassword)

	email := "test@email.com"
	password := "securepassword123"
	hashedPassword := []byte("$2a$10$hashedpasswordexample")
	query := `
  INSERT INTO managers (email, password)
  VALUES ($1, $2)
  RETURNING id`

	mockEncrypt.On("GenerateFromPassword", []byte(password), bcrypt.DefaultCost).Return(hashedPassword, nil)

	mockRow.On("Scan", mock.AnythingOfType("*int")).Return(errors.New("Email already registered"))

	mockDB.On(
		"QueryRow",
		query,
		email,
		string(hashedPassword),
	).Return(mockRow)

	id, error := repo.Create(email, password)

	utils.Error(t, error)
	assert.Equal(t, 0, id)

	mockEncrypt.AssertExpectations(t)
	mockRow.AssertExpectations(t)
	mockDB.AssertExpectations(t)
}

func TestManagerRepository_Create_DatabaseError(t *testing.T) {
	mockDB := &mocksDatabase.DB{}
	mockRow := &mocksDatabase.Row{}
	mockEncrypt := &mocksUtils.Encryption{}

	repo := repository.NewManagerRepository(mockDB, mockEncrypt.GenerateFromPassword)

	email := "test@email.com"
	password := "securepassword123"
	hashedPassword := []byte("$2a$10$hashedpasswordexample")
	query := `
  INSERT INTO managers (email, password)
  VALUES ($1, $2)
  RETURNING id`

	mockEncrypt.On("GenerateFromPassword", []byte(password), bcrypt.DefaultCost).Return(hashedPassword, nil)

	mockRow.On("Scan", mock.AnythingOfType("*int")).Return(errors.New("Database error"))

	mockDB.On(
		"QueryRow",
		query,
		email,
		string(hashedPassword),
	).Return(mockRow)

	id, error := repo.Create(email, password)

	utils.Error(t, error)
	assert.Equal(t, 0, id)

	mockEncrypt.AssertExpectations(t)
	mockRow.AssertExpectations(t)
	mockDB.AssertExpectations(t)
}

func TestManagerRepository_GetAll_Success(t *testing.T) {
	mockDB := &mocksDatabase.DB{}
	mockRows := &mocksDatabase.Rows{}
	mockEncrypt := &mocksUtils.Encryption{}

	repo := repository.NewManagerRepository(mockDB, mockEncrypt.GenerateFromPassword)

	expectedManagers := []models.Manager{
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

	var callIndex int

	mockRows.On("Next").Return(true).Times(len(expectedManagers))
	mockRows.On("Next").Return(false).Once()
	mockRows.On(
		"Scan",
		mock.AnythingOfType("*int"),
		mock.AnythingOfType("*string"),
		mock.AnythingOfType("*string"),
		mock.AnythingOfType("**string"),
		mock.AnythingOfType("**string"),
		mock.AnythingOfType("**string"),
		mock.AnythingOfType("**string"),
		mock.AnythingOfType("*time.Time"),
		mock.AnythingOfType("*time.Time"),
		mock.AnythingOfType("**time.Time"),
	).Run(func(args mock.Arguments) {
		// Use callIndex to fetch the current manager
		manager := expectedManagers[callIndex]

		// Map arguments to manager fields
		*(args[0].(*int)) = manager.ID
		*(args[1].(*string)) = manager.Email
		*(args[2].(*string)) = manager.Password
		*(args[3].(**string)) = manager.Name
		*(args[4].(**string)) = manager.UserImageUri
		*(args[5].(**string)) = manager.CompanyName
		*(args[6].(**string)) = manager.CompanyImageUri
		*(args[7].(*time.Time)) = manager.CreatedAt
		*(args[8].(*time.Time)) = manager.UpdatedAt
		*(args[9].(**time.Time)) = manager.DeletedAt

		// Increment the callIndex
		callIndex++
	}).Return(nil).Times(len(expectedManagers))

	mockRows.On("Err").Return(nil)
	mockRows.On("Close").Return(nil)

	query := `
  SELECT id, email, password, name, user_image_uri, company_name, company_image_uri, created_at, updated_at, deleted_at
  FROM managers`

	mockDB.On("Query", query).Return(mockRows, nil)

	managers, err := repo.GetAll()

	utils.NoError(t, err)
	assert.Equal(t, expectedManagers, managers)

	mockRows.AssertExpectations(t)
	mockDB.AssertExpectations(t)
	mockEncrypt.AssertExpectations(t)
}

func TestManagerRepository_GetAll_ScanError(t *testing.T) {
	mockDB := &mocksDatabase.DB{}
	mockRows := &mocksDatabase.Rows{}
	mockEncrypt := &mocksUtils.Encryption{}

	repo := repository.NewManagerRepository(mockDB, mockEncrypt.GenerateFromPassword)

	mockRows.On("Next").Return(true).Once()
	mockRows.On(
		"Scan",
		mock.AnythingOfType("*int"),
		mock.AnythingOfType("*string"),
		mock.AnythingOfType("*string"),
		mock.AnythingOfType("**string"),
		mock.AnythingOfType("**string"),
		mock.AnythingOfType("**string"),
		mock.AnythingOfType("**string"),
		mock.AnythingOfType("*time.Time"),
		mock.AnythingOfType("*time.Time"),
		mock.AnythingOfType("**time.Time"),
	).Return(errors.New("Failed to scan row"))

	mockRows.On("Close").Return(nil)

	query := `
  SELECT id, email, password, name, user_image_uri, company_name, company_image_uri, created_at, updated_at, deleted_at
  FROM managers`

	mockDB.On("Query", query).Return(mockRows, nil)

	managers, err := repo.GetAll()

	utils.Error(t, err)
	assert.Nil(t, managers)
	assert.Equal(t, "Failed to scan row", err.Err.Error())

	mockRows.AssertExpectations(t)
	mockDB.AssertExpectations(t)
	mockEncrypt.AssertExpectations(t)
}

func TestManagerRepository_GetAll_RowsError(t *testing.T) {
	mockDB := &mocksDatabase.DB{}
	mockRows := &mocksDatabase.Rows{}
	mockEncrypt := &mocksUtils.Encryption{}

	repo := repository.NewManagerRepository(mockDB, mockEncrypt.GenerateFromPassword)

	mockRows.On("Next").Return(false).Once()
	mockRows.On("Err").Return(errors.New("Error scanning row"))
	mockRows.On("Close").Return(nil)

	query := `
  SELECT id, email, password, name, user_image_uri, company_name, company_image_uri, created_at, updated_at, deleted_at
  FROM managers`

	mockDB.On("Query", query).Return(mockRows, nil)

	managers, err := repo.GetAll()

	utils.Error(t, err)
	assert.Nil(t, managers)
	assert.Equal(t, "Error scanning row", err.Err.Error())

	mockRows.AssertExpectations(t)
	mockDB.AssertExpectations(t)
	mockEncrypt.AssertExpectations(t)
}

func TestManagerRepository_GetAll_DatabaseError(t *testing.T) {
	mockDB := &mocksDatabase.DB{}
	mockRows := &mocksDatabase.Rows{}
	mockEncrypt := &mocksUtils.Encryption{}

	repo := repository.NewManagerRepository(mockDB, mockEncrypt.GenerateFromPassword)

	query := `
  SELECT id, email, password, name, user_image_uri, company_name, company_image_uri, created_at, updated_at, deleted_at
  FROM managers`

	mockDB.On("Query", query).Return(nil, errors.New("Failed to execute query"))

	managers, err := repo.GetAll()

	utils.Error(t, err)
	assert.Nil(t, managers)

	mockRows.AssertExpectations(t)
	mockDB.AssertExpectations(t)
	mockEncrypt.AssertExpectations(t)
}

func TestManagerRepository_GetByID_Success(t *testing.T) {
	mockDB := &mocksDatabase.DB{}
	mockRow := &mocksDatabase.Row{}
	mockEncrypt := &mocksUtils.Encryption{}

	repo := repository.NewManagerRepository(mockDB, mockEncrypt.GenerateFromPassword)

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

	mockRow.On(
		"Scan",
		mock.AnythingOfType("*int"),
		mock.AnythingOfType("*string"),
		mock.AnythingOfType("*string"),
		mock.AnythingOfType("**string"),
		mock.AnythingOfType("**string"),
		mock.AnythingOfType("**string"),
		mock.AnythingOfType("**string"),
		mock.AnythingOfType("*time.Time"),
		mock.AnythingOfType("*time.Time"),
		mock.AnythingOfType("**time.Time"),
	).Run(func(args mock.Arguments) {
		*(args[0].(*int)) = manager.ID
		*(args[1].(*string)) = manager.Email
		*(args[2].(*string)) = manager.Password
		*(args[3].(**string)) = manager.Name
		*(args[4].(**string)) = manager.UserImageUri
		*(args[5].(**string)) = manager.CompanyName
		*(args[6].(**string)) = manager.CompanyImageUri
		*(args[7].(*time.Time)) = manager.CreatedAt
		*(args[8].(*time.Time)) = manager.UpdatedAt
		*(args[9].(**time.Time)) = manager.DeletedAt
	}).Return(nil)

	query := `
  SELECT id, email, password, name, user_image_uri, company_name, company_image_uri, created_at, updated_at, deleted_at 
  FROM managers 
  WHERE id = $1`

	mockDB.On("QueryRow", query, manager.ID).Return(mockRow)

	actualManager, err := repo.GetByID(manager.ID)

	utils.NoError(t, err)
	assert.Equal(t, manager, actualManager)

	mockRow.AssertExpectations(t)
	mockDB.AssertExpectations(t)
	mockEncrypt.AssertExpectations(t)
}

func TestManagerRepository_GetByID_DatabaseError(t *testing.T) {
	mockDB := &mocksDatabase.DB{}
	mockRow := &mocksDatabase.Row{}
	mockEncrypt := &mocksUtils.Encryption{}

	repo := repository.NewManagerRepository(mockDB, mockEncrypt.GenerateFromPassword)

	id := 1

	mockRow.On(
		"Scan",
		mock.AnythingOfType("*int"),
		mock.AnythingOfType("*string"),
		mock.AnythingOfType("*string"),
		mock.AnythingOfType("**string"),
		mock.AnythingOfType("**string"),
		mock.AnythingOfType("**string"),
		mock.AnythingOfType("**string"),
		mock.AnythingOfType("*time.Time"),
		mock.AnythingOfType("*time.Time"),
		mock.AnythingOfType("**time.Time"),
	).Return(errors.New("Error scannin row"))

	query := `
  SELECT id, email, password, name, user_image_uri, company_name, company_image_uri, created_at, updated_at, deleted_at 
  FROM managers 
  WHERE id = $1`

	mockDB.On("QueryRow", query, id).Return(mockRow)

	actualManager, err := repo.GetByID(id)

	utils.Error(t, err)
	assert.Nil(t, actualManager)

	mockRow.AssertExpectations(t)
	mockDB.AssertExpectations(t)
	mockEncrypt.AssertExpectations(t)
}

func TestManagerRepository_GetByEmail_Success(t *testing.T) {
	mockDB := &mocksDatabase.DB{}
	mockRow := &mocksDatabase.Row{}
	mockEncrypt := &mocksUtils.Encryption{}

	repo := repository.NewManagerRepository(mockDB, mockEncrypt.GenerateFromPassword)

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

	mockRow.On(
		"Scan",
		mock.AnythingOfType("*int"),
		mock.AnythingOfType("*string"),
		mock.AnythingOfType("*string"),
		mock.AnythingOfType("**string"),
		mock.AnythingOfType("**string"),
		mock.AnythingOfType("**string"),
		mock.AnythingOfType("**string"),
		mock.AnythingOfType("*time.Time"),
		mock.AnythingOfType("*time.Time"),
		mock.AnythingOfType("**time.Time"),
	).Run(func(args mock.Arguments) {
		*(args[0].(*int)) = manager.ID
		*(args[1].(*string)) = manager.Email
		*(args[2].(*string)) = manager.Password
		*(args[3].(**string)) = manager.Name
		*(args[4].(**string)) = manager.UserImageUri
		*(args[5].(**string)) = manager.CompanyName
		*(args[6].(**string)) = manager.CompanyImageUri
		*(args[7].(*time.Time)) = manager.CreatedAt
		*(args[8].(*time.Time)) = manager.UpdatedAt
		*(args[9].(**time.Time)) = manager.DeletedAt
	}).Return(nil)

	query := `
  SELECT id, email, password, name, user_image_uri, company_name, company_image_uri, created_at, updated_at, deleted_at 
  FROM managers 
  WHERE email = $1`

	mockDB.On("QueryRow", query, manager.Email).Return(mockRow)

	actualManager, err := repo.GetByEmail(manager.Email)

	utils.NoError(t, err)
	assert.Equal(t, manager, actualManager)

	mockRow.AssertExpectations(t)
	mockDB.AssertExpectations(t)
	mockEncrypt.AssertExpectations(t)
}

func TestManagerRepository_GetByEmail_DatabaseError(t *testing.T) {
	mockDB := &mocksDatabase.DB{}
	mockRow := &mocksDatabase.Row{}
	mockEncrypt := &mocksUtils.Encryption{}

	repo := repository.NewManagerRepository(mockDB, mockEncrypt.GenerateFromPassword)

	email := "test1@example.com"

	mockRow.On(
		"Scan",
		mock.AnythingOfType("*int"),
		mock.AnythingOfType("*string"),
		mock.AnythingOfType("*string"),
		mock.AnythingOfType("**string"),
		mock.AnythingOfType("**string"),
		mock.AnythingOfType("**string"),
		mock.AnythingOfType("**string"),
		mock.AnythingOfType("*time.Time"),
		mock.AnythingOfType("*time.Time"),
		mock.AnythingOfType("**time.Time"),
	).Return(errors.New("Error scanning row"))

	query := `
  SELECT id, email, password, name, user_image_uri, company_name, company_image_uri, created_at, updated_at, deleted_at 
  FROM managers 
  WHERE email = $1`

	mockDB.On("QueryRow", query, email).Return(mockRow)

	actualManager, err := repo.GetByEmail(email)

	utils.Error(t, err)
	assert.Nil(t, actualManager)

	mockRow.AssertExpectations(t)
	mockDB.AssertExpectations(t)
	mockEncrypt.AssertExpectations(t)
}

func TestManagerRepository_Update_Success(t *testing.T) {
	mockDB := &mocksDatabase.DB{}
	mockEncrypt := &mocksUtils.Encryption{}

	repo := repository.NewManagerRepository(mockDB, mockEncrypt.GenerateFromPassword)

	hashedPassword := []byte("$2a$10$hashedpasswordexample")
	manager := &utils.ManagerRequest{
		ID:              1,
		Email:           ptr("test1@example.com"),
		Password:        ptr("hashedpassword1"),
		Name:            ptr("Manager One"),
		UserImageUri:    ptr("http://aws-s3.com/image1.png"),
		CompanyName:     ptr("Company A"),
		CompanyImageUri: ptr("http://aws-s3.com/company1.png"),
	}

	query := `UPDATE managers SET email = $1, password = $2, name = $3, user_image_uri = $4, company_name = $5, company_image_uri = $6, updated_at = $7 WHERE id = $8`

	mockEncrypt.On("GenerateFromPassword", []byte(*manager.Password), bcrypt.DefaultCost).Return(hashedPassword, nil)

	mockDB.On("Exec",
		query,
		mock.AnythingOfType("string"),
		mock.AnythingOfType("[]uint8"),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("time.Time"),
		mock.AnythingOfType("int"),
	).Return(nil, nil)

	err := repo.Update(manager)

	utils.NoError(t, err)

	mockDB.AssertExpectations(t)
	mockEncrypt.AssertExpectations(t)
}

func TestManagerRepository_Update_NoField(t *testing.T) {
	mockDB := &mocksDatabase.DB{}
	mockEncrypt := &mocksUtils.Encryption{}

	repo := repository.NewManagerRepository(mockDB, mockEncrypt.GenerateFromPassword)

	manager := &utils.ManagerRequest{
		ID:              1,
		Email:           ptr("old_email@example.com"),
		Password:        nil,
		Name:            nil,
		UserImageUri:    nil,
		CompanyName:     nil,
		CompanyImageUri: nil,
	}

	mockDB.On("Exec",
		`UPDATE managers SET email = $1, updated_at = $2 WHERE id = $3`,
		mock.AnythingOfType("string"),
		mock.AnythingOfType("time.Time"),
		mock.AnythingOfType("int"),
	).Return(nil, nil)

	err := repo.Update(manager)

	utils.NoError(t, err)

	mockDB.AssertExpectations(t)
	mockEncrypt.AssertExpectations(t)
}

func TestManagerRepository_Update_ExecError(t *testing.T) {
	mockDB := &mocksDatabase.DB{}
	mockEncrypt := &mocksUtils.Encryption{}

	repo := repository.NewManagerRepository(mockDB, mockEncrypt.GenerateFromPassword)

	hashedPassword := []byte("$2a$10$hashedpasswordexample")
	manager := &utils.ManagerRequest{
		ID:              1,
		Email:           ptr("test1@example.com"),
		Password:        ptr("hashedpassword1"),
		Name:            ptr("Manager One"),
		UserImageUri:    ptr("http://aws-s3.com/image1.png"),
		CompanyName:     ptr("Company A"),
		CompanyImageUri: ptr("http://aws-s3.com/company1.png"),
	}

	query := `UPDATE managers SET email = $1, password = $2, name = $3, user_image_uri = $4, company_name = $5, company_image_uri = $6, updated_at = $7 WHERE id = $8`

	mockEncrypt.On("GenerateFromPassword", []byte(*manager.Password), bcrypt.DefaultCost).Return(hashedPassword, nil)

	mockDB.On("Exec",
		query,
		mock.AnythingOfType("string"),
		mock.AnythingOfType("[]uint8"),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("time.Time"),
		mock.AnythingOfType("int"),
	).Return(nil, errors.New("Database error"))

	err := repo.Update(manager)

	utils.Error(t, err)

	mockDB.AssertExpectations(t)
	mockEncrypt.AssertExpectations(t)
}

func TestManagerRepository_Update_PartialUpdate(t *testing.T) {
	mockDB := &mocksDatabase.DB{}
	mockEncrypt := &mocksUtils.Encryption{}

	repo := repository.NewManagerRepository(mockDB, mockEncrypt.GenerateFromPassword)

	manager := &utils.ManagerRequest{
		ID:              1,
		Email:           ptr("new_email@example.com"),
		Password:        nil,
		Name:            ptr("New Name"),
		UserImageUri:    nil,
		CompanyName:     nil,
		CompanyImageUri: nil,
	}

	mockDB.On("Exec",
		`UPDATE managers SET email = $1, name = $2, updated_at = $3 WHERE id = $4`,
		mock.AnythingOfType("string"),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("time.Time"),
		mock.AnythingOfType("int"),
	).Return(nil, errors.New("Database error"))

	err := repo.Update(manager)

	utils.Error(t, err)

	mockDB.AssertExpectations(t)
	mockEncrypt.AssertExpectations(t)
}

// Helper

func ptr(s string) *string {
	return &s
}
