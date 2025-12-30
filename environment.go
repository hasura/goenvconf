// Package goenvconf contains reusable structures and utilities for configuration with environment variables.
package goenvconf

import (
	"errors"
	"os"
	"strconv"
)

// GetEnvFunc abstracts a custom function to get the value of an environment variable.
type GetEnvFunc func(string) (string, error)

// EnvString represents either a literal string or an environment reference.
type EnvString struct {
	Value    *string `json:"value,omitempty" jsonschema:"anyof_required=value,description=Default literal value if the env is empty" mapstructure:"value" yaml:"value,omitempty"`
	Variable *string `json:"env,omitempty"                                                                                           mapstructure:"env"   yaml:"env,omitempty"   hema:"anyof_required=env,description=Environment variable to be evaluated"`
}

// NewEnvString creates an EnvString instance.
func NewEnvString(env string, value string) EnvString {
	return EnvString{
		Variable: &env,
		Value:    &value,
	}
}

// NewEnvStringValue creates an EnvString with a literal value.
func NewEnvStringValue(value string) EnvString {
	return EnvString{
		Value: &value,
	}
}

// NewEnvStringVariable creates an EnvString with a variable name.
func NewEnvStringVariable(name string) EnvString {
	return EnvString{
		Variable: &name,
	}
}

// IsZero checks if the instance is empty.
func (ev EnvString) IsZero() bool {
	return (ev.Variable == nil || *ev.Variable == "") &&
		ev.Value == nil
}

// Equal checks if this instance equals the target value.
func (ev EnvString) Equal(target EnvString) bool {
	isSameValue := (ev.Value == nil && target.Value == nil) ||
		(ev.Value != nil && target.Value != nil && *ev.Value == *target.Value)
	if !isSameValue {
		return false
	}

	return (ev.Variable == nil && target.Variable == nil) ||
		(ev.Variable != nil && target.Variable != nil && *ev.Variable == *target.Variable)
}

// Get gets literal value or from system environment.
func (ev EnvString) Get() (string, error) {
	if ev.IsZero() {
		return "", ErrEnvironmentValueRequired
	}

	var value string

	var envExisted bool

	if ev.Variable != nil && *ev.Variable != "" {
		value, envExisted = os.LookupEnv(*ev.Variable)
		if value != "" {
			return value, nil
		}
	}

	if ev.Value != nil {
		return *ev.Value, nil
	}

	if envExisted {
		return "", nil
	}

	return "", getEnvVariableValueRequiredError(ev.Variable)
}

// GetOrDefault returns the default value if the environment value is empty.
func (ev EnvString) GetOrDefault(defaultValue string) (string, error) {
	result, err := ev.Get()
	if err != nil {
		if errors.Is(err, ErrEnvironmentVariableValueRequired) {
			return defaultValue, nil
		}

		return "", err
	} else if result == "" {
		result = defaultValue
	}

	return result, nil
}

// GetCustom gets literal value or from system environment by a custom function.
func (ev EnvString) GetCustom(getFunc GetEnvFunc) (string, error) {
	if ev.IsZero() {
		return "", ErrEnvironmentValueRequired
	}

	if ev.Variable != nil && *ev.Variable != "" {
		return getFunc(*ev.Variable)
	}

	if ev.Value != nil {
		return *ev.Value, nil
	}

	return "", getEnvVariableValueRequiredError(ev.Variable)
}

// EnvInt represents either a literal integer or an environment reference.
type EnvInt struct {
	Value    *int64  `json:"value,omitempty" jsonschema:"anyof_required=value,description=Default literal value if the env is empty" mapstructure:"value" yaml:"value,omitempty"`
	Variable *string `json:"env,omitempty"                                                                                           mapstructure:"env"   yaml:"env,omitempty"   hema:"anyof_required=env,description=Environment variable to be evaluated"`
}

// NewEnvInt creates an EnvInt instance.
func NewEnvInt(env string, value int64) EnvInt {
	return EnvInt{
		Variable: &env,
		Value:    &value,
	}
}

// NewEnvIntValue creates an EnvInt with a literal value.
func NewEnvIntValue(value int64) EnvInt {
	return EnvInt{
		Value: &value,
	}
}

// NewEnvIntVariable creates an EnvInt with a variable name.
func NewEnvIntVariable(name string) EnvInt {
	return EnvInt{
		Variable: &name,
	}
}

