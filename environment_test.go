package goenvconf

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestEnvString(t *testing.T) {
	t.Setenv("SOME_FOO", "bar")
	testCases := []struct {
		Input    EnvString
		Expected string
		ErrorMsg string
	}{
		{
			Input:    NewEnvStringValue("foo"),
			Expected: "foo",
		},
		{
			Input:    NewEnvStringVariable("SOME_FOO"),
			Expected: "bar",
		},
		{
			Input:    EnvString{},
			ErrorMsg: ErrEnvironmentValueRequired.Error(),
		},
		{
			Input:    NewEnvString("SOME_BAR", "bar"),
			Expected: "bar",
		},
		{
			Input: EnvString{
				Variable: toPtr(""),
			},
			ErrorMsg: ErrEnvironmentValueRequired.Error(),
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
				assertDeepEqual(t, tc.Input.IsZero(), false)
			}
		})
	}

	t.Run("json_decode", func(t *testing.T) {
		var ev EnvString
		assertNilError(t, json.Unmarshal([]byte(`{"env": "SOME_FOO"}`), &ev))
		result, err := ev.Get()
		assertNilError(t, err)
		assertDeepEqual(t, "bar", result)
	})

	t.Run("get_default", func(t *testing.T) {
		result, err := NewEnvStringVariable("SOME_BAZ").GetOrDefault("baz")
		assertNilError(t, err)
		assertDeepEqual(t, "baz", result)
	})
}

func TestEnvBool(t *testing.T) {
	t.Setenv("SOME_FOO", "true")
	testCases := []struct {
		Input    EnvBool
		Expected bool
		ErrorMsg string
	}{
		{
			Input:    NewEnvBoolValue(true),
			Expected: true,
		},
		{
			Input:    NewEnvBoolVariable("SOME_FOO"),
			Expected: true,
		},
		{
			Input:    NewEnvBoolVariable("SOME_FOO_2"),
			ErrorMsg: getEnvVariableValueRequiredError(toPtr("SOME_FOO_2")).Error(),
		},
		{
			Input:    EnvBool{},
			ErrorMsg: ErrEnvironmentValueRequired.Error(),
		},
		{
			Input:    NewEnvBool("SOME_FOO_2", true),
			Expected: true,
		},
		{
			Input: EnvBool{
				Variable: toPtr(""),
			},
			ErrorMsg: ErrEnvironmentValueRequired.Error(),
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			result, err := tc.Input.Get()
			if tc.ErrorMsg != "" {
				assertErrorContains(t, err, tc.ErrorMsg)
				if tc.ErrorMsg == ErrEnvironmentVariableValueRequired.Error() {
					newValue, err := tc.Input.GetOrDefault(true)
					assertNilError(t, err)
					assertDeepEqual(t, newValue, true)
				}
			} else {
				assertNilError(t, err)
				assertDeepEqual(t, result, tc.Expected)

				newValue, err := tc.Input.GetOrDefault(true)
				assertNilError(t, err)
				assertDeepEqual(t, newValue, tc.Expected)
				assertDeepEqual(t, tc.Input.IsZero(), false)
			}
		})
	}

	t.Run("json_decode", func(t *testing.T) {
		var ev EnvBool
		assertNilError(t, json.Unmarshal([]byte(`{"env": "SOME_FOO"}`), &ev))
		result, err := ev.Get()
		assertNilError(t, err)
		assertDeepEqual(t, true, result)
	})

	t.Run("get_default", func(t *testing.T) {
		result, err := NewEnvBoolVariable("SOME_TRUE").GetOrDefault(true)
		assertNilError(t, err)
		assertDeepEqual(t, true, result)

		assertDeepEqual(t, EnvBool{}.IsZero(), true)
	})
}

