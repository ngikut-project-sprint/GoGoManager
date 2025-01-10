package handlers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/ngikut-project-sprint/GoGoManager/internal/config"
	"github.com/ngikut-project-sprint/GoGoManager/internal/handlers"
	"github.com/ngikut-project-sprint/GoGoManager/internal/middleware"
	"github.com/ngikut-project-sprint/GoGoManager/internal/models"
	"github.com/ngikut-project-sprint/GoGoManager/internal/utils"
	mocksServices "github.com/ngikut-project-sprint/GoGoManager/mocks/services"
	mocksUtils "github.com/ngikut-project-sprint/GoGoManager/mocks/utils"
)

type MockRequestBody struct {
	Data  string
	Error error
}

func (m *MockRequestBody) Read(p []byte) (n int, err error) {
	copy(p, m.Data)
	return len(m.Data), m.Error
}

func (m *MockRequestBody) Close() error {
	return nil
}

func TestAuthHandler_RegisterManager_Success(t *testing.T) {
	mockService := new(mocksServices.ManagerService)
	mockJWTGen := &mocksUtils.JWTGenerator{}
	mockBCrypt := &mocksUtils.Encryption{}
	handler := handlers.NewAuthHandler(mockService, mockJWTGen.GenerateJWT, mockBCrypt.CompareHashAndPassword)

	manager_id := 1
	email := "random@name.com"
	password := "cobalagi"
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MjgsImVtYWlsIjoiYW5kb21AbmFtZS5jb20iLCJleHAiOjE3MzY1MTgzMjMsImlhdCI6MTczNjQzMTkyM30.Hc91AimPJHnXZmwUxAd25g_YyJabZmZMuwtss0wt5zg"

	mockBody := &MockRequestBody{
		Data:  `{"email": "random@name.com", "password": "cobalagi", "action": "create"}`,
		Error: nil,
	}

	cfg := &config.Config{
		Database: config.DatabaseConfig{},
		JWT: config.JWTConfig{
			Secret: "i-am-amazing-huntsman",
		},
	}

	req := httptest.NewRequest(http.MethodPost, "/v1/auth", mockBody)
	res := httptest.NewRecorder()

	mockService.On("Create", email, password).Return(manager_id, nil)
	mockJWTGen.On("GenerateJWT", cfg.JWT.Secret, manager_id, email).Return(token, nil)

	middleware.ConfigMiddleware(cfg, http.HandlerFunc(handler.Auth)).ServeHTTP(res, req)

	assert.Equal(t, http.StatusCreated, res.Code)

	mockService.AssertExpectations(t)
	mockJWTGen.AssertExpectations(t)
	mockBCrypt.AssertExpectations(t)
}

func TestAuthHandler_RegisterManager_WrongMethod(t *testing.T) {
	mockService := new(mocksServices.ManagerService)
	mockJWTGen := &mocksUtils.JWTGenerator{}
	mockBCrypt := &mocksUtils.Encryption{}
	handler := handlers.NewAuthHandler(mockService, mockJWTGen.GenerateJWT, mockBCrypt.CompareHashAndPassword)

	mockBody := &MockRequestBody{
		Data:  `{"email": "random@name.com", "password": "cobalagi", "action": "create"}`,
		Error: nil,
	}

	cfg := &config.Config{
		Database: config.DatabaseConfig{},
		JWT: config.JWTConfig{
			Secret: "i-am-amazing-huntsman",
		},
	}

	req := httptest.NewRequest(http.MethodDelete, "/v1/auth", mockBody)
	res := httptest.NewRecorder()

	middleware.ConfigMiddleware(cfg, http.HandlerFunc(handler.Auth)).ServeHTTP(res, req)

	assert.Equal(t, http.StatusMethodNotAllowed, res.Code)

	mockService.AssertExpectations(t)
	mockJWTGen.AssertExpectations(t)
	mockBCrypt.AssertExpectations(t)
}

