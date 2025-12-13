package goenvconf

import (
	"testing"
)

func TestEnvMapString_Equal(t *testing.T) {
	testCases := []struct {
		Name     string
		Input    EnvMapString
		Target   EnvMapString
		Expected bool
	}{
		{
			Name:     "both_nil",
			Input:    EnvMapString{},
			Target:   EnvMapString{},
			Expected: true,
		},
		{
			Name:     "same_literal_values",
			Input:    NewEnvMapStringValue(map[string]string{"key1": "value1", "key2": "value2"}),
			Target:   NewEnvMapStringValue(map[string]string{"key1": "value1", "key2": "value2"}),
			Expected: true,
		},
		{
			Name:     "different_literal_values",
			Input:    NewEnvMapStringValue(map[string]string{"key1": "value1"}),
			Target:   NewEnvMapStringValue(map[string]string{"key2": "value2"}),
			Expected: false,
		},
		{
			Name:     "same_variable_names",
			Input:    NewEnvMapStringVariable("MY_VAR"),
			Target:   NewEnvMapStringVariable("MY_VAR"),
			Expected: true,
		},
		{
			Name:     "different_variable_names",
			Input:    NewEnvMapStringVariable("VAR1"),
			Target:   NewEnvMapStringVariable("VAR2"),
			Expected: false,
		},
		{
			Name:     "same_value_and_variable",
			Input:    NewEnvMapString("MY_VAR", map[string]string{"key": "value"}),
			Target:   NewEnvMapString("MY_VAR", map[string]string{"key": "value"}),
			Expected: true,
		},
		{
			Name:     "value_vs_variable",
			Input:    NewEnvMapStringValue(map[string]string{"key": "value"}),
			Target:   NewEnvMapStringVariable("MY_VAR"),
			Expected: false,
		},
		{
			Name:     "nil_vs_empty_map",
			Input:    EnvMapString{Value: nil},
			Target:   NewEnvMapStringValue(map[string]string{}),
			Expected: false,
		},
		{
			Name:     "both_empty_maps",
			Input:    NewEnvMapStringValue(map[string]string{}),
			Target:   NewEnvMapStringValue(map[string]string{}),
			Expected: true,
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

func TestEnvMapInt_Equal(t *testing.T) {
	testCases := []struct {
		Name     string
		Input    EnvMapInt
		Target   EnvMapInt
		Expected bool
	}{
		{
			Name:     "both_nil",
			Input:    EnvMapInt{},
			Target:   EnvMapInt{},
			Expected: true,
		},
		{
			Name:     "same_literal_values",
			Input:    NewEnvMapIntValue(map[string]int64{"key1": 1, "key2": 2}),
			Target:   NewEnvMapIntValue(map[string]int64{"key1": 1, "key2": 2}),
			Expected: true,
		},
		{
			Name:     "different_literal_values",
			Input:    NewEnvMapIntValue(map[string]int64{"key1": 1}),
			Target:   NewEnvMapIntValue(map[string]int64{"key2": 2}),
			Expected: false,
		},
		{
			Name:     "same_variable_names",
			Input:    NewEnvMapIntVariable("MY_VAR"),
			Target:   NewEnvMapIntVariable("MY_VAR"),
			Expected: true,
		},
		{
			Name:     "different_variable_names",
			Input:    NewEnvMapIntVariable("VAR1"),
			Target:   NewEnvMapIntVariable("VAR2"),
			Expected: false,
		},
		{
			Name:     "same_value_and_variable",
			Input:    NewEnvMapInt("MY_VAR", map[string]int64{"key": 42}),
			Target:   NewEnvMapInt("MY_VAR", map[string]int64{"key": 42}),
			Expected: true,
		},
		{
			Name:     "value_vs_variable",
			Input:    NewEnvMapIntValue(map[string]int64{"key": 42}),
			Target:   NewEnvMapIntVariable("MY_VAR"),
			Expected: false,
		},
		{
			Name:     "nil_vs_empty_map",
			Input:    EnvMapInt{Value: nil},
			Target:   NewEnvMapIntValue(map[string]int64{}),
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

func TestEnvMapFloat_Equal(t *testing.T) {
	testCases := []struct {
		Name     string
		Input    EnvMapFloat
		Target   EnvMapFloat
		Expected bool
	}{
		{
			Name:     "both_nil",
			Input:    EnvMapFloat{},
			Target:   EnvMapFloat{},
			Expected: true,
		},
		{
			Name:     "same_literal_values",
			Input:    NewEnvMapFloatValue(map[string]float64{"key1": 3.14, "key2": 2.718}),
			Target:   NewEnvMapFloatValue(map[string]float64{"key1": 3.14, "key2": 2.718}),
			Expected: true,
		},
		{
			Name:     "different_literal_values",
			Input:    NewEnvMapFloatValue(map[string]float64{"key1": 3.14}),
			Target:   NewEnvMapFloatValue(map[string]float64{"key2": 2.718}),
			Expected: false,
		},
		{
			Name:     "same_variable_names",
			Input:    NewEnvMapFloatVariable("MY_VAR"),
			Target:   NewEnvMapFloatVariable("MY_VAR"),
			Expected: true,
		},
		{
			Name:     "different_variable_names",
			Input:    NewEnvMapFloatVariable("VAR1"),
			Target:   NewEnvMapFloatVariable("VAR2"),
			Expected: false,
		},
		{
			Name:     "same_value_and_variable",
			Input:    NewEnvMapFloat("MY_VAR", map[string]float64{"key": 3.14}),
			Target:   NewEnvMapFloat("MY_VAR", map[string]float64{"key": 3.14}),
			Expected: true,
		},
		{
			Name:     "value_vs_variable",
			Input:    NewEnvMapFloatValue(map[string]float64{"key": 3.14}),
			Target:   NewEnvMapFloatVariable("MY_VAR"),
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

func TestEnvMapBool_Equal(t *testing.T) {
	testCases := []struct {
		Name     string
		Input    EnvMapBool
		Target   EnvMapBool
		Expected bool
	}{
		{
			Name:     "both_nil",
			Input:    EnvMapBool{},
			Target:   EnvMapBool{},
			Expected: true,
		},
		{
			Name:     "same_literal_values",
			Input:    NewEnvMapBoolValue(map[string]bool{"key1": true, "key2": false}),
			Target:   NewEnvMapBoolValue(map[string]bool{"key1": true, "key2": false}),
			Expected: true,
		},
		{
			Name:     "different_literal_values",
			Input:    NewEnvMapBoolValue(map[string]bool{"key1": true}),
			Target:   NewEnvMapBoolValue(map[string]bool{"key2": false}),
			Expected: false,
		},
		{
			Name:     "same_variable_names",
			Input:    NewEnvMapBoolVariable("MY_VAR"),
			Target:   NewEnvMapBoolVariable("MY_VAR"),
			Expected: true,
		},
		{
			Name:     "different_variable_names",
			Input:    NewEnvMapBoolVariable("VAR1"),
			Target:   NewEnvMapBoolVariable("VAR2"),
			Expected: false,
		},
		{
			Name:     "same_value_and_variable",
			Input:    NewEnvMapBool("MY_VAR", map[string]bool{"key": true}),
			Target:   NewEnvMapBool("MY_VAR", map[string]bool{"key": true}),
			Expected: true,
		},
		{
			Name:     "value_vs_variable",
			Input:    NewEnvMapBoolValue(map[string]bool{"key": true}),
			Target:   NewEnvMapBoolVariable("MY_VAR"),
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
