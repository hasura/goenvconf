// Package goenvconf contains reusable structures and utilities for configuration with environment variables.
package goenvconf

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"
)

// GetEnvFunc abstracts a custom function to get the value of an environment variable.
type GetEnvFunc func(string) (string, error)

// EnvString represents either a literal string or an environment reference.
type EnvString struct {
	Value    *string `json:"value,omitempty" jsonschema:"anyof_required=value" mapstructure:"value" yaml:"value,omitempty"`
	Variable *string `json:"env,omitempty"   jsonschema:"anyof_required=env"   mapstructure:"env"   yaml:"env,omitempty"`
}

// NewEnvString creates an EnvString instance.
func NewEnvString(env string, value string) EnvString {
	return EnvString{
		Variable: &env,
		Value:    &value,
	}
}

// NewEnvStringValue creates an EnvString with a literal value.
func NewEnvStringValue(value string) EnvString {
	return EnvString{
		Value: &value,
	}
}

// NewEnvStringVariable creates an EnvString with a variable name.
func NewEnvStringVariable(name string) EnvString {
	return EnvString{
		Variable: &name,
	}
}

// UnmarshalJSON implements json.Unmarshaler.
func (ev *EnvString) UnmarshalJSON(b []byte) error {
	type Plain EnvString

	var rawValue Plain

	err := json.Unmarshal(b, &rawValue)
	if err != nil {
		return err
	}

	value := EnvString(rawValue)
	if value.IsZero() {
		return ErrEnvironmentValueRequired
	}

	*ev = value

	return nil
}

// IsZero checks if the instance is empty.
func (ev EnvString) IsZero() bool {
	return (ev.Variable == nil || *ev.Variable == "") &&
		ev.Value == nil
}

// Get gets literal value or from system environment.
func (ev EnvString) Get() (string, error) {
	if ev.IsZero() {
		return "", ErrEnvironmentValueRequired
	}

	var value string

	var envExisted bool

	if ev.Variable != nil && *ev.Variable != "" {
		value, envExisted = os.LookupEnv(*ev.Variable)
		if value != "" {
			return value, nil
		}
	}

	if ev.Value != nil {
		return *ev.Value, nil
	}

	if envExisted {
		return "", nil
	}

	return "", getEnvVariableValueRequiredError(ev.Variable)
}

// GetOrDefault returns the default value if the environment value is empty.
func (ev EnvString) GetOrDefault(defaultValue string) (string, error) {
	result, err := ev.Get()
	if err != nil {
		if errors.Is(err, ErrEnvironmentVariableValueRequired) {
			return defaultValue, nil
		}

		return "", err
	} else if result == "" {
		result = defaultValue
	}

	return result, nil
}

// GetCustom gets literal value or from system environment by a custom function.
func (ev EnvString) GetCustom(getFunc GetEnvFunc) (string, error) {
	if ev.IsZero() {
		return "", ErrEnvironmentValueRequired
	}

	if ev.Variable != nil && *ev.Variable != "" {
		return getFunc(*ev.Variable)
	}

	if ev.Value != nil {
		return *ev.Value, nil
	}

	return "", getEnvVariableValueRequiredError(ev.Variable)
}

// EnvInt represents either a literal integer or an environment reference.
type EnvInt struct {
	Value    *int64  `json:"value,omitempty" jsonschema:"anyof_required=value" mapstructure:"value" yaml:"value,omitempty"`
	Variable *string `json:"env,omitempty"   jsonschema:"anyof_required=env"   mapstructure:"env"   yaml:"env,omitempty"`
}

// NewEnvInt creates an EnvInt instance.
func NewEnvInt(env string, value int64) EnvInt {
	return EnvInt{
		Variable: &env,
		Value:    &value,
	}
}

// NewEnvIntValue creates an EnvInt with a literal value.
func NewEnvIntValue(value int64) EnvInt {
	return EnvInt{
		Value: &value,
	}
}

// NewEnvIntVariable creates an EnvInt with a variable name.
func NewEnvIntVariable(name string) EnvInt {
	return EnvInt{
		Variable: &name,
	}
}

