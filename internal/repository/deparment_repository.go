package repository

import (
	"database/sql"
	"fmt"
	
	"github.com/ngikut-project-sprint/GoGoManager/internal/models"
)

// repositories/department.go
type DepartmentRepository interface {
    Create(name string, managerID int) (*models.Department, error)
    FindAll(limit, offset int, name string) ([]models.Department, error)
    FindByID(id int) (*models.Department, error)  // Added
    Update(id int, name string) (*models.Department, error) 
    Delete(id int) error                         // Added
    HasEmployees(id int) (bool, error)           // Added
}

// Define the implementation struct
type departmentRepository struct {
    db *sql.DB
}

// Constructor
func NewDepartmentRepository(db *sql.DB) DepartmentRepository {
    return &departmentRepository{db: db}
}

// Implement Create method
func (r *departmentRepository) Create(name string, managerID int) (*models.Department, error) {
    var dept models.Department
    
    query := `
        INSERT INTO departments (name, manager_id, created_at, updated_at)
        VALUES ($1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        RETURNING department_id, name, manager_id, created_at, updated_at`
    
    err := r.db.QueryRow(query, name, managerID).Scan(
        &dept.ID,
        &dept.Name,
        &dept.ManagerID,
        &dept.CreatedAt,
        &dept.UpdatedAt,
    )
    
    if err != nil {
        return nil, fmt.Errorf("error creating department: %v", err)
    }
    
    return &dept, nil
}

// Implement FindAll method
func (r *departmentRepository) FindAll(limit, offset int, name string) ([]models.Department, error) {
    query := `
        SELECT department_id, name
        FROM departments 
        WHERE ($1 = '' OR name ILIKE $1 || '%')
        ORDER BY created_at DESC
        LIMIT $2 OFFSET $3`

    rows, err := r.db.Query(query, name, limit, offset)
    if err != nil {
        return nil, fmt.Errorf("error querying departments: %v", err)
    }
    defer rows.Close()

    var departments []models.Department
    for rows.Next() {
        var dept models.Department
        err := rows.Scan(
            &dept.ID,
            &dept.Name,
        )
        if err != nil {
            return nil, fmt.Errorf("error scanning department: %v", err)
        }
        departments = append(departments, dept)
    }

    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating departments: %v", err)
    }

    return departments, nil
}

func (r *departmentRepository) FindByID(id int) (*models.Department, error) {
    var dept models.Department
    query := `SELECT department_id, name, manager_id FROM departments WHERE department_id = $1`
    
    err := r.db.QueryRow(query, id).Scan(&dept.ID, &dept.Name, &dept.ManagerID)
    if err == sql.ErrNoRows {
        return nil, fmt.Errorf("department not found")
    }
    if err != nil {
        return nil, fmt.Errorf("error finding department: %v", err)
    }
    
    return &dept, nil
}

func (r *departmentRepository) Update(id int, name string) (*models.Department, error) {
    var dept models.Department
    query := `
        UPDATE departments 
        SET name = $1, updated_at = CURRENT_TIMESTAMP 
        WHERE department_id = $2 
        RETURNING department_id, name`

    err := r.db.QueryRow(query, name, id).Scan(&dept.ID, &dept.Name)
    if err == sql.ErrNoRows {
        return nil, fmt.Errorf("department not found")
    }
    if err != nil {
        return nil, fmt.Errorf("error updating department: %v", err)
    }

    return &dept, nil
}

func (r *departmentRepository) Delete(id int) error {
    query := `DELETE FROM departments WHERE department_id = $1`
    
    result, err := r.db.Exec(query, id)
    if err != nil {
        return fmt.Errorf("error deleting department: %v", err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("error getting rows affected: %v", err)
    }

    if rowsAffected == 0 {
        return fmt.Errorf("department not found")
    }

    return nil
}

func (r *departmentRepository) HasEmployees(id int) (bool, error) {
    var count int
    query := `SELECT COUNT(*) FROM employees WHERE department_id = $1`
    
    err := r.db.QueryRow(query, id).Scan(&count)
    if err != nil {
        return false, err
    }
    
    return count > 0, nil
}

