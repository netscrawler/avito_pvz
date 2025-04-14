package domain

import "errors"

var (
	ErrNotFound      = errors.New("NotFound")
	ErrPVZNotExist   = errors.New("PVZNotFound")
	ErrInvalidEmail  = errors.New("InvalidEmail")
	ErrInternal      = errors.New("InternalError")
	ErrInvalidRole   = errors.New("InvalidRole")
	ErrAlreadyExists = errors.New("UserAlreadyExist")
)