// IsZero checks if the instance is empty.
func (ev EnvInt) IsZero() bool {
	return (ev.Variable == nil || *ev.Variable == "") &&
		ev.Value == nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (ev *EnvInt) UnmarshalJSON(b []byte) error {
	type Plain EnvInt

	var rawValue Plain

	err := json.Unmarshal(b, &rawValue)
	if err != nil {
		return err
	}

	value := EnvInt(rawValue)
	if value.IsZero() {
		return ErrEnvironmentValueRequired
	}

	*ev = value

	return nil
}

// Get gets literal value or from system environment.
func (ev EnvInt) Get() (int64, error) {
	if ev.IsZero() {
		return 0, ErrEnvironmentValueRequired
	}

	if ev.Variable != nil && *ev.Variable != "" {
		rawValue := os.Getenv(*ev.Variable)
		if rawValue != "" {
			return strconv.ParseInt(rawValue, 10, 64)
		}
	}

	if ev.Value != nil {
		return *ev.Value, nil
	}

	return 0, getEnvVariableValueRequiredError(ev.Variable)
}

// GetOrDefault returns the default value if the environment value is empty.
func (ev EnvInt) GetOrDefault(defaultValue int64) (int64, error) {
	result, err := ev.Get()
	if err != nil {
		if errors.Is(err, ErrEnvironmentVariableValueRequired) {
			return defaultValue, nil
		}

		return 0, err
	}

	return result, nil
}

// GetCustom gets literal value or from system environment by a custom function.
func (ev EnvInt) GetCustom(getFunc GetEnvFunc) (int64, error) {
	if ev.IsZero() {
		return 0, ErrEnvironmentValueRequired
	}

	if ev.Variable != nil && *ev.Variable != "" {
		rawValue, err := getFunc(*ev.Variable)
		if err != nil {
			return 0, err
		}

		if rawValue != "" {
			return strconv.ParseInt(rawValue, 10, 64)
		}
	}

	if ev.Value != nil {
		return *ev.Value, nil
	}

	return 0, getEnvVariableValueRequiredError(ev.Variable)
}

// EnvBool represents either a literal boolean or an environment reference.
type EnvBool struct {
	Value    *bool   `json:"value,omitempty" jsonschema:"anyof_required=value" mapstructure:"value" yaml:"value,omitempty"`
	Variable *string `json:"env,omitempty"   jsonschema:"anyof_required=env"   mapstructure:"env"   yaml:"env,omitempty"`
}

// NewEnvBool creates an EnvBool instance.
func NewEnvBool(env string, value bool) EnvBool {
	return EnvBool{
		Variable: &env,
		Value:    &value,
	}
}

// NewEnvBoolValue creates an EnvBool with a literal value.
func NewEnvBoolValue(value bool) EnvBool {
	return EnvBool{
		Value: &value,
	}
}

// NewEnvBoolVariable creates an EnvBool with a variable name.
func NewEnvBoolVariable(name string) EnvBool {
	return EnvBool{
		Variable: &name,
	}
}

// IsZero checks if the instance is empty.
func (ev EnvBool) IsZero() bool {
	return (ev.Variable == nil || *ev.Variable == "") &&
		ev.Value == nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (ev *EnvBool) UnmarshalJSON(b []byte) error {
	type Plain EnvBool

	var rawValue Plain

	err := json.Unmarshal(b, &rawValue)
	if err != nil {
		return err
	}

	value := EnvBool(rawValue)
	if value.IsZero() {
		return ErrEnvironmentValueRequired
	}

	*ev = value

	return nil
}

// Get gets literal value or from system environment.
func (ev EnvBool) Get() (bool, error) {
	if ev.IsZero() {
		return false, ErrEnvironmentValueRequired
	}

	if ev.Variable != nil && *ev.Variable != "" {
		rawValue := os.Getenv(*ev.Variable)
		if rawValue != "" {
			return strconv.ParseBool(rawValue)
		}
	}

	if ev.Value != nil {
		return *ev.Value, nil
	}

	return false, getEnvVariableValueRequiredError(ev.Variable)
}

// GetOrDefault returns the default value if the environment value is empty.
func (ev EnvBool) GetOrDefault(defaultValue bool) (bool, error) {
	result, err := ev.Get()
	if err != nil {
		if errors.Is(err, ErrEnvironmentVariableValueRequired) {
			return defaultValue, nil
		}

		return false, err
	}

	return result, nil
}

// GetCustom gets literal value or from system environment with custom function.
func (ev EnvBool) GetCustom(getFunc GetEnvFunc) (bool, error) {
	if ev.IsZero() {
		return false, ErrEnvironmentValueRequired
	}

	if ev.Variable != nil && *ev.Variable != "" {
		rawValue, err := getFunc(*ev.Variable)
		if err != nil {
			return false, err
		}

		if rawValue != "" {
			return strconv.ParseBool(rawValue)
		}
	}

	if ev.Value != nil {
		return *ev.Value, nil
	}

	return false, getEnvVariableValueRequiredError(ev.Variable)
}

// EnvFloat represents either a literal floating point number or an environment reference.
type EnvFloat struct {
	Value    *float64 `json:"value,omitempty" jsonschema:"anyof_required=value" mapstructure:"value" yaml:"value,omitempty"`
	Variable *string  `json:"env,omitempty"   jsonschema:"anyof_required=env"   mapstructure:"env"   yaml:"env,omitempty"`
}

// NewEnvFloat creates an EnvFloat instance.
func NewEnvFloat(env string, value float64) EnvFloat {
	return EnvFloat{
		Variable: &env,
		Value:    &value,
	}
}

// NewEnvFloatValue creates an EnvFloat with a literal value.
func NewEnvFloatValue(value float64) EnvFloat {
	return EnvFloat{
		Value: &value,
	}
}

// NewEnvFloatVariable creates an EnvFloat with a variable name.
func NewEnvFloatVariable(name string) EnvFloat {
	return EnvFloat{
		Variable: &name,
	}
}

// IsZero checks if the instance is empty.
func (ev EnvFloat) IsZero() bool {
	return (ev.Variable == nil || *ev.Variable == "") &&
		ev.Value == nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (ev *EnvFloat) UnmarshalJSON(b []byte) error {
	type Plain EnvFloat

	var rawValue Plain

	err := json.Unmarshal(b, &rawValue)
	if err != nil {
		return err
	}

	value := EnvFloat(rawValue)
	if value.IsZero() {
		return ErrEnvironmentValueRequired
	}

	*ev = value

	return nil
}

// Get gets literal value or from system environment.
func (ev EnvFloat) Get() (float64, error) {
	if ev.IsZero() {
		return 0, ErrEnvironmentValueRequired
	}

	if ev.Variable != nil && *ev.Variable != "" {
		rawValue := os.Getenv(*ev.Variable)
		if rawValue != "" {
			return strconv.ParseFloat(rawValue, 64)
		}
	}

	if ev.Value != nil {
		return *ev.Value, nil
	}

	return 0, getEnvVariableValueRequiredError(ev.Variable)
}

// GetOrDefault returns the default value if the environment value is empty.
func (ev EnvFloat) GetOrDefault(defaultValue float64) (float64, error) {
	result, err := ev.Get()
	if err != nil {
		if errors.Is(err, ErrEnvironmentVariableValueRequired) {
			return defaultValue, nil
		}

		return 0, err
	}

	return result, nil
}

// GetCustom gets literal value or from system environment by a custom function.
func (ev EnvFloat) GetCustom(getFunc GetEnvFunc) (float64, error) {
	if ev.IsZero() {
		return 0, ErrEnvironmentValueRequired
	}

	if ev.Variable != nil && *ev.Variable != "" {
		rawValue, err := getFunc(*ev.Variable)
		if err != nil {
			return 0, err
		}

		if rawValue != "" {
			return strconv.ParseFloat(rawValue, 64)
		}
	}

	if ev.Value != nil {
		return *ev.Value, nil
	}

	return 0, getEnvVariableValueRequiredError(ev.Variable)
}

// EnvMapString represents either a literal string map or an environment reference.
type EnvMapString struct {
	Value    map[string]string `json:"value,omitempty" jsonschema:"anyof_required=value" mapstructure:"value" yaml:"value,omitempty"`
	Variable *string           `json:"env,omitempty"   jsonschema:"anyof_required=env"   mapstructure:"env"   yaml:"env,omitempty"`
}

// NewEnvMapString creates an EnvMapString instance.
func NewEnvMapString(env string, value map[string]string) EnvMapString {
	return EnvMapString{
		Variable: &env,
		Value:    value,
	}
}

// NewEnvMapStringValue creates an EnvMapString with a literal value.
func NewEnvMapStringValue(value map[string]string) EnvMapString {
	return EnvMapString{
		Value: value,
	}
}

// NewEnvMapStringVariable creates an EnvMapString with a variable name.
func NewEnvMapStringVariable(name string) EnvMapString {
	return EnvMapString{
		Variable: &name,
	}
}

// IsZero checks if the instance is empty.
func (ev EnvMapString) IsZero() bool {
	return (ev.Variable == nil || *ev.Variable == "") &&
		ev.Value == nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (ev *EnvMapString) UnmarshalJSON(b []byte) error {
	type Plain EnvMapString

	var rawValue Plain

	err := json.Unmarshal(b, &rawValue)
	if err != nil {
		return err
	}

	*ev = EnvMapString(rawValue)

	return nil
}

// Get gets literal value or from system environment.
func (ev EnvMapString) Get() (map[string]string, error) {
	if ev.Variable != nil && *ev.Variable != "" {
		rawValue := os.Getenv(*ev.Variable)
		if rawValue != "" {
			return ParseStringMapFromString(rawValue)
		}
	}

	return ev.Value, nil
}

// GetCustom gets literal value or from system environment by a custom function.
func (ev EnvMapString) GetCustom(getFunc GetEnvFunc) (map[string]string, error) {
	if ev.Variable != nil && *ev.Variable != "" {
		rawValue, err := getFunc(*ev.Variable)
		if err != nil {
			return nil, err
		}

		if rawValue != "" {
			return ParseStringMapFromString(rawValue)
		}
	}

	return ev.Value, nil
}

// EnvMapInt represents either a literal int map or an environment reference.
type EnvMapInt struct {
	Value    map[string]int64 `json:"value,omitempty" jsonschema:"anyof_required=value" mapstructure:"value" yaml:"value,omitempty"`
	Variable *string          `json:"env,omitempty"   jsonschema:"anyof_required=env"   mapstructure:"env"   yaml:"env,omitempty"`
}

// NewEnvMapInt creates an EnvMapInt instance.
func NewEnvMapInt(env string, value map[string]int64) EnvMapInt {
	return EnvMapInt{
		Variable: &env,
		Value:    value,
	}
}

// NewEnvMapIntValue creates an EnvMapInt with a literal value.
func NewEnvMapIntValue(value map[string]int64) EnvMapInt {
	return EnvMapInt{
		Value: value,
	}
}

// NewEnvMapIntVariable creates an EnvMapInt with a variable name.
func NewEnvMapIntVariable(name string) EnvMapInt {
	return EnvMapInt{
		Variable: &name,
	}
}

// IsZero checks if the instance is empty.
func (ev EnvMapInt) IsZero() bool {
	return (ev.Variable == nil || *ev.Variable == "") &&
		ev.Value == nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (ev *EnvMapInt) UnmarshalJSON(b []byte) error {
	type Plain EnvMapInt

	var rawValue Plain

	err := json.Unmarshal(b, &rawValue)
	if err != nil {
		return err
	}

	*ev = EnvMapInt(rawValue)

	return nil
}

// Get gets literal value or from system environment.
func (ev EnvMapInt) Get() (map[string]int64, error) {
	if ev.Variable != nil && *ev.Variable != "" {
		rawValue := os.Getenv(*ev.Variable)
		if rawValue != "" {
			return ParseIntegerMapFromString[int64](rawValue)
		}
	}

	return ev.Value, nil
}

// GetCustom gets literal value or from system environment by a custom function.
func (ev EnvMapInt) GetCustom(getFunc GetEnvFunc) (map[string]int64, error) {
	if ev.Variable != nil && *ev.Variable != "" {
		rawValue, err := getFunc(*ev.Variable)
		if err != nil {
			return nil, err
		}

		if rawValue != "" {
			return ParseIntegerMapFromString[int64](rawValue)
		}
	}

	return ev.Value, nil
}

// EnvMapFloat represents either a literal float map or an environment reference.
type EnvMapFloat struct {
	Value    map[string]float64 `json:"value,omitempty" jsonschema:"anyof_required=value" mapstructure:"value" yaml:"value,omitempty"`
	Variable *string            `json:"env,omitempty"   jsonschema:"anyof_required=env"   mapstructure:"env"   yaml:"env,omitempty"`
}

// NewEnvMapFloat creates an EnvMapFloat instance.
func NewEnvMapFloat(env string, value map[string]float64) EnvMapFloat {
	return EnvMapFloat{
		Variable: &env,
		Value:    value,
	}
}

// NewEnvMapFloatValue creates an EnvMapFloat with a literal value.
func NewEnvMapFloatValue(value map[string]float64) EnvMapFloat {
	return EnvMapFloat{
		Value: value,
	}
}

// NewEnvMapFloatVariable creates an EnvMapFloat with a variable name.
func NewEnvMapFloatVariable(name string) EnvMapFloat {
	return EnvMapFloat{
		Variable: &name,
	}
}

// IsZero checks if the instance is empty.
func (ev EnvMapFloat) IsZero() bool {
	return (ev.Variable == nil || *ev.Variable == "") &&
		ev.Value == nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (ev *EnvMapFloat) UnmarshalJSON(b []byte) error {
	type Plain EnvMapFloat

	var rawValue Plain

	err := json.Unmarshal(b, &rawValue)
	if err != nil {
		return err
	}

	*ev = EnvMapFloat(rawValue)

	return nil
}

// Get gets literal value or from system environment.
func (ev EnvMapFloat) Get() (map[string]float64, error) {
	if ev.Variable != nil && *ev.Variable != "" {
		rawValue := os.Getenv(*ev.Variable)
		if rawValue != "" {
			return ParseFloatMapFromString[float64](rawValue)
		}
	}

	return ev.Value, nil
}

// GetCustom gets literal value or from system environment by a custom function.
func (ev EnvMapFloat) GetCustom(getFunc GetEnvFunc) (map[string]float64, error) {
	if ev.Variable != nil && *ev.Variable != "" {
		rawValue, err := getFunc(*ev.Variable)
		if err != nil {
			return nil, err
		}

		if rawValue != "" {
			return ParseFloatMapFromString[float64](rawValue)
		}
	}

	return ev.Value, nil
}

// EnvMapBool represents either a literal bool map or an environment reference.
type EnvMapBool struct {
	Value    map[string]bool `json:"value,omitempty" jsonschema:"anyof_required=value" mapstructure:"value" yaml:"value,omitempty"`
	Variable *string         `json:"env,omitempty"   jsonschema:"anyof_required=env"   mapstructure:"env"   yaml:"env,omitempty"`
}

// NewEnvMapBool creates an EnvMapBool instance.
func NewEnvMapBool(env string, value map[string]bool) EnvMapBool {
	return EnvMapBool{
		Variable: &env,
		Value:    value,
	}
}

// NewEnvMapBoolValue creates an EnvMapBool with a literal value.
func NewEnvMapBoolValue(value map[string]bool) EnvMapBool {
	return EnvMapBool{
		Value: value,
	}
}

// NewEnvMapBoolVariable creates an EnvMapBool with a variable name.
func NewEnvMapBoolVariable(name string) EnvMapBool {
	return EnvMapBool{
		Variable: &name,
	}
}

// IsZero checks if the instance is empty.
func (ev EnvMapBool) IsZero() bool {
	return (ev.Variable == nil || *ev.Variable == "") &&
		ev.Value == nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (ev *EnvMapBool) UnmarshalJSON(b []byte) error {
	type Plain EnvMapBool

	var rawValue Plain

	err := json.Unmarshal(b, &rawValue)
	if err != nil {
		return err
	}

	*ev = EnvMapBool(rawValue)

	return nil
}

// Get gets literal value or from system environment.
func (ev EnvMapBool) Get() (map[string]bool, error) {
	if ev.Variable != nil && *ev.Variable != "" {
		rawValue := os.Getenv(*ev.Variable)
		if rawValue != "" {
			return ParseBoolMapFromString(rawValue)
		}
	}

	return ev.Value, nil
}

// GetCustom gets literal value or from system environment by a custom function.
func (ev EnvMapBool) GetCustom(getFunc GetEnvFunc) (map[string]bool, error) {
	if ev.Variable != nil && *ev.Variable != "" {
		rawValue, err := getFunc(*ev.Variable)
		if err != nil {
			return nil, err
		}

		if rawValue != "" {
			return ParseBoolMapFromString(rawValue)
		}
	}

	return ev.Value, nil
}
