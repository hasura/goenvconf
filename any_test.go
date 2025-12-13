package goenvconf

import (
	"errors"
	"fmt"
	"testing"
)

func TestEnvAny(t *testing.T) {
	t.Setenv("SOME_FOO", "2.2")

	testCases := []struct {
		Input    EnvAny
		Expected any
		ErrorMsg string
	}{
		{
			Input: NewEnvAnyValue(map[string]float64{
				"foo": 1.1,
			}),
			Expected: map[string]float64{
				"foo": 1.1,
			},
		},
		{
			Input:    NewEnvAnyVariable("SOME_FOO"),
			Expected: float64(2.2),
		},
		{
			Input:    EnvAny{},
			Expected: nil,
		},
		{
			Input:    NewEnvAny("SOME_FOO_2", "baz"),
			Expected: "baz",
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			result, err := tc.Input.Get()

			if tc.ErrorMsg != "" {
				assertErrorContains(t, err, tc.ErrorMsg)
			} else {
				assertNilError(t, err)
				assertDeepEqual(t, result, tc.Expected)
			}

			assertDeepEqual(t, tc.Input.IsZero(), tc.Expected == nil)
		})
	}
}

// mockGetEnvFuncForAny creates a mock GetEnvFunc for EnvAny tests
func mockGetEnvFuncForAny(values map[string]string, returnError bool) GetEnvFunc {
	return func(key string) (string, error) {
		if returnError {
			return "", errors.New("mock error: failed to get environment variable")
		}
		if val, ok := values[key]; ok {
			return val, nil
		}
		return "", nil
	}
}

func TestEnvAny_GetCustom(t *testing.T) {
	testCases := []struct {
		Name     string
		Input    EnvAny
		GetFunc  GetEnvFunc
		Expected any
		ErrorMsg string
	}{
		{
			Name:     "literal_string_value",
			Input:    NewEnvAnyValue("hello"),
			GetFunc:  mockGetEnvFuncForAny(map[string]string{}, false),
			Expected: "hello",
		},
		{
			Name:     "literal_number_value",
			Input:    NewEnvAnyValue(42.5),
			GetFunc:  mockGetEnvFuncForAny(map[string]string{}, false),
			Expected: 42.5,
		},
		{
			Name:     "literal_map_value",
			Input:    NewEnvAnyValue(map[string]any{"key": "value"}),
			GetFunc:  mockGetEnvFuncForAny(map[string]string{}, false),
			Expected: map[string]any{"key": "value"},
		},
		{
			Name:     "variable_from_custom_func_string",
			Input:    NewEnvAnyVariable("CUSTOM_VAR"),
			GetFunc:  mockGetEnvFuncForAny(map[string]string{"CUSTOM_VAR": `"test_string"`}, false),
			Expected: "test_string",
		},
		{
			Name:     "variable_from_custom_func_number",
			Input:    NewEnvAnyVariable("CUSTOM_VAR"),
			GetFunc:  mockGetEnvFuncForAny(map[string]string{"CUSTOM_VAR": "123.45"}, false),
			Expected: 123.45,
		},
		{
			Name:     "variable_from_custom_func_json_object",
			Input:    NewEnvAnyVariable("CUSTOM_VAR"),
			GetFunc:  mockGetEnvFuncForAny(map[string]string{"CUSTOM_VAR": `{"foo":"bar","num":42}`}, false),
			Expected: map[string]any{"foo": "bar", "num": float64(42)},
		},
		{
			Name:     "variable_from_custom_func_json_array",
			Input:    NewEnvAnyVariable("CUSTOM_VAR"),
			GetFunc:  mockGetEnvFuncForAny(map[string]string{"CUSTOM_VAR": `[1,2,3]`}, false),
			Expected: []any{float64(1), float64(2), float64(3)},
		},
		{
			Name:     "variable_with_fallback_value",
			Input:    NewEnvAny("CUSTOM_VAR", "fallback"),
			GetFunc:  mockGetEnvFuncForAny(map[string]string{"CUSTOM_VAR": `"custom"`}, false),
			Expected: "custom",
		},
		{
			Name:     "empty_variable_uses_fallback",
			Input:    NewEnvAny("EMPTY_VAR", "fallback"),
			GetFunc:  mockGetEnvFuncForAny(map[string]string{"EMPTY_VAR": ""}, false),
			Expected: "fallback",
		},
		{
			Name:     "nil_value_and_no_variable",
			Input:    EnvAny{},
			GetFunc:  mockGetEnvFuncForAny(map[string]string{}, false),
			Expected: nil,
		},
		{
			Name:     "custom_func_error",
			Input:    NewEnvAnyVariable("SOME_VAR"),
			GetFunc:  mockGetEnvFuncForAny(map[string]string{}, true),
			ErrorMsg: "mock error",
		},
		{
			Name:     "invalid_json_format",
			Input:    NewEnvAnyVariable("INVALID_VAR"),
			GetFunc:  mockGetEnvFuncForAny(map[string]string{"INVALID_VAR": `{invalid json`}, false),
			ErrorMsg: "invalid character",
		},
		{
			Name:     "missing_variable_returns_nil",
			Input:    NewEnvAnyVariable("MISSING_VAR"),
			GetFunc:  mockGetEnvFuncForAny(map[string]string{}, false),
			Expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			result, err := tc.Input.GetCustom(tc.GetFunc)
			if tc.ErrorMsg != "" {
				assertErrorContains(t, err, tc.ErrorMsg)
			} else {
				assertNilError(t, err)
				assertDeepEqual(t, tc.Expected, result)
			}
		})
	}
}

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
