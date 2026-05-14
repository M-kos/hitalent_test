package domain

import "time"

type Employee struct {
	ID           int
	DepartmentID *int
	FullName     string
	Position     string
	HiredAt      *time.Time
	CreatedAt    time.Time
}
