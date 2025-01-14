package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/ngikut-project-sprint/GoGoManager/internal/constants"
	"github.com/ngikut-project-sprint/GoGoManager/internal/models"
	"github.com/ngikut-project-sprint/GoGoManager/internal/utils"
)

type EmployeeRepository interface {
	List(ctx context.Context, filter models.FilterOptions) ([]models.Employee, error)
	Create(ctx context.Context, employee *models.Employee) (*models.Employee, error)
	Update(ctx context.Context, identityNumber string, req models.UpdateEmployeeRequest) (*models.Employee, error)
	Delete(ctx context.Context, identityNumber string) error
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
	claims, ok := ctx.Value(constants.JWTKey).(*utils.Claims)
	if !ok {
		return nil, fmt.Errorf("unauthorized: missing or invalid JWT claims")
	}

	query := `
			SELECT e.id, e.identity_number, e.name, e.employee_image_uri, e.gender, e.department_id, 
						 e.created_at, e.updated_at, e.deleted_at
			FROM employees e
			JOIN departments d ON e.department_id = d.department_id
			WHERE e.deleted_at IS NULL
			AND d.manager_id = $1
	`
	args := []interface{}{claims.ID} // Current manager's ID from JWT
	argCount := 2                    // Start from 2 since we used $1 for manager_id

	if filter.IdentityNumber != nil {
		query += fmt.Sprintf(" AND e.identity_number LIKE $%d", argCount)
		args = append(args, "%"+*filter.IdentityNumber+"%")
		argCount++
	}

	if filter.Gender != nil {
		query += fmt.Sprintf(" AND e.gender = $%d", argCount)
		args = append(args, *filter.Gender)
		argCount++
	}

	if filter.DepartmentID != nil {
		query += fmt.Sprintf(" AND e.department_id = $%d", argCount)
		args = append(args, *filter.DepartmentID)
		argCount++
	}

	if filter.Name != nil {
		query += fmt.Sprintf(" AND e.name LIKE $%d", argCount)
		args = append(args, "%"+*filter.Name+"%")
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

func (r *employeeRepository) Create(ctx context.Context, employee *models.Employee) (*models.Employee, error) {
	query := `
			INSERT INTO employees (
					identity_number, name, employee_image_uri, gender, department_id,
					created_at, updated_at
			) VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
			RETURNING id, identity_number, name, employee_image_uri, gender, department_id, 
								created_at, updated_at, deleted_at
	`

	row := r.db.QueryRowContext(
		ctx,
		query,
		employee.IdentityNumber,
		employee.Name,
		employee.EmployeeImageURI,
		employee.Gender,
		employee.DepartmentID,
	)

	err := row.Scan(
		&employee.ID,
		&employee.IdentityNumber,
		&employee.Name,
		&employee.EmployeeImageURI,
		&employee.Gender,
		&employee.DepartmentID,
		&employee.CreatedAt,
		&employee.UpdatedAt,
		&employee.DeletedAt,
	)
	if err != nil {
		log.Println("error repo: ", err.Error())
		return nil, fmt.Errorf("error creating employee: %w", err)
	}

	return employee, nil
}

func (r *employeeRepository) Update(ctx context.Context, identityNumber string, req models.UpdateEmployeeRequest) (*models.Employee, error) {
	claims, ok := ctx.Value(constants.JWTKey).(*utils.Claims)
	if !ok {
		return nil, fmt.Errorf("unauthorized: missing or invalid JWT claims")
	}

	var existingDeptID int
	err := r.db.QueryRowContext(ctx, `
			SELECT e.department_id 
			FROM employees e
			JOIN departments d ON e.department_id = d.department_id
			WHERE e.identity_number = $1 
			AND e.deleted_at IS NULL 
			AND d.manager_id = $2`,
		identityNumber, claims.ID,
	).Scan(&existingDeptID)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("employee not found or unauthorized")
	}
	if err != nil {
		return nil, fmt.Errorf("error verifying employee: %w", err)
	}

	if req.DepartmentID != nil {
		var count int
		err := r.db.QueryRowContext(ctx,
			"SELECT COUNT(*) FROM departments WHERE department_id = $1 AND manager_id = $2 AND deleted_at IS NULL",
			*req.DepartmentID, claims.ID,
		).Scan(&count)

		if err != nil {
			return nil, fmt.Errorf("error verifying new department: %w", err)
		}

		if count == 0 {
			return nil, fmt.Errorf("unauthorized: new department does not belong to the current manager")
		}
	}

	// Build dynamic update query
	query := "UPDATE employees SET updated_at = NOW()"
	args := []interface{}{}
	argCount := 1

	// Only include fields that are provided in the request
	if req.Name != nil {
		query += fmt.Sprintf(", name = $%d", argCount)
		args = append(args, *req.Name)
		argCount++
	}
	if req.EmployeeImageURI != nil {
		query += fmt.Sprintf(", employee_image_uri = $%d", argCount)
		args = append(args, *req.EmployeeImageURI)
		argCount++
	}
	if req.Gender != nil {
		query += fmt.Sprintf(", gender = $%d", argCount)
		args = append(args, *req.Gender)
		argCount++
	}
	if req.DepartmentID != nil {
		query += fmt.Sprintf(", department_id = $%d", argCount)
		args = append(args, *req.DepartmentID)
		argCount++
	}
	if req.IdentityNumber != nil {
		query += fmt.Sprintf(", identity_number = $%d", argCount)
		args = append(args, *req.IdentityNumber)
		argCount++
	}

	query += fmt.Sprintf(" WHERE identity_number = $%d AND deleted_at IS NULL", argCount)
	args = append(args, identityNumber)
	query += " RETURNING id, identity_number, name, employee_image_uri, gender, department_id, created_at, updated_at, deleted_at"

	var employee models.Employee
	err = r.db.QueryRowContext(ctx, query, args...).Scan(
		&employee.ID,
		&employee.IdentityNumber,
		&employee.Name,
		&employee.EmployeeImageURI,
		&employee.Gender,
		&employee.DepartmentID,
		&employee.CreatedAt,
		&employee.UpdatedAt,
		&employee.DeletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("employee not found")
		}
		return nil, fmt.Errorf("error updating employee: %w", err)
	}

	return &employee, nil
}

func (r *employeeRepository) Delete(ctx context.Context, identityNumber string) error {
	claims, ok := ctx.Value(constants.JWTKey).(*utils.Claims)
	if !ok {
		return fmt.Errorf("unauthorized: missing or invalid JWT claims")
	}

	query := `
			UPDATE employees e
			SET deleted_at = NOW()
			FROM departments d
			WHERE e.department_id = d.department_id
			AND e.identity_number = $1 
			AND e.deleted_at IS NULL
			AND d.manager_id = $2
	`

	result, err := r.db.ExecContext(ctx, query, identityNumber, claims.ID)
	if err != nil {
		return fmt.Errorf("error deleting employee: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking deletion result: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("employee not found") // Changed this line
	}

	return nil
}
