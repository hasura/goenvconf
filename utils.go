package goenvconf

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	keyValueLength = 2
)

var (
	errEnvironmentValueRequired = errors.New("require either value or env")
	// ErrEnvironmentVariableRequired the error happens when the name of environment variable is empty.
	ErrEnvironmentVariableRequired = errors.New("the environment variable name is empty")
	// ErrEnvironmentVariableValueRequired the error happens when the value from environment variable is empty.
	ErrEnvironmentVariableValueRequired = errors.New("the environment variable value is empty")
	// ErrParseStringFailed is the error when failed to parse a string to another type.
	ErrParseStringFailed = errors.New("ParseStringFailed")
)

// ParseStringMapFromString parses a string map from a string with format:
//
//	<key1>=<value1>;<key2>=<value2>
func ParseStringMapFromString(input string) (map[string]string, error) {
	result := make(map[string]string)
	if input == "" {
		return result, nil
	}

	rawItems := strings.SplitSeq(input, ";")

	for rawItem := range rawItems {
		keyValue := strings.Split(rawItem, "=")

		if len(keyValue) != keyValueLength {
			return nil, fmt.Errorf(
				"%w: invalid int map from string, expected: <key1>=<value1>;<key2>=<value2>, got: %s",
				ErrParseStringFailed,
				input,
			)
		}

		result[keyValue[0]] = keyValue[1]
	}

	return result, nil
}

// ParseIntMapFromString parses an integer map from a string with format:
//
//	<key1>=<value1>;<key2>=<value2>
func ParseIntMapFromString(input string) (map[string]int, error) {
	return ParseIntegerMapFromString[int](input)
}

// ParseIntegerMapFromString parses an integer map from a string with format:
//
//	<key1>=<value1>;<key2>=<value2>
func ParseIntegerMapFromString[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64](
	input string,
) (map[string]T, error) {
	rawValues, err := ParseStringMapFromString(input)
	if err != nil {
		return nil, err
	}

	result := make(map[string]T)

	for key, value := range rawValues {
		intValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, fmt.Errorf(
				"%w: invalid integer value %s in item %s",
				ErrParseStringFailed,
				value,
				key,
			)
		}

		result[key] = T(intValue)
	}

	return result, nil
}

// ParseFloatMapFromString parses a float map from a string with format:
//
//	<key1>=<value1>;<key2>=<value2>
func ParseFloatMapFromString[T float32 | float64](input string) (map[string]T, error) {
	rawValues, err := ParseStringMapFromString(input)
	if err != nil {
		return nil, err
	}

	result := make(map[string]T)

	for key, value := range rawValues {
		floatValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, fmt.Errorf(
				"%w: invalid float value %s in item %s",
				ErrParseStringFailed,
				value,
				key,
			)
		}

		result[key] = T(floatValue)
	}

	return result, nil
}

// ParseBoolMapFromString parses a bool map from a string with format:
//
//	<key1>=<value1>;<key2>=<value2>
func ParseBoolMapFromString(input string) (map[string]bool, error) {
	rawValues, err := ParseStringMapFromString(input)
	if err != nil {
		return nil, err
	}

	result := make(map[string]bool)

	for key, value := range rawValues {
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return nil, fmt.Errorf(
				"%w: invalid bool value %s in item %s",
				ErrParseStringFailed,
				value,
				key,
			)
		}

		result[key] = boolValue
	}

	return result, nil
}

func validateEnvironmentValue[T any](value *T, variable *string) error {
	if value == nil && variable == nil {
		return errEnvironmentValueRequired
	}

	if variable != nil && *variable == "" {
		return ErrEnvironmentVariableRequired
	}

	return nil
}

func validateEnvironmentMapValue(variable *string) error {
	if variable != nil && *variable == "" {
		return ErrEnvironmentVariableRequired
	}

	return nil
}

func getEnvVariableValueRequiredError(envName *string) error {
	if envName != nil {
		return fmt.Errorf("%s: %w", *envName, ErrEnvironmentVariableValueRequired)
	}

	return ErrEnvironmentVariableValueRequired
}
