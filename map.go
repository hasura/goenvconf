package goenvconf

import (
	"maps"
	"os"
)

// EnvMapString represents either a literal string map or an environment reference.
type EnvMapString struct {
	Value    map[string]string `json:"value,omitempty" jsonschema:"anyof_required=value" mapstructure:"value" yaml:"value,omitempty"`
	Variable *string           `json:"env,omitempty"   jsonschema:"anyof_required=env"   mapstructure:"env"   yaml:"env,omitempty"`
}

// NewEnvMapString creates an EnvMapString instance.
func NewEnvMapString(env string, value map[string]string) EnvMapString {
	return EnvMapString{
		Variable: &env,
		Value:    value,
	}
}

// NewEnvMapStringValue creates an EnvMapString with a literal value.
func NewEnvMapStringValue(value map[string]string) EnvMapString {
	return EnvMapString{
		Value: value,
	}
}

// NewEnvMapStringVariable creates an EnvMapString with a variable name.
func NewEnvMapStringVariable(name string) EnvMapString {
	return EnvMapString{
		Variable: &name,
	}
}

// IsZero checks if the instance is empty.
func (ev EnvMapString) IsZero() bool {
	return (ev.Variable == nil || *ev.Variable == "") &&
		ev.Value == nil
}

// Equal checks if this instance equals the target value.
func (ev EnvMapString) Equal(target EnvMapString) bool {
	isSameEnv := (ev.Variable == nil && target.Variable == nil) ||
		(ev.Variable != nil && target.Variable != nil && *ev.Variable == *target.Variable)
	if !isSameEnv {
		return false
	}

	return (ev.Value == nil && target.Value == nil) ||
		(ev.Value != nil && target.Value != nil && maps.Equal(ev.Value, target.Value))
}

// Get gets literal value or from system environment.
func (ev EnvMapString) Get() (map[string]string, error) {
	if ev.Variable != nil && *ev.Variable != "" {
		rawValue := os.Getenv(*ev.Variable)
		if rawValue != "" {
			return ParseStringMapFromString(rawValue)
		}
	}

	return ev.Value, nil
}

// GetCustom gets literal value or from system environment by a custom function.
func (ev EnvMapString) GetCustom(getFunc GetEnvFunc) (map[string]string, error) {
	if ev.Variable != nil && *ev.Variable != "" {
		rawValue, err := getFunc(*ev.Variable)
		if err != nil {
			return nil, err
		}

		if rawValue != "" {
			return ParseStringMapFromString(rawValue)
		}
	}

	return ev.Value, nil
}

// EnvMapInt represents either a literal int map or an environment reference.
type EnvMapInt struct {
	Value    map[string]int64 `json:"value,omitempty" jsonschema:"anyof_required=value" mapstructure:"value" yaml:"value,omitempty"`
	Variable *string          `json:"env,omitempty"   jsonschema:"anyof_required=env"   mapstructure:"env"   yaml:"env,omitempty"`
}

// NewEnvMapInt creates an EnvMapInt instance.
func NewEnvMapInt(env string, value map[string]int64) EnvMapInt {
	return EnvMapInt{
		Variable: &env,
		Value:    value,
	}
}

// NewEnvMapIntValue creates an EnvMapInt with a literal value.
func NewEnvMapIntValue(value map[string]int64) EnvMapInt {
	return EnvMapInt{
		Value: value,
	}
}

// NewEnvMapIntVariable creates an EnvMapInt with a variable name.
func NewEnvMapIntVariable(name string) EnvMapInt {
	return EnvMapInt{
		Variable: &name,
	}
}

// IsZero checks if the instance is empty.
func (ev EnvMapInt) IsZero() bool {
	return (ev.Variable == nil || *ev.Variable == "") &&
		ev.Value == nil
}

// Equal checks if this instance equals the target value.
func (ev EnvMapInt) Equal(target EnvMapInt) bool {
	isSameEnv := (ev.Variable == nil && target.Variable == nil) ||
		(ev.Variable != nil && target.Variable != nil && *ev.Variable == *target.Variable)
	if !isSameEnv {
		return false
	}

	return (ev.Value == nil && target.Value == nil) ||
		(ev.Value != nil && target.Value != nil && maps.Equal(ev.Value, target.Value))
}

// Get gets literal value or from system environment.
func (ev EnvMapInt) Get() (map[string]int64, error) {
	if ev.Variable != nil && *ev.Variable != "" {
		rawValue := os.Getenv(*ev.Variable)
		if rawValue != "" {
			return ParseIntegerMapFromString[int64](rawValue)
		}
	}

	return ev.Value, nil
}

// GetCustom gets literal value or from system environment by a custom function.
func (ev EnvMapInt) GetCustom(getFunc GetEnvFunc) (map[string]int64, error) {
	if ev.Variable != nil && *ev.Variable != "" {
		rawValue, err := getFunc(*ev.Variable)
		if err != nil {
			return nil, err
		}

		if rawValue != "" {
			return ParseIntegerMapFromString[int64](rawValue)
		}
	}

	return ev.Value, nil
}

