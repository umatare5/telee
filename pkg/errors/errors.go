package errors

import "errors"

// Used for error handling
var (
	ErrMissingUsername     = errors.New("TELEE_USERNAME must be set")
	ErrMissingPassword     = errors.New("TELEE_PASSWORD must be set")
	ErrMissingPrivPassword = errors.New("TELEE_PRIVPASSWORD must be set")
	ErrMissingHostname     = errors.New("hostname must be set")
	ErrMissingCommand      = errors.New("command must be set")
	ErrInvalidPlatform     = errors.New("platform is not defined")
)
