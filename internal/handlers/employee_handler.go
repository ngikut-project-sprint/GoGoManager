package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/ngikut-project-sprint/GoGoManager/internal/models"
	"github.com/ngikut-project-sprint/GoGoManager/internal/services"
	"github.com/ngikut-project-sprint/GoGoManager/internal/utils"
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
		utils.SendErrorResponse(w, "Internal server error", http.StatusInternalServerError)
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
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(utils.Response{
		Data:    response,
		Message: fmt.Sprintf("Successfully retrieved %d employees", len(response)),
	}); err != nil {
		log.Printf("Error encoding response: %v", err)
		utils.SendErrorResponse(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *EmployeeHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateEmployeeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Create new employee
	employee, err := h.service.Create(r.Context(), req)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "unique_identity_number"):
			utils.SendErrorResponse(w, fmt.Sprintf("Identity number %s is already registered", req.IdentityNumber), http.StatusConflict)
		default:
			utils.SendErrorResponse(w, "Failed to create employee", http.StatusInternalServerError)
		}
		return
	}

	// Prepare response
	response := EmployeeResponse{
		IdentityNumber:   employee.IdentityNumber,
		Name:             employee.Name,
		EmployeeImageUri: employee.EmployeeImageURI,
		Gender:           string(employee.Gender),
		DepartmentId:     employee.DepartmentID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(utils.Response{
		Data:    response,
		Message: fmt.Sprintf("Employee with ID %s created successfully", response.IdentityNumber),
	}); err != nil {
		log.Printf("Error encoding response: %v", err)
		utils.SendErrorResponse(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *EmployeeHandler) Update(w http.ResponseWriter, r *http.Request, identityNumber string) {
	// Decode request body
	var req models.UpdateEmployeeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Update employee
	employee, err := h.service.Update(r.Context(), identityNumber, req)
	if err != nil {
		utils.SendErrorResponse(w, "Failed to update employee", http.StatusInternalServerError)
		return
	}

	// Prepare response
	response := EmployeeResponse{
		IdentityNumber:   employee.IdentityNumber,
		Name:             employee.Name,
		EmployeeImageUri: employee.EmployeeImageURI,
		Gender:           string(employee.Gender),
		DepartmentId:     employee.DepartmentID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(utils.Response{
		Data:    response,
		Message: fmt.Sprintf("Employee with ID %s updated successfully", response.IdentityNumber),
	}); err != nil {
		log.Printf("Error encoding response: %v", err)
		utils.SendErrorResponse(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *EmployeeHandler) Delete(w http.ResponseWriter, r *http.Request, identityNumber string) {
	err := h.service.Delete(r.Context(), identityNumber)
	if err != nil {
		switch {
		case err.Error() == "employee not found":
			utils.SendErrorResponse(w, fmt.Sprintf("identityNumber %s is not found", identityNumber), http.StatusNotFound)
		case strings.Contains(err.Error(), "unauthorized") ||
			strings.Contains(err.Error(), "expired token") ||
			strings.Contains(err.Error(), "invalid token"):
			utils.SendErrorResponse(w, "expired / invalid / missing request token", http.StatusUnauthorized)
		default:
			utils.SendErrorResponse(w, "Server Error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(utils.Response{
		Message: fmt.Sprintf("Employee with ID %s deleted successfully", identityNumber),
	}); err != nil {
		log.Printf("Error encoding response: %v", err)
		utils.SendErrorResponse(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
