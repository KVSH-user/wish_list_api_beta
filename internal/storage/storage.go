package storage

import "errors"

var (
	ErrListNotFound = errors.New("list not found")
	ErrListExists   = errors.New("list exists")
)
