package domain

import (
	"errors"
)

var (
	ErrDepartmentNotFound      = errors.New("department not found")
	ErrWrongParentId           = errors.New("wrong parent id")
	ErrWrongMode               = errors.New("wrong delete mode")
	ErrDepartmentAlreadyExists = errors.New("department with the same name already exists")
)
