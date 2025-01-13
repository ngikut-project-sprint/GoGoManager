package services

import (
	"errors"
	"fmt"

	"github.com/ngikut-project-sprint/GoGoManager/internal/repository"
)

var (
	ErrDepartmentNotFound     = errors.New("department not found")
	ErrDepartmentHasEmployees = errors.New("department has employees")
)

type DepartmentService interface {
	CreateDepartment(name string, managerID int) (*DepartmentResponse, error)
	GetDepartments(limit, offset int, name string) ([]DepartmentResponse, error)
	UpdateDepartment(id int, name string, managerID int) (*DepartmentResponse, error)
	DeleteDepartment(id int, managerID int) error
}

type departmentService struct {
	repo repository.DepartmentRepository
}

// First, make sure your DepartmentResponse struct is defined correctly
type DepartmentResponse struct {
	DepartmentId int    `json:"departmentId"` // Change type to int to match Department.ID
	Name         string `json:"name"`
}

// Constructor
func NewDepartmentService(repo repository.DepartmentRepository) DepartmentService {
	return &departmentService{
		repo: repo,
	}
}

// Implement all interface methods
func (s *departmentService) CreateDepartment(name string, managerID int) (*DepartmentResponse, error) {
	dept, err := s.repo.Create(name, managerID)
	if err != nil {
		return nil, err
	}

	return &DepartmentResponse{
		DepartmentId: dept.ID,
		Name:         dept.Name,
	}, nil
}

func (s *departmentService) GetDepartments(limit, offset int, name string) ([]DepartmentResponse, error) {
	departments, err := s.repo.FindAll(limit, offset, name)
	if err != nil {
		return nil, fmt.Errorf("failed to find departments: %v", err)
	}

	response := make([]DepartmentResponse, len(departments))
	for i, dept := range departments {
		response[i] = DepartmentResponse{
			DepartmentId: dept.ID,
			Name:         dept.Name,
		}
	}

	return response, nil
}

func (s *departmentService) UpdateDepartment(departmentID int, name string, managerID int) (*DepartmentResponse, error) {
    // Check if department exists and belongs to the manager
    existing, err := s.repo.FindByID(departmentID)
    if err != nil {
        return nil, fmt.Errorf("failed to find department: %v", err)
    }

    if existing.ManagerID != managerID {
        return nil, fmt.Errorf("unauthorized: department does not belong to this manager")
    }

    // Update department
    dept, err := s.repo.Update(departmentID, name)
    if err != nil {
        return nil, fmt.Errorf("failed to update department: %v", err)
    }

    return &DepartmentResponse{
        DepartmentId: dept.ID,
        Name:        dept.Name,
    }, nil
}



func (s *departmentService) DeleteDepartment(departmentID int, managerID int) error {
    // Check if department exists and belongs to the manager
    existing, err := s.repo.FindByID(departmentID)
    if err != nil {
        return fmt.Errorf("failed to find department: %v", err)
    }

    if existing.ManagerID != managerID {
        return fmt.Errorf("unauthorized: department does not belong to this manager")
    }

    // Delete department
    err = s.repo.Delete(departmentID)
    if err != nil {
        return fmt.Errorf("failed to delete department: %v", err)
    }

    return nil
}