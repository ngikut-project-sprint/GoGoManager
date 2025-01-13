package models

import "time"

type Department struct {
	ID        int       `json:"department_id"`
	Name      string    `json:"name"`
	ManagerID int       `json:"manager_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateDepartmentRequest struct {
	Name string `json:"name" validate:"required,min=4,max=33"`
}

type UpdateDepartmentRequest struct {
	Name string `json:"name" validate:"required,min=4,max=33"`
}

type DepartmentResponse struct {
	DepartmentId string `json:"departmentId"`
	Name         string `json:"name"`
}

type GetDepartmentQuery struct {
	Limit  int    `query:"limit"`
	Offset int    `query:"offset"`
	Name   string `query:"name"`
}