func TestAuthHandler_RegisterManager_CorruptRequestBody(t *testing.T) {
	mockService := new(mocksServices.ManagerService)
	mockJWTGen := &mocksUtils.JWTGenerator{}
	mockBCrypt := &mocksUtils.Encryption{}
	handler := handlers.NewAuthHandler(mockService, mockJWTGen.GenerateJWT, mockBCrypt.CompareHashAndPassword)

	mockBody := &MockRequestBody{
		Data:  `{"email": "rando`,
		Error: nil,
	}

	cfg := &config.Config{
		Database: config.DatabaseConfig{},
		JWT: config.JWTConfig{
			Secret: "i-am-amazing-huntsman",
		},
	}

	req := httptest.NewRequest(http.MethodPost, "/v1/auth", mockBody)
	res := httptest.NewRecorder()

	middleware.ConfigMiddleware(cfg, http.HandlerFunc(handler.Auth)).ServeHTTP(res, req)

	assert.Equal(t, http.StatusBadRequest, res.Code)

	mockService.AssertExpectations(t)
	mockJWTGen.AssertExpectations(t)
	mockBCrypt.AssertExpectations(t)
}

func TestAuthHandler_RegisterManager_ConfigNotFound(t *testing.T) {
	mockService := new(mocksServices.ManagerService)
	mockJWTGen := &mocksUtils.JWTGenerator{}
	mockBCrypt := &mocksUtils.Encryption{}
	handler := handlers.NewAuthHandler(mockService, mockJWTGen.GenerateJWT, mockBCrypt.CompareHashAndPassword)

	mockBody := &MockRequestBody{
		Data:  `{"email": "random@name.com", "password": "cobalagi", "action": "create"}`,
		Error: nil,
	}

	req := httptest.NewRequest(http.MethodPost, "/v1/auth", mockBody)
	res := httptest.NewRecorder()

	handler.Auth(res, req)

	assert.Equal(t, http.StatusInternalServerError, res.Code)

	mockService.AssertExpectations(t)
	mockJWTGen.AssertExpectations(t)
	mockBCrypt.AssertExpectations(t)
}

func TestAuthHandler_RegisterManager_EmailAlreadyRegistered(t *testing.T) {
	mockService := new(mocksServices.ManagerService)
	mockJWTGen := &mocksUtils.JWTGenerator{}
	mockBCrypt := &mocksUtils.Encryption{}
	handler := handlers.NewAuthHandler(mockService, mockJWTGen.GenerateJWT, mockBCrypt.CompareHashAndPassword)

	email := "random@name.com"
	password := "cobalagi"

	mockBody := &MockRequestBody{
		Data:  `{"email": "random@name.com", "password": "cobalagi", "action": "create"}`,
		Error: nil,
	}

	cfg := &config.Config{
		Database: config.DatabaseConfig{},
		JWT: config.JWTConfig{
			Secret: "i-am-amazing-huntsman",
		},
	}

	e := errors.New("Email already registered")
	error := &utils.GoGoError{
		Type:    utils.SQLUniqueViolated,
		Message: "Email already registered",
		Err:     e,
	}

	req := httptest.NewRequest(http.MethodPost, "/v1/auth", mockBody)
	res := httptest.NewRecorder()

	mockService.On("Create", email, password).Return(0, error)

	middleware.ConfigMiddleware(cfg, http.HandlerFunc(handler.Auth)).ServeHTTP(res, req)

	assert.Equal(t, http.StatusConflict, res.Code)

	mockService.AssertExpectations(t)
	mockJWTGen.AssertExpectations(t)
	mockBCrypt.AssertExpectations(t)
}

