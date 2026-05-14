package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/M-kos/hitalent_test/internal/config"
	"github.com/M-kos/hitalent_test/internal/constants"
	"github.com/M-kos/hitalent_test/internal/domain"
	"github.com/M-kos/hitalent_test/internal/handler/dto"
	"github.com/M-kos/hitalent_test/pkg/logger"
)

type DepartmentService interface {
	CreateDepartment(ctx context.Context, department *domain.Department) (*domain.Department, error)
	CreateEmployee(ctx context.Context, employee *domain.Employee) (*domain.Employee, error)
	Department(
		ctx context.Context,
		id int,
		depth int,
		includeEmployees bool) (*domain.DepartmentTree, error)
	UpdateDepartment(ctx context.Context, department *domain.Department) (*domain.Department, error)
	DeleteDepartment(ctx context.Context, departmentId int, mode string, reassignId int) error
}

const (
	CreateDepartmentUrl = "POST /departments"
	CreateEmployeeUrl   = "POST /departments/{id}/employees"
	GetDepartmentUrl    = "GET /departments/{id}"
	UpdateDepartmentUrl = "PATCH /departments/{id}"
	DeleteDepartmentUrl = "DELETE /departments/{id}"
)

type Department struct {
	config  *config.Config
	logger  *logger.Logger
	service DepartmentService
}

