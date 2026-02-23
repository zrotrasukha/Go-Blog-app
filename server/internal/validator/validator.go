package validator

import (
	"slices"
	"strings"
	"unicode/utf8"
)

type Validator struct {
	FieldError map[string]string
}

func (v *Validator) Valid() bool {
	return len(v.FieldError) == 0
}

func (v *Validator) AddError(key, error string) {
	if v.FieldError == nil {
		v.FieldError = make(map[string]string)
	}

	if _, exists := v.FieldError[key]; !exists {
		v.FieldError[key] = error
	}
}

func (v *Validator) CheckField(ok bool, key, error string) {
	if !ok {
		v.AddError(key, error)
	}
}

func (v *Validator) NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func (v *Validator) MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

func (v *Validator) PermittedInt(value int, permittedValues ...int) bool {
	if slices.Contains(permittedValues, value) {
		return true
	}
	return false
}