func TestAuthHandler_RegisterManager_InvalidEmailFormat(t *testing.T) {
	mockService := new(mocksServices.ManagerService)
	mockJWTGen := &mocksUtils.JWTGenerator{}
	mockBCrypt := &mocksUtils.Encryption{}
	handler := handlers.NewAuthHandler(mockService, mockJWTGen.GenerateJWT, mockBCrypt.CompareHashAndPassword)

	email := "random@name.com"
	password := "cobalagi"

	mockBody := &MockRequestBody{
		Data:  `{"email": "random@name.com", "password": "cobalagi", "action": "create"}`,
		Error: nil,
	}

	cfg := &config.Config{
		Database: config.DatabaseConfig{},
		JWT: config.JWTConfig{
			Secret: "i-am-amazing-huntsman",
		},
	}

	e := errors.New("Invalid email format")
	error := &utils.GoGoError{
		Type:    utils.InvalidEmailFormat,
		Message: "Invalid email format",
		Err:     e,
	}

	req := httptest.NewRequest(http.MethodPost, "/v1/auth", mockBody)
	res := httptest.NewRecorder()

	mockService.On("Create", email, password).Return(0, error)

	middleware.ConfigMiddleware(cfg, http.HandlerFunc(handler.Auth)).ServeHTTP(res, req)

	assert.Equal(t, http.StatusBadRequest, res.Code)

	mockService.AssertExpectations(t)
	mockJWTGen.AssertExpectations(t)
	mockBCrypt.AssertExpectations(t)
}

func TestAuthHandler_RegisterManager_InvalidPasswordLeght(t *testing.T) {
	mockService := new(mocksServices.ManagerService)
	mockJWTGen := &mocksUtils.JWTGenerator{}
	mockBCrypt := &mocksUtils.Encryption{}
	handler := handlers.NewAuthHandler(mockService, mockJWTGen.GenerateJWT, mockBCrypt.CompareHashAndPassword)

	email := "random@name.com"
	password := "cobalagi"

	mockBody := &MockRequestBody{
		Data:  `{"email": "random@name.com", "password": "cobalagi", "action": "create"}`,
		Error: nil,
	}

	cfg := &config.Config{
		Database: config.DatabaseConfig{},
		JWT: config.JWTConfig{
			Secret: "i-am-amazing-huntsman",
		},
	}

	e := errors.New("Invalid password length (min length: 8, max length: 32)")
	error := &utils.GoGoError{
		Type:    utils.InvalidEmailFormat,
		Message: "Invalid password length",
		Err:     e,
	}

	req := httptest.NewRequest(http.MethodPost, "/v1/auth", mockBody)
	res := httptest.NewRecorder()

	mockService.On("Create", email, password).Return(0, error)

	middleware.ConfigMiddleware(cfg, http.HandlerFunc(handler.Auth)).ServeHTTP(res, req)

	assert.Equal(t, http.StatusBadRequest, res.Code)

	mockService.AssertExpectations(t)
	mockJWTGen.AssertExpectations(t)
	mockBCrypt.AssertExpectations(t)
}

func TestAuthHandler_RegisterManager_DatabaseError(t *testing.T) {
	mockService := new(mocksServices.ManagerService)
	mockJWTGen := &mocksUtils.JWTGenerator{}
	mockBCrypt := &mocksUtils.Encryption{}
	handler := handlers.NewAuthHandler(mockService, mockJWTGen.GenerateJWT, mockBCrypt.CompareHashAndPassword)

	email := "random@name.com"
	password := "cobalagi"

	mockBody := &MockRequestBody{
		Data:  `{"email": "random@name.com", "password": "cobalagi", "action": "create"}`,
		Error: nil,
	}

	cfg := &config.Config{
		Database: config.DatabaseConfig{},
		JWT: config.JWTConfig{
			Secret: "i-am-amazing-huntsman",
		},
	}

	e := errors.New("Database error")
	error := &utils.GoGoError{
		Type:    utils.SQLError,
		Message: "Database error",
		Err:     e,
	}

	req := httptest.NewRequest(http.MethodPost, "/v1/auth", mockBody)
	res := httptest.NewRecorder()

	mockService.On("Create", email, password).Return(0, error)

	middleware.ConfigMiddleware(cfg, http.HandlerFunc(handler.Auth)).ServeHTTP(res, req)

	assert.Equal(t, http.StatusInternalServerError, res.Code)

	mockService.AssertExpectations(t)
	mockJWTGen.AssertExpectations(t)
	mockBCrypt.AssertExpectations(t)
}

