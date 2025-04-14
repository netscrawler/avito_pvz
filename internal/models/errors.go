package models

import "errors"

var (
	ErrCheckUser       = errors.New("ErrCheckUserData")
	ErrInvalidPassword = errors.New("ErrInvalidPassword")
	ErrInvalidEmail    = errors.New("ErrInvalidPassword")
	ErrUserNotFoud     = errors.New("NotFoundUser")
)

var (
	ErrInvalidToken         = errors.New("ErrInvalidToken")
	ErrUnexpectedSignMethod = errors.New("ErrUnexpectedSignMethod")
	ErrInvalidTokenClaims   = errors.New("ErrInvalidTokenClaims")
	ErrInternalCodeGen      = errors.New("ErrInternalCodeGen")
)

var (
	ErrReceptionAlreadyClosed = errors.New("ReceptionAlreadyClosed")
	ErrReceptionDontExist     = errors.New("ReceptionDontExist")
	ErrReceptionAlreadyExist  = errors.New("ReceptionAlreadyExist")
	ErrProductNotFound        = errors.New("ProductNotFound")
	ErrPVZNotFound            = errors.New("PVZNotFound")
	ErrReceptionAlreadyExists = errors.New("ReceptionAlreadyClosed")
	ErrInvalidCity            = errors.New("InvalidCity")
	ErrInternal               = errors.New("InternalError")
	ErrInvalidProductType     = errors.New("InvalidProductType")
	ErrUserAlreadyExist       = errors.New("UserWithThisEmailAlreadyExist")
)
