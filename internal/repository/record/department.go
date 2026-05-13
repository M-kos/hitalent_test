package record

import (
	"time"

	"github.com/M-kos/hitalent_test/internal/domain"
)

type DepartmentRecord struct {
	ID        int    `db:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Name      string `db:"name" gorm:"column:name;not null"`
	ParentID  *int   `db:"parent_id" gorm:"column:parent_id"`
	CreatedAt string `db:"created_at" gorm:"column:created_at"`
}

func (d *DepartmentRecord) ToDomain() (*domain.Department, error) {
	createdAt, err := time.Parse(time.RFC3339, d.CreatedAt)
	if err != nil {
		return nil, err
	}

	department := &domain.Department{
		ID:        d.ID,
		Name:      d.Name,
		CreatedAt: createdAt,
	}

	if d.ParentID != nil {
		department.Parent = &domain.Department{
			ID: *d.ParentID,
		}
	}

	return department, nil
}

func (d *DepartmentRecord) FromDomain(department *domain.Department) {
	d.ID = department.ID
	d.Name = department.Name
	d.CreatedAt = department.CreatedAt.Format(time.RFC3339)

	if department.Parent != nil {
		d.ParentID = &department.Parent.ID
	}
}
