package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ngikut-project-sprint/GoGoManager/internal/models"
)

type EmployeeRepository interface {
	List(ctx context.Context, filter models.FilterOptions) ([]models.Employee, error)
}

type employeeRepository struct {
	db *sql.DB
}

func NewEmployeeRepository(db *sql.DB) EmployeeRepository {
	return &employeeRepository{
		db: db,
	}
}

func (r *employeeRepository) List(ctx context.Context, filter models.FilterOptions) ([]models.Employee, error) {
	query := `
        SELECT id, identity_number, name, employee_image_uri, gender, department_id, 
               created_at, updated_at, deleted_at
        FROM employees
        WHERE deleted_at IS NULL
    `
	args := []interface{}{}
	argCount := 1

	if filter.IdentityNumber != nil {
		query += fmt.Sprintf(" AND identity_number LIKE $%d", argCount)
		args = append(args, "%"+*filter.IdentityNumber+"%")
		argCount++
	}

	if filter.Gender != nil {
		query += fmt.Sprintf(" AND gender = $%d", argCount)
		args = append(args, *filter.Gender)
		argCount++
	}

	if filter.DepartmentID != nil {
		query += fmt.Sprintf(" AND department_id = $%d", argCount)
		args = append(args, *filter.DepartmentID)
		argCount++
	}

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		fmt.Printf("Database error: %v\n", err) // Add this line for debugging
		return nil, fmt.Errorf("error querying employees: %w", err)
	}
	defer rows.Close()

	var employees []models.Employee
	for rows.Next() {
		var emp models.Employee
		err := rows.Scan(
			&emp.ID,
			&emp.IdentityNumber,
			&emp.Name,
			&emp.EmployeeImageURI,
			&emp.Gender,
			&emp.DepartmentID,
			&emp.CreatedAt,
			&emp.UpdatedAt,
			&emp.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning employee: %w", err)
		}
		employees = append(employees, emp)
	}

	return employees, nil
}
