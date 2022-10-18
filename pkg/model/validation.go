package model

import (
	"fmt"
	"github.com/jmoiron/modl"
)

// NewValidator create a new validator
func NewValidator(executor modl.SqlExecutor) *Validator {
	return &Validator{
		errors: map[string][]string{},
		ex:     executor,
	}
}

// Validator is an object who can test properties and register error messages
// in an error map [property]errorMessages
type Validator struct {
	ex     modl.SqlExecutor
	errors map[string][]string
}

type ValidationErrors map[string][]string

// Errors return the errors map
func (v *Validator) Errors() ValidationErrors {
	return v.errors
}

// AddError add the message in msg to the error map at the index key
func (v *Validator) AddError(key, msg string) {
	if _, ok := v.errors[key]; !ok {
		v.errors[key] = []string{}
	}
	v.errors[key] = append(v.errors[key], msg)
}

func (v ValidationErrors) Error() string {
	return fmt.Sprintf("Validation errors: %v", map[string][]string(v))
}
