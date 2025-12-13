package goenvconf

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
)

// mockGetEnvFuncForMaps creates a mock GetEnvFunc for map tests
func mockGetEnvFuncForMaps(values map[string]string, returnError bool) GetEnvFunc {
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

func TestEnvMapString_GetCustom(t *testing.T) {
	testCases := []struct {
		Name     string
		Input    EnvMapString
		GetFunc  GetEnvFunc
		Expected map[string]string
		ErrorMsg string
	}{
		{
			Name:     "literal_value",
			Input:    NewEnvMapStringValue(map[string]string{"key1": "val1", "key2": "val2"}),
			GetFunc:  mockGetEnvFuncForMaps(map[string]string{}, false),
			Expected: map[string]string{"key1": "val1", "key2": "val2"},
		},
		{
			Name:     "variable_from_custom_func",
			Input:    NewEnvMapStringVariable("CUSTOM_MAP"),
			GetFunc:  mockGetEnvFuncForMaps(map[string]string{"CUSTOM_MAP": "foo=bar;baz=qux"}, false),
			Expected: map[string]string{"foo": "bar", "baz": "qux"},
		},
		{
			Name:     "variable_with_fallback_value",
			Input:    NewEnvMapString("CUSTOM_MAP", map[string]string{"default": "value"}),
			GetFunc:  mockGetEnvFuncForMaps(map[string]string{"CUSTOM_MAP": "custom=data"}, false),
			Expected: map[string]string{"custom": "data"},
		},
		{
			Name:     "empty_variable_uses_fallback",
			Input:    NewEnvMapString("EMPTY_MAP", map[string]string{"fallback": "value"}),
			GetFunc:  mockGetEnvFuncForMaps(map[string]string{"EMPTY_MAP": ""}, false),
			Expected: map[string]string{"fallback": "value"},
		},
		{
			Name:     "nil_value_and_no_variable",
			Input:    EnvMapString{},
			GetFunc:  mockGetEnvFuncForMaps(map[string]string{}, false),
			Expected: nil,
		},
		{
			Name:     "custom_func_error",
			Input:    NewEnvMapStringVariable("SOME_MAP"),
			GetFunc:  mockGetEnvFuncForMaps(map[string]string{}, true),
			ErrorMsg: "mock error",
		},
		{
			Name:     "invalid_map_format",
			Input:    NewEnvMapStringVariable("INVALID_MAP"),
			GetFunc:  mockGetEnvFuncForMaps(map[string]string{"INVALID_MAP": "invalid_format_no_equals"}, false),
			ErrorMsg: ErrParseStringFailed.Error(),
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

func TestEnvMapInt_GetCustom(t *testing.T) {
	testCases := []struct {
		Name     string
		Input    EnvMapInt
		GetFunc  GetEnvFunc
		Expected map[string]int64
		ErrorMsg string
	}{
		{
			Name:     "literal_value",
			Input:    NewEnvMapIntValue(map[string]int64{"key1": 10, "key2": 20}),
			GetFunc:  mockGetEnvFuncForMaps(map[string]string{}, false),
			Expected: map[string]int64{"key1": 10, "key2": 20},
		},
		{
			Name:     "variable_from_custom_func",
			Input:    NewEnvMapIntVariable("CUSTOM_MAP"),
			GetFunc:  mockGetEnvFuncForMaps(map[string]string{"CUSTOM_MAP": "foo=100;bar=200"}, false),
			Expected: map[string]int64{"foo": 100, "bar": 200},
		},
		{
			Name:     "variable_with_fallback_value",
			Input:    NewEnvMapInt("CUSTOM_MAP", map[string]int64{"default": 99}),
			GetFunc:  mockGetEnvFuncForMaps(map[string]string{"CUSTOM_MAP": "custom=42"}, false),
			Expected: map[string]int64{"custom": 42},
		},
		{
			Name:     "empty_variable_uses_fallback",
			Input:    NewEnvMapInt("EMPTY_MAP", map[string]int64{"fallback": 123}),
			GetFunc:  mockGetEnvFuncForMaps(map[string]string{"EMPTY_MAP": ""}, false),
			Expected: map[string]int64{"fallback": 123},
		},
		{
			Name:     "nil_value_and_no_variable",
			Input:    EnvMapInt{},
			GetFunc:  mockGetEnvFuncForMaps(map[string]string{}, false),
			Expected: nil,
		},
		{
			Name:     "custom_func_error",
			Input:    NewEnvMapIntVariable("SOME_MAP"),
			GetFunc:  mockGetEnvFuncForMaps(map[string]string{}, true),
			ErrorMsg: "mock error",
		},
		{
			Name:     "invalid_int_value",
			Input:    NewEnvMapIntVariable("INVALID_MAP"),
			GetFunc:  mockGetEnvFuncForMaps(map[string]string{"INVALID_MAP": "key=not_a_number"}, false),
			ErrorMsg: ErrParseStringFailed.Error(),
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

func TestEnvMapFloat_GetCustom(t *testing.T) {
	testCases := []struct {
		Name     string
		Input    EnvMapFloat
		GetFunc  GetEnvFunc
		Expected map[string]float64
		ErrorMsg string
	}{
		{
			Name:     "literal_value",
			Input:    NewEnvMapFloatValue(map[string]float64{"key1": 1.5, "key2": 2.5}),
			GetFunc:  mockGetEnvFuncForMaps(map[string]string{}, false),
			Expected: map[string]float64{"key1": 1.5, "key2": 2.5},
		},
		{
			Name:     "variable_from_custom_func",
			Input:    NewEnvMapFloatVariable("CUSTOM_MAP"),
			GetFunc:  mockGetEnvFuncForMaps(map[string]string{"CUSTOM_MAP": "foo=3.14;bar=2.718"}, false),
			Expected: map[string]float64{"foo": 3.14, "bar": 2.718},
		},
		{
			Name:     "variable_with_fallback_value",
			Input:    NewEnvMapFloat("CUSTOM_MAP", map[string]float64{"default": 9.99}),
			GetFunc:  mockGetEnvFuncForMaps(map[string]string{"CUSTOM_MAP": "custom=1.23"}, false),
			Expected: map[string]float64{"custom": 1.23},
		},
		{
			Name:     "empty_variable_uses_fallback",
			Input:    NewEnvMapFloat("EMPTY_MAP", map[string]float64{"fallback": 0.5}),
			GetFunc:  mockGetEnvFuncForMaps(map[string]string{"EMPTY_MAP": ""}, false),
			Expected: map[string]float64{"fallback": 0.5},
		},
		{
			Name:     "nil_value_and_no_variable",
			Input:    EnvMapFloat{},
			GetFunc:  mockGetEnvFuncForMaps(map[string]string{}, false),
			Expected: nil,
		},
		{
			Name:     "custom_func_error",
			Input:    NewEnvMapFloatVariable("SOME_MAP"),
			GetFunc:  mockGetEnvFuncForMaps(map[string]string{}, true),
			ErrorMsg: "mock error",
		},
		{
			Name:     "invalid_float_value",
			Input:    NewEnvMapFloatVariable("INVALID_MAP"),
			GetFunc:  mockGetEnvFuncForMaps(map[string]string{"INVALID_MAP": "key=not_a_float"}, false),
			ErrorMsg: ErrParseStringFailed.Error(),
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

func TestEnvMapBool_GetCustom(t *testing.T) {
	testCases := []struct {
		Name     string
		Input    EnvMapBool
		GetFunc  GetEnvFunc
		Expected map[string]bool
		ErrorMsg string
	}{
		{
			Name:     "literal_value",
			Input:    NewEnvMapBoolValue(map[string]bool{"key1": true, "key2": false}),
			GetFunc:  mockGetEnvFuncForMaps(map[string]string{}, false),
			Expected: map[string]bool{"key1": true, "key2": false},
		},
		{
			Name:     "variable_from_custom_func",
			Input:    NewEnvMapBoolVariable("CUSTOM_MAP"),
			GetFunc:  mockGetEnvFuncForMaps(map[string]string{"CUSTOM_MAP": "foo=true;bar=false"}, false),
			Expected: map[string]bool{"foo": true, "bar": false},
		},
		{
			Name:     "variable_with_fallback_value",
			Input:    NewEnvMapBool("CUSTOM_MAP", map[string]bool{"default": true}),
			GetFunc:  mockGetEnvFuncForMaps(map[string]string{"CUSTOM_MAP": "custom=false"}, false),
			Expected: map[string]bool{"custom": false},
		},
		{
			Name:     "empty_variable_uses_fallback",
			Input:    NewEnvMapBool("EMPTY_MAP", map[string]bool{"fallback": true}),
			GetFunc:  mockGetEnvFuncForMaps(map[string]string{"EMPTY_MAP": ""}, false),
			Expected: map[string]bool{"fallback": true},
		},
		{
			Name:     "nil_value_and_no_variable",
			Input:    EnvMapBool{},
			GetFunc:  mockGetEnvFuncForMaps(map[string]string{}, false),
			Expected: nil,
		},
		{
			Name:     "custom_func_error",
			Input:    NewEnvMapBoolVariable("SOME_MAP"),
			GetFunc:  mockGetEnvFuncForMaps(map[string]string{}, true),
			ErrorMsg: "mock error",
		},
		{
			Name:     "invalid_bool_value",
			Input:    NewEnvMapBoolVariable("INVALID_MAP"),
			GetFunc:  mockGetEnvFuncForMaps(map[string]string{"INVALID_MAP": "key=not_a_bool"}, false),
			ErrorMsg: ErrParseStringFailed.Error(),
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

func TestEnvMapBool(t *testing.T) {
	t.Setenv("SOME_FOO", "foo=true;bar=false")
	testCases := []struct {
		Input    EnvMapBool
		Expected map[string]bool
		ErrorMsg string
	}{
		{
			Input: NewEnvMapBoolValue(map[string]bool{
				"foo": true,
			}),
			Expected: map[string]bool{
				"foo": true,
			},
		},
		{
			Input: NewEnvMapBoolVariable("SOME_FOO"),
			Expected: map[string]bool{
				"foo": true,
				"bar": false,
			},
		},
		{
			Input:    EnvMapBool{},
			Expected: nil,
		},
		{
			Input:    NewEnvMapBool("SOME_FOO_2", map[string]bool{}),
			Expected: map[string]bool{},
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
				assertDeepEqual(t, tc.Input.IsZero(), tc.Expected == nil)
			}
		})
	}

	t.Run("json_decode", func(t *testing.T) {
		var ev EnvMapBool
		assertNilError(t, json.Unmarshal([]byte(`{"env": "SOME_FOO"}`), &ev))
		result, err := ev.Get()
		assertNilError(t, err)
		assertDeepEqual(t, map[string]bool{
			"foo": true,
			"bar": false,
		}, result)
	})
}

func TestEnvMapInt(t *testing.T) {
	t.Setenv("SOME_FOO", "foo=2;bar=3")
	testCases := []struct {
		Input    EnvMapInt
		Expected map[string]int64
		ErrorMsg string
	}{
		{
			Input: NewEnvMapIntValue(map[string]int64{
				"foo": 1,
			}),
			Expected: map[string]int64{
				"foo": 1,
			},
		},
		{
			Input: NewEnvMapIntVariable("SOME_FOO"),
			Expected: map[string]int64{
				"foo": 2,
				"bar": 3,
			},
		},
		{
			Input:    EnvMapInt{},
			Expected: nil,
		},
		{
			Input:    NewEnvMapInt("SOME_FOO_2", map[string]int64{}),
			Expected: map[string]int64{},
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
				assertDeepEqual(t, tc.Input.IsZero(), tc.Expected == nil)
			}
		})
	}

	t.Run("json_decode", func(t *testing.T) {
		var ev EnvMapInt
		assertNilError(t, json.Unmarshal([]byte(`{"env": "SOME_FOO"}`), &ev))
		result, err := ev.Get()
		assertNilError(t, err)
		assertDeepEqual(t, map[string]int64{
			"foo": 2,
			"bar": 3,
		}, result)
	})
}

func TestEnvMapFloat(t *testing.T) {
	t.Setenv("SOME_FOO", "foo=2.2;bar=3.3")
	testCases := []struct {
		Input    EnvMapFloat
		Expected map[string]float64
		ErrorMsg string
	}{
		{
			Input: NewEnvMapFloatValue(map[string]float64{
				"foo": 1.1,
			}),
			Expected: map[string]float64{
				"foo": 1.1,
			},
		},
		{
			Input: NewEnvMapFloatVariable("SOME_FOO"),
			Expected: map[string]float64{
				"foo": 2.2,
				"bar": 3.3,
			},
		},
		{
			Input:    EnvMapFloat{},
			Expected: nil,
		},
		{
			Input:    NewEnvMapFloat("SOME_FOO_2", map[string]float64{}),
			Expected: map[string]float64{},
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
				assertDeepEqual(t, tc.Input.IsZero(), tc.Expected == nil)
			}
		})
	}

	t.Run("json_decode", func(t *testing.T) {
		var ev EnvMapFloat
		assertNilError(t, json.Unmarshal([]byte(`{"env": "SOME_FOO"}`), &ev))
		result, err := ev.Get()
		assertNilError(t, err)
		assertDeepEqual(t, map[string]float64{
			"foo": 2.2,
			"bar": 3.3,
		}, result)
	})
}

func TestEnvMapString(t *testing.T) {
	t.Setenv("SOME_FOO", "foo=2.2;bar=3.3")

	testCases := []struct {
		Input    EnvMapString
		Expected map[string]string
		ErrorMsg string
	}{
		{
			Input: NewEnvMapStringValue(map[string]string{
				"foo": "1.1",
			}),
			Expected: map[string]string{
				"foo": "1.1",
			},
		},
		{
			Input: NewEnvMapStringVariable("SOME_FOO"),
			Expected: map[string]string{
				"foo": "2.2",
				"bar": "3.3",
			},
		},
		{
			Input:    EnvMapString{},
			Expected: nil,
		},
		{
			Input:    NewEnvMapString("SOME_FOO_2", map[string]string{}),
			Expected: map[string]string{},
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
				assertDeepEqual(t, tc.Input.IsZero(), tc.Expected == nil)
			}
		})
	}

	t.Run("json_decode", func(t *testing.T) {
		var ev EnvMapString
		assertNilError(t, json.Unmarshal([]byte(`{"env": "SOME_FOO"}`), &ev))
		result, err := ev.Get()
		assertNilError(t, err)
		assertDeepEqual(t, map[string]string{
			"foo": "2.2",
			"bar": "3.3",
		}, result)
	})
}

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