// IsZero checks if the instance is empty.
func (ev EnvInt) IsZero() bool {
	return (ev.Variable == nil || *ev.Variable == "") &&
		ev.Value == nil
}

// Equal checks if this instance equals the target value.
func (ev EnvInt) Equal(target EnvInt) bool {
	isSameValue := (ev.Value == nil && target.Value == nil) ||
		(ev.Value != nil && target.Value != nil && *ev.Value == *target.Value)
	if !isSameValue {
		return false
	}

	return (ev.Variable == nil && target.Variable == nil) ||
		(ev.Variable != nil && target.Variable != nil && *ev.Variable == *target.Variable)
}

// Get gets literal value or from system environment.
func (ev EnvInt) Get() (int64, error) {
	if ev.IsZero() {
		return 0, ErrEnvironmentValueRequired
	}

	if ev.Variable != nil && *ev.Variable != "" {
		rawValue := os.Getenv(*ev.Variable)
		if rawValue != "" {
			return strconv.ParseInt(rawValue, 10, 64)
		}
	}

	if ev.Value != nil {
		return *ev.Value, nil
	}

	return 0, getEnvVariableValueRequiredError(ev.Variable)
}

// GetOrDefault returns the default value if the environment value is empty.
func (ev EnvInt) GetOrDefault(defaultValue int64) (int64, error) {
	result, err := ev.Get()
	if err != nil {
		if errors.Is(err, ErrEnvironmentVariableValueRequired) {
			return defaultValue, nil
		}

		return 0, err
	}

	return result, nil
}

// GetCustom gets literal value or from system environment by a custom function.
func (ev EnvInt) GetCustom(getFunc GetEnvFunc) (int64, error) {
	if ev.IsZero() {
		return 0, ErrEnvironmentValueRequired
	}

	if ev.Variable != nil && *ev.Variable != "" {
		rawValue, err := getFunc(*ev.Variable)
		if err != nil {
			return 0, err
		}

		if rawValue != "" {
			return strconv.ParseInt(rawValue, 10, 64)
		}
	}

	if ev.Value != nil {
		return *ev.Value, nil
	}

	return 0, getEnvVariableValueRequiredError(ev.Variable)
}

// EnvBool represents either a literal boolean or an environment reference.
type EnvBool struct {
	Value    *bool   `json:"value,omitempty" jsonschema:"anyof_required=value,description=Default literal value if the env is empty" mapstructure:"value" yaml:"value,omitempty"`
	Variable *string `json:"env,omitempty"                                                                                           mapstructure:"env"   yaml:"env,omitempty"   hema:"anyof_required=env,description=Environment variable to be evaluated"`
}

// NewEnvBool creates an EnvBool instance.
func NewEnvBool(env string, value bool) EnvBool {
	return EnvBool{
		Variable: &env,
		Value:    &value,
	}
}

// NewEnvBoolValue creates an EnvBool with a literal value.
func NewEnvBoolValue(value bool) EnvBool {
	return EnvBool{
		Value: &value,
	}
}

// NewEnvBoolVariable creates an EnvBool with a variable name.
func NewEnvBoolVariable(name string) EnvBool {
	return EnvBool{
		Variable: &name,
	}
}

// IsZero checks if the instance is empty.
func (ev EnvBool) IsZero() bool {
	return (ev.Variable == nil || *ev.Variable == "") &&
		ev.Value == nil
}

// Equal checks if this instance equals the target value.
func (ev EnvBool) Equal(target EnvBool) bool {
	isSameValue := (ev.Value == nil && target.Value == nil) ||
		(ev.Value != nil && target.Value != nil && *ev.Value == *target.Value)
	if !isSameValue {
		return false
	}

	return (ev.Variable == nil && target.Variable == nil) ||
		(ev.Variable != nil && target.Variable != nil && *ev.Variable == *target.Variable)
}

// Get gets literal value or from system environment.
func (ev EnvBool) Get() (bool, error) {
	if ev.IsZero() {
		return false, ErrEnvironmentValueRequired
	}

	if ev.Variable != nil && *ev.Variable != "" {
		rawValue := os.Getenv(*ev.Variable)
		if rawValue != "" {
			return strconv.ParseBool(rawValue)
		}
	}

	if ev.Value != nil {
		return *ev.Value, nil
	}

	return false, getEnvVariableValueRequiredError(ev.Variable)
}

