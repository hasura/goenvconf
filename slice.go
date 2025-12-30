package goenvconf

import (
	"fmt"
	"os"
	"slices"
)

// EnvStringSlice represents either a literal string slice or an environment reference.
type EnvStringSlice struct {
	Value    []string `json:"value,omitempty" jsonschema:"anyof_required=value,description=Default literal value if the env is empty" mapstructure:"value" yaml:"value,omitempty"`
	Variable *string  `json:"env,omitempty"                                                                                           mapstructure:"env"   yaml:"env,omitempty"   hema:"anyof_required=env,description=Environment variable to be evaluated"`
}

// NewEnvStringSlice creates an EnvStringSlice instance.
func NewEnvStringSlice(env string, value []string) EnvStringSlice {
	return EnvStringSlice{
		Variable: &env,
		Value:    value,
	}
}

// NewEnvStringSliceValue creates an EnvStringSlice with a literal value.
func NewEnvStringSliceValue(value []string) EnvStringSlice {
	return EnvStringSlice{
		Value: value,
	}
}

// NewEnvStringSliceVariable creates an EnvStringSlice with a variable name.
func NewEnvStringSliceVariable(name string) EnvStringSlice {
	return EnvStringSlice{
		Variable: &name,
	}
}

// IsZero checks if the instance is empty.
func (ev EnvStringSlice) IsZero() bool {
	return (ev.Variable == nil || *ev.Variable == "") &&
		ev.Value == nil
}

// Equal checks if this instance equals the target value.
func (ev EnvStringSlice) Equal(target EnvStringSlice) bool {
	isSameValue := slices.Equal(ev.Value, target.Value)
	if !isSameValue {
		return false
	}

	return (ev.Variable == nil && target.Variable == nil) ||
		(ev.Variable != nil && target.Variable != nil && *ev.Variable == *target.Variable)
}

// Get gets literal value or from system environment.
func (ev EnvStringSlice) Get() ([]string, error) {
	if ev.IsZero() {
		return nil, ErrEnvironmentValueRequired
	}

	var value string

	var envExisted bool

	if ev.Variable != nil && *ev.Variable != "" {
		value, envExisted = os.LookupEnv(*ev.Variable)
		if value != "" {
			return ParseStringSliceFromString(value), nil
		}
	}

	if ev.Value != nil {
		return ev.Value, nil
	}

	if envExisted {
		return []string{}, nil
	}

	return nil, getEnvVariableValueRequiredError(ev.Variable)
}

// GetCustom gets literal value or from system environment by a custom function.
func (ev EnvStringSlice) GetCustom(getFunc GetEnvFunc) ([]string, error) {
	if ev.IsZero() {
		return nil, ErrEnvironmentValueRequired
	}

	if ev.Variable != nil && *ev.Variable != "" {
		value, err := getFunc(*ev.Variable)
		if err != nil {
			return nil, err
		}

		if value != "" {
			return ParseStringSliceFromString(value), nil
		}
	}

	if ev.Value != nil {
		return ev.Value, nil
	}

	return nil, getEnvVariableValueRequiredError(ev.Variable)
}

// EnvIntSlice represents either a literal integer slice or an environment reference.
type EnvIntSlice struct {
	Value    []int64 `json:"value,omitempty" jsonschema:"anyof_required=value,description=Default literal value if the env is empty" mapstructure:"value" yaml:"value,omitempty"`
	Variable *string `json:"env,omitempty"                                                                                           mapstructure:"env"   yaml:"env,omitempty"   hema:"anyof_required=env,description=Environment variable to be evaluated"`
}

// NewEnvIntSlice creates an EnvIntSlice instance.
func NewEnvIntSlice(env string, value []int64) EnvIntSlice {
	return EnvIntSlice{
		Variable: &env,
		Value:    value,
	}
}

// NewEnvIntSliceValue creates an EnvIntSlice with a literal value.
func NewEnvIntSliceValue(value []int64) EnvIntSlice {
	return EnvIntSlice{
		Value: value,
	}
}

// NewEnvIntSliceVariable creates an EnvIntSlice with a variable name.
func NewEnvIntSliceVariable(name string) EnvIntSlice {
	return EnvIntSlice{
		Variable: &name,
	}
}

