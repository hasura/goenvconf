package goenvconf

import (
	"fmt"
	"testing"
)

func TestEnvJSON(t *testing.T) {
	t.Setenv("SOME_FOO", "2.2")

	testCases := []struct {
		Input    EnvJSON
		Expected any
		ErrorMsg string
	}{
		{
			Input: NewEnvJSONValue(map[string]float64{
				"foo": 1.1,
			}),
			Expected: map[string]float64{
				"foo": 1.1,
			},
		},
		{
			Input:    NewEnvJSONVariable("SOME_FOO"),
			Expected: float64(2.2),
		},
		{
			Input:    EnvJSON{},
			Expected: nil,
		},
		{
			Input:    NewEnvJSON("SOME_FOO_2", "baz"),
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

func TestEnvJSON_GetCustom(t *testing.T) {
	testCases := []struct {
		Name     string
		Input    EnvJSON
		GetFunc  GetEnvFunc
		Expected any
		ErrorMsg string
	}{
		{
			Name:     "literal_string_value",
			Input:    NewEnvJSONValue("hello"),
			GetFunc:  mockGetEnvFuncForAny(map[string]string{}, false),
			Expected: "hello",
		},
		{
			Name:     "literal_number_value",
			Input:    NewEnvJSONValue(42.5),
			GetFunc:  mockGetEnvFuncForAny(map[string]string{}, false),
			Expected: 42.5,
		},
		{
			Name:     "literal_map_value",
			Input:    NewEnvJSONValue(map[string]any{"key": "value"}),
			GetFunc:  mockGetEnvFuncForAny(map[string]string{}, false),
			Expected: map[string]any{"key": "value"},
		},
		{
			Name:     "variable_from_custom_func_string",
			Input:    NewEnvJSONVariable("CUSTOM_VAR"),
			GetFunc:  mockGetEnvFuncForAny(map[string]string{"CUSTOM_VAR": `"test_string"`}, false),
			Expected: "test_string",
		},
		{
			Name:     "variable_from_custom_func_number",
			Input:    NewEnvJSONVariable("CUSTOM_VAR"),
			GetFunc:  mockGetEnvFuncForAny(map[string]string{"CUSTOM_VAR": "123.45"}, false),
			Expected: 123.45,
		},
		{
			Name:     "variable_from_custom_func_json_object",
			Input:    NewEnvJSONVariable("CUSTOM_VAR"),
			GetFunc:  mockGetEnvFuncForAny(map[string]string{"CUSTOM_VAR": `{"foo":"bar","num":42}`}, false),
			Expected: map[string]any{"foo": "bar", "num": float64(42)},
		},
		{
			Name:     "variable_from_custom_func_json_array",
			Input:    NewEnvJSONVariable("CUSTOM_VAR"),
			GetFunc:  mockGetEnvFuncForAny(map[string]string{"CUSTOM_VAR": `[1,2,3]`}, false),
			Expected: []any{float64(1), float64(2), float64(3)},
		},
		{
			Name:     "variable_with_fallback_value",
			Input:    NewEnvJSON("CUSTOM_VAR", "fallback"),
			GetFunc:  mockGetEnvFuncForAny(map[string]string{"CUSTOM_VAR": `"custom"`}, false),
			Expected: "custom",
		},
		{
			Name:     "empty_variable_uses_fallback",
			Input:    NewEnvJSON("EMPTY_VAR", "fallback"),
			GetFunc:  mockGetEnvFuncForAny(map[string]string{"EMPTY_VAR": ""}, false),
			Expected: "fallback",
		},
		{
			Name:     "nil_value_and_no_variable",
			Input:    EnvJSON{},
			GetFunc:  mockGetEnvFuncForAny(map[string]string{}, false),
			Expected: nil,
		},
		{
			Name:     "custom_func_error",
			Input:    NewEnvJSONVariable("SOME_VAR"),
			GetFunc:  mockGetEnvFuncForAny(map[string]string{}, true),
			ErrorMsg: "mock error",
		},
		{
			Name:     "invalid_json_format",
			Input:    NewEnvJSONVariable("INVALID_VAR"),
			GetFunc:  mockGetEnvFuncForAny(map[string]string{"INVALID_VAR": `{invalid json`}, false),
			ErrorMsg: "invalid character",
		},
		{
			Name:     "missing_variable_returns_nil",
			Input:    NewEnvJSONVariable("MISSING_VAR"),
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

func TestEnvJSON_Equal(t *testing.T) {
	testCases := []struct {
		Name     string
		Input    EnvJSON
		Target   EnvJSON
		Expected bool
	}{
		{
			Name:     "both_nil",
			Input:    EnvJSON{},
			Target:   EnvJSON{},
			Expected: true,
		},
		{
			Name:     "same_string_values",
			Input:    NewEnvJSONValue("hello"),
			Target:   NewEnvJSONValue("hello"),
			Expected: true,
		},
		{
			Name:     "different_string_values",
			Input:    NewEnvJSONValue("hello"),
			Target:   NewEnvJSONValue("world"),
			Expected: false,
		},
		{
			Name:     "same_int_values",
			Input:    NewEnvJSONValue(42),
			Target:   NewEnvJSONValue(42),
			Expected: true,
		},
		{
			Name:     "different_int_values",
			Input:    NewEnvJSONValue(42),
			Target:   NewEnvJSONValue(100),
			Expected: false,
		},
		{
			Name:     "same_float_values",
			Input:    NewEnvJSONValue(3.14),
			Target:   NewEnvJSONValue(3.14),
			Expected: true,
		},
		{
			Name:     "different_float_values",
			Input:    NewEnvJSONValue(3.14),
			Target:   NewEnvJSONValue(2.718),
			Expected: false,
		},
		{
			Name:     "same_bool_values",
			Input:    NewEnvJSONValue(true),
			Target:   NewEnvJSONValue(true),
			Expected: true,
		},
		{
			Name:     "different_bool_values",
			Input:    NewEnvJSONValue(true),
			Target:   NewEnvJSONValue(false),
			Expected: false,
		},
		{
			Name:     "same_map_values",
			Input:    NewEnvJSONValue(map[string]any{"key": "value"}),
			Target:   NewEnvJSONValue(map[string]any{"key": "value"}),
			Expected: true,
		},
		{
			Name:     "different_map_values",
			Input:    NewEnvJSONValue(map[string]any{"key1": "value1"}),
			Target:   NewEnvJSONValue(map[string]any{"key2": "value2"}),
			Expected: false,
		},
		{
			Name:     "same_slice_values",
			Input:    NewEnvJSONValue([]any{1, 2, 3}),
			Target:   NewEnvJSONValue([]any{1, 2, 3}),
			Expected: true,
		},
		{
			Name:     "different_slice_values",
			Input:    NewEnvJSONValue([]any{1, 2, 3}),
			Target:   NewEnvJSONValue([]any{4, 5, 6}),
			Expected: false,
		},
		{
			Name:     "same_variable_names",
			Input:    NewEnvJSONVariable("MY_VAR"),
			Target:   NewEnvJSONVariable("MY_VAR"),
			Expected: true,
		},
		{
			Name:     "different_variable_names",
			Input:    NewEnvJSONVariable("VAR1"),
			Target:   NewEnvJSONVariable("VAR2"),
			Expected: false,
		},
		{
			Name:     "same_value_and_variable",
			Input:    NewEnvJSON("MY_VAR", "default"),
			Target:   NewEnvJSON("MY_VAR", "default"),
			Expected: true,
		},
		{
			Name:     "same_variable_different_value",
			Input:    NewEnvJSON("MY_VAR", "value1"),
			Target:   NewEnvJSON("MY_VAR", "value2"),
			Expected: false,
		},
		{
			Name:     "different_variable_same_value",
			Input:    NewEnvJSON("VAR1", "value"),
			Target:   NewEnvJSON("VAR2", "value"),
			Expected: false,
		},
		{
			Name:     "value_vs_variable",
			Input:    NewEnvJSONValue("hello"),
			Target:   NewEnvJSONVariable("MY_VAR"),
			Expected: false,
		},
		{
			Name:     "different_types",
			Input:    NewEnvJSONValue("42"),
			Target:   NewEnvJSONValue(42),
			Expected: false,
		},
		{
			Name:     "nil_value_vs_non_nil",
			Input:    NewEnvJSONVariable("MY_VAR"),
			Target:   NewEnvJSONValue("hello"),
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