func TestAuthHandler_RegisterManager_JWTGenerateError(t *testing.T) {
	mockService := new(mocksServices.ManagerService)
	mockJWTGen := &mocksUtils.JWTGenerator{}
	mockBCrypt := &mocksUtils.Encryption{}
	handler := handlers.NewAuthHandler(mockService, mockJWTGen.GenerateJWT, mockBCrypt.CompareHashAndPassword)

	manager_id := 1
	email := "random@name.com"
	password := "cobalagi"

	mockBody := &MockRequestBody{
		Data:  `{"email": "random@name.com", "password": "cobalagi", "action": "create"}`,
		Error: nil,
	}

	cfg := &config.Config{
		Database: config.DatabaseConfig{},
		JWT: config.JWTConfig{
			Secret: "i-am-amazing-huntsman",
		},
	}

	req := httptest.NewRequest(http.MethodPost, "/v1/auth", mockBody)
	res := httptest.NewRecorder()

	mockService.On("Create", email, password).Return(manager_id, nil)
	mockJWTGen.On("GenerateJWT", cfg.JWT.Secret, manager_id, email).Return("", errors.New("Failed to generate JWT"))

	middleware.ConfigMiddleware(cfg, http.HandlerFunc(handler.Auth)).ServeHTTP(res, req)

	assert.Equal(t, http.StatusInternalServerError, res.Code)

	mockService.AssertExpectations(t)
	mockJWTGen.AssertExpectations(t)
	mockBCrypt.AssertExpectations(t)
}

func TestAuthHandler_LoginManager_Success(t *testing.T) {
	mockService := new(mocksServices.ManagerService)
	mockJWTGen := &mocksUtils.JWTGenerator{}
	mockBCrypt := &mocksUtils.Encryption{}
	handler := handlers.NewAuthHandler(mockService, mockJWTGen.GenerateJWT, mockBCrypt.CompareHashAndPassword)

	manager_id := 1
	email := "random@name.com"
	password := "cobalagi"
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MjgsImVtYWlsIjoiYW5kb21AbmFtZS5jb20iLCJleHAiOjE3MzY1MTgzMjMsImlhdCI6MTczNjQzMTkyM30.Hc91AimPJHnXZmwUxAd25g_YyJabZmZMuwtss0wt5zg"

	manager := &models.Manager{
		ID:              1,
		Email:           "random@name.com",
		Password:        "cobalagi",
		Name:            ptr("Manager One"),
		UserImageUri:    ptr("image1.png"),
		CompanyName:     ptr("Company A"),
		CompanyImageUri: ptr("company1.png"),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		DeletedAt:       nil,
	}

	mockBody := &MockRequestBody{
		Data:  `{"email": "random@name.com", "password": "cobalagi", "action": "login"}`,
		Error: nil,
	}

	cfg := &config.Config{
		Database: config.DatabaseConfig{},
		JWT: config.JWTConfig{
			Secret: "i-am-amazing-huntsman",
		},
	}

	req := httptest.NewRequest(http.MethodPost, "/v1/auth", mockBody)
	res := httptest.NewRecorder()

	mockService.On("GetByEmail", email).Return(manager, nil)
	mockBCrypt.On("CompareHashAndPassword", []byte(manager.Password), []byte(password)).Return(nil)
	mockJWTGen.On("GenerateJWT", cfg.JWT.Secret, manager_id, email).Return(token, nil)

	middleware.ConfigMiddleware(cfg, http.HandlerFunc(handler.Auth)).ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)

	mockService.AssertExpectations(t)
	mockJWTGen.AssertExpectations(t)
	mockBCrypt.AssertExpectations(t)
}

