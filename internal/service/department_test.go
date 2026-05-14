package service

import (
	"context"
	"testing"

	"github.com/M-kos/hitalent_test/internal/constants"
	"github.com/M-kos/hitalent_test/internal/domain"
	"github.com/stretchr/testify/assert"
)

type StubDepartmentRepository struct {
	CreateDepartmentFunc            func(ctx context.Context, department *domain.Department) (*domain.Department, error)
	CreateEmployeeFunc              func(ctx context.Context, employee *domain.Employee) (*domain.Employee, error)
	ListEmployeesByDepartmentIdFunc func(ctx context.Context, ids []int) ([]*domain.Employee, error)
	DepartmentTreeFunc              func(ctx context.Context, id int, depth int) (*domain.DepartmentTree, error)
	UpdateDepartmentFunc            func(ctx context.Context, department *domain.Department) (*domain.Department, error)
	DeleteCascadeDepartmentFunc     func(ctx context.Context, departmentId int) error
	DeleteAndReassignDepartmentFunc func(ctx context.Context, departmentId int, reassignDepartmentId int) error
}

func (s *StubDepartmentRepository) CreateDepartment(ctx context.Context, department *domain.Department) (*domain.Department, error) {
	if s.CreateDepartmentFunc != nil {
		return s.CreateDepartmentFunc(ctx, department)
	}
	return nil, nil
}

func (s *StubDepartmentRepository) CreateEmployee(ctx context.Context, employee *domain.Employee) (*domain.Employee, error) {
	if s.CreateEmployeeFunc != nil {
		return s.CreateEmployeeFunc(ctx, employee)
	}
	return nil, nil
}

func (s *StubDepartmentRepository) ListEmployeesByDepartmentId(ctx context.Context, ids []int) ([]*domain.Employee, error) {
	if s.ListEmployeesByDepartmentIdFunc != nil {
		return s.ListEmployeesByDepartmentIdFunc(ctx, ids)
	}
	return nil, nil
}

func (s *StubDepartmentRepository) DepartmentTree(ctx context.Context, id int, depth int) (*domain.DepartmentTree, error) {
	if s.DepartmentTreeFunc != nil {
		return s.DepartmentTreeFunc(ctx, id, depth)
	}
	return nil, nil
}

func (s *StubDepartmentRepository) UpdateDepartment(ctx context.Context, department *domain.Department) (*domain.Department, error) {
	if s.UpdateDepartmentFunc != nil {
		return s.UpdateDepartmentFunc(ctx, department)
	}
	return nil, nil
}

func (s *StubDepartmentRepository) DeleteCascadeDepartment(ctx context.Context, departmentId int) error {
	if s.DeleteCascadeDepartmentFunc != nil {
		return s.DeleteCascadeDepartmentFunc(ctx, departmentId)
	}
	return nil
}

func (s *StubDepartmentRepository) DeleteAndReassignDepartment(ctx context.Context, departmentId int, reassignDepartmentId int) error {
	if s.DeleteAndReassignDepartmentFunc != nil {
		return s.DeleteAndReassignDepartmentFunc(ctx, departmentId, reassignDepartmentId)
	}
	return nil
}

func TestDepartmentService_CreateDepartment(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	dept := &domain.Department{Name: "Backend"}

	tests := []struct {
		name           string
		setup          func() DepartmentRepository
		expectedError  error
		expectedResult *domain.Department
	}{
		{
			name: "success",
			setup: func() DepartmentRepository {
				return &StubDepartmentRepository{
					CreateDepartmentFunc: func(ctx context.Context, department *domain.Department) (*domain.Department, error) {
						return dept, nil
					},
				}
			},
			expectedError:  nil,
			expectedResult: dept,
		},
		{
			name: "error from repo",
			setup: func() DepartmentRepository {
				return &StubDepartmentRepository{
					CreateDepartmentFunc: func(ctx context.Context, department *domain.Department) (*domain.Department, error) {
						return nil, domain.ErrDepartmentAlreadyExists
					},
				}
			},
			expectedError:  domain.ErrDepartmentAlreadyExists,
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := tt.setup()
			service := NewDepartmentService(repo)

			result, err := service.CreateDepartment(ctx, dept)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
		})
	}
}

func TestDepartmentService_CreateEmployee(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	empl := &domain.Employee{FullName: "Aaa Bbb", Position: "Developer"}

	tests := []struct {
		name           string
		setup          func() DepartmentRepository
		expectedError  error
		expectedResult *domain.Employee
	}{
		{
			name: "success",
			setup: func() DepartmentRepository {
				return &StubDepartmentRepository{
					CreateEmployeeFunc: func(ctx context.Context, employee *domain.Employee) (*domain.Employee, error) {
						return empl, nil
					},
				}
			},
			expectedError:  nil,
			expectedResult: empl,
		},
		{
			name: "error from repo",
			setup: func() DepartmentRepository {
				return &StubDepartmentRepository{
					CreateEmployeeFunc: func(ctx context.Context, employee *domain.Employee) (*domain.Employee, error) {
						return nil, domain.ErrDepartmentNotFound
					},
				}
			},
			expectedError:  domain.ErrDepartmentNotFound,
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := tt.setup()
			service := NewDepartmentService(repo)

			result, err := service.CreateEmployee(ctx, empl)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
		})
	}
}

