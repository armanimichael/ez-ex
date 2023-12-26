package security

import "errors"

var (
	ErrInvalidEncodingFormat = errors.New("invalid encoding format")
	ErrInvalidJWTToken       = errors.New("jwt token is not valid")
)
