package validator

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

var EmailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// package object for methods
type Validator struct {
	FieldErrors map[string]string
}

// returns if the FieldErrors map is empty
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

// Adds error message to the FieldErrors map; private to this module
func (v *Validator) AddFieldError(key, msg string) {
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
func (v *Validator) CheckField(ok bool, key, msg string) {
	if !ok {
		v.AddFieldError(key, msg)
	}
}

// Returns true if a string is not empty
func (v *Validator) NotBlank(entry string) bool {
	return len(strings.TrimSpace(entry)) > 0
}

// Returns true if a string length is within a range
func (v *Validator) WithinRange(min, max int, entry string) bool {
	return len(strings.TrimSpace(entry)) >= min && len(strings.TrimSpace(entry)) <= max
}

// Returns true if a string is at least X chars long
func (v *Validator) LongEnough(min int, entry string) bool {
	return len(strings.TrimSpace(entry)) >= min
}

// Returns true if a string is at the most X chars long
func (v *Validator) ShortEnough(max int, entry string) bool {
	return len(strings.TrimSpace(entry)) <= max
}

func (v *Validator) ValidExpiration(value string) bool {
	permittedValues := []string{
		`1 day`,
		`7 days`,
		`1 month`,
		`1 year`,
	}

	for i := range permittedValues {
		if value == permittedValues[i] {
			return true
		}
	}
	return false
}

// Return if an Integer can be found on a list of integers
func (v *Validator) PermittedInt(value int, permittedValues ...int) bool {
	for i := range permittedValues {
		if value == permittedValues[i] {
			return true
		}
	}
	return false
}

// returns true if a string contains at least n string chars
func (v *Validator) MinChars(s string, n int) bool {
	return utf8.RuneCountInString(s) >= 8
}

// returns true if a string matches a pre-compiled regex
func (v *Validator) Matches(s string, rx *regexp.Regexp) bool {
	return rx.MatchString(s)
}

// vim: ts=4 sw=4 fdm=indent
