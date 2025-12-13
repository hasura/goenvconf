package goenvconf

import (
	"testing"
)

func TestEnvAny_Equal(t *testing.T) {
	testCases := []struct {
		Name     string
		Input    EnvAny
		Target   EnvAny
		Expected bool
	}{
		{
			Name:     "both_nil",
			Input:    EnvAny{},
			Target:   EnvAny{},
			Expected: true,
		},
		{
			Name:     "same_string_values",
			Input:    NewEnvAnyValue("hello"),
			Target:   NewEnvAnyValue("hello"),
			Expected: true,
		},
		{
			Name:     "different_string_values",
			Input:    NewEnvAnyValue("hello"),
			Target:   NewEnvAnyValue("world"),
			Expected: false,
		},
		{
			Name:     "same_int_values",
			Input:    NewEnvAnyValue(42),
			Target:   NewEnvAnyValue(42),
			Expected: true,
		},
		{
			Name:     "different_int_values",
			Input:    NewEnvAnyValue(42),
			Target:   NewEnvAnyValue(100),
			Expected: false,
		},
		{
			Name:     "same_float_values",
			Input:    NewEnvAnyValue(3.14),
			Target:   NewEnvAnyValue(3.14),
			Expected: true,
		},
		{
			Name:     "different_float_values",
			Input:    NewEnvAnyValue(3.14),
			Target:   NewEnvAnyValue(2.718),
			Expected: false,
		},
		{
			Name:     "same_bool_values",
			Input:    NewEnvAnyValue(true),
			Target:   NewEnvAnyValue(true),
			Expected: true,
		},
		{
			Name:     "different_bool_values",
			Input:    NewEnvAnyValue(true),
			Target:   NewEnvAnyValue(false),
			Expected: false,
		},
		{
			Name:     "same_map_values",
			Input:    NewEnvAnyValue(map[string]any{"key": "value"}),
			Target:   NewEnvAnyValue(map[string]any{"key": "value"}),
			Expected: true,
		},
		{
			Name:     "different_map_values",
			Input:    NewEnvAnyValue(map[string]any{"key1": "value1"}),
			Target:   NewEnvAnyValue(map[string]any{"key2": "value2"}),
			Expected: false,
		},
		{
			Name:     "same_slice_values",
			Input:    NewEnvAnyValue([]any{1, 2, 3}),
			Target:   NewEnvAnyValue([]any{1, 2, 3}),
			Expected: true,
		},
		{
			Name:     "different_slice_values",
			Input:    NewEnvAnyValue([]any{1, 2, 3}),
			Target:   NewEnvAnyValue([]any{4, 5, 6}),
			Expected: false,
		},
		{
			Name:     "same_variable_names",
			Input:    NewEnvAnyVariable("MY_VAR"),
			Target:   NewEnvAnyVariable("MY_VAR"),
			Expected: true,
		},
		{
			Name:     "different_variable_names",
			Input:    NewEnvAnyVariable("VAR1"),
			Target:   NewEnvAnyVariable("VAR2"),
			Expected: false,
		},
		{
			Name:     "same_value_and_variable",
			Input:    NewEnvAny("MY_VAR", "default"),
			Target:   NewEnvAny("MY_VAR", "default"),
			Expected: true,
		},
		{
			Name:     "same_variable_different_value",
			Input:    NewEnvAny("MY_VAR", "value1"),
			Target:   NewEnvAny("MY_VAR", "value2"),
			Expected: false,
		},
		{
			Name:     "different_variable_same_value",
			Input:    NewEnvAny("VAR1", "value"),
			Target:   NewEnvAny("VAR2", "value"),
			Expected: false,
		},
		{
			Name:     "value_vs_variable",
			Input:    NewEnvAnyValue("hello"),
			Target:   NewEnvAnyVariable("MY_VAR"),
			Expected: false,
		},
		{
			Name:     "different_types",
			Input:    NewEnvAnyValue("42"),
			Target:   NewEnvAnyValue(42),
			Expected: false,
		},
		{
			Name:     "nil_value_vs_non_nil",
			Input:    NewEnvAnyVariable("MY_VAR"),
			Target:   NewEnvAnyValue("hello"),
			Expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			result := tc.Input.Equal(tc.Target)
			if result != tc.Expected {
				t.Errorf("Expected %v, got %v", tc.Expected, result)
			}
		})
	}
}
