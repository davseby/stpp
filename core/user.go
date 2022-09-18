package core

import (
	"foodie/server/apierr"
	"time"

	"github.com/rs/xid"
)

const (
	// RootAdminName specifies the name of the admin user that shouldn't
	// be deleted.
	RootAdminName = "admin"
)

// User contains user information.
type User struct {
	// ID specifies the id of the user.
	ID xid.ID `json:"id"`

	// Name specifies the name of the user.
	Name string `json:"name"`

	// Admin specifies whether the user has administrator permissions.
	Admin bool `json:"admin"`

	// PasswordHash contains password hash information. It shouldn't be
	// exposed to the clients.
	PasswordHash []byte `json:"-"`

	// CreatedAt specifies a time at which the object was created.
	CreatedAt time.Time `json:"created_at"`
}

// UserInput contains core user information that is used only when creating
// or updating users.
type UserInput struct {
	// Name specifies the name of the user.
	Name string `json:"name"`

	// Password specifies the password of the user.
	Password string `json:"password"`
}

// Validate checks whether user input contains valid attributes.
func (ui *UserInput) Validate() *apierr.Error {
	if ui.Name == "" {
		return apierr.InvalidAttribute("name", "cannot be empty")
	}

	return ValidatePassword(ui.Password)
}

// ValidatePassword validates the password.
func ValidatePassword(pass string) *apierr.Error {
	if len(pass) < 4 {
		return apierr.InvalidAttribute("password", "must be at least 4 characters long")
	}

	return nil
}
