package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/M-kos/hitalent_test/internal/config"
	"github.com/M-kos/hitalent_test/internal/domain"
	"github.com/M-kos/hitalent_test/internal/handler/dto"
	"github.com/M-kos/hitalent_test/pkg/logger"
	"github.com/stretchr/testify/assert"
)

type StubDepartmentService struct {
	CreateDepartmentFunc func(ctx context.Context, department *domain.Department) (*domain.Department, error)
	CreateEmployeeFunc   func(ctx context.Context, employee *domain.Employee) (*domain.Employee, error)
	DepartmentFunc       func(ctx context.Context, id int, depth int, includeEmployees bool) (*domain.DepartmentTree, error)
	UpdateDepartmentFunc func(ctx context.Context, department *domain.Department) (*domain.Department, error)
	DeleteDepartmentFunc func(ctx context.Context, departmentId int, mode string, reassignId int) error
}

func (s *StubDepartmentService) CreateDepartment(ctx context.Context, department *domain.Department) (*domain.Department, error) {
	if s.CreateDepartmentFunc != nil {
		return s.CreateDepartmentFunc(ctx, department)
	}
	return nil, nil
}

func (s *StubDepartmentService) CreateEmployee(ctx context.Context, employee *domain.Employee) (*domain.Employee, error) {
	if s.CreateEmployeeFunc != nil {
		return s.CreateEmployeeFunc(ctx, employee)
	}
	return nil, nil
}

func (s *StubDepartmentService) Department(ctx context.Context, id int, depth int, includeEmployees bool) (*domain.DepartmentTree, error) {
	if s.DepartmentFunc != nil {
		return s.DepartmentFunc(ctx, id, depth, includeEmployees)
	}
	return nil, nil
}

func (s *StubDepartmentService) UpdateDepartment(ctx context.Context, department *domain.Department) (*domain.Department, error) {
	if s.UpdateDepartmentFunc != nil {
		return s.UpdateDepartmentFunc(ctx, department)
	}
	return nil, nil
}

func (s *StubDepartmentService) DeleteDepartment(ctx context.Context, departmentId int, mode string, reassignId int) error {
	if s.DeleteDepartmentFunc != nil {
		return s.DeleteDepartmentFunc(ctx, departmentId, mode, reassignId)
	}
	return nil
}

func TestDepartmentHandler_CreateDepartment(t *testing.T) {
	t.Parallel()

	logger := logger.New(&config.Config{LogLevel: "debug"})

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		dept := &domain.Department{ID: 1, Name: "Backend"}
		stubService := &StubDepartmentService{
			CreateDepartmentFunc: func(ctx context.Context, department *domain.Department) (*domain.Department, error) {
				return dept, nil
			},
		}
		handler := &Department{service: stubService, logger: logger, config: &config.Config{}}

		reqBody := `{"name": "Backend"}`
		req := httptest.NewRequest("POST", "/departments", strings.NewReader(reqBody))
		w := httptest.NewRecorder()

		handler.createDepartment(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response dto.CreateDepartmentResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 1, response.ID)
		assert.Equal(t, "Backend", response.Name)
	})

	t.Run("validation error", func(t *testing.T) {
		t.Parallel()
		stubService := &StubDepartmentService{}
		handler := &Department{service: stubService, logger: logger, config: &config.Config{}}

		reqBody := `{"name": ""}`
		req := httptest.NewRequest("POST", "/departments", strings.NewReader(reqBody))
		w := httptest.NewRecorder()

		handler.createDepartment(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "validation error")
	})

	t.Run("wrong parent id", func(t *testing.T) {
		t.Parallel()
		stubService := &StubDepartmentService{
			CreateDepartmentFunc: func(ctx context.Context, department *domain.Department) (*domain.Department, error) {
				return nil, domain.ErrWrongParentId
			},
		}
		handler := &Department{service: stubService, logger: logger, config: &config.Config{}}

		reqBody := `{"name": "Backend", "parent_id": 1}`
		req := httptest.NewRequest("POST", "/departments", strings.NewReader(reqBody))
		w := httptest.NewRecorder()

		handler.createDepartment(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
		assert.Contains(t, w.Body.String(), "wrong parent id")
	})
}

