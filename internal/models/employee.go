package models

import (
	"time"
)

type Gender string

const (
	Male   Gender = "male"
	Female Gender = "female"
)

type Employee struct {
	ID               int        `json:"id" db:"id"`
	IdentityNumber   string     `json:"identityNumber" db:"identity_number"`
	Name             string     `json:"name" db:"name"`
	EmployeeImageURI string     `json:"employeeImageUri" db:"employee_image_uri"`
	Gender           Gender     `json:"gender" db:"gender"`
	DepartmentID     int        `json:"departmentId" db:"department_id"`
	CreatedAt        time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt        time.Time  `json:"updatedAt" db:"updated_at"`
	DeletedAt        *time.Time `json:"deletedAt,omitempty" db:"deleted_at"`
}

type CreateEmployeeRequest struct {
	IdentityNumber   string `json:"identityNumber" validate:"required,min=5,max=33"`
	Name             string `json:"name" validate:"required,min=4,max=33"`
	EmployeeImageURI string `json:"employeeImageUri" validate:"required,url"`
	Gender           Gender `json:"gender" validate:"required,oneof=male female"`
	DepartmentID     int    `json:"departmentId" validate:"required"`
}

type UpdateEmployeeRequest struct {
	IdentityNumber   *string `json:"identityNumber,omitempty" validate:"omitempty,min=5,max=33"`
	Name             *string `json:"name,omitempty" validate:"omitempty,min=4,max=33"`
	EmployeeImageURI *string `json:"employeeImageUri,omitempty" validate:"omitempty,url"`
	Gender           *Gender `json:"gender,omitempty" validate:"omitempty,oneof=male female"`
	DepartmentID     *int    `json:"departmentId,omitempty" validate:"omitempty"`
}

type FilterOptions struct {
	IdentityNumber *string `json:"identityNumber,omitempty"`
	Gender         *Gender `json:"gender,omitempty"`
	DepartmentID   *int    `json:"departmentId,omitempty"`
	Limit          int     `json:"limit" validate:"required,min=1" default:"10"`
	Offset         int     `json:"offset" validate:"min=0" default:"0"`
}
