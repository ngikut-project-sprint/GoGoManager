package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"
    "log"

	"github.com/ngikut-project-sprint/GoGoManager/internal/constants"
    "github.com/ngikut-project-sprint/GoGoManager/internal/services"
    "github.com/ngikut-project-sprint/GoGoManager/internal/utils"
)

var req struct {
    Name string `json:"name"`
}

type DepartmentHandler struct {
    service services.DepartmentService
}

func NewDepartmentHandler(service services.DepartmentService) *DepartmentHandler {
    return &DepartmentHandler{service: service}
}

func (h *DepartmentHandler) HandleDepartment(w http.ResponseWriter, r *http.Request) {
    log.Printf("HandleDepartment called with method: %s and path: %s\n", r.Method, r.URL.Path)

    switch r.Method {
    case http.MethodGet:
        h.ListDepartments(w, r)
    case http.MethodPost:
        h.CreateDepartment(w, r)
    default:
        utils.MethodNotAllowed(w, r.Method)
    }
}

func (h *DepartmentHandler) HandleDepartmentWithID(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    if idStr == "" {
        utils.BadRequest(w, "Department ID is required")
        return    }
    
    departmentID, err := strconv.Atoi(idStr)
    if err != nil {
        utils.BadRequest(w, "Department ID is required")
        return   
     }

    switch r.Method {
    case http.MethodPatch:
        h.UpdateDepartment(w, r, departmentID)  // Pass departmentID here
    case http.MethodDelete:
        h.DeleteDepartment(w, r, departmentID)
    default:
        utils.MethodNotAllowed(w, r.Method)
    }
}

func (h *DepartmentHandler) CreateDepartment(w http.ResponseWriter, r *http.Request) {

    w.Header().Set("Content-Type", "application/json; charset=utf-8")

    // Get manager ID from context
    claims, ok := r.Context().Value(constants.JWTKey).(*utils.Claims)
	if !ok {
		http.Error(w, "User not aunthenticated", http.StatusUnauthorized)
		return
	} 
    if !ok {
        utils.SendErrorResponse(w, 
            "Invalid user ID format",
            http.StatusUnauthorized)
        return
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        utils.SendErrorResponse(w, 
            "Invalid request body",
            http.StatusBadRequest)
        return
    }

    // Validate name
    if len(req.Name) < 4 || len(req.Name) > 33 {
        utils.SendErrorResponse(w, 
            "Name must be between 4 and 33 characters",
            http.StatusBadRequest)
        return
    }

    dept, err := h.service.CreateDepartment(req.Name, claims.ID)
    if err != nil {
        utils.SendErrorResponse(w, 
            "Failed to create department",
            http.StatusInternalServerError)
        return
    }

    // Success response
    w.WriteHeader(http.StatusCreated)

    if err := json.NewEncoder(w).Encode(dept); err != nil {
        utils.SendErrorResponse(w, "Failed create department", http.StatusBadRequest)
        return
    }
}


func (h *DepartmentHandler) ListDepartments(w http.ResponseWriter, r *http.Request) {
    // Parse query parameters
    query := r.URL.Query()
    limit, _ := strconv.Atoi(query.Get("limit"))
    offset, _ := strconv.Atoi(query.Get("offset"))
    name := query.Get("name")

    // Set defaults
    if limit <= 0 {
        limit = 5
    }
    if offset < 0 {
        offset = 0
    }

    departments, err := h.service.GetDepartments(limit, offset, name)
    if err != nil {
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    if err := json.NewEncoder(w).Encode(departments); err != nil {
        utils.SendErrorResponse(w, "Failed get department list", http.StatusBadRequest)
        return
    }
}

func (h *DepartmentHandler) UpdateDepartment(w http.ResponseWriter, r *http.Request, departmentID int) { // Add departmentID parameter
    w.Header().Set("Content-Type", "application/json; charset=utf-8")

    // Get manager ID from token
    userID, ok := r.Context().Value("user_id").(int)
    if !ok {
        utils.SendErrorResponse(w, "User ID not found in token", http.StatusUnauthorized)
        return
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        utils.SendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Validate name
    if len(req.Name) < 4 || len(req.Name) > 33 {
        utils.SendErrorResponse(w, "Name must be between 4 and 33 characters", http.StatusBadRequest)
        return
    }

    // Update department
    dept, err := h.service.UpdateDepartment(departmentID, req.Name, userID)
    if err != nil {
        utils.SendErrorResponse(w, "Failed to update department", http.StatusInternalServerError)
        return
    }

    utils.WriteJSON(w, http.StatusOK, dept)
}

func (h *DepartmentHandler) DeleteDepartment(w http.ResponseWriter, r *http.Request, departmentID int) {
    w.Header().Set("Content-Type", "application/json; charset=utf-8")

    // Get manager ID from token
    userID, ok := r.Context().Value("user_id").(int)
    if !ok {
        utils.SendErrorResponse(w, "User ID not found in token", http.StatusUnauthorized)
        return
    }

    // Delete department
    err := h.service.DeleteDepartment(departmentID, userID)
    if err != nil {       
        utils.SendErrorResponse(w, "Failed to delete department", http.StatusUnauthorized)
        return
    }

    // Return success response
    utils.WriteJSON(w, http.StatusOK, map[string]string{
        "message": "Department deleted successfully",
    })
}