func TestAuthHandler_LoginManager_WrongMethod(t *testing.T) {
	mockService := new(mocksServices.ManagerService)
	mockJWTGen := &mocksUtils.JWTGenerator{}
	mockBCrypt := &mocksUtils.Encryption{}
	handler := handlers.NewAuthHandler(mockService, mockJWTGen.GenerateJWT, mockBCrypt.CompareHashAndPassword)

	mockBody := &MockRequestBody{
		Data:  `{"email": "random@name.com", "password": "cobalagi", "action": "login"}`,
		Error: nil,
	}

	cfg := &config.Config{
		Database: config.DatabaseConfig{},
		JWT: config.JWTConfig{
			Secret: "i-am-amazing-huntsman",
		},
	}

	req := httptest.NewRequest(http.MethodDelete, "/v1/auth", mockBody)
	res := httptest.NewRecorder()

	middleware.ConfigMiddleware(cfg, http.HandlerFunc(handler.Auth)).ServeHTTP(res, req)

	assert.Equal(t, http.StatusMethodNotAllowed, res.Code)

	mockService.AssertExpectations(t)
	mockJWTGen.AssertExpectations(t)
	mockBCrypt.AssertExpectations(t)
}

func TestAuthHandler_LoginManager_CorruptRequestBody(t *testing.T) {
	mockService := new(mocksServices.ManagerService)
	mockJWTGen := &mocksUtils.JWTGenerator{}
	mockBCrypt := &mocksUtils.Encryption{}
	handler := handlers.NewAuthHandler(mockService, mockJWTGen.GenerateJWT, mockBCrypt.CompareHashAndPassword)

	mockBody := &MockRequestBody{
		Data:  `{"email": "rando`,
		Error: nil,
	}

	cfg := &config.Config{
		Database: config.DatabaseConfig{},
		JWT: config.JWTConfig{
			Secret: "i-am-amazing-huntsman",
		},
	}

	req := httptest.NewRequest(http.MethodPost, "/v1/auth", mockBody)
	res := httptest.NewRecorder()

	middleware.ConfigMiddleware(cfg, http.HandlerFunc(handler.Auth)).ServeHTTP(res, req)

	assert.Equal(t, http.StatusBadRequest, res.Code)

	mockService.AssertExpectations(t)
	mockJWTGen.AssertExpectations(t)
	mockBCrypt.AssertExpectations(t)
}

func TestAuthHandler_LoginManager_ConfigNotFound(t *testing.T) {
	mockService := new(mocksServices.ManagerService)
	mockJWTGen := &mocksUtils.JWTGenerator{}
	mockBCrypt := &mocksUtils.Encryption{}
	handler := handlers.NewAuthHandler(mockService, mockJWTGen.GenerateJWT, mockBCrypt.CompareHashAndPassword)

	mockBody := &MockRequestBody{
		Data:  `{"email": "random@name.com", "password": "cobalagi", "action": "login"}`,
		Error: nil,
	}

	req := httptest.NewRequest(http.MethodPost, "/v1/auth", mockBody)
	res := httptest.NewRecorder()

	handler.Auth(res, req)

	assert.Equal(t, http.StatusInternalServerError, res.Code)

	mockService.AssertExpectations(t)
	mockJWTGen.AssertExpectations(t)
	mockBCrypt.AssertExpectations(t)
}

func TestAuthHandler_LoginManager_UserNotFound(t *testing.T) {
	mockService := new(mocksServices.ManagerService)
	mockJWTGen := &mocksUtils.JWTGenerator{}
	mockBCrypt := &mocksUtils.Encryption{}
	handler := handlers.NewAuthHandler(mockService, mockJWTGen.GenerateJWT, mockBCrypt.CompareHashAndPassword)

	email := "random@name.com"

	mockBody := &MockRequestBody{
		Data:  `{"email": "random@name.com", "password": "cobalagi", "action": "login"}`,
		Error: nil,
	}

	cfg := &config.Config{
		Database: config.DatabaseConfig{},
		JWT: config.JWTConfig{
			Secret: "i-am-amazing-huntsman",
		},
	}

	e := errors.New("User not found")
	error := &utils.GoGoError{
		Type:    utils.SQLError,
		Message: "Email not found",
		Err:     e,
	}

	req := httptest.NewRequest(http.MethodPost, "/v1/auth", mockBody)
	res := httptest.NewRecorder()

	mockService.On("GetByEmail", email).Return(nil, error)

	middleware.ConfigMiddleware(cfg, http.HandlerFunc(handler.Auth)).ServeHTTP(res, req)

	assert.Equal(t, http.StatusNotFound, res.Code)

	mockService.AssertExpectations(t)
	mockJWTGen.AssertExpectations(t)
	mockBCrypt.AssertExpectations(t)
}

