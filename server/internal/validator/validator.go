package validator

import (
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"
)

type Validator struct {
	FieldError map[string]string
}

var EmailRX = regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)+$`)

func (v *Validator) Valid() bool {
	return len(v.FieldError) == 0
}

func (v *Validator) AddFieldError(key, error string) {
	if v.FieldError == nil {
		v.FieldError = make(map[string]string)
	}

	if _, exists := v.FieldError[key]; !exists {
		v.FieldError[key] = error
	}
}

func (v *Validator) CheckField(ok bool, key, error string) {
	if !ok {
		v.AddFieldError(key, error)
	}
}

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

func PermittedInt(value int, permittedValues ...int) bool {
	if slices.Contains(permittedValues, value) {
		return true
	}
	return false
}

func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}
