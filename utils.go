package goenvconf

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	keyValueLength = 2
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

		if len(keyValue) != keyValueLength || keyValue[0] == "" {
			return nil, NewParseEnvFailedError(
				"invalid string map syntax, expected: <key1>=<value1>;<key2>=<value2>",
				keyValue[0],
			)
		}

		result[keyValue[0]] = keyValue[1]
	}

	return result, nil
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
			return nil, NewParseEnvFailedError("invalid integer map syntax", key)
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
			return nil, NewParseEnvFailedError("invalid float map syntax", key)
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
			return nil, NewParseEnvFailedError("invalid boolean map syntax", key)
		}

		result[key] = boolValue
	}

	return result, nil
}

// ParseStringSliceFromString parses a string slice from a comma-separated string.
func ParseStringSliceFromString(input string) []string {
	return strings.Split(input, ",")
}

// ParseIntSliceFromString parses an integer slice from a comma-separated string.
func ParseIntSliceFromString[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64](
	input string,
) ([]T, error) {
	return parseIntSliceFromStringWithErrorPrefix[T](input, "")
}

func parseIntSliceFromStringWithErrorPrefix[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64](
	input string,
	errorPrefix string,
) ([]T, error) {
	rawValues := ParseStringSliceFromString(input)
	results := make([]T, len(rawValues))

	for index, val := range rawValues {
		intVal, err := strconv.ParseInt(strings.TrimSpace(val), 10, 64)
		if err != nil {
			return nil, NewParseEnvFailedError(
				errorPrefix+"invalid integer slice syntax",
				strconv.Itoa(index),
			)
		}

		results[index] = T(intVal)
	}

	return results, nil
}

// ParseFloatSliceFromString parses a floating-point number slice from a comma-separated string.
func ParseFloatSliceFromString[T float32 | float64](input string) ([]T, error) {
	return parseFloatSliceFromStringWithErrorPrefix[T](input, "")
}

func parseFloatSliceFromStringWithErrorPrefix[T float32 | float64](
	input string,
	errorPrefix string,
) ([]T, error) {
	rawValues := ParseStringSliceFromString(input)
	results := make([]T, len(rawValues))

	for index, val := range rawValues {
		floatVal, err := strconv.ParseFloat(strings.TrimSpace(val), 64)
		if err != nil {
			return nil, NewParseEnvFailedError(
				errorPrefix+"invalid floating-point number slice syntax",
				strconv.Itoa(index),
			)
		}

		results[index] = T(floatVal)
	}

	return results, nil
}

// ParseBoolSliceFromString parses a boolean slice from a comma-separated string.
func ParseBoolSliceFromString(input string) ([]bool, error) {
	return parseBoolSliceFromStringWithErrorPrefix(input, "")
}

func parseBoolSliceFromStringWithErrorPrefix(input string, errorPrefix string) ([]bool, error) {
	rawValues := ParseStringSliceFromString(input)
	results := make([]bool, len(rawValues))

	for index, val := range rawValues {
		boolVal, err := strconv.ParseBool(strings.TrimSpace(val))
		if err != nil {
			return nil, NewParseEnvFailedError(
				errorPrefix+"invalid boolean slice syntax",
				strconv.Itoa(index),
			)
		}

		results[index] = boolVal
	}

	return results, nil
}

// OSEnvGetter wraps the GetOSEnv function with context.
func OSEnvGetter(_ context.Context) GetEnvFunc {
	return GetOSEnv
}

// GetOSEnv implements the GetEnvFunc with OS environment.
func GetOSEnv(s string) (string, error) {
	value, ok := os.LookupEnv(s)
	if !ok {
		return value, ErrEnvironmentVariableValueRequired
	}

	return value, nil
}

func getEnvVariableValueRequiredError(envName *string) error {
	if envName != nil {
		return fmt.Errorf("%s: %w", *envName, ErrEnvironmentVariableValueRequired)
	}

	return ErrEnvironmentVariableValueRequired
}