func TestAuthHandler_LoginManager_Unauthorized(t *testing.T) {
	mockService := new(mocksServices.ManagerService)
	mockJWTGen := &mocksUtils.JWTGenerator{}
	mockBCrypt := &mocksUtils.Encryption{}
	handler := handlers.NewAuthHandler(mockService, mockJWTGen.GenerateJWT, mockBCrypt.CompareHashAndPassword)

	email := "random@name.com"
	password := "cobalagi"

	manager := &models.Manager{
		ID:              1,
		Email:           "random@name.com",
		Password:        "cobalagi",
		Name:            ptr("Manager One"),
		UserImageUri:    ptr("image1.png"),
		CompanyName:     ptr("Company A"),
		CompanyImageUri: ptr("company1.png"),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		DeletedAt:       nil,
	}

	mockBody := &MockRequestBody{
		Data:  `{"email": "random@name.com", "password": "cobalagi", "action": "login"}`,
		Error: nil,
	}

	cfg := &config.Config{
		Database: config.DatabaseConfig{},
		JWT: config.JWTConfig{
			Secret: "i-am-amazing-huntsman",
		},
	}

	e := errors.New("Invalid credential")

	req := httptest.NewRequest(http.MethodPost, "/v1/auth", mockBody)
	res := httptest.NewRecorder()

	mockService.On("GetByEmail", email).Return(manager, nil)
	mockBCrypt.On("CompareHashAndPassword", []byte(manager.Password), []byte(password)).Return(e)

	middleware.ConfigMiddleware(cfg, http.HandlerFunc(handler.Auth)).ServeHTTP(res, req)

	assert.Equal(t, http.StatusUnauthorized, res.Code)

	mockService.AssertExpectations(t)
	mockJWTGen.AssertExpectations(t)
	mockBCrypt.AssertExpectations(t)
}

func TestAuthHandler_LoginManager_FailedGenerateJWT(t *testing.T) {
	mockService := new(mocksServices.ManagerService)
	mockJWTGen := &mocksUtils.JWTGenerator{}
	mockBCrypt := &mocksUtils.Encryption{}
	handler := handlers.NewAuthHandler(mockService, mockJWTGen.GenerateJWT, mockBCrypt.CompareHashAndPassword)

	manager_id := 1
	email := "random@name.com"
	password := "cobalagi"

	manager := &models.Manager{
		ID:              1,
		Email:           "random@name.com",
		Password:        "cobalagi",
		Name:            ptr("Manager One"),
		UserImageUri:    ptr("image1.png"),
		CompanyName:     ptr("Company A"),
		CompanyImageUri: ptr("company1.png"),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		DeletedAt:       nil,
	}

	mockBody := &MockRequestBody{
		Data:  `{"email": "random@name.com", "password": "cobalagi", "action": "login"}`,
		Error: nil,
	}

	cfg := &config.Config{
		Database: config.DatabaseConfig{},
		JWT: config.JWTConfig{
			Secret: "i-am-amazing-huntsman",
		},
	}

	e := errors.New("Failed generate JWT")

	req := httptest.NewRequest(http.MethodPost, "/v1/auth", mockBody)
	res := httptest.NewRecorder()

	mockService.On("GetByEmail", email).Return(manager, nil)
	mockBCrypt.On("CompareHashAndPassword", []byte(manager.Password), []byte(password)).Return(nil)
	mockJWTGen.On("GenerateJWT", cfg.JWT.Secret, manager_id, email).Return("", e)

	middleware.ConfigMiddleware(cfg, http.HandlerFunc(handler.Auth)).ServeHTTP(res, req)

	assert.Equal(t, http.StatusInternalServerError, res.Code)

	mockService.AssertExpectations(t)
	mockJWTGen.AssertExpectations(t)
	mockBCrypt.AssertExpectations(t)
}

// Helper

func ptr(s string) *string {
	return &s
}