func TestEnvInt(t *testing.T) {
	t.Setenv("SOME_FOO", "10")
	testCases := []struct {
		Input    EnvInt
		Expected int64
		ErrorMsg string
	}{
		{
			Input:    NewEnvIntValue(1),
			Expected: 1,
		},
		{
			Input:    NewEnvIntVariable("SOME_FOO"),
			Expected: 10,
		},
		{
			Input:    NewEnvIntVariable("SOME_FOO_2"),
			ErrorMsg: ErrEnvironmentVariableValueRequired.Error(),
		},
		{
			Input:    EnvInt{},
			ErrorMsg: ErrEnvironmentValueRequired.Error(),
		},
		{
			Input:    NewEnvInt("SOME_FOO_2", 10),
			Expected: 10,
		},
		{
			Input: EnvInt{
				Variable: toPtr(""),
			},
			ErrorMsg: ErrEnvironmentValueRequired.Error(),
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			result, err := tc.Input.Get()
			if tc.ErrorMsg != "" {
				assertErrorContains(t, err, tc.ErrorMsg)
				if tc.ErrorMsg == ErrEnvironmentVariableValueRequired.Error() {
					newValue, err := tc.Input.GetOrDefault(100)
					assertNilError(t, err)
					assertDeepEqual(t, newValue, int64(100))
				}

			} else {
				assertNilError(t, err)
				assertDeepEqual(t, result, tc.Expected)

				newValue, err := tc.Input.GetOrDefault(100)
				assertNilError(t, err)
				assertDeepEqual(t, newValue, tc.Expected)

				assertDeepEqual(t, tc.Input.IsZero(), false)
			}
		})
	}

	t.Run("json_decode", func(t *testing.T) {
		var ev EnvInt
		assertNilError(t, json.Unmarshal([]byte(`{"env": "SOME_FOO"}`), &ev))
		result, err := ev.Get()
		assertNilError(t, err)
		assertDeepEqual(t, int64(10), result)
	})
}

func TestEnvFloat(t *testing.T) {
	t.Setenv("SOME_FOO", "10.5")
	testCases := []struct {
		Input    EnvFloat
		Expected float64
		ErrorMsg string
	}{
		{
			Input:    NewEnvFloatValue(1.1),
			Expected: 1.1,
		},
		{
			Input:    NewEnvFloatVariable("SOME_FOO"),
			Expected: 10.5,
		},
		{
			Input:    NewEnvFloatVariable("SOME_FOO_2"),
			ErrorMsg: ErrEnvironmentVariableValueRequired.Error(),
		},
		{
			Input:    EnvFloat{},
			ErrorMsg: ErrEnvironmentValueRequired.Error(),
		},
		{
			Input:    NewEnvFloat("SOME_FOO_1", 10),
			Expected: 10,
		},
		{
			Input: EnvFloat{
				Variable: toPtr(""),
			},
			ErrorMsg: ErrEnvironmentValueRequired.Error(),
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			result, err := tc.Input.Get()
			if tc.ErrorMsg != "" {
				assertErrorContains(t, err, tc.ErrorMsg)
				if tc.ErrorMsg == ErrEnvironmentVariableValueRequired.Error() {
					newValue, err := tc.Input.GetOrDefault(100.5)
					assertNilError(t, err)
					assertDeepEqual(t, newValue, float64(100.5))
				}
			} else {
				assertNilError(t, err)
				assertDeepEqual(t, result, tc.Expected)

				newValue, err := tc.Input.GetOrDefault(100)
				assertNilError(t, err)
				assertDeepEqual(t, newValue, tc.Expected)
				assertDeepEqual(t, tc.Input.IsZero(), false)
			}
		})
	}

	t.Run("json_decode", func(t *testing.T) {
		var ev EnvFloat
		assertNilError(t, json.Unmarshal([]byte(`{"env": "SOME_FOO"}`), &ev))
		result, err := ev.Get()
		assertNilError(t, err)
		assertDeepEqual(t, float64(10.5), result)
	})
}

