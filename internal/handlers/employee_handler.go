package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ngikut-project-sprint/GoGoManager/internal/models"
	"github.com/ngikut-project-sprint/GoGoManager/internal/services"
)

type EmployeeHandler struct {
	service services.EmployeeService
}

type EmployeeResponse struct {
	IdentityNumber   string `json:"identityNumber"`
	Name             string `json:"name"`
	EmployeeImageUri string `json:"employeeImageUri"`
	Gender           string `json:"gender"`
	DepartmentId     int    `json:"departmentId"`
}

func NewEmployeeHandler(service services.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{
		service: service,
	}
}

func (h *EmployeeHandler) List(w http.ResponseWriter, r *http.Request) {
	filter := models.FilterOptions{
		Limit:  5, // default
		Offset: 0, // default
	}

	// Parse query parameters
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			filter.Limit = limit
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			filter.Offset = offset
		}
	}

	if identityNumber := r.URL.Query().Get("identityNumber"); identityNumber != "" {
		filter.IdentityNumber = &identityNumber
	}

	if gender := r.URL.Query().Get("gender"); gender != "" {
		if gender == "male" || gender == "female" {
			g := models.Gender(gender)
			filter.Gender = &g
		}
	}

	if deptID := r.URL.Query().Get("departmentId"); deptID != "" {
		if id, err := strconv.Atoi(deptID); err == nil {
			filter.DepartmentID = &id
		}
	}

	employees, err := h.service.List(r.Context(), filter)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := make([]EmployeeResponse, len(employees))
	for i, emp := range employees {
		response[i] = EmployeeResponse{
			IdentityNumber:   emp.IdentityNumber,
			Name:             emp.Name,
			EmployeeImageUri: emp.EmployeeImageURI,
			Gender:           string(emp.Gender),
			DepartmentId:     emp.DepartmentID,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
