package goenvconf

import (
	"fmt"
	"testing"
)

// TestEnvStringSlice tests the EnvStringSlice Get method
func TestEnvStringSlice(t *testing.T) {
	t.Setenv("STRING_SLICE_VAR", "foo,bar,baz")
	t.Setenv("EMPTY_STRING_SLICE", "")

	testCases := []struct {
		Input    EnvStringSlice
		Expected []string
		ErrorMsg string
	}{
		{
			Input:    NewEnvStringSliceValue([]string{"foo", "bar"}),
			Expected: []string{"foo", "bar"},
		},
		{
			Input:    NewEnvStringSliceVariable("STRING_SLICE_VAR"),
			Expected: []string{"foo", "bar", "baz"},
		},
		{
			Input:    EnvStringSlice{},
			ErrorMsg: ErrEnvironmentValueRequired.Error(),
		},
		{
			Input:    NewEnvStringSlice("SOME_MISSING_VAR", []string{"fallback"}),
			Expected: []string{"fallback"},
		},
		{
			Input: EnvStringSlice{
				Variable: toPtr(""),
			},
			ErrorMsg: ErrEnvironmentValueRequired.Error(),
		},
		{
			Input:    NewEnvStringSliceVariable("EMPTY_STRING_SLICE"),
			Expected: []string{},
		},
		{
			Input:    NewEnvStringSliceVariable("MISSING_VAR"),
			ErrorMsg: ErrEnvironmentVariableValueRequired.Error(),
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			result, err := tc.Input.Get()
			if tc.ErrorMsg != "" {
				assertErrorContains(t, err, tc.ErrorMsg)
			} else {
				assertNilError(t, err)
				assertDeepEqual(t, tc.Expected, result)
			}
		})
	}
}

