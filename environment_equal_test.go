package goenvconf

import (
	"testing"
)

func TestEnvString_Equal(t *testing.T) {
	testCases := []struct {
		Name     string
		Input    EnvString
		Target   EnvString
		Expected bool
	}{
		{
			Name:     "both_nil_values_and_variables",
			Input:    EnvString{},
			Target:   EnvString{},
			Expected: true,
		},
		{
			Name:     "same_literal_values",
			Input:    NewEnvStringValue("hello"),
			Target:   NewEnvStringValue("hello"),
			Expected: true,
		},
		{
			Name:     "different_literal_values",
			Input:    NewEnvStringValue("hello"),
			Target:   NewEnvStringValue("world"),
			Expected: false,
		},
		{
			Name:     "same_variable_names",
			Input:    NewEnvStringVariable("MY_VAR"),
			Target:   NewEnvStringVariable("MY_VAR"),
			Expected: true,
		},
		{
			Name:     "different_variable_names",
			Input:    NewEnvStringVariable("VAR1"),
			Target:   NewEnvStringVariable("VAR2"),
			Expected: false,
		},
		{
			Name:     "same_value_and_variable",
			Input:    NewEnvString("MY_VAR", "default"),
			Target:   NewEnvString("MY_VAR", "default"),
			Expected: true,
		},
		{
			Name:     "same_variable_different_value",
			Input:    NewEnvString("MY_VAR", "value1"),
			Target:   NewEnvString("MY_VAR", "value2"),
			Expected: false,
		},
		{
			Name:     "different_variable_same_value",
			Input:    NewEnvString("VAR1", "value"),
			Target:   NewEnvString("VAR2", "value"),
			Expected: false,
		},
		{
			Name:     "value_vs_variable",
			Input:    NewEnvStringValue("hello"),
			Target:   NewEnvStringVariable("MY_VAR"),
			Expected: false,
		},
		{
			Name:     "nil_value_vs_non_nil",
			Input:    NewEnvStringVariable("MY_VAR"),
			Target:   NewEnvStringValue("hello"),
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

func TestEnvInt_Equal(t *testing.T) {
	testCases := []struct {
		Name     string
		Input    EnvInt
		Target   EnvInt
		Expected bool
	}{
		{
			Name:     "both_nil",
			Input:    EnvInt{},
			Target:   EnvInt{},
			Expected: true,
		},
		{
			Name:     "same_literal_values",
			Input:    NewEnvIntValue(42),
			Target:   NewEnvIntValue(42),
			Expected: true,
		},
		{
			Name:     "different_literal_values",
			Input:    NewEnvIntValue(42),
			Target:   NewEnvIntValue(100),
			Expected: false,
		},
		{
			Name:     "same_variable_names",
			Input:    NewEnvIntVariable("MY_VAR"),
			Target:   NewEnvIntVariable("MY_VAR"),
			Expected: true,
		},
		{
			Name:     "different_variable_names",
			Input:    NewEnvIntVariable("VAR1"),
			Target:   NewEnvIntVariable("VAR2"),
			Expected: false,
		},
		{
			Name:     "same_value_and_variable",
			Input:    NewEnvInt("MY_VAR", 42),
			Target:   NewEnvInt("MY_VAR", 42),
			Expected: true,
		},
		{
			Name:     "value_vs_variable",
			Input:    NewEnvIntValue(42),
			Target:   NewEnvIntVariable("MY_VAR"),
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

func TestEnvBool_Equal(t *testing.T) {
	testCases := []struct {
		Name     string
		Input    EnvBool
		Target   EnvBool
		Expected bool
	}{
		{
			Name:     "both_nil",
			Input:    EnvBool{},
			Target:   EnvBool{},
			Expected: true,
		},
		{
			Name:     "same_literal_true",
			Input:    NewEnvBoolValue(true),
			Target:   NewEnvBoolValue(true),
			Expected: true,
		},
		{
			Name:     "same_literal_false",
			Input:    NewEnvBoolValue(false),
			Target:   NewEnvBoolValue(false),
			Expected: true,
		},
		{
			Name:     "different_literal_values",
			Input:    NewEnvBoolValue(true),
			Target:   NewEnvBoolValue(false),
			Expected: false,
		},
		{
			Name:     "same_variable_names",
			Input:    NewEnvBoolVariable("MY_VAR"),
			Target:   NewEnvBoolVariable("MY_VAR"),
			Expected: true,
		},
		{
			Name:     "different_variable_names",
			Input:    NewEnvBoolVariable("VAR1"),
			Target:   NewEnvBoolVariable("VAR2"),
			Expected: false,
		},
		{
			Name:     "same_value_and_variable",
			Input:    NewEnvBool("MY_VAR", true),
			Target:   NewEnvBool("MY_VAR", true),
			Expected: true,
		},
		{
			Name:     "value_vs_variable",
			Input:    NewEnvBoolValue(true),
			Target:   NewEnvBoolVariable("MY_VAR"),
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

func TestEnvFloat_Equal(t *testing.T) {
	testCases := []struct {
		Name     string
		Input    EnvFloat
		Target   EnvFloat
		Expected bool
	}{
		{
			Name:     "both_nil",
			Input:    EnvFloat{},
			Target:   EnvFloat{},
			Expected: true,
		},
		{
			Name:     "same_literal_values",
			Input:    NewEnvFloatValue(3.14),
			Target:   NewEnvFloatValue(3.14),
			Expected: true,
		},
		{
			Name:     "different_literal_values",
			Input:    NewEnvFloatValue(3.14),
			Target:   NewEnvFloatValue(2.718),
			Expected: false,
		},
		{
			Name:     "same_variable_names",
			Input:    NewEnvFloatVariable("MY_VAR"),
			Target:   NewEnvFloatVariable("MY_VAR"),
			Expected: true,
		},
		{
			Name:     "different_variable_names",
			Input:    NewEnvFloatVariable("VAR1"),
			Target:   NewEnvFloatVariable("VAR2"),
			Expected: false,
		},
		{
			Name:     "same_value_and_variable",
			Input:    NewEnvFloat("MY_VAR", 3.14),
			Target:   NewEnvFloat("MY_VAR", 3.14),
			Expected: true,
		},
		{
			Name:     "value_vs_variable",
			Input:    NewEnvFloatValue(3.14),
			Target:   NewEnvFloatVariable("MY_VAR"),
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
