package service

import (
	"context"

	"github.com/M-kos/hitalent_test/internal/constants"
	"github.com/M-kos/hitalent_test/internal/domain"
)

type DepartmentRepository interface {
	CreateDepartment(ctx context.Context, department *domain.Department) (*domain.Department, error)
	CreateEmployee(ctx context.Context, employee *domain.Employee) (*domain.Employee, error)
	ListEmployeesByDepartmentId(ctx context.Context, ids []int) ([]*domain.Employee, error)
	DepartmentTree(ctx context.Context, id int, depth int) (*domain.DepartmentTree, error)
	UpdateDepartment(ctx context.Context, department *domain.Department) (*domain.Department, error)
	DeleteCascadeDepartment(ctx context.Context, departmentId int) error
	DeleteAndReassignDepartment(ctx context.Context, departmentId int, reassignDepartmentId int) error
}

type DepartmentService struct {
	repo DepartmentRepository
}

func NewDepartmentService(repo DepartmentRepository) *DepartmentService {
	return &DepartmentService{repo: repo}
}

func (d *DepartmentService) CreateDepartment(ctx context.Context, department *domain.Department) (*domain.Department, error) {
	return d.repo.CreateDepartment(ctx, department)
}

func (d *DepartmentService) CreateEmployee(ctx context.Context, employee *domain.Employee) (*domain.Employee, error) {
	return d.repo.CreateEmployee(ctx, employee)
}

func (d *DepartmentService) Department(ctx context.Context,
	id int,
	depth int,
	includeEmployees bool,
) (*domain.DepartmentTree, error) {
	tree, err := d.repo.DepartmentTree(ctx, id, depth)
	if err != nil {
		return nil, err
	}

	if includeEmployees {
		employeea, err := d.repo.ListEmployeesByDepartmentId(ctx, []int{id})
		if err != nil {
			return nil, err
		}
		tree.Employees = employeea
	}

	return tree, nil
}

func (d *DepartmentService) UpdateDepartment(ctx context.Context, department *domain.Department) (*domain.Department, error) {
	return d.repo.UpdateDepartment(ctx, department)
}

func (d *DepartmentService) DeleteDepartment(ctx context.Context, departmentId int, mode string, reassignId int) error {
	switch mode {
	case constants.ModeCascade:
		return d.repo.DeleteCascadeDepartment(ctx, departmentId)
	case constants.ModeReassign:
		return d.repo.DeleteAndReassignDepartment(ctx, departmentId, reassignId)
	default:
		return domain.ErrWrongMode
	}
}
