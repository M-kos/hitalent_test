package repository

import (
	"context"
	"errors"

	"github.com/M-kos/hitalent_test/internal/domain"
	"github.com/M-kos/hitalent_test/internal/repository/query"
	"github.com/M-kos/hitalent_test/internal/repository/record"
	"gorm.io/gorm"
)

const MaxDepth = 999

type DepartmentRepository struct {
	db *gorm.DB
}

func NewDepartmentRepository(db *gorm.DB) *DepartmentRepository {
	return &DepartmentRepository{
		db: db,
	}
}

func (dr *DepartmentRepository) CreateDepartment(ctx context.Context, department *domain.Department) (*domain.Department, error) {
	var rec record.DepartmentRecord
	rec.FromDomain(department)

	err := dr.db.WithContext(ctx).Raw(query.CreateDepartment, rec.Name, rec.ParentID).Scan(&rec).Error
	if err != nil {
		return nil, err
	}

	dep, err := rec.ToDomain()
	if err != nil {
		return nil, err
	}

	return dep, nil
}

func (dr *DepartmentRepository) CreateEmployee(ctx context.Context, employee *domain.Employee) (*domain.Employee, error) {
	var rec record.EmployeeRecord
	rec.FromDomain(employee)

	err := dr.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var exists bool
		tx.Raw(query.CheckDepartment, employee.DepartmentID).Scan(&exists)
		if !exists {
			return domain.ErrDepartmentNotFound
		}

		err := tx.Raw(query.CreateEmployee, rec.DepartmentID, rec.FullName, rec.Position, rec.HiredAt).Scan(&rec).Error
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	emp, err := rec.ToDomain()
	if err != nil {
		return nil, err
	}

	return emp, nil
}

func (dr *DepartmentRepository) ListEmployeesByDepartmentId(ctx context.Context, ids []int) ([]*domain.Employee, error) {
	employeeRecords := make([]record.EmployeeRecord, 0)
	err := dr.db.WithContext(ctx).Raw(query.ListEmployeesByDepartmentId, ids).Scan(&employeeRecords).Error
	if err != nil {
		return nil, err
	}

	employees := make([]*domain.Employee, 0, len(employeeRecords))

	for _, rec := range employeeRecords {
		employee, err := rec.ToDomain()
		if err != nil {
			return nil, err
		}
		employees = append(employees, employee)
	}

	return employees, nil
}

func (dr *DepartmentRepository) DepartmentTree(ctx context.Context, id int, depth int) (*domain.DepartmentTree, error) {
	departmentRecords := make([]record.DepartmentTreeRecord, 0)
	err := dr.db.WithContext(ctx).Raw(query.DepartmentTree, id, depth).Scan(&departmentRecords).Error
	if err != nil {
		return nil, err
	}

	if len(departmentRecords) == 0 {
		return nil, domain.ErrDepartmentNotFound
	}

	recDepartment := departmentRecords[0]
	department, err := recDepartment.ToDomain()
	if err != nil {
		return nil, err
	}

	childrenDepartment := make([]*domain.Department, 0, len(departmentRecords)-1)

	for i := 1; i < len(departmentRecords); i++ {
		rec := departmentRecords[i]
		dep, err := rec.ToDomain()
		if err != nil {
			return nil, err
		}

		childrenDepartment = append(childrenDepartment, dep)
	}

	return &domain.DepartmentTree{
		Department:         department,
		ChildrenDepartment: childrenDepartment,
	}, nil
}

func (dr *DepartmentRepository) UpdateDepartment(ctx context.Context, department *domain.Department) (*domain.Department, error) {
	var rec record.DepartmentRecord
	rec.FromDomain(department)

	err := dr.db.WithContext(ctx).Raw(query.UpdateDepartment, rec.Name, rec.ParentID).Scan(&rec).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrDepartmentNotFound
		}
		return nil, err
	}

	dep, err := rec.ToDomain()
	if err != nil {
		return nil, err
	}

	return dep, nil
}

func (dr *DepartmentRepository) DeleteCascadeDepartment(ctx context.Context, departmentId int) error {
	return dr.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		departmentRecords := make([]record.DepartmentTreeRecord, 0)
		err := tx.Raw(query.DepartmentTree, departmentId, MaxDepth).Scan(&departmentRecords).Error
		if err != nil {
			return err
		}

		if len(departmentRecords) == 0 {
			return domain.ErrDepartmentNotFound
		}

		departmentIds := make([]int, 0, len(departmentRecords))
		for _, rec := range departmentRecords {
			departmentIds = append(departmentIds, rec.ID)
		}

		err = tx.Raw(query.DeleteDepartments, departmentIds).Error
		if err != nil {
			return err
		}

		err = tx.Raw(query.DeleteEmployeesByDepartmentId, departmentIds).Error
		if err != nil {
			return err
		}

		return nil
	})
}

func (dr *DepartmentRepository) DeleteAndReassignDepartment(ctx context.Context, departmentId int, reassignDepartmentId int) error {
	return dr.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var exists bool
		tx.Raw(query.CheckDepartment, departmentId).Scan(&exists)
		if !exists {
			return domain.ErrDepartmentNotFound
		}

		err := tx.Raw(query.DeleteDepartments, []int{departmentId}).Error
		if err != nil {
			return err
		}

		err = tx.Raw(query.UpdateDepartmentForEmployees, reassignDepartmentId, departmentId).Error
		if err != nil {
			return err
		}

		return nil
	})
}
