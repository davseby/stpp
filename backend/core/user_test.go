package core

import (
	"foodie/server/apierr"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_UserInput_Validate(t *testing.T) {
	tests := map[string]struct {
		UserInput UserInput
		Error     *apierr.Error
	}{
		"Invalid name": {
			UserInput: UserInput{
				Password: "1234",
			},
			Error: apierr.InvalidAttribute("name", "cannot be empty"),
		},
		"Invalid password": {
			UserInput: UserInput{
				Name: "1234",
			},
			Error: apierr.InvalidAttribute("password", "must be at least 4 characters long"),
		},
		"Valid user input": {
			UserInput: UserInput{
				Name:     "1234",
				Password: "1234",
			},
		},
	}

	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, test.Error, test.UserInput.Validate())
		})
	}
}

func Test_ValidatePassword(t *testing.T) {
	tests := map[string]struct {
		Password string
		Error    *apierr.Error
	}{
		"Invalid password": {
			Error: apierr.InvalidAttribute("password", "must be at least 4 characters long"),
		},
		"Valid password": {
			Password: "1234",
		},
	}

	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, test.Error, ValidatePassword(test.Password))
		})
	}
}
