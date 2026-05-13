package dto

import (
	"time"

	"github.com/M-kos/hitalent_test/internal/domain"
	"github.com/go-playground/validator/v10"
)

type CreateEmployeeRequest struct {
	FullName string  `json:"full_name" validate:"required,min=1,max=200"`
	Position string  `json:"position" validate:"required,min=1,max=200"`
	HiredAt  *string `json:"hired_at" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
}

func (e *CreateEmployeeRequest) Validate() error {
	validate := validator.New()

	return validate.Struct(e)
}

func (e *CreateEmployeeRequest) ToDomain(id int) (*domain.Employee, error) {
	employee := &domain.Employee{
		DepartmentID: &id,
		FullName:     e.FullName,
		Position:     e.Position,
		CreatedAt:    time.Now(),
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
