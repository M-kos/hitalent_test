package record

import (
	"time"

	"github.com/M-kos/hitalent_test/internal/domain"
)

type EmployeeRecord struct {
	ID           int     `db:"id" gorm:"column:id"`
	DepartmentID *int    `db:"department_id" gorm:"column:department_id"`
	FullName     string  `db:"full_name" gorm:"column:full_name"`
	Position     string  `db:"position" gorm:"column:position"`
	HiredAt      *string `db:"hired_at" gorm:"column:hired_at"`
	CreatedAt    string  `db:"created_at" gorm:"column:created_at"`
}

func (e *EmployeeRecord) ToDomain() (*domain.Employee, error) {
	createdAt, err := time.Parse(time.RFC3339, e.CreatedAt)
	if err != nil {
		return nil, err
	}

	employee := &domain.Employee{
		ID:           e.ID,
		DepartmentID: e.DepartmentID,
		FullName:     e.FullName,
		Position:     e.Position,
		CreatedAt:    createdAt,
	}

	if e.HiredAt != nil {
		hiredAt, err := time.Parse(time.RFC3339, *e.HiredAt)
		if err != nil {
			return nil, err
		}
		employee.HiredAt = &hiredAt
	}

	return employee, nil
}

func (e *EmployeeRecord) FromDomain(employee *domain.Employee) {
	e.ID = employee.ID
	e.FullName = employee.FullName
	e.Position = employee.Position
	e.CreatedAt = employee.CreatedAt.Format(time.RFC3339)

	if employee.HiredAt != nil {
		hiredAt := *employee.HiredAt
		hired := hiredAt.Format(time.RFC3339)
		e.HiredAt = &hired
	}
}
