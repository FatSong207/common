package id

import "errors"

var (
	ErrInvalidDatacenter = errors.New("invalid datacenter id: must be between 0 and 31")
	ErrInvalidMachine    = errors.New("invalid machine id: must be between 0 and 31")
)
