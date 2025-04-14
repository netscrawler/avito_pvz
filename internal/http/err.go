package httpserver

import "errors"

var (
	ErrJsonUnmarshal = errors.New("ErrorWithJsonDecode")
	ErrInvalidRole   = errors.New("InvalidRole")
	ErrInvalidEmail  = errors.New("InvalidEmail")
	ErrJsonMarshal   = errors.New("ErrorWithJsonEncode")
	ErrInvalidPVZID  = errors.New("InvalidPVZId")
)
