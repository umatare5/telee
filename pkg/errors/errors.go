package errors

import "errors"

// Used for config validation
var (
	ErrMissingUsername         = errors.New("TELEE_USERNAME must be set")
	ErrMissingPassword         = errors.New("TELEE_PASSWORD must be set")
	ErrMissingPrivPassword     = errors.New("TELEE_PRIVPASSWORD must be set")
	ErrMissingHostname         = errors.New("hostname must be set")
	ErrMissingCommand          = errors.New("command must be set")
	ErrInvalidPlatform         = errors.New("exec-platform is not supported")
	ErrUnsupportedHAMode       = errors.New("ha-mode is not supported in this platform")
	ErrUnsupportedSecureMode   = errors.New("secure-mode is not supported in this platform")
	ErrUnsupportedUnsecureMode = errors.New("non secure-mode is not supported in this platform")
)

// Used for useca
var (
	ErrTermLengthIsEnforced = errors.New("EnableMode must be set. Terminal length expansion in user-level is not supporting.") // nolint
)