// TestEnvStringSlice_GetCustom tests the EnvStringSlice GetCustom method
func TestEnvStringSlice_GetCustom(t *testing.T) {
	testCases := []struct {
		Name     string
		Input    EnvStringSlice
		GetFunc  GetEnvFunc
		Expected []string
		ErrorMsg string
	}{
		{
			Name:     "literal value",
			Input:    NewEnvStringSliceValue([]string{"foo", "bar"}),
			GetFunc:  mockGetEnvFunc(map[string]string{}, false),
			Expected: []string{"foo", "bar"},
		},
		{
			Name:     "variable from custom func",
			Input:    NewEnvStringSliceVariable("STRING_SLICE_VAR"),
			GetFunc:  mockGetEnvFunc(map[string]string{"STRING_SLICE_VAR": "a,b,c"}, false),
			Expected: []string{"a", "b", "c"},
		},
		{
			Name:     "variable with fallback value",
			Input:    NewEnvStringSlice("MISSING_VAR", []string{"fallback"}),
			GetFunc:  mockGetEnvFunc(map[string]string{}, false),
			Expected: []string{"fallback"},
		},
		{
			Name:     "empty variable uses fallback",
			Input:    NewEnvStringSlice("EMPTY_VAR", []string{"fallback"}),
			GetFunc:  mockGetEnvFunc(map[string]string{"EMPTY_VAR": ""}, false),
			Expected: []string{"fallback"},
		},
		{
			Name:     "nil value and no variable",
			Input:    EnvStringSlice{},
			GetFunc:  mockGetEnvFunc(map[string]string{}, false),
			ErrorMsg: ErrEnvironmentValueRequired.Error(),
		},
		{
			Name:     "custom func error",
			Input:    NewEnvStringSliceVariable("ERROR_VAR"),
			GetFunc:  mockGetEnvFunc(map[string]string{}, true),
			ErrorMsg: "mock error",
		},
		{
			Name:     "missing variable returns error",
			Input:    NewEnvStringSliceVariable("MISSING_VAR"),
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

// TestEnvStringSlice_IsZero tests the EnvStringSlice IsZero method
func TestEnvStringSlice_IsZero(t *testing.T) {
	testCases := []struct {
		Name     string
		Input    EnvStringSlice
		Expected bool
	}{
		{
			Name:     "empty struct",
			Input:    EnvStringSlice{},
			Expected: true,
		},
		{
			Name:     "with value",
			Input:    NewEnvStringSliceValue([]string{"foo"}),
			Expected: false,
		},
		{
			Name:     "with variable",
			Input:    NewEnvStringSliceVariable("VAR"),
			Expected: false,
		},
		{
			Name:     "with empty variable",
			Input:    EnvStringSlice{Variable: toPtr("")},
			Expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			result := tc.Input.IsZero()
			if result != tc.Expected {
				t.Errorf("Expected %v, got %v", tc.Expected, result)
			}
		})
	}
}

// TestEnvIntSlice tests the EnvIntSlice Get method
func TestEnvIntSlice(t *testing.T) {
	t.Setenv("INT_SLICE_VAR", "1,2,3")
	t.Setenv("EMPTY_INT_SLICE", "")
	t.Setenv("INVALID_INT_SLICE", "1,abc,3")

	testCases := []struct {
		Input    EnvIntSlice
		Expected []int64
		ErrorMsg string
	}{
		{
			Input:    NewEnvIntSliceValue([]int64{10, 20, 30}),
			Expected: []int64{10, 20, 30},
		},
		{
			Input:    NewEnvIntSliceVariable("INT_SLICE_VAR"),
			Expected: []int64{1, 2, 3},
		},
		{
			Input:    EnvIntSlice{},
			ErrorMsg: ErrEnvironmentValueRequired.Error(),
		},
		{
			Input:    NewEnvIntSlice("SOME_MISSING_VAR", []int64{100}),
			Expected: []int64{100},
		},
		{
			Input: EnvIntSlice{
				Variable: toPtr(""),
			},
			ErrorMsg: ErrEnvironmentValueRequired.Error(),
		},
		{
			Input:    NewEnvIntSliceVariable("EMPTY_INT_SLICE"),
			Expected: []int64{},
		},
		{
			Input:    NewEnvIntSliceVariable("MISSING_VAR"),
			ErrorMsg: ErrEnvironmentVariableValueRequired.Error(),
		},
		{
			Input:    NewEnvIntSliceVariable("INVALID_INT_SLICE"),
			ErrorMsg: "failed to convert INVALID_INT_SLICE variable to integers",
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			result, err := tc.Input.Get()
			if tc.ErrorMsg != "" {
				assertErrorContains(t, err, tc.ErrorMsg)
			} else {
				assertNilError(t, err)
				assertDeepEqual(t, tc.Expected, result)
			}
		})
	}
}

// TestEnvIntSlice_GetCustom tests the EnvIntSlice GetCustom method
func TestEnvIntSlice_GetCustom(t *testing.T) {
	testCases := []struct {
		Name     string
		Input    EnvIntSlice
		GetFunc  GetEnvFunc
		Expected []int64
		ErrorMsg string
	}{
		{
			Name:     "literal value",
			Input:    NewEnvIntSliceValue([]int64{10, 20}),
			GetFunc:  mockGetEnvFunc(map[string]string{}, false),
			Expected: []int64{10, 20},
		},
		{
			Name:     "variable from custom func",
			Input:    NewEnvIntSliceVariable("INT_SLICE_VAR"),
			GetFunc:  mockGetEnvFunc(map[string]string{"INT_SLICE_VAR": "5,10,15"}, false),
			Expected: []int64{5, 10, 15},
		},
		{
			Name:     "variable with fallback value",
			Input:    NewEnvIntSlice("MISSING_VAR", []int64{99}),
			GetFunc:  mockGetEnvFunc(map[string]string{}, false),
			Expected: []int64{99},
		},
		{
			Name:     "empty variable uses fallback",
			Input:    NewEnvIntSlice("EMPTY_VAR", []int64{42}),
			GetFunc:  mockGetEnvFunc(map[string]string{"EMPTY_VAR": ""}, false),
			Expected: []int64{42},
		},
		{
			Name:     "nil value and no variable",
			Input:    EnvIntSlice{},
			GetFunc:  mockGetEnvFunc(map[string]string{}, false),
			ErrorMsg: ErrEnvironmentValueRequired.Error(),
		},
		{
			Name:     "custom func error",
			Input:    NewEnvIntSliceVariable("ERROR_VAR"),
			GetFunc:  mockGetEnvFunc(map[string]string{}, true),
			ErrorMsg: "mock error",
		},
		{
			Name:     "invalid int format",
			Input:    NewEnvIntSliceVariable("INVALID_VAR"),
			GetFunc:  mockGetEnvFunc(map[string]string{"INVALID_VAR": "1,abc,3"}, false),
			ErrorMsg: "failed to convert INVALID_VAR variable to integers",
		},
		{
			Name:     "missing variable no fallback",
			Input:    NewEnvIntSliceVariable("MISSING_VAR"),
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

// TestEnvIntSlice_IsZero tests the EnvIntSlice IsZero method
func TestEnvIntSlice_IsZero(t *testing.T) {
	testCases := []struct {
		Name     string
		Input    EnvIntSlice
		Expected bool
	}{
		{
			Name:     "empty struct",
			Input:    EnvIntSlice{},
			Expected: true,
		},
		{
			Name:     "with value",
			Input:    NewEnvIntSliceValue([]int64{1}),
			Expected: false,
		},
		{
			Name:     "with variable",
			Input:    NewEnvIntSliceVariable("VAR"),
			Expected: false,
		},
		{
			Name:     "with empty variable",
			Input:    EnvIntSlice{Variable: toPtr("")},
			Expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			result := tc.Input.IsZero()
			if result != tc.Expected {
				t.Errorf("Expected %v, got %v", tc.Expected, result)
			}
		})
	}
}

// TestEnvFloatSlice tests the EnvFloatSlice Get method
func TestEnvFloatSlice(t *testing.T) {
	t.Setenv("FLOAT_SLICE_VAR", "1.5,2.5,3.5")
	t.Setenv("EMPTY_FLOAT_SLICE", "")
	t.Setenv("INVALID_FLOAT_SLICE", "1.5,abc,3.5")

	testCases := []struct {
		Input    EnvFloatSlice
		Expected []float64
		ErrorMsg string
	}{
		{
			Input:    NewEnvFloatSliceValue([]float64{10.5, 20.5, 30.5}),
			Expected: []float64{10.5, 20.5, 30.5},
		},
		{
			Input:    NewEnvFloatSliceVariable("FLOAT_SLICE_VAR"),
			Expected: []float64{1.5, 2.5, 3.5},
		},
		{
			Input:    EnvFloatSlice{},
			ErrorMsg: ErrEnvironmentValueRequired.Error(),
		},
		{
			Input:    NewEnvFloatSlice("SOME_MISSING_VAR", []float64{100.5}),
			Expected: []float64{100.5},
		},
		{
			Input: EnvFloatSlice{
				Variable: toPtr(""),
			},
			ErrorMsg: ErrEnvironmentValueRequired.Error(),
		},
		{
			Input:    NewEnvFloatSliceVariable("EMPTY_FLOAT_SLICE"),
			Expected: []float64{},
		},
		{
			Input:    NewEnvFloatSliceVariable("MISSING_VAR"),
			ErrorMsg: ErrEnvironmentVariableValueRequired.Error(),
		},
		{
			Input:    NewEnvFloatSliceVariable("INVALID_FLOAT_SLICE"),
			ErrorMsg: "failed to convert INVALID_FLOAT_SLICE variable to integers",
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			result, err := tc.Input.Get()
			if tc.ErrorMsg != "" {
				assertErrorContains(t, err, tc.ErrorMsg)
			} else {
				assertNilError(t, err)
				assertDeepEqual(t, tc.Expected, result)
			}
		})
	}
}

// TestEnvFloatSlice_GetCustom tests the EnvFloatSlice GetCustom method
func TestEnvFloatSlice_GetCustom(t *testing.T) {
	testCases := []struct {
		Name     string
		Input    EnvFloatSlice
		GetFunc  GetEnvFunc
		Expected []float64
		ErrorMsg string
	}{
		{
			Name:     "literal value",
			Input:    NewEnvFloatSliceValue([]float64{10.5, 20.5}),
			GetFunc:  mockGetEnvFunc(map[string]string{}, false),
			Expected: []float64{10.5, 20.5},
		},
		{
			Name:     "variable from custom func",
			Input:    NewEnvFloatSliceVariable("FLOAT_SLICE_VAR"),
			GetFunc:  mockGetEnvFunc(map[string]string{"FLOAT_SLICE_VAR": "5.5,10.5,15.5"}, false),
			Expected: []float64{5.5, 10.5, 15.5},
		},
		{
			Name:     "variable with fallback value",
			Input:    NewEnvFloatSlice("MISSING_VAR", []float64{99.9}),
			GetFunc:  mockGetEnvFunc(map[string]string{}, false),
			Expected: []float64{99.9},
		},
		{
			Name:     "empty variable uses fallback",
			Input:    NewEnvFloatSlice("EMPTY_VAR", []float64{42.5}),
			GetFunc:  mockGetEnvFunc(map[string]string{"EMPTY_VAR": ""}, false),
			Expected: []float64{42.5},
		},
		{
			Name:     "nil value and no variable",
			Input:    EnvFloatSlice{},
			GetFunc:  mockGetEnvFunc(map[string]string{}, false),
			ErrorMsg: ErrEnvironmentValueRequired.Error(),
		},
		{
			Name:     "custom func error",
			Input:    NewEnvFloatSliceVariable("ERROR_VAR"),
			GetFunc:  mockGetEnvFunc(map[string]string{}, true),
			ErrorMsg: "mock error",
		},
		{
			Name:     "invalid float format",
			Input:    NewEnvFloatSliceVariable("INVALID_VAR"),
			GetFunc:  mockGetEnvFunc(map[string]string{"INVALID_VAR": "1.5,abc,3.5"}, false),
			ErrorMsg: "failed to convert INVALID_VAR variable to integers",
		},
		{
			Name:     "missing variable no fallback",
			Input:    NewEnvFloatSliceVariable("MISSING_VAR"),
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

// TestEnvFloatSlice_IsZero tests the EnvFloatSlice IsZero method
func TestEnvFloatSlice_IsZero(t *testing.T) {
	testCases := []struct {
		Name     string
		Input    EnvFloatSlice
		Expected bool
	}{
		{
			Name:     "empty struct",
			Input:    EnvFloatSlice{},
			Expected: true,
		},
		{
			Name:     "with value",
			Input:    NewEnvFloatSliceValue([]float64{1.5}),
			Expected: false,
		},
		{
			Name:     "with variable",
			Input:    NewEnvFloatSliceVariable("VAR"),
			Expected: false,
		},
		{
			Name:     "with empty variable",
			Input:    EnvFloatSlice{Variable: toPtr("")},
			Expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			result := tc.Input.IsZero()
			if result != tc.Expected {
				t.Errorf("Expected %v, got %v", tc.Expected, result)
			}
		})
	}
}

// TestEnvBoolSlice tests the EnvBoolSlice Get method
func TestEnvBoolSlice(t *testing.T) {
	t.Setenv("BOOL_SLICE_VAR", "true,false,true")
	t.Setenv("EMPTY_BOOL_SLICE", "")
	t.Setenv("INVALID_BOOL_SLICE", "true,invalid,false")

	testCases := []struct {
		Input    EnvBoolSlice
		Expected []bool
		ErrorMsg string
	}{
		{
			Input:    NewEnvBoolSliceValue([]bool{true, false, true}),
			Expected: []bool{true, false, true},
		},
		{
			Input:    NewEnvBoolSliceVariable("BOOL_SLICE_VAR"),
			Expected: []bool{true, false, true},
		},
		{
			Input:    EnvBoolSlice{},
			ErrorMsg: ErrEnvironmentValueRequired.Error(),
		},
		{
			Input:    NewEnvBoolSlice("SOME_MISSING_VAR", []bool{true}),
			Expected: []bool{true},
		},
		{
			Input: EnvBoolSlice{
				Variable: toPtr(""),
			},
			ErrorMsg: ErrEnvironmentValueRequired.Error(),
		},
		{
			Input:    NewEnvBoolSliceVariable("EMPTY_BOOL_SLICE"),
			Expected: []bool{},
		},
		{
			Input:    NewEnvBoolSliceVariable("MISSING_VAR"),
			ErrorMsg: ErrEnvironmentVariableValueRequired.Error(),
		},
		{
			Input:    NewEnvBoolSliceVariable("INVALID_BOOL_SLICE"),
			ErrorMsg: "failed to convert INVALID_BOOL_SLICE variable to integers",
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			result, err := tc.Input.Get()
			if tc.ErrorMsg != "" {
				assertErrorContains(t, err, tc.ErrorMsg)
			} else {
				assertNilError(t, err)
				assertDeepEqual(t, tc.Expected, result)
			}
		})
	}
}

// TestEnvBoolSlice_GetCustom tests the EnvBoolSlice GetCustom method
func TestEnvBoolSlice_GetCustom(t *testing.T) {
	testCases := []struct {
		Name     string
		Input    EnvBoolSlice
		GetFunc  GetEnvFunc
		Expected []bool
		ErrorMsg string
	}{
		{
			Name:     "literal value",
			Input:    NewEnvBoolSliceValue([]bool{true, false}),
			GetFunc:  mockGetEnvFunc(map[string]string{}, false),
			Expected: []bool{true, false},
		},
		{
			Name:     "variable from custom func",
			Input:    NewEnvBoolSliceVariable("BOOL_SLICE_VAR"),
			GetFunc:  mockGetEnvFunc(map[string]string{"BOOL_SLICE_VAR": "true,false,true"}, false),
			Expected: []bool{true, false, true},
		},
		{
			Name:     "variable with fallback value",
			Input:    NewEnvBoolSlice("MISSING_VAR", []bool{false}),
			GetFunc:  mockGetEnvFunc(map[string]string{}, false),
			Expected: []bool{false},
		},
		{
			Name:     "empty variable uses fallback",
			Input:    NewEnvBoolSlice("EMPTY_VAR", []bool{true}),
			GetFunc:  mockGetEnvFunc(map[string]string{"EMPTY_VAR": ""}, false),
			Expected: []bool{true},
		},
		{
			Name:     "nil value and no variable",
			Input:    EnvBoolSlice{},
			GetFunc:  mockGetEnvFunc(map[string]string{}, false),
			ErrorMsg: ErrEnvironmentValueRequired.Error(),
		},
		{
			Name:     "custom func error",
			Input:    NewEnvBoolSliceVariable("ERROR_VAR"),
			GetFunc:  mockGetEnvFunc(map[string]string{}, true),
			ErrorMsg: "mock error",
		},
		{
			Name:     "invalid bool format",
			Input:    NewEnvBoolSliceVariable("INVALID_VAR"),
			GetFunc:  mockGetEnvFunc(map[string]string{"INVALID_VAR": "true,invalid,false"}, false),
			ErrorMsg: "failed to convert INVALID_VAR variable to integers",
		},
		{
			Name:     "missing variable no fallback",
			Input:    NewEnvBoolSliceVariable("MISSING_VAR"),
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

// TestEnvBoolSlice_IsZero tests the EnvBoolSlice IsZero method
func TestEnvBoolSlice_IsZero(t *testing.T) {
	testCases := []struct {
		Name     string
		Input    EnvBoolSlice
		Expected bool
	}{
		{
			Name:     "empty struct",
			Input:    EnvBoolSlice{},
			Expected: true,
		},
		{
			Name:     "with value",
			Input:    NewEnvBoolSliceValue([]bool{true}),
			Expected: false,
		},
		{
			Name:     "with variable",
			Input:    NewEnvBoolSliceVariable("VAR"),
			Expected: false,
		},
		{
			Name:     "with empty variable",
			Input:    EnvBoolSlice{Variable: toPtr("")},
			Expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			result := tc.Input.IsZero()
			if result != tc.Expected {
				t.Errorf("Expected %v, got %v", tc.Expected, result)
			}
		})
	}
}

// TestEnvStringSlice_Equal tests the EnvStringSlice Equal method
func TestEnvStringSlice_Equal(t *testing.T) {
	testCases := []struct {
		Name     string
		A        EnvStringSlice
		B        EnvStringSlice
		Expected bool
	}{
		{
			Name:     "both nil",
			A:        EnvStringSlice{},
			B:        EnvStringSlice{},
			Expected: true,
		},
		{
			Name:     "same literal values",
			A:        NewEnvStringSliceValue([]string{"foo", "bar"}),
			B:        NewEnvStringSliceValue([]string{"foo", "bar"}),
			Expected: true,
		},
		{
			Name:     "different literal values",
			A:        NewEnvStringSliceValue([]string{"foo", "bar"}),
			B:        NewEnvStringSliceValue([]string{"baz", "qux"}),
			Expected: false,
		},
		{
			Name:     "different order",
			A:        NewEnvStringSliceValue([]string{"foo", "bar"}),
			B:        NewEnvStringSliceValue([]string{"bar", "foo"}),
			Expected: false,
		},
		{
			Name:     "same variable names",
			A:        NewEnvStringSliceVariable("VAR1"),
			B:        NewEnvStringSliceVariable("VAR1"),
			Expected: true,
		},
		{
			Name:     "different variable names",
			A:        NewEnvStringSliceVariable("VAR1"),
			B:        NewEnvStringSliceVariable("VAR2"),
			Expected: false,
		},
		{
			Name:     "same value and variable",
			A:        NewEnvStringSlice("VAR1", []string{"foo"}),
			B:        NewEnvStringSlice("VAR1", []string{"foo"}),
			Expected: true,
		},
		{
			Name:     "value vs variable",
			A:        NewEnvStringSliceValue([]string{"foo"}),
			B:        NewEnvStringSliceVariable("VAR1"),
			Expected: false,
		},
		{
			Name:     "nil vs empty slice",
			A:        EnvStringSlice{Value: nil},
			B:        NewEnvStringSliceValue([]string{}),
			Expected: true, // slices.Equal considers nil and empty slices as equal
		},
		{
			Name:     "both empty slices",
			A:        NewEnvStringSliceValue([]string{}),
			B:        NewEnvStringSliceValue([]string{}),
			Expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			result := tc.A.Equal(tc.B)
			if result != tc.Expected {
				t.Errorf("Expected %v, got %v", tc.Expected, result)
			}
		})
	}
}

// TestEnvIntSlice_Equal tests the EnvIntSlice Equal method
func TestEnvIntSlice_Equal(t *testing.T) {
	testCases := []struct {
		Name     string
		A        EnvIntSlice
		B        EnvIntSlice
		Expected bool
	}{
		{
			Name:     "both nil",
			A:        EnvIntSlice{},
			B:        EnvIntSlice{},
			Expected: true,
		},
		{
			Name:     "same literal values",
			A:        NewEnvIntSliceValue([]int64{1, 2, 3}),
			B:        NewEnvIntSliceValue([]int64{1, 2, 3}),
			Expected: true,
		},
		{
			Name:     "different literal values",
			A:        NewEnvIntSliceValue([]int64{1, 2, 3}),
			B:        NewEnvIntSliceValue([]int64{4, 5, 6}),
			Expected: false,
		},
		{
			Name:     "different order",
			A:        NewEnvIntSliceValue([]int64{1, 2, 3}),
			B:        NewEnvIntSliceValue([]int64{3, 2, 1}),
			Expected: false,
		},
		{
			Name:     "same variable names",
			A:        NewEnvIntSliceVariable("VAR1"),
			B:        NewEnvIntSliceVariable("VAR1"),
			Expected: true,
		},
		{
			Name:     "different variable names",
			A:        NewEnvIntSliceVariable("VAR1"),
			B:        NewEnvIntSliceVariable("VAR2"),
			Expected: false,
		},
		{
			Name:     "same value and variable",
			A:        NewEnvIntSlice("VAR1", []int64{10}),
			B:        NewEnvIntSlice("VAR1", []int64{10}),
			Expected: true,
		},
		{
			Name:     "value vs variable",
			A:        NewEnvIntSliceValue([]int64{10}),
			B:        NewEnvIntSliceVariable("VAR1"),
			Expected: false,
		},
		{
			Name:     "nil vs empty slice",
			A:        EnvIntSlice{Value: nil},
			B:        NewEnvIntSliceValue([]int64{}),
			Expected: true, // slices.Equal considers nil and empty slices as equal
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			result := tc.A.Equal(tc.B)
			if result != tc.Expected {
				t.Errorf("Expected %v, got %v", tc.Expected, result)
			}
		})
	}
}

// TestEnvFloatSlice_Equal tests the EnvFloatSlice Equal method
func TestEnvFloatSlice_Equal(t *testing.T) {
	testCases := []struct {
		Name     string
		A        EnvFloatSlice
		B        EnvFloatSlice
		Expected bool
	}{
		{
			Name:     "both nil",
			A:        EnvFloatSlice{},
			B:        EnvFloatSlice{},
			Expected: true,
		},
		{
			Name:     "same literal values",
			A:        NewEnvFloatSliceValue([]float64{1.5, 2.5, 3.5}),
			B:        NewEnvFloatSliceValue([]float64{1.5, 2.5, 3.5}),
			Expected: true,
		},
		{
			Name:     "different literal values",
			A:        NewEnvFloatSliceValue([]float64{1.5, 2.5, 3.5}),
			B:        NewEnvFloatSliceValue([]float64{4.5, 5.5, 6.5}),
			Expected: false,
		},
		{
			Name:     "different order",
			A:        NewEnvFloatSliceValue([]float64{1.5, 2.5, 3.5}),
			B:        NewEnvFloatSliceValue([]float64{3.5, 2.5, 1.5}),
			Expected: false,
		},
		{
			Name:     "same variable names",
			A:        NewEnvFloatSliceVariable("VAR1"),
			B:        NewEnvFloatSliceVariable("VAR1"),
			Expected: true,
		},
		{
			Name:     "different variable names",
			A:        NewEnvFloatSliceVariable("VAR1"),
			B:        NewEnvFloatSliceVariable("VAR2"),
			Expected: false,
		},
		{
			Name:     "same value and variable",
			A:        NewEnvFloatSlice("VAR1", []float64{10.5}),
			B:        NewEnvFloatSlice("VAR1", []float64{10.5}),
			Expected: true,
		},
		{
			Name:     "value vs variable",
			A:        NewEnvFloatSliceValue([]float64{10.5}),
			B:        NewEnvFloatSliceVariable("VAR1"),
			Expected: false,
		},
		{
			Name:     "nil vs empty slice",
			A:        EnvFloatSlice{Value: nil},
			B:        NewEnvFloatSliceValue([]float64{}),
			Expected: true, // slices.Equal considers nil and empty slices as equal
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			result := tc.A.Equal(tc.B)
			if result != tc.Expected {
				t.Errorf("Expected %v, got %v", tc.Expected, result)
			}
		})
	}
}

// TestEnvBoolSlice_Equal tests the EnvBoolSlice Equal method
func TestEnvBoolSlice_Equal(t *testing.T) {
	testCases := []struct {
		Name     string
		A        EnvBoolSlice
		B        EnvBoolSlice
		Expected bool
	}{
		{
			Name:     "both nil",
			A:        EnvBoolSlice{},
			B:        EnvBoolSlice{},
			Expected: true,
		},
		{
			Name:     "same literal values",
			A:        NewEnvBoolSliceValue([]bool{true, false, true}),
			B:        NewEnvBoolSliceValue([]bool{true, false, true}),
			Expected: true,
		},
		{
			Name:     "different literal values",
			A:        NewEnvBoolSliceValue([]bool{true, false, true}),
			B:        NewEnvBoolSliceValue([]bool{false, true, false}),
			Expected: false,
		},
		{
			Name:     "different order",
			A:        NewEnvBoolSliceValue([]bool{true, false}),
			B:        NewEnvBoolSliceValue([]bool{false, true}),
			Expected: false,
		},
		{
			Name:     "same variable names",
			A:        NewEnvBoolSliceVariable("VAR1"),
			B:        NewEnvBoolSliceVariable("VAR1"),
			Expected: true,
		},
		{
			Name:     "different variable names",
			A:        NewEnvBoolSliceVariable("VAR1"),
			B:        NewEnvBoolSliceVariable("VAR2"),
			Expected: false,
		},
		{
			Name:     "same value and variable",
			A:        NewEnvBoolSlice("VAR1", []bool{true}),
			B:        NewEnvBoolSlice("VAR1", []bool{true}),
			Expected: true,
		},
		{
			Name:     "value vs variable",
			A:        NewEnvBoolSliceValue([]bool{true}),
			B:        NewEnvBoolSliceVariable("VAR1"),
			Expected: false,
		},
		{
			Name:     "nil vs empty slice",
			A:        EnvBoolSlice{Value: nil},
			B:        NewEnvBoolSliceValue([]bool{}),
			Expected: true, // slices.Equal considers nil and empty slices as equal
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			result := tc.A.Equal(tc.B)
			if result != tc.Expected {
				t.Errorf("Expected %v, got %v", tc.Expected, result)
			}
		})
	}
}
