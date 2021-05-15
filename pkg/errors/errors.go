package errors

import "errors"

// Used for config validation
var (
	ErrMissingUsername            = errors.New("TELEE_USERNAME must be set")
	ErrMissingPassword            = errors.New("TELEE_PASSWORD must be set")
	ErrMissingPrivPassword        = errors.New("TELEE_PRIVPASSWORD must be set")
	ErrMissingHostname            = errors.New("hostname must be set")
	ErrMissingCommand             = errors.New("command must be set")
	ErrInvalidPlatform            = errors.New("exec-platform is not supported")
	ErrUnsupportedRedundantMode   = errors.New("redundant-mode is not supported in this platform")
	ErrUnsupportedModeSet         = errors.New("enable-mode and default-priv-mode cannot use at once")
	ErrUnsupportedSecureMode      = errors.New("secure-mode is not supported in this platform")
	ErrUnsupportedDefaultPrivMode = errors.New("default-privilege-mode is not supported in this platform")
	ErrUnsupportedUnsecureMode    = errors.New("non secure-mode is not supported in this platform")
)

// Used for useca
var (
	ErrTermLengthIsEnforced = errors.New("EnableMode must be set. Terminal length expansion in user-level is not supporting.") // nolint
)
