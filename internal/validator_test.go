package internal_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"realworld-fiber-htmx/internal"
)

type testStruct struct {
	Name  string `validate:"required"`
	Email string `validate:"required,email"`
}

func TestValidate_Valid(t *testing.T) {
	v := internal.NewValidator()
	err := v.Validate(testStruct{Name: "test", Email: "test@example.com"})
	assert.NoError(t, err)
}

func TestValidate_MissingRequired(t *testing.T) {
	v := internal.NewValidator()
	err := v.Validate(testStruct{Name: "", Email: "test@example.com"})
	assert.Error(t, err)
}

func TestValidate_InvalidEmail(t *testing.T) {
	v := internal.NewValidator()
	err := v.Validate(testStruct{Name: "test", Email: "not-an-email"})
	assert.Error(t, err)
}