// IsZero checks if the instance is empty.
func (ev EnvIntSlice) IsZero() bool {
	return (ev.Variable == nil || *ev.Variable == "") &&
		ev.Value == nil
}

// Equal checks if this instance equals the target value.
func (ev EnvIntSlice) Equal(target EnvIntSlice) bool {
	isSameValue := slices.Equal(ev.Value, target.Value)
	if !isSameValue {
		return false
	}

	return (ev.Variable == nil && target.Variable == nil) ||
		(ev.Variable != nil && target.Variable != nil && *ev.Variable == *target.Variable)
}

// Get gets literal value or from system environment.
func (ev EnvIntSlice) Get() ([]int64, error) {
	if ev.IsZero() {
		return nil, ErrEnvironmentValueRequired
	}

	var value string

	var envExisted bool

	if ev.Variable != nil && *ev.Variable != "" {
		value, envExisted = os.LookupEnv(*ev.Variable)
		if value != "" {
			return parseIntSliceFromStringWithErrorPrefix[int64](
				value,
				fmt.Sprintf("failed to parse %s: ", *ev.Variable),
			)
		}
	}

	if ev.Value != nil {
		return ev.Value, nil
	}

	if envExisted {
		return []int64{}, nil
	}

	return nil, getEnvVariableValueRequiredError(ev.Variable)
}

// GetCustom gets literal value or from system environment by a custom function.
func (ev EnvIntSlice) GetCustom(getFunc GetEnvFunc) ([]int64, error) {
	if ev.IsZero() {
		return nil, ErrEnvironmentValueRequired
	}

	if ev.Variable != nil && *ev.Variable != "" {
		value, err := getFunc(*ev.Variable)
		if err != nil {
			return nil, err
		}

		if value != "" {
			return parseIntSliceFromStringWithErrorPrefix[int64](
				value,
				fmt.Sprintf("failed to parse %s: ", *ev.Variable),
			)
		}
	}

	if ev.Value != nil {
		return ev.Value, nil
	}

	return nil, getEnvVariableValueRequiredError(ev.Variable)
}

// EnvFloatSlice represents either a literal floating-point number slice or an environment reference.
type EnvFloatSlice struct {
	Value    []float64 `json:"value,omitempty" jsonschema:"anyof_required=value,description=Default literal value if the env is empty" mapstructure:"value" yaml:"value,omitempty"`
	Variable *string   `json:"env,omitempty"                                                                                           mapstructure:"env"   yaml:"env,omitempty"   hema:"anyof_required=env,description=Environment variable to be evaluated"`
}

// NewEnvFloatSlice creates an EnvFloatSlice instance.
func NewEnvFloatSlice(env string, value []float64) EnvFloatSlice {
	return EnvFloatSlice{
		Variable: &env,
		Value:    value,
	}
}

// NewEnvFloatSliceValue creates an EnvFloatSlice with a literal value.
func NewEnvFloatSliceValue(value []float64) EnvFloatSlice {
	return EnvFloatSlice{
		Value: value,
	}
}

// NewEnvFloatSliceVariable creates an EnvFloatSlice with a variable name.
func NewEnvFloatSliceVariable(name string) EnvFloatSlice {
	return EnvFloatSlice{
		Variable: &name,
	}
}

// IsZero checks if the instance is empty.
func (ev EnvFloatSlice) IsZero() bool {
	return (ev.Variable == nil || *ev.Variable == "") &&
		ev.Value == nil
}

// Equal checks if this instance equals the target value.
func (ev EnvFloatSlice) Equal(target EnvFloatSlice) bool {
	isSameValue := slices.Equal(ev.Value, target.Value)
	if !isSameValue {
		return false
	}

	return (ev.Variable == nil && target.Variable == nil) ||
		(ev.Variable != nil && target.Variable != nil && *ev.Variable == *target.Variable)
}

// Get gets literal value or from system environment.
func (ev EnvFloatSlice) Get() ([]float64, error) {
	if ev.IsZero() {
		return nil, ErrEnvironmentValueRequired
	}

	var value string

	var envExisted bool

	if ev.Variable != nil && *ev.Variable != "" {
		value, envExisted = os.LookupEnv(*ev.Variable)
		if value != "" {
			return parseFloatSliceFromStringWithErrorPrefix[float64](
				value,
				fmt.Sprintf("failed to parse %s: ", *ev.Variable),
			)
		}
	}

	if ev.Value != nil {
		return ev.Value, nil
	}

	if envExisted {
		return []float64{}, nil
	}

	return nil, getEnvVariableValueRequiredError(ev.Variable)
}

