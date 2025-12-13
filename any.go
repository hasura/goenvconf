package goenvconf

import (
	"encoding/json"
	"os"
	"reflect"
)

// EnvAny represents either arbitrary value or an environment reference.
type EnvAny struct {
	Value    any     `json:"value,omitempty" jsonschema:"anyof_required=value" mapstructure:"value" yaml:"value,omitempty"`
	Variable *string `json:"env,omitempty"   jsonschema:"anyof_required=env"   mapstructure:"env"   yaml:"env,omitempty"`
}

// NewEnvAny creates an EnvAny instance.
func NewEnvAny(env string, value any) EnvAny {
	return EnvAny{
		Variable: &env,
		Value:    value,
	}
}

// NewEnvAnyValue creates an EnvAny with a literal value.
func NewEnvAnyValue(value any) EnvAny {
	return EnvAny{
		Value: value,
	}
}

// NewEnvAnyVariable creates an EnvAny with a variable name.
func NewEnvAnyVariable(name string) EnvAny {
	return EnvAny{
		Variable: &name,
	}
}

// UnmarshalJSON implements json.Unmarshaler.
func (ev *EnvAny) UnmarshalJSON(b []byte) error {
	type Plain EnvAny

	var rawValue Plain

	err := json.Unmarshal(b, &rawValue)
	if err != nil {
		return err
	}

	*ev = EnvAny(rawValue)

	return nil
}

// IsZero checks if the instance is empty.
func (ev EnvAny) IsZero() bool {
	return (ev.Variable == nil || *ev.Variable == "") &&
		ev.Value == nil
}

// Get gets literal value or from system environment.
func (ev EnvAny) Get() (any, error) {
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
func (ev EnvAny) GetCustom(getFunc GetEnvFunc) (any, error) {
	if ev.Variable != nil && *ev.Variable != "" {
		rawValue, err := getFunc(*ev.Variable)
		if err != nil {
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
func (ev EnvAny) Equal(target EnvAny) bool {
	isSameValue := (ev.Value == nil && target.Value == nil) ||
		(ev.Value != nil && target.Value != nil && reflect.DeepEqual(ev.Value, target.Value))
	if !isSameValue {
		return false
	}

	return (ev.Variable == nil && target.Variable == nil) ||
		(ev.Variable != nil && target.Variable != nil && *ev.Variable == *target.Variable)
}