// mockGetEnvFunc creates a mock GetEnvFunc that returns predefined values
func mockGetEnvFunc(values map[string]string, returnError bool) GetEnvFunc {
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

func TestEnvString_GetCustom(t *testing.T) {
	t.Setenv("TEST_VAR", "bar")

	testCases := []struct {
		Name     string
		Input    EnvString
		GetFunc  GetEnvFunc
		Expected string
		ErrorMsg string
	}{
		{
			Name:     "literal_value",
			Input:    NewEnvStringValue("foo"),
			GetFunc:  OSEnvGetter(context.TODO()),
			Expected: "foo",
		},
		{
			Name:     "variable_from_custom_func",
			Input:    NewEnvStringVariable("TEST_VAR"),
			GetFunc:  OSEnvGetter(context.TODO()),
			Expected: "bar",
		},
		{
			Name:     "variable_with_fallback_value",
			Input:    NewEnvString("CUSTOM_VAR", "fallback"),
			GetFunc:  mockGetEnvFunc(map[string]string{"CUSTOM_VAR": "custom"}, false),
			Expected: "custom",
		},
		{
			Name:     "empty_variable_returns_empty_string",
			Input:    NewEnvString("EMPTY_VAR", "fallback"),
			GetFunc:  mockGetEnvFunc(map[string]string{"EMPTY_VAR": ""}, false),
			Expected: "",
		},
		{
			Name:     "zero_value_error",
			Input:    EnvString{},
			GetFunc:  mockGetEnvFunc(map[string]string{}, false),
			ErrorMsg: ErrEnvironmentValueRequired.Error(),
		},
		{
			Name:     "custom_func_error",
			Input:    NewEnvStringVariable("SOME_VAR"),
			GetFunc:  OSEnvGetter(context.TODO()),
			ErrorMsg: ErrEnvironmentVariableValueRequired.Error(),
		},
		{
			Name:     "missing_variable_returns_empty_string",
			Input:    NewEnvStringVariable("MISSING_VAR"),
			GetFunc:  mockGetEnvFunc(map[string]string{}, false),
			Expected: "",
		},
		{
			Name:     "only_value_no_variable",
			Input:    NewEnvStringValue("only_value"),
			GetFunc:  mockGetEnvFunc(map[string]string{}, false),
			Expected: "only_value",
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

func TestEnvInt_GetCustom(t *testing.T) {
	testCases := []struct {
		Name     string
		Input    EnvInt
		GetFunc  GetEnvFunc
		Expected int64
		ErrorMsg string
	}{
		{
			Name:     "literal_value",
			Input:    NewEnvIntValue(42),
			GetFunc:  mockGetEnvFunc(map[string]string{}, false),
			Expected: 42,
		},
		{
			Name:     "variable_from_custom_func",
			Input:    NewEnvIntVariable("CUSTOM_INT"),
			GetFunc:  mockGetEnvFunc(map[string]string{"CUSTOM_INT": "100"}, false),
			Expected: 100,
		},
		{
			Name:     "variable_with_fallback_value",
			Input:    NewEnvInt("CUSTOM_INT", 50),
			GetFunc:  mockGetEnvFunc(map[string]string{"CUSTOM_INT": "200"}, false),
			Expected: 200,
		},
		{
			Name:     "empty_variable_uses_fallback",
			Input:    NewEnvInt("EMPTY_INT", 99),
			GetFunc:  mockGetEnvFunc(map[string]string{"EMPTY_INT": ""}, false),
			Expected: 99,
		},
		{
			Name:     "zero_value_error",
			Input:    EnvInt{},
			GetFunc:  mockGetEnvFunc(map[string]string{}, false),
			ErrorMsg: ErrEnvironmentValueRequired.Error(),
		},
		{
			Name:     "custom_func_error",
			Input:    NewEnvIntVariable("SOME_INT"),
			GetFunc:  mockGetEnvFunc(map[string]string{}, true),
			ErrorMsg: "mock error",
		},
		{
			Name:     "invalid_int_format",
			Input:    NewEnvIntVariable("INVALID_INT"),
			GetFunc:  mockGetEnvFunc(map[string]string{"INVALID_INT": "not_a_number"}, false),
			ErrorMsg: "invalid syntax",
		},
		{
			Name:     "missing_variable_no_fallback",
			Input:    NewEnvIntVariable("MISSING_INT"),
			GetFunc:  mockGetEnvFunc(map[string]string{}, false),
			ErrorMsg: ErrEnvironmentVariableValueRequired.Error(),
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

func TestEnvBool_GetCustom(t *testing.T) {
	testCases := []struct {
		Name     string
		Input    EnvBool
		GetFunc  GetEnvFunc
		Expected bool
		ErrorMsg string
	}{
		{
			Name:     "literal_value_true",
			Input:    NewEnvBoolValue(true),
			GetFunc:  mockGetEnvFunc(map[string]string{}, false),
			Expected: true,
		},
		{
			Name:     "literal_value_false",
			Input:    NewEnvBoolValue(false),
			GetFunc:  mockGetEnvFunc(map[string]string{}, false),
			Expected: false,
		},
		{
			Name:     "variable_from_custom_func_true",
			Input:    NewEnvBoolVariable("CUSTOM_BOOL"),
			GetFunc:  mockGetEnvFunc(map[string]string{"CUSTOM_BOOL": "true"}, false),
			Expected: true,
		},
		{
			Name:     "variable_from_custom_func_false",
			Input:    NewEnvBoolVariable("CUSTOM_BOOL"),
			GetFunc:  mockGetEnvFunc(map[string]string{"CUSTOM_BOOL": "false"}, false),
			Expected: false,
		},
		{
			Name:     "variable_with_fallback_value",
			Input:    NewEnvBool("CUSTOM_BOOL", false),
			GetFunc:  mockGetEnvFunc(map[string]string{"CUSTOM_BOOL": "true"}, false),
			Expected: true,
		},
		{
			Name:     "empty_variable_uses_fallback",
			Input:    NewEnvBool("EMPTY_BOOL", true),
			GetFunc:  mockGetEnvFunc(map[string]string{"EMPTY_BOOL": ""}, false),
			Expected: true,
		},
		{
			Name:     "zero_value_error",
			Input:    EnvBool{},
			GetFunc:  mockGetEnvFunc(map[string]string{}, false),
			ErrorMsg: ErrEnvironmentValueRequired.Error(),
		},
		{
			Name:     "custom_func_error",
			Input:    NewEnvBoolVariable("SOME_BOOL"),
			GetFunc:  mockGetEnvFunc(map[string]string{}, true),
			ErrorMsg: "mock error",
		},
		{
			Name:     "invalid_bool_format",
			Input:    NewEnvBoolVariable("INVALID_BOOL"),
			GetFunc:  mockGetEnvFunc(map[string]string{"INVALID_BOOL": "not_a_bool"}, false),
			ErrorMsg: "invalid syntax",
		},
		{
			Name:     "missing_variable_no_fallback",
			Input:    NewEnvBoolVariable("MISSING_BOOL"),
			GetFunc:  mockGetEnvFunc(map[string]string{}, false),
			ErrorMsg: ErrEnvironmentVariableValueRequired.Error(),
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

func TestEnvFloat_GetCustom(t *testing.T) {
	testCases := []struct {
		Name     string
		Input    EnvFloat
		GetFunc  GetEnvFunc
		Expected float64
		ErrorMsg string
	}{
		{
			Name:     "literal_value",
			Input:    NewEnvFloatValue(3.14),
			GetFunc:  mockGetEnvFunc(map[string]string{}, false),
			Expected: 3.14,
		},
		{
			Name:     "variable_from_custom_func",
			Input:    NewEnvFloatVariable("CUSTOM_FLOAT"),
			GetFunc:  mockGetEnvFunc(map[string]string{"CUSTOM_FLOAT": "2.718"}, false),
			Expected: 2.718,
		},
		{
			Name:     "variable_with_fallback_value",
			Input:    NewEnvFloat("CUSTOM_FLOAT", 1.5),
			GetFunc:  mockGetEnvFunc(map[string]string{"CUSTOM_FLOAT": "9.81"}, false),
			Expected: 9.81,
		},
		{
			Name:     "empty_variable_uses_fallback",
			Input:    NewEnvFloat("EMPTY_FLOAT", 0.5),
			GetFunc:  mockGetEnvFunc(map[string]string{"EMPTY_FLOAT": ""}, false),
			Expected: 0.5,
		},
		{
			Name:     "zero_value_error",
			Input:    EnvFloat{},
			GetFunc:  mockGetEnvFunc(map[string]string{}, false),
			ErrorMsg: ErrEnvironmentValueRequired.Error(),
		},
		{
			Name:     "custom_func_error",
			Input:    NewEnvFloatVariable("SOME_FLOAT"),
			GetFunc:  mockGetEnvFunc(map[string]string{}, true),
			ErrorMsg: "mock error",
		},
		{
			Name:     "invalid_float_format",
			Input:    NewEnvFloatVariable("INVALID_FLOAT"),
			GetFunc:  mockGetEnvFunc(map[string]string{"INVALID_FLOAT": "not_a_float"}, false),
			ErrorMsg: "invalid syntax",
		},
		{
			Name:     "missing_variable_no_fallback",
			Input:    NewEnvFloatVariable("MISSING_FLOAT"),
			GetFunc:  mockGetEnvFunc(map[string]string{}, false),
			ErrorMsg: ErrEnvironmentVariableValueRequired.Error(),
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

func assertNilError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Errorf("expected nil error, got: %s", err)
		t.FailNow()
	}
}

func assertErrorContains(t *testing.T, err error, msg string) {
	t.Helper()

	if err == nil {
		t.Errorf("expected error with content: `%s`, got: nil", msg)
		t.FailNow()
	}

	if !strings.Contains(err.Error(), msg) {
		t.Errorf("expected error with content: %s, got: %s", msg, err)
		t.FailNow()
	}
}

func assertDeepEqual(t *testing.T, expected, reality any) {
	t.Helper()

	if !reflect.DeepEqual(expected, reality) {
		t.Errorf("%v != %v", expected, reality)
		t.FailNow()
	}
}

func toPtr[T any](input T) *T {
	return &input
}

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
