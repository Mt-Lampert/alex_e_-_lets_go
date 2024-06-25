package validator

import "strings"

// package object for methods
type validator struct {
	FieldErrors map[string]string
}

// returns if the FieldErrors map is empty
func (v *validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

// Adds error message to the FieldErrors map; private to this module
func (v *validator) addFieldError(key, msg string) {
	// initialize map, if necessary
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	// if the error to be added does NOT exist ...
	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = msg
	}
}

// Adds error message if an entry fails its validation check
func (v *validator) CheckField(ok bool, key, msg string) {
	if !ok {
		v.addFieldError(key, msg)
	}
}

// Returns true if a string is not empty
func (v *validator) NotBlank(entry string) bool {
	return len(strings.TrimSpace(entry)) > 0
}

// Returns true if a string length is within a range
func (v *validator) WithinRange(min, max int, entry string) bool {
	return len(strings.TrimSpace(entry)) >= min && len(strings.TrimSpace(entry)) <= max
}

// Returns true if a string is at least X chars long
func (v *validator) LongEnough(min int, entry string) bool {
	return len(strings.TrimSpace(entry)) >= min
}

// Returns true if a string is at the most X chars long
func (v *validator) ShortEnough(max int, entry string) bool {
	return len(strings.TrimSpace(entry)) <= max
}

// Return if an Integer can be found on a list of integers
func PermittedInt(value int, permittedValues ...int) bool {
	for i := range permittedValues {
		if value == permittedValues[i] {
			return true
		}
	}
	return false
}

// vim: ts=4 sw=4 fdm=indent
