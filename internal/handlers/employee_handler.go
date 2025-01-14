package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	DepartmentId     string `json:"departmentId"`
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

	if name := r.URL.Query().Get("name"); name != "" {
		filter.Name = &name
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
			DepartmentId:     strconv.Itoa(emp.DepartmentID),
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		utils.SendErrorResponse(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *EmployeeHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		utils.SendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	r.Body = ioutil.NopCloser(bytes.NewReader(bodyBytes))
	// Decode JSON body into a map
	var body map[string]interface{}
	var decoder = json.NewDecoder(bytes.NewReader(bodyBytes))
	if err := decoder.Decode(&body); err != nil {
		utils.SendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Example checking for a specific field, e.g., "value"
	value, exists := body["departmentId"]
	if !exists {
		utils.SendErrorResponse(w, "Missing Depatment ID", http.StatusBadRequest)
		return
	}

	if _, isString := value.(string); !isString {
		utils.SendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	fmt.Println("masuk")
	fmt.Println()
	var req models.CreateEmployeeRequest
	decoder1 := json.NewDecoder(bytes.NewReader(bodyBytes))
	decoder.DisallowUnknownFields()
	if err := decoder1.Decode(&req); err != nil {
		fmt.Println("decode ", err.Error())
		utils.SendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if !req.ValidIdentityNumber() {
		utils.SendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if !req.ValidName() {
		utils.SendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if !req.ValidImageURI() {
		utils.SendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if !req.ValidGender() {
		utils.SendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if !req.ValidDepartmentId() {
		utils.SendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fmt.Println("masuk sini")

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
		DepartmentId:     strconv.Itoa(employee.DepartmentID),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		utils.SendErrorResponse(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *EmployeeHandler) Update(w http.ResponseWriter, r *http.Request, identityNumber string) {
	if r.Header.Get("Content-Type") != "application/json" {
		utils.SendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Decode request body
	var req models.UpdateEmployeeRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		utils.SendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if !req.ValidIdentityNumber() {
		utils.SendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if !req.ValidName() {
		utils.SendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if !req.ValidImageURI() {
		utils.SendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if !req.ValidGender() {
		utils.SendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if !req.ValidDepartmentId() {
		utils.SendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update employee
	employee, err := h.service.Update(r.Context(), identityNumber, req)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "employee not found"):
			utils.SendErrorResponse(w, "employee not found", http.StatusNotFound)
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
		DepartmentId:     strconv.Itoa(employee.DepartmentID),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		utils.SendErrorResponse(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *EmployeeHandler) Delete(w http.ResponseWriter, r *http.Request, identityNumber string) {
	err := h.service.Delete(r.Context(), identityNumber)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "missing JWT claims"):
			utils.SendErrorResponse(w, "expired / invalid / missing request token", http.StatusUnauthorized)
		case err.Error() == "employee not found":
			utils.SendErrorResponse(w, fmt.Sprintf("identityNumber %s is not found", identityNumber), http.StatusNotFound)
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
