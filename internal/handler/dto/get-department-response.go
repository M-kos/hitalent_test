package dto

import (
	"time"

	"github.com/M-kos/hitalent_test/internal/domain"
)

type GetDepartmentDto struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	ParentID  *int   `json:"parent_id,omitempty"`
	CreatedAt string `json:"created_at"`
}

type GetDepartmentEmployeeDto struct {
	ID           int    `json:"id"`
	DepartmentID *int   `json:"department_id,omitempty"`
	FullName     string `json:"full_name"`
	Position     string `json:"position"`
	HiredAt      string `json:"hired_at"`
	CreatedAt    string `json:"created_at"`
}

type GetDepartmentChildrenDto struct {
	ID        int                        `json:"id"`
	Name      string                     `json:"name"`
	ParentID  *int                       `json:"parent_id,omitempty"`
	CreatedAt string                     `json:"created_at"`
	Children  []GetDepartmentChildrenDto `json:"children"`
}

type DepartmentResponse struct {
	Department GetDepartmentDto           `json:"department"`
	Employees  []GetDepartmentEmployeeDto `json:"employees"`
	Children   []GetDepartmentChildrenDto `json:"children"`
}

func (d *DepartmentResponse) FromDomain(department *domain.Department, employees []*domain.Employee, childrenDepartment []*domain.Department) {
	d.Department.ID = department.ID
	d.Department.Name = department.Name

	if department.Parent != nil {
		d.Department.ParentID = &department.Parent.ID
	}

	d.Department.CreatedAt = department.CreatedAt.Format(time.RFC3339)

	d.Employees = make([]GetDepartmentEmployeeDto, len(employees))
	for i, employee := range employees {
		empl := GetDepartmentEmployeeDto{
			ID:           employee.ID,
			DepartmentID: employee.DepartmentID,
			FullName:     employee.FullName,
			Position:     employee.Position,
			HiredAt:      employee.HiredAt.Format(time.RFC3339),
			CreatedAt:    employee.CreatedAt.Format(time.RFC3339),
		}

		if employee.HiredAt != nil {
			empl.HiredAt = employee.HiredAt.Format(time.RFC3339)
		}

		d.Employees[i] = empl
	}

	d.Children = makeDepartmentChildren(d.Department.ID, childrenDepartment)
}

func makeDepartmentChildren(startParentId int, departments []*domain.Department) []GetDepartmentChildrenDto {
	if len(departments) == 0 {
		return nil
	}

	parrentIdToDepatments := map[int][]*domain.Department{
		startParentId: make([]*domain.Department, 0, len(departments)),
	}

	deeperChildren := make([]*domain.Department, 0, len(departments)/2)

	for _, department := range departments {
		if startParentId == department.Parent.ID {
			parrentIdToDepatments[startParentId] = append(parrentIdToDepatments[startParentId], department)
			continue
		}

		deeperChildren = append(deeperChildren, department)
	}

	result := make([]GetDepartmentChildrenDto, 0, len(parrentIdToDepatments[startParentId]))

	for _, department := range parrentIdToDepatments[startParentId] {
		dto := GetDepartmentChildrenDto{
			ID:        department.ID,
			Name:      department.Name,
			CreatedAt: department.CreatedAt.Format(time.RFC3339),
			Children:  makeDepartmentChildren(department.ID, deeperChildren),
		}

		if department.Parent != nil {
			dto.ParentID = &department.Parent.ID
		}
		result = append(result, dto)
	}

	return result
}
