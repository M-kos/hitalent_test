package domain

import "time"

type Department struct {
	ID        int
	Name      string
	Parent    *Department
	CreatedAt time.Time
}