func (d *Department) createDepartment(w http.ResponseWriter, r *http.Request) {
	var createDepartmentDto dto.CreateDepartmentRequest
	if err := json.NewDecoder(r.Body).Decode(&createDepartmentDto); err != nil {
		d.logger.Handler.Error("[Create Department] error decoding create department", "error", err)
		writeErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := createDepartmentDto.Validate(); err != nil {
		d.logger.Handler.Error("[Create Department] error validating create department", "error", err)
		writeErrorJSON(w, fmt.Sprintf("validation error: %s", err.Error()), http.StatusBadRequest)
		return
	}

	department, err := d.service.CreateDepartment(r.Context(), createDepartmentDto.ToDomain())
	if err != nil {
		if errors.Is(err, domain.ErrWrongParentId) {
			d.logger.Handler.Error("[Create Department] wrong parent id", "error", err, "parentId", createDepartmentDto.ParentID)
			writeErrorJSON(w, err.Error(), http.StatusConflict)
			return
		}

		d.logger.Handler.Error("[Create Department] error creating department", "error", err)
		writeErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var departmentResponse dto.CreateDepartmentResponse
	departmentResponse.FromDomain(department)

	err = writeJSON(w, http.StatusOK, departmentResponse)
	if err != nil {
		d.logger.Handler.Error("[Create Department] write response error", "error", err)
		writeErrorJSON(w, err.Error(), http.StatusInternalServerError)
	}
}

func (d *Department) createEmployee(w http.ResponseWriter, r *http.Request) {
	queryId := r.PathValue("id")
	departmentId, err := strconv.Atoi(queryId)
	if err != nil {
		d.logger.Handler.Error("[Create Employee] error converting query department id to int", "error", err)
		writeErrorJSON(w, "invalid department id", http.StatusBadRequest)
		return
	}

	var createEmployeeDto dto.CreateEmployeeRequest
	if err := json.NewDecoder(r.Body).Decode(&createEmployeeDto); err != nil {
		d.logger.Handler.Error("[Create Employee] error decoding create employee", "error", err)
		writeErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := createEmployeeDto.Validate(); err != nil {
		d.logger.Handler.Error("[Create Employee] error validating create employee", "error", err)
		writeErrorJSON(w, fmt.Sprintf("validation error: %s", err.Error()), http.StatusBadRequest)
		return
	}

	emplDto, err := createEmployeeDto.ToDomain(departmentId)
	if err != nil {
		d.logger.Handler.Error("[Create Employee] error creating employee", "error", err)
		writeErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	employee, err := d.service.CreateEmployee(r.Context(), emplDto)
	if err != nil {
		if errors.Is(err, domain.ErrDepartmentNotFound) {
			d.logger.Handler.Error("[Create Employee] wrong department", "error", err, "departmentId", departmentId)
			writeErrorJSON(w, err.Error(), http.StatusNotFound)
			return
		}

		d.logger.Handler.Error("[Create Employee] error creating employee", "error", err)
		writeErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var employeeResponse dto.CreateEmployeeResponse
	employeeResponse.FromDomain(employee)

	err = writeJSON(w, http.StatusOK, employeeResponse)
	if err != nil {
		d.logger.Handler.Error("[Create Employee] write response error", "error", err)
		writeErrorJSON(w, err.Error(), http.StatusInternalServerError)
	}
}

func (d *Department) department(w http.ResponseWriter, r *http.Request) {
	queryId := r.PathValue("id")
	departmentId, err := strconv.Atoi(queryId)
	if err != nil {
		d.logger.Handler.Error("[Get Department] error converting query department id to int", "error", err)
		writeErrorJSON(w, "invalid department id", http.StatusBadRequest)
		return
	}

	queryDepth := r.URL.Query().Get("depth")
	queryIncludeEmployees := r.URL.Query().Get("include_employees")

	depth, err := strconv.Atoi(queryDepth)
	if err != nil || depth < 1 {
		depth = 1
	}

	if depth > 5 {
		depth = 5
	}

	includeEmployees, err := strconv.ParseBool(queryIncludeEmployees)
	if err != nil {
		includeEmployees = true
	}

	tree, err := d.service.Department(r.Context(), departmentId, depth, includeEmployees)
	if err != nil {
		if errors.Is(err, domain.ErrDepartmentNotFound) {
			d.logger.Handler.Error("[Get Department] wrong department", "error", err)
			writeErrorJSON(w, err.Error(), http.StatusNotFound)
			return
		}

		d.logger.Handler.Error("[Get Department] error getting department", "error", err)
		writeErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var departmentResponse dto.DepartmentResponse
	departmentResponse.FromDomain(tree.Department, tree.Employees, tree.ChildrenDepartment)

	err = writeJSON(w, http.StatusOK, departmentResponse)
	if err != nil {
		d.logger.Handler.Error("[Get Department] write response error", "error", err)
		writeErrorJSON(w, err.Error(), http.StatusInternalServerError)
	}
}

func (d *Department) updateDepartment(w http.ResponseWriter, r *http.Request) {
	queryId := r.PathValue("id")
	departmentId, err := strconv.Atoi(queryId)
	if err != nil {
		d.logger.Handler.Error("[Update Department] error converting query department id to int", "error", err)
		writeErrorJSON(w, "invalid department id", http.StatusBadRequest)
		return
	}

	var updateDepartmentDto dto.UpdateDepartmentRequest
	if err := json.NewDecoder(r.Body).Decode(&updateDepartmentDto); err != nil {
		d.logger.Handler.Error("[Update Department] error decoding update department", "error", err)
		writeErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := updateDepartmentDto.Validate(); err != nil {
		d.logger.Handler.Error("[Update Department] error validating update department", "error", err)
		writeErrorJSON(w, fmt.Sprintf("validation error: %s", err.Error()), http.StatusBadRequest)
		return
	}

	department, err := d.service.UpdateDepartment(r.Context(), updateDepartmentDto.ToDomain(departmentId))
	if err != nil {
		if errors.Is(err, domain.ErrWrongParentId) {
			d.logger.Handler.Error("[Update Department] wrong parent id", "error", err, "parentId", updateDepartmentDto.ParentID)
			writeErrorJSON(w, err.Error(), http.StatusConflict)
			return
		}

		d.logger.Handler.Error("[Update Department] error updating department", "error", err)
		writeErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var departmentResponse dto.UpdateDepartmentResponse
	departmentResponse.FromDomain(department)

	err = writeJSON(w, http.StatusOK, departmentResponse)
	if err != nil {
		d.logger.Handler.Error("[Update Department] write response error", "error", err)
		writeErrorJSON(w, err.Error(), http.StatusInternalServerError)
	}
}

func (d *Department) deleteDepartment(w http.ResponseWriter, r *http.Request) {
	queryId := r.PathValue("id")
	departmentId, err := strconv.Atoi(queryId)
	if err != nil {
		d.logger.Handler.Error("[Update Department] error converting query department id to int", "error", err)
		writeErrorJSON(w, "invalid department id", http.StatusBadRequest)
		return
	}

	mode := r.URL.Query().Get("mode")
	queryReassignDepartmentId := r.URL.Query().Get("reassign_to_department_id")

	if ok := modeValidate(mode); !ok {
		d.logger.Handler.Error("[Delete Department] wrong delete mode", "mode", mode)
		writeErrorJSON(w, "invalid mode", http.StatusBadRequest)
		return
	}

	reassignDepartmentId, err := strconv.Atoi(queryReassignDepartmentId)
	if err != nil && mode == constants.ModeReassign {
		d.logger.Handler.Error("[Delete Department] wrong reassign department id", "error", err)
		writeErrorJSON(w, "invalid reassign department id", http.StatusBadRequest)
		return
	}

	err = d.service.DeleteDepartment(r.Context(), departmentId, mode, reassignDepartmentId)
	if err != nil {
		if errors.Is(err, domain.ErrDepartmentNotFound) {
			d.logger.Handler.Error("[Delete Department] wrong department id", "error", err)
			writeErrorJSON(w, err.Error(), http.StatusNotFound)
			return
		}

		d.logger.Handler.Error("[Delete Department] error deleting department", "error", err)
		writeErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func NewDepartmentHandler(router *http.ServeMux, config *config.Config, logger *logger.Logger, service DepartmentService) {
	handler := &Department{
		config:  config,
		logger:  logger,
		service: service,
	}

	router.HandleFunc(CreateDepartmentUrl, handler.createDepartment)
	router.HandleFunc(CreateEmployeeUrl, handler.createEmployee)
	router.HandleFunc(GetDepartmentUrl, handler.department)
	router.HandleFunc(UpdateDepartmentUrl, handler.updateDepartment)
	router.HandleFunc(DeleteDepartmentUrl, handler.deleteDepartment)
}

func writeErrorJSON(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(dto.ErrorResponse{Error: message}) //nolint:errcheck
}

func writeJSON(w http.ResponseWriter, statusCode int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(data)
}

func modeValidate(mode string) bool {
	switch mode {
	case constants.ModeCascade:
		return true
	case constants.ModeReassign:
		return true
	default:
		return false
	}
}
