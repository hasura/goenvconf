package goenvconf

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
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

	if rawValue.Variable != nil && *rawValue.Variable == "" {
		return fmt.Errorf("EnvAny: %w", ErrEnvironmentVariableRequired)
	}

	*ev = EnvAny(rawValue)

	return nil
}

// Get gets literal value or from system environment.
func (ev EnvAny) Get() (any, error) {
	if ev.Variable != nil {
		if *ev.Variable == "" {
			return nil, fmt.Errorf("EnvAny: %w", ErrEnvironmentVariableRequired)
		}

		rawValue := os.Getenv(*ev.Variable)
		if rawValue != "" {
			var result any

			err := json.Unmarshal([]byte(rawValue), &result)

			return result, err
		}
	}

	return ev.Value, nil
}

// GetOrDefault returns the default value if the environment value is empty.
func (ev EnvAny) GetOrDefault(defaultValue any) (any, error) {
	result, err := ev.Get()
	if err != nil {
		if errors.Is(err, ErrEnvironmentVariableValueRequired) {
			return defaultValue, nil
		}

		return false, err
	}

	return result, nil
}
