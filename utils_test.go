package goenvconf

import (
	"testing"
)

func TestParseIntMapFromString(t *testing.T) {
	testCases := []struct {
		Input    string
		Expected map[string]int
		ErrorMsg string
	}{
		{
			Expected: map[string]int{},
		},
		{
			Input: "a=1;b=2;c=3",
			Expected: map[string]int{
				"a": 1,
				"b": 2,
				"c": 3,
			},
		},
		{
			Input:    "a;b=2",
			ErrorMsg: "ParseStringFailed: invalid int map from string, expected: <key1>=<value1>;<key2>=<value2>, got: a;b=2",
		},
		{
			Input:    "a=c;b=2",
			ErrorMsg: "invalid integer value c in item a",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Input, func(t *testing.T) {
			result, err := ParseIntMapFromString(tc.Input)
			if tc.ErrorMsg != "" {
				assertErrorContains(t, err, tc.ErrorMsg)
			} else {
				assertNilError(t, err)
				assertDeepEqual(t, result, tc.Expected)
			}
		})
	}
}
