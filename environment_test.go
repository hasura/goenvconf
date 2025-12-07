package goenvconf

import (
	"encoding/json"
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
