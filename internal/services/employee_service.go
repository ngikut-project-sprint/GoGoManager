package services

import (
	"context"

	"github.com/ngikut-project-sprint/GoGoManager/internal/models"
	"github.com/ngikut-project-sprint/GoGoManager/internal/repositories"
)

type EmployeeService interface {
	List(ctx context.Context, filter models.FilterOptions) ([]models.Employee, error)
	Create(ctx context.Context, req models.CreateEmployeeRequest) (*models.Employee, error)
	Update(ctx context.Context, identityNumber string, req models.UpdateEmployeeRequest) (*models.Employee, error)
	Delete(ctx context.Context, identityNumber string) error
}

type employeeService struct {
	repo repositories.EmployeeRepository
}

func NewEmployeeService(repo repositories.EmployeeRepository) EmployeeService {
	return &employeeService{
		repo: repo,
	}
}

func (s *employeeService) List(ctx context.Context, filter models.FilterOptions) ([]models.Employee, error) {
	// Set defaults if not provided
	if filter.Limit == 0 {
		filter.Limit = 5
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}

	return s.repo.List(ctx, filter)
}

func (s *employeeService) Create(ctx context.Context, req models.CreateEmployeeRequest) (*models.Employee, error) {
	employee := &models.Employee{
		IdentityNumber:   req.IdentityNumber,
		Name:             req.Name,
		EmployeeImageURI: req.EmployeeImageURI,
		Gender:           req.Gender,
		DepartmentID:     req.DepartmentID,
	}

	return s.repo.Create(ctx, employee)
}

func (s *employeeService) Update(ctx context.Context, identityNumber string, req models.UpdateEmployeeRequest) (*models.Employee, error) {
	return s.repo.Update(ctx, identityNumber, req)
}

func (s *employeeService) Delete(ctx context.Context, identityNumber string) error {
	return s.repo.Delete(ctx, identityNumber)
}
