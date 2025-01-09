package services

import (
	"context"

	"github.com/ngikut-project-sprint/GoGoManager/internal/models"
	"github.com/ngikut-project-sprint/GoGoManager/internal/repositories"
)

type EmployeeService interface {
	List(ctx context.Context, filter models.FilterOptions) ([]models.Employee, error)
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
