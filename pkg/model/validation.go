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
	v.errors[key] = append(v.errors[key], msg)
}

func (v ValidationErrors) Error() string {
	return fmt.Sprintf("Validation errors: %v", map[string][]string(v))
}

// Exists verifies value exists in table.field
func (v *Validator) Exists(table, field string, value interface{}, errKey, msg string) bool {
	query := fmt.Sprintf("SELECT count(1) FROM %#v WHERE %#v = $1", table, field)
	res := struct{ Count int }{-1}
	v.ex.SelectOne(&res, query, value)
	if res.Count < 1 {
		v.AddError(errKey, msg)
		return false
	}
	return true
}
