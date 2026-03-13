package internal_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"realworld-fiber-htmx/internal"
)

func TestErrorMessage(t *testing.T) {
	tests := []struct {
		name      string
		fieldName string
		fieldType string
		want      string
	}{
		{"required field", "Name", "required", "Name is required."},
		{"email field", "Email", "email", "Email must be an email."},
		{"unknown type", "Field", "unknown", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := internal.ErrorMessage(tt.fieldName, tt.fieldType)
			assert.Equal(t, tt.want, result)
		})
	}
}
