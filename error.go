package goenvconf

import "fmt"

var (
	// ErrEnvironmentValueRequired occurs when both value and env fields are null or empty.
	ErrEnvironmentValueRequired = ParseEnvError{
		Code:   "EmptyEnv",
		Detail: "require either value or env",
	}

	// ErrEnvironmentVariableValueRequired the error occurs when the value from environment variable is empty.
	ErrEnvironmentVariableValueRequired = ParseEnvError{
		Code:   "EmptyVar",
		Detail: "the environment variable value is empty",
	}
)

const (
	// ErrCodeParseEnvFailed is the error code when parsing environment variable failed.
	ErrCodeParseEnvFailed = "ParseEnvFailed"
)

// ParseEnvError structures a detailed error for parsed env.
type ParseEnvError struct {
	Code   string `json:"code"           jsonschema:"enum=EmptyEnv,enum=EmptyVar,enum=ParseEnvFailed"`
	Detail string `json:"detail"`
	Hint   string `json:"hint,omitempty"`
}

// NewParseEnvFailedError creates a [ParseEnvError] for parsing env variable errors.
func NewParseEnvFailedError(detail string, hint string) ParseEnvError {
	return ParseEnvError{
		Code:   ErrCodeParseEnvFailed,
		Detail: detail,
		Hint:   hint,
	}
}

// Error returns the error message.
func (pee ParseEnvError) Error() string {
	if pee.Hint != "" {
		return fmt.Sprintf("%s: %s. Hint: %s", pee.Code, pee.Detail, pee.Hint)
	}

	return pee.Code + ": " + pee.Detail
}