// GetOrDefault returns the default value if the environment value is empty.
func (ev EnvBool) GetOrDefault(defaultValue bool) (bool, error) {
	result, err := ev.Get()
	if err != nil {
		if errors.Is(err, ErrEnvironmentVariableValueRequired) {
			return defaultValue, nil
		}

		return false, err
	}

	return result, nil
}

// GetCustom gets literal value or from system environment with custom function.
func (ev EnvBool) GetCustom(getFunc GetEnvFunc) (bool, error) {
	if ev.IsZero() {
		return false, ErrEnvironmentValueRequired
	}

	if ev.Variable != nil && *ev.Variable != "" {
		rawValue, err := getFunc(*ev.Variable)
		if err != nil {
			return false, err
		}

		if rawValue != "" {
			return strconv.ParseBool(rawValue)
		}
	}

	if ev.Value != nil {
		return *ev.Value, nil
	}

	return false, getEnvVariableValueRequiredError(ev.Variable)
}

// EnvFloat represents either a literal floating point number or an environment reference.
type EnvFloat struct {
	Value    *float64 `json:"value,omitempty" jsonschema:"anyof_required=value,description=Default literal value if the env is empty" mapstructure:"value" yaml:"value,omitempty"`
	Variable *string  `json:"env,omitempty"                                                                                           mapstructure:"env"   yaml:"env,omitempty"   hema:"anyof_required=env,description=Environment variable to be evaluated"`
}

// NewEnvFloat creates an EnvFloat instance.
func NewEnvFloat(env string, value float64) EnvFloat {
	return EnvFloat{
		Variable: &env,
		Value:    &value,
	}
}

// NewEnvFloatValue creates an EnvFloat with a literal value.
func NewEnvFloatValue(value float64) EnvFloat {
	return EnvFloat{
		Value: &value,
	}
}

// NewEnvFloatVariable creates an EnvFloat with a variable name.
func NewEnvFloatVariable(name string) EnvFloat {
	return EnvFloat{
		Variable: &name,
	}
}

// IsZero checks if the instance is empty.
func (ev EnvFloat) IsZero() bool {
	return (ev.Variable == nil || *ev.Variable == "") &&
		ev.Value == nil
}

// Equal checks if this instance equals the target value.
func (ev EnvFloat) Equal(target EnvFloat) bool {
	isSameValue := (ev.Value == nil && target.Value == nil) ||
		(ev.Value != nil && target.Value != nil && *ev.Value == *target.Value)
	if !isSameValue {
		return false
	}

	return (ev.Variable == nil && target.Variable == nil) ||
		(ev.Variable != nil && target.Variable != nil && *ev.Variable == *target.Variable)
}

// Get gets literal value or from system environment.
func (ev EnvFloat) Get() (float64, error) {
	if ev.IsZero() {
		return 0, ErrEnvironmentValueRequired
	}

	if ev.Variable != nil && *ev.Variable != "" {
		rawValue := os.Getenv(*ev.Variable)
		if rawValue != "" {
			return strconv.ParseFloat(rawValue, 64)
		}
	}

	if ev.Value != nil {
		return *ev.Value, nil
	}

	return 0, getEnvVariableValueRequiredError(ev.Variable)
}

// GetOrDefault returns the default value if the environment value is empty.
func (ev EnvFloat) GetOrDefault(defaultValue float64) (float64, error) {
	result, err := ev.Get()
	if err != nil {
		if errors.Is(err, ErrEnvironmentVariableValueRequired) {
			return defaultValue, nil
		}

		return 0, err
	}

	return result, nil
}

// GetCustom gets literal value or from system environment by a custom function.
func (ev EnvFloat) GetCustom(getFunc GetEnvFunc) (float64, error) {
	if ev.IsZero() {
		return 0, ErrEnvironmentValueRequired
	}

	if ev.Variable != nil && *ev.Variable != "" {
		rawValue, err := getFunc(*ev.Variable)
		if err != nil {
			return 0, err
		}

		if rawValue != "" {
			return strconv.ParseFloat(rawValue, 64)
		}
	}

	if ev.Value != nil {
		return *ev.Value, nil
	}

	return 0, getEnvVariableValueRequiredError(ev.Variable)
}