func TestDepartmentHandler_CreateEmployee(t *testing.T) {
	t.Parallel()

	logger := logger.New(&config.Config{LogLevel: "debug"})

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		empl := &domain.Employee{ID: 1, FullName: "Aaa Bbb", Position: "Developer"}
		stubService := &StubDepartmentService{
			CreateEmployeeFunc: func(ctx context.Context, employee *domain.Employee) (*domain.Employee, error) {
				return empl, nil
			},
		}
		handler := &Department{service: stubService, logger: logger, config: &config.Config{}}

		reqBody := `{"full_name": "Aaa Bbb", "position": "Developer"}`
		req := httptest.NewRequest("POST", "/departments/1/employees", strings.NewReader(reqBody))
		req.SetPathValue("id", "1")
		w := httptest.NewRecorder()

		handler.createEmployee(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response dto.CreateEmployeeResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 1, response.ID)
		assert.Equal(t, "Aaa Bbb", response.FullName)
	})

	t.Run("invalid department id", func(t *testing.T) {
		t.Parallel()
		stubService := &StubDepartmentService{}
		handler := &Department{service: stubService, logger: logger, config: &config.Config{}}

		req := httptest.NewRequest("POST", "/departments/invalid/employees", nil)
		w := httptest.NewRecorder()

		handler.createEmployee(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid department id")
	})

	t.Run("department not found", func(t *testing.T) {
		t.Parallel()
		stubService := &StubDepartmentService{
			CreateEmployeeFunc: func(ctx context.Context, employee *domain.Employee) (*domain.Employee, error) {
				return nil, domain.ErrDepartmentNotFound
			},
		}
		handler := &Department{service: stubService, logger: logger, config: &config.Config{}}

		reqBody := `{"full_name": "Aaa Bbb", "position": "Developer"}`
		req := httptest.NewRequest("POST", "/departments/999/employees", strings.NewReader(reqBody))
		req.SetPathValue("id", "999")
		w := httptest.NewRecorder()

		handler.createEmployee(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "department not found")
	})
}

