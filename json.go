package goenvconf

import (
	"encoding/json"
	"errors"
	"os"
	"reflect"
)

// EnvJSON represents either a JSON value or an environment reference.
type EnvJSON struct {
	Value    any     `json:"value,omitempty" jsonschema:"anyof_required=value,description=Default literal value if the env is empty" mapstructure:"value" yaml:"value,omitempty"`
	Variable *string `json:"env,omitempty"   jsonschema:"anyof_required=env,description=Environment variable to be evaluated"        mapstructure:"env"   yaml:"env,omitempty"`
}

// NewEnvJSON creates an EnvJSON instance.
func NewEnvJSON(env string, value any) EnvJSON {
	return EnvJSON{
		Variable: &env,
		Value:    value,
	}
}

// NewEnvJSONValue creates an EnvJSON with a literal value.
func NewEnvJSONValue(value any) EnvJSON {
	return EnvJSON{
		Value: value,
	}
}

// NewEnvJSONVariable creates an EnvJSON with a variable name.
func NewEnvJSONVariable(name string) EnvJSON {
	return EnvJSON{
		Variable: &name,
	}
}

// IsZero checks if the instance is empty.
func (ev EnvJSON) IsZero() bool {
	return (ev.Variable == nil || *ev.Variable == "") &&
		ev.Value == nil
}

// Get gets literal value or from system environment.
func (ev EnvJSON) Get() (any, error) {
	if ev.Variable != nil && *ev.Variable != "" {
		rawValue := os.Getenv(*ev.Variable)
		if rawValue != "" {
			var result any

			err := json.Unmarshal([]byte(rawValue), &result)

			return result, err
		}
	}

	return ev.Value, nil
}

// GetCustom gets literal value or from system environment by a custom function.
func (ev EnvJSON) GetCustom(getFunc GetEnvFunc) (any, error) {
	if ev.Variable != nil && *ev.Variable != "" {
		rawValue, err := getFunc(*ev.Variable)
		if err != nil && !errors.Is(err, ErrEnvironmentVariableValueRequired) {
			return nil, err
		}

		if rawValue != "" {
			var result any

			err := json.Unmarshal([]byte(rawValue), &result)

			return result, err
		}
	}

	return ev.Value, nil
}

// Equal checks if this instance equals the target value.
func (ev EnvJSON) Equal(target EnvJSON) bool {
	isSameValue := (ev.Value == nil && target.Value == nil) ||
		(ev.Value != nil && target.Value != nil && reflect.DeepEqual(ev.Value, target.Value))
	if !isSameValue {
		return false
	}

	return (ev.Variable == nil && target.Variable == nil) ||
		(ev.Variable != nil && target.Variable != nil && *ev.Variable == *target.Variable)
}