// GetCustom gets literal value or from system environment by a custom function.
func (ev EnvFloatSlice) GetCustom(getFunc GetEnvFunc) ([]float64, error) {
	if ev.IsZero() {
		return nil, ErrEnvironmentValueRequired
	}

	if ev.Variable != nil && *ev.Variable != "" {
		value, err := getFunc(*ev.Variable)
		if err != nil {
			return nil, err
		}

		if value != "" {
			return parseFloatSliceFromStringWithErrorPrefix[float64](
				value,
				fmt.Sprintf("failed to parse %s: ", *ev.Variable),
			)
		}
	}

	if ev.Value != nil {
		return ev.Value, nil
	}

	return nil, getEnvVariableValueRequiredError(ev.Variable)
}

// EnvBoolSlice represents either a literal boolean slice or an environment reference.
type EnvBoolSlice struct {
	Value    []bool  `json:"value,omitempty" jsonschema:"anyof_required=value,description=Default literal value if the env is empty" mapstructure:"value" yaml:"value,omitempty"`
	Variable *string `json:"env,omitempty"                                                                                           mapstructure:"env"   yaml:"env,omitempty"   hema:"anyof_required=env,description=Environment variable to be evaluated"`
}

// NewEnvBoolSlice creates an EnvBoolSlice instance.
func NewEnvBoolSlice(env string, value []bool) EnvBoolSlice {
	return EnvBoolSlice{
		Variable: &env,
		Value:    value,
	}
}

// NewEnvBoolSliceValue creates an EnvBoolSlice with a literal value.
func NewEnvBoolSliceValue(value []bool) EnvBoolSlice {
	return EnvBoolSlice{
		Value: value,
	}
}

// NewEnvBoolSliceVariable creates an EnvBoolSlice with a variable name.
func NewEnvBoolSliceVariable(name string) EnvBoolSlice {
	return EnvBoolSlice{
		Variable: &name,
	}
}

// IsZero checks if the instance is empty.
func (ev EnvBoolSlice) IsZero() bool {
	return (ev.Variable == nil || *ev.Variable == "") &&
		ev.Value == nil
}

// Equal checks if this instance equals the target value.
func (ev EnvBoolSlice) Equal(target EnvBoolSlice) bool {
	isSameValue := slices.Equal(ev.Value, target.Value)
	if !isSameValue {
		return false
	}

	return (ev.Variable == nil && target.Variable == nil) ||
		(ev.Variable != nil && target.Variable != nil && *ev.Variable == *target.Variable)
}

// Get gets literal value or from system environment.
func (ev EnvBoolSlice) Get() ([]bool, error) {
	if ev.IsZero() {
		return nil, ErrEnvironmentValueRequired
	}

	var value string

	var envExisted bool

	if ev.Variable != nil && *ev.Variable != "" {
		value, envExisted = os.LookupEnv(*ev.Variable)
		if value != "" {
			return parseBoolSliceFromStringWithErrorPrefix(
				value,
				fmt.Sprintf("failed to parse %s: ", *ev.Variable),
			)
		}
	}

	if ev.Value != nil {
		return ev.Value, nil
	}

	if envExisted {
		return []bool{}, nil
	}

	return nil, getEnvVariableValueRequiredError(ev.Variable)
}

// GetCustom gets literal value or from system environment by a custom function.
func (ev EnvBoolSlice) GetCustom(getFunc GetEnvFunc) ([]bool, error) {
	if ev.IsZero() {
		return nil, ErrEnvironmentValueRequired
	}

	if ev.Variable != nil && *ev.Variable != "" {
		value, err := getFunc(*ev.Variable)
		if err != nil {
			return nil, err
		}

		if value != "" {
			return parseBoolSliceFromStringWithErrorPrefix(
				value,
				fmt.Sprintf("failed to parse %s: ", *ev.Variable),
			)
		}
	}

	if ev.Value != nil {
		return ev.Value, nil
	}

	return nil, getEnvVariableValueRequiredError(ev.Variable)
}
