package record

import (
	"time"

	"github.com/M-kos/hitalent_test/internal/domain"
)

type CreateDepartmentRecord struct {
	ID        int       `db:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Name      string    `db:"name" gorm:"column:name;not null"`
	ParentID  *int      `db:"parent_id" gorm:"column:parent_id"`
	CreatedAt time.Time `db:"created_at" gorm:"column:created_at"`
}

func (d *CreateDepartmentRecord) TableName() string {
	return "department"
}

func (d *CreateDepartmentRecord) ToDomain() (*domain.Department, error) {
	department := &domain.Department{
		ID:        d.ID,
		Name:      d.Name,
		CreatedAt: d.CreatedAt,
	}

	if d.ParentID != nil {
		department.Parent = &domain.Department{
			ID: *d.ParentID,
		}
	}

	return department, nil
}

func (d *CreateDepartmentRecord) FromDomain(department *domain.Department) {
	d.ID = department.ID
	if department.Name != "" {
		d.Name = department.Name
	}

	if department.Parent != nil {
		d.ParentID = &department.Parent.ID
	}
}

type UpdateDepartmentRecord struct {
	ID        int       `db:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Name      *string   `db:"name" gorm:"column:name"`
	ParentID  *int      `db:"parent_id" gorm:"column:parent_id"`
	CreatedAt time.Time `db:"created_at" gorm:"column:created_at"`
}

func (d *UpdateDepartmentRecord) TableName() string {
	return "department"
}

func (d *UpdateDepartmentRecord) ToDomain() (*domain.Department, error) {
	department := &domain.Department{
		ID:        d.ID,
		Name:      *d.Name,
		CreatedAt: d.CreatedAt,
	}

	if d.ParentID != nil {
		department.Parent = &domain.Department{
			ID: *d.ParentID,
		}
	}

	return department, nil
}

func (d *UpdateDepartmentRecord) FromDomain(department *domain.Department) {
	d.ID = department.ID
	if department.Name != "" {
		d.Name = &department.Name
	}

	if department.Parent != nil {
		d.ParentID = &department.Parent.ID
	}
}
