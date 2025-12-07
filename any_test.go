package goenvconf

import (
	"encoding/json"
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

	t.Run("json_decode", func(t *testing.T) {
		var ev EnvAny
		assertNilError(t, json.Unmarshal([]byte(`{"env": "SOME_FOO"}`), &ev))
	})
}
