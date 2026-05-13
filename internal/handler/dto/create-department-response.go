package dto

import (
	"time"

	"github.com/M-kos/hitalent_test/internal/domain"
)

type CreateDepartmentResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	ParentID  int    `json:"parent_id"`
	CreatedAt string `json:"created_at"`
}

func (d *CreateDepartmentResponse) FromDomain(department *domain.Department) {
	d.ID = department.ID
	d.Name = department.Name
	if department.Parent != nil {
		d.ParentID = department.Parent.ID
	}
	d.CreatedAt = department.CreatedAt.Format(time.RFC3339)
}
