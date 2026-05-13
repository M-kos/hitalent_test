package domain

import "errors"

var (
	ErrDepartmentNotFound = errors.New("department not found")
	ErrWrongParentId      = errors.New("wrong parent id")
	ErrWrongMode          = errors.New("wrong delete mode")
)
