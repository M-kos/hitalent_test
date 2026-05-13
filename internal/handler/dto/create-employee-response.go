package dto

import (
	"time"

	"github.com/M-kos/hitalent_test/internal/domain"
)

type CreateEmployeeResponse struct {
	ID           int    `json:"id"`
	DepartmentID *int   `json:"department_id,omitempty"`
	FullName     string `json:"full_name"`
	Position     string `json:"position"`
	HiredAt      string `json:"hired_at"`
	CreatedAt    string `json:"created_at"`
}

func (e *CreateEmployeeResponse) FromDomain(employee *domain.Employee) {
	e.ID = employee.ID
	e.DepartmentID = employee.DepartmentID
	e.FullName = employee.FullName
	e.Position = employee.Position

	if employee.HiredAt != nil {
		e.HiredAt = employee.HiredAt.Format(time.RFC3339)
	}

	e.CreatedAt = employee.CreatedAt.Format(time.RFC3339)
}
