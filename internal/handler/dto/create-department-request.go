package dto

import (
	"strings"

	"github.com/M-kos/hitalent_test/internal/domain"
	"github.com/go-playground/validator/v10"
)

type CreateDepartmentRequest struct {
	Name     string `json:"name" validate:"required,min=1,max=200"`
	ParentID *int   `json:"parent_id" validate:"omitempty"`
}

func (d *CreateDepartmentRequest) Validate() error {
	validate := validator.New()
	d.Name = strings.TrimSpace(d.Name)

	return validate.Struct(d)
}

func (d *CreateDepartmentRequest) ToDomain() *domain.Department {
	return &domain.Department{
		Name: strings.TrimSpace(d.Name),
		Parent: &domain.Department{
			ID: *d.ParentID,
		},
	}
}