func TestDepartmentHandler_GetDepartment(t *testing.T) {
	t.Parallel()

	logger := logger.New(&config.Config{LogLevel: "debug"})

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		depId := 1
		tree := &domain.DepartmentTree{
			Department:         &domain.Department{ID: depId, Name: "Developers"},
			ChildrenDepartment: []*domain.Department{{ID: 2, Name: "Backend", Parent: &domain.Department{ID: depId, Name: "Developers"}}},
			Employees:          []*domain.Employee{{ID: 1, FullName: "Aaa Bbb", DepartmentID: &depId}},
		}
		stubService := &StubDepartmentService{
			DepartmentFunc: func(ctx context.Context, id int, depth int, includeEmployees bool) (*domain.DepartmentTree, error) {
				return tree, nil
			},
		}
		handler := &Department{service: stubService, logger: logger, config: &config.Config{}}

		req := httptest.NewRequest("GET", "/departments/1?depth=1&include_employees=true", nil)
		req.SetPathValue("id", "1")
		w := httptest.NewRecorder()

		handler.department(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response dto.DepartmentResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 1, response.Department.ID)
		assert.Len(t, response.Children, 1)
		assert.Len(t, response.Employees, 1)
	})

	t.Run("invalid id", func(t *testing.T) {
		t.Parallel()
		stubService := &StubDepartmentService{}
		handler := &Department{service: stubService, logger: logger, config: &config.Config{}}

		req := httptest.NewRequest("GET", "/departments/invalid", nil)
		w := httptest.NewRecorder()

		handler.department(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("department not found", func(t *testing.T) {
		t.Parallel()
		stubService := &StubDepartmentService{
			DepartmentFunc: func(ctx context.Context, id int, depth int, includeEmployees bool) (*domain.DepartmentTree, error) {
				return nil, domain.ErrDepartmentNotFound
			},
		}
		handler := &Department{service: stubService, logger: logger, config: &config.Config{}}

		req := httptest.NewRequest("GET", "/departments/999", nil)
		req.SetPathValue("id", "999")
		w := httptest.NewRecorder()

		handler.department(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestDepartmentHandler_UpdateDepartment(t *testing.T) {
	t.Parallel()

	logger := logger.New(&config.Config{LogLevel: "debug"})

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		dept := &domain.Department{ID: 1, Name: "Backend"}
		stubService := &StubDepartmentService{
			UpdateDepartmentFunc: func(ctx context.Context, department *domain.Department) (*domain.Department, error) {
				return dept, nil
			},
		}
		handler := &Department{service: stubService, logger: logger, config: &config.Config{}}

		reqBody := `{"name": "Backend"}`
		req := httptest.NewRequest("PATCH", "/departments/1", strings.NewReader(reqBody))
		req.SetPathValue("id", "1")
		w := httptest.NewRecorder()

		handler.updateDepartment(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response dto.UpdateDepartmentResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 1, response.ID)
		assert.Equal(t, "Backend", response.Name)
	})

	t.Run("wrong parent id", func(t *testing.T) {
		t.Parallel()
		stubService := &StubDepartmentService{
			UpdateDepartmentFunc: func(ctx context.Context, department *domain.Department) (*domain.Department, error) {
				return nil, domain.ErrWrongParentId
			},
		}
		handler := &Department{service: stubService, logger: logger, config: &config.Config{}}

		reqBody := `{"parent_id": 1}`
		req := httptest.NewRequest("PATCH", "/departments/1", strings.NewReader(reqBody))
		req.SetPathValue("id", "1")
		w := httptest.NewRecorder()

		handler.updateDepartment(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
	})
}

func TestDepartmentHandler_DeleteDepartment(t *testing.T) {
	t.Parallel()

	logger := logger.New(&config.Config{LogLevel: "debug"})

	t.Run("cascade success", func(t *testing.T) {
		t.Parallel()
		stubService := &StubDepartmentService{
			DeleteDepartmentFunc: func(ctx context.Context, departmentId int, mode string, reassignId int) error {
				return nil
			},
		}
		handler := &Department{service: stubService, logger: logger, config: &config.Config{}}

		req := httptest.NewRequest("DELETE", "/departments/1?mode=cascade", nil)
		req.SetPathValue("id", "1")
		w := httptest.NewRecorder()

		handler.deleteDepartment(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("reassign success", func(t *testing.T) {
		t.Parallel()
		stubService := &StubDepartmentService{
			DeleteDepartmentFunc: func(ctx context.Context, departmentId int, mode string, reassignId int) error {
				return nil
			},
		}
		handler := &Department{service: stubService, logger: logger, config: &config.Config{}}

		req := httptest.NewRequest("DELETE", "/departments/1?mode=reassign&reassign_to_department_id=2", nil)
		req.SetPathValue("id", "1")
		w := httptest.NewRecorder()

		handler.deleteDepartment(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("invalid mode", func(t *testing.T) {
		t.Parallel()
		stubService := &StubDepartmentService{}
		handler := &Department{service: stubService, logger: logger, config: &config.Config{}}

		req := httptest.NewRequest("DELETE", "/departments/1?mode=invalid", nil)
		req.SetPathValue("id", "1")
		w := httptest.NewRecorder()

		handler.deleteDepartment(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("department not found", func(t *testing.T) {
		t.Parallel()
		stubService := &StubDepartmentService{
			DeleteDepartmentFunc: func(ctx context.Context, departmentId int, mode string, reassignId int) error {
				return domain.ErrDepartmentNotFound
			},
		}
		handler := &Department{service: stubService, logger: logger, config: &config.Config{}}

		req := httptest.NewRequest("DELETE", "/departments/999?mode=cascade", nil)
		req.SetPathValue("id", "999")
		w := httptest.NewRecorder()

		handler.deleteDepartment(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
