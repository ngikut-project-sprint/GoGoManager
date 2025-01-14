package repository

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"

	"github.com/ngikut-project-sprint/GoGoManager/internal/database"
	"github.com/ngikut-project-sprint/GoGoManager/internal/models"
	"github.com/ngikut-project-sprint/GoGoManager/internal/utils"
)

type ManagerRepository interface {
	Create(email string, password string) (int, *utils.GoGoError)
	GetAll() ([]models.Manager, *utils.GoGoError)
	GetByID(id int) (*models.Manager, *utils.GoGoError)
	GetByEmail(email string) (*models.Manager, *utils.GoGoError)
	Update(manager *utils.ManagerRequest) *utils.GoGoError
}

type managerRepository struct {
	db           database.DB
	hashPassword utils.HashPassword
}

func NewManagerRepository(db database.DB, hashPassword utils.HashPassword) ManagerRepository {
	return &managerRepository{db: db, hashPassword: hashPassword}
}

func (r *managerRepository) Create(email string, password string) (int, *utils.GoGoError) {
	// Hash password
	hashedPassword, err := r.hashPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, utils.WrapError(err, utils.PasswordHashFailed, "Failed to hash password")
	}

	// Insert new manager
	query := `
  INSERT INTO managers (email, password)
  VALUES ($1, $2)
  RETURNING id`

	row := r.db.QueryRow(query, email, string(hashedPassword))

	var id int
	error := row.Scan(&id)
	if error != nil {
		if uniqueErr := utils.UniqueConstraintError(error); uniqueErr != nil {
			return 0, utils.WrapError(error, utils.SQLUniqueViolated, "Email already registered")
		}

		return 0, utils.WrapError(error, utils.SQLError, "Failed to create user")
	}

	return id, nil
}

func (r *managerRepository) GetAll() ([]models.Manager, *utils.GoGoError) {
	var managers []models.Manager

	query := `
  SELECT id, email, password, name, user_image_uri, company_name, company_image_uri, created_at, updated_at, deleted_at
  FROM managers`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, utils.WrapError(err, utils.SQLError, "Error querying managers")
	}
	defer rows.Close()

	for rows.Next() {
		var manager models.Manager
		err := rows.Scan(&manager.ID, &manager.Email, &manager.Password, &manager.Name, &manager.UserImageUri, &manager.CompanyName, &manager.CompanyImageUri, &manager.CreatedAt, &manager.UpdatedAt, &manager.DeletedAt)
		if err != nil {
			return nil, utils.WrapError(err, utils.SQLError, "Error scanning row")
		}
		managers = append(managers, manager)
	}

	if err := rows.Err(); err != nil {
		return nil, utils.WrapError(err, utils.SQLError, "Error scanning row")
	}

	return managers, nil
}

func (r *managerRepository) GetByID(id int) (*models.Manager, *utils.GoGoError) {
	var manager models.Manager

	query := `
  SELECT id, email, password, name, user_image_uri, company_name, company_image_uri, created_at, updated_at, deleted_at 
  FROM managers 
  WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(&manager.ID, &manager.Email, &manager.Password, &manager.Name, &manager.UserImageUri, &manager.CompanyName, &manager.CompanyImageUri, &manager.CreatedAt, &manager.UpdatedAt, &manager.DeletedAt)
	if err != nil {
		return nil, utils.WrapError(err, utils.SQLError, "Error querying manager by id")
	}

	return &manager, nil
}

func (r *managerRepository) GetByEmail(email string) (*models.Manager, *utils.GoGoError) {
	var manager models.Manager

	query := `
  SELECT id, email, password, name, user_image_uri, company_name, company_image_uri, created_at, updated_at, deleted_at 
  FROM managers 
  WHERE email = $1`

	err := r.db.QueryRow(query, email).Scan(&manager.ID, &manager.Email, &manager.Password, &manager.Name, &manager.UserImageUri, &manager.CompanyName, &manager.CompanyImageUri, &manager.CreatedAt, &manager.UpdatedAt, &manager.DeletedAt)
	if err != nil {
		log.Println(err)
		return nil, utils.WrapError(err, utils.SQLError, "Error querying manager by email")
	}

	return &manager, nil
}

func (r *managerRepository) Update(manager *utils.ManagerRequest) *utils.GoGoError {
	query := "UPDATE managers SET "
	var params []interface{}
	var setClauses []string
	paramCounter := 1

	if manager.Email != nil {
		if !manager.ValidEmail() {
			err := errors.New("invalid email format")
			return utils.WrapError(err, utils.InvalidEmailFormat, "Invalid email")
		}
		setClauses = append(setClauses, fmt.Sprintf("email = $%d", paramCounter))
		params = append(params, *manager.Email)
		paramCounter++
	}

	if manager.Password != nil {
		if !manager.ValidPassword() {
			err := errors.New("invalid password format")
			return utils.WrapError(err, utils.InvalidPasswordLength, "Invalid password")
		}
		hashedPassword, err := r.hashPassword([]byte(*manager.Password), bcrypt.DefaultCost)
		if err != nil {
			return utils.WrapError(err, utils.SQLError, "Error hashing password")
		}
		setClauses = append(setClauses, fmt.Sprintf("password = $%d", paramCounter))
		params = append(params, hashedPassword)
		paramCounter++
	}

	if manager.Name != nil {
		if !manager.ValidName() {
			err := errors.New("invalid name length")
			return utils.WrapError(err, utils.InvalidNameLength, "Invalid name")
		}
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", paramCounter))
		params = append(params, *manager.Name)
		paramCounter++
	}

	if manager.UserImageUri != nil {
		if !manager.ValidImageURI() {
			err := errors.New("invalid user image uri format")
			return utils.WrapError(err, utils.InvalidURIFormat, "Invalid user image uri")
		}
		setClauses = append(setClauses, fmt.Sprintf("user_image_uri = $%d", paramCounter))
		params = append(params, *manager.UserImageUri)
		paramCounter++
	}

	if manager.CompanyName != nil {
		if !manager.ValidCompanyName() {
			err := errors.New("invalid company name length")
			return utils.WrapError(err, utils.InvalidNameLength, "Invalid company name")
		}
		setClauses = append(setClauses, fmt.Sprintf("company_name = $%d", paramCounter))
		params = append(params, *manager.CompanyName)
		paramCounter++
	}

	if manager.CompanyImageUri != nil {
		if !manager.ValidCompanyImageURI() {
			err := errors.New("invalid company image uri format")
			return utils.WrapError(err, utils.InvalidURIFormat, "Invalid company image uri")
		}
		setClauses = append(setClauses, fmt.Sprintf("company_image_uri = $%d", paramCounter))
		params = append(params, *manager.CompanyImageUri)
		paramCounter++
	}

	setClauses = append(setClauses, fmt.Sprintf("updated_at = $%d", paramCounter))
	params = append(params, time.Now())
	paramCounter++

	query += strings.Join(setClauses, ", ") + fmt.Sprintf(" WHERE id = $%d", paramCounter)
	params = append(params, manager.ID)

	_, err := r.db.Exec(query, params...)
	if err != nil {
		if uniqueErr := utils.UniqueConstraintError(err); uniqueErr != nil {
			return utils.WrapError(err, utils.SQLUniqueViolated, "Email already registered")
		}
		return utils.WrapError(err, utils.SQLError, "Error updating manager")
	}
	return nil
}