// EnvMapFloat represents either a literal float map or an environment reference.
type EnvMapFloat struct {
	Value    map[string]float64 `json:"value,omitempty" jsonschema:"anyof_required=value" mapstructure:"value" yaml:"value,omitempty"`
	Variable *string            `json:"env,omitempty"   jsonschema:"anyof_required=env"   mapstructure:"env"   yaml:"env,omitempty"`
}

// NewEnvMapFloat creates an EnvMapFloat instance.
func NewEnvMapFloat(env string, value map[string]float64) EnvMapFloat {
	return EnvMapFloat{
		Variable: &env,
		Value:    value,
	}
}

// NewEnvMapFloatValue creates an EnvMapFloat with a literal value.
func NewEnvMapFloatValue(value map[string]float64) EnvMapFloat {
	return EnvMapFloat{
		Value: value,
	}
}

// NewEnvMapFloatVariable creates an EnvMapFloat with a variable name.
func NewEnvMapFloatVariable(name string) EnvMapFloat {
	return EnvMapFloat{
		Variable: &name,
	}
}

// IsZero checks if the instance is empty.
func (ev EnvMapFloat) IsZero() bool {
	return (ev.Variable == nil || *ev.Variable == "") &&
		ev.Value == nil
}

// Equal checks if this instance equals the target value.
func (ev EnvMapFloat) Equal(target EnvMapFloat) bool {
	isSameEnv := (ev.Variable == nil && target.Variable == nil) ||
		(ev.Variable != nil && target.Variable != nil && *ev.Variable == *target.Variable)
	if !isSameEnv {
		return false
	}

	return (ev.Value == nil && target.Value == nil) ||
		(ev.Value != nil && target.Value != nil && maps.Equal(ev.Value, target.Value))
}

// Get gets literal value or from system environment.
func (ev EnvMapFloat) Get() (map[string]float64, error) {
	if ev.Variable != nil && *ev.Variable != "" {
		rawValue := os.Getenv(*ev.Variable)
		if rawValue != "" {
			return ParseFloatMapFromString[float64](rawValue)
		}
	}

	return ev.Value, nil
}

// GetCustom gets literal value or from system environment by a custom function.
func (ev EnvMapFloat) GetCustom(getFunc GetEnvFunc) (map[string]float64, error) {
	if ev.Variable != nil && *ev.Variable != "" {
		rawValue, err := getFunc(*ev.Variable)
		if err != nil {
			return nil, err
		}

		if rawValue != "" {
			return ParseFloatMapFromString[float64](rawValue)
		}
	}

	return ev.Value, nil
}

// EnvMapBool represents either a literal bool map or an environment reference.
type EnvMapBool struct {
	Value    map[string]bool `json:"value,omitempty" jsonschema:"anyof_required=value" mapstructure:"value" yaml:"value,omitempty"`
	Variable *string         `json:"env,omitempty"   jsonschema:"anyof_required=env"   mapstructure:"env"   yaml:"env,omitempty"`
}

// NewEnvMapBool creates an EnvMapBool instance.
func NewEnvMapBool(env string, value map[string]bool) EnvMapBool {
	return EnvMapBool{
		Variable: &env,
		Value:    value,
	}
}

// NewEnvMapBoolValue creates an EnvMapBool with a literal value.
func NewEnvMapBoolValue(value map[string]bool) EnvMapBool {
	return EnvMapBool{
		Value: value,
	}
}

// NewEnvMapBoolVariable creates an EnvMapBool with a variable name.
func NewEnvMapBoolVariable(name string) EnvMapBool {
	return EnvMapBool{
		Variable: &name,
	}
}

// IsZero checks if the instance is empty.
func (ev EnvMapBool) IsZero() bool {
	return (ev.Variable == nil || *ev.Variable == "") &&
		ev.Value == nil
}

// Equal checks if this instance equals the target value.
func (ev EnvMapBool) Equal(target EnvMapBool) bool {
	isSameEnv := (ev.Variable == nil && target.Variable == nil) ||
		(ev.Variable != nil && target.Variable != nil && *ev.Variable == *target.Variable)
	if !isSameEnv {
		return false
	}

	return (ev.Value == nil && target.Value == nil) ||
		(ev.Value != nil && target.Value != nil && maps.Equal(ev.Value, target.Value))
}

// Get gets literal value or from system environment.
func (ev EnvMapBool) Get() (map[string]bool, error) {
	if ev.Variable != nil && *ev.Variable != "" {
		rawValue := os.Getenv(*ev.Variable)
		if rawValue != "" {
			return ParseBoolMapFromString(rawValue)
		}
	}

	return ev.Value, nil
}

// GetCustom gets literal value or from system environment by a custom function.
func (ev EnvMapBool) GetCustom(getFunc GetEnvFunc) (map[string]bool, error) {
	if ev.Variable != nil && *ev.Variable != "" {
		rawValue, err := getFunc(*ev.Variable)
		if err != nil {
			return nil, err
		}

		if rawValue != "" {
			return ParseBoolMapFromString(rawValue)
		}
	}

	return ev.Value, nil
}
