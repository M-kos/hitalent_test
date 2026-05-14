package dto

import (
	"strings"

	"github.com/M-kos/hitalent_test/internal/domain"
	"github.com/go-playground/validator/v10"
)

type UpdateDepartmentRequest struct {
	Name     string `json:"name,omitempty" validate:"omitempty,min=1,max=200"`
	ParentID *int   `json:"parent_id,omitempty" validate:"omitempty"`
}

func (d *UpdateDepartmentRequest) Validate() error {
	validate := validator.New()
	d.Name = strings.TrimSpace(d.Name)

	return validate.Struct(d)
}

func (d *UpdateDepartmentRequest) ToDomain(id int) *domain.Department {
	return &domain.Department{
		ID:   id,
		Name: strings.TrimSpace(d.Name),
		Parent: &domain.Department{
			ID: *d.ParentID,
		},
	}
}