func TestDepartmentService_Department(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	tree := &domain.DepartmentTree{
		Department:         &domain.Department{Name: "Developers"},
		ChildrenDepartment: []*domain.Department{{Name: "Backend"}},
	}

	tests := []struct {
		name           string
		setup          func() DepartmentRepository
		id             int
		depth          int
		include        bool
		expectedError  error
		expectedResult *domain.DepartmentTree
	}{
		{
			name: "success with employees",
			setup: func() DepartmentRepository {
				return &StubDepartmentRepository{
					DepartmentTreeFunc: func(ctx context.Context, id int, depth int) (*domain.DepartmentTree, error) {
						return tree, nil
					},
					ListEmployeesByDepartmentIdFunc: func(ctx context.Context, ids []int) ([]*domain.Employee, error) {
						return []*domain.Employee{{FullName: "Aaa Bbb"}}, nil
					},
				}
			},
			id:             1,
			depth:          1,
			include:        true,
			expectedError:  nil,
			expectedResult: tree,
		},
		{
			name: "success without employees",
			setup: func() DepartmentRepository {
				return &StubDepartmentRepository{
					DepartmentTreeFunc: func(ctx context.Context, id int, depth int) (*domain.DepartmentTree, error) {
						return tree, nil
					},
				}
			},
			id:             1,
			depth:          1,
			include:        false,
			expectedError:  nil,
			expectedResult: tree,
		},
		{
			name: "error from repo",
			setup: func() DepartmentRepository {
				return &StubDepartmentRepository{
					DepartmentTreeFunc: func(ctx context.Context, id int, depth int) (*domain.DepartmentTree, error) {
						return nil, domain.ErrDepartmentNotFound
					},
				}
			},
			id:            1,
			depth:         1,
			include:       true,
			expectedError: domain.ErrDepartmentNotFound,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := tt.setup()
			service := NewDepartmentService(repo)

			result, err := service.Department(ctx, tt.id, tt.depth, tt.include)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
				if tt.include {
					assert.NotEmpty(t, result.Employees)
				} else {
					assert.Nil(t, result.Employees)
				}
			}
		})
	}
}

func TestDepartmentService_UpdateDepartment(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	dept := &domain.Department{Name: "Backend"}

	tests := []struct {
		name           string
		setup          func() DepartmentRepository
		expectedError  error
		expectedResult *domain.Department
	}{
		{
			name: "success",
			setup: func() DepartmentRepository {
				return &StubDepartmentRepository{
					UpdateDepartmentFunc: func(ctx context.Context, department *domain.Department) (*domain.Department, error) {
						return dept, nil
					},
				}
			},
			expectedError:  nil,
			expectedResult: dept,
		},
		{
			name: "error from repo",
			setup: func() DepartmentRepository {
				return &StubDepartmentRepository{
					UpdateDepartmentFunc: func(ctx context.Context, department *domain.Department) (*domain.Department, error) {
						return nil, domain.ErrWrongParentId
					},
				}
			},
			expectedError:  domain.ErrWrongParentId,
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := tt.setup()
			service := NewDepartmentService(repo)

			result, err := service.UpdateDepartment(ctx, dept)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
		})
	}
}

func TestDepartmentService_DeleteDepartment(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name          string
		setup         func() DepartmentRepository
		departmentId  int
		mode          string
		reassignId    int
		expectedError error
	}{
		{
			name: "cascade mode success",
			setup: func() DepartmentRepository {
				return &StubDepartmentRepository{
					DeleteCascadeDepartmentFunc: func(ctx context.Context, departmentId int) error {
						return nil
					},
				}
			},
			departmentId:  1,
			mode:          constants.ModeCascade,
			reassignId:    0,
			expectedError: nil,
		},
		{
			name: "reassign mode success",
			setup: func() DepartmentRepository {
				return &StubDepartmentRepository{
					DeleteAndReassignDepartmentFunc: func(ctx context.Context, departmentId int, reassignDepartmentId int) error {
						return nil
					},
				}
			},
			departmentId:  1,
			mode:          constants.ModeReassign,
			reassignId:    2,
			expectedError: nil,
		},
		{
			name:          "wrong mode",
			setup:         func() DepartmentRepository { return &StubDepartmentRepository{} },
			departmentId:  1,
			mode:          "invalid",
			reassignId:    0,
			expectedError: domain.ErrWrongMode,
		},
		{
			name: "error from repo in cascade",
			setup: func() DepartmentRepository {
				return &StubDepartmentRepository{
					DeleteCascadeDepartmentFunc: func(ctx context.Context, departmentId int) error {
						return domain.ErrDepartmentNotFound
					},
				}
			},
			departmentId:  1,
			mode:          constants.ModeCascade,
			reassignId:    0,
			expectedError: domain.ErrDepartmentNotFound,
		},
		{
			name: "error from repo in reassign",
			setup: func() DepartmentRepository {
				return &StubDepartmentRepository{
					DeleteAndReassignDepartmentFunc: func(ctx context.Context, departmentId int, reassignDepartmentId int) error {
						return domain.ErrDepartmentNotFound
					},
				}
			},
			departmentId:  1,
			mode:          constants.ModeReassign,
			reassignId:    2,
			expectedError: domain.ErrDepartmentNotFound,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := tt.setup()
			service := NewDepartmentService(repo)

			err := service.DeleteDepartment(ctx, tt.departmentId, tt.mode, tt.reassignId)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
