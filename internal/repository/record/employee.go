package record

import (
	"time"

	"github.com/M-kos/hitalent_test/internal/domain"
)

type EmployeeRecord struct {
	ID           int        `db:"id" gorm:"column:id"`
	DepartmentID *int       `db:"department_id" gorm:"column:department_id"`
	FullName     string     `db:"full_name" gorm:"column:full_name"`
	Position     string     `db:"position" gorm:"column:position"`
	HiredAt      *time.Time `db:"hired_at" gorm:"column:hired_at"`
	CreatedAt    time.Time  `db:"created_at" gorm:"column:created_at"`
}

func (e *EmployeeRecord) ToDomain() (*domain.Employee, error) {
	employee := &domain.Employee{
		ID:           e.ID,
		DepartmentID: e.DepartmentID,
		FullName:     e.FullName,
		Position:     e.Position,
		CreatedAt:    e.CreatedAt,
	}

	if e.HiredAt != nil {
		employee.HiredAt = e.HiredAt
	}

	return employee, nil
}

func (e *EmployeeRecord) FromDomain(employee *domain.Employee) {
	e.ID = employee.ID
	e.FullName = employee.FullName
	e.Position = employee.Position
	e.CreatedAt = employee.CreatedAt

	if employee.DepartmentID != nil {
		e.DepartmentID = employee.DepartmentID
	}

	if employee.HiredAt != nil {
		e.HiredAt = employee.HiredAt
	}
}
