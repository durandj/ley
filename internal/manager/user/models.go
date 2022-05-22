package user

import "time"

// User defines the backend view of what a user is.
type User struct {
	id         string
	username   string
	status     Status
	createdOn  time.Time
	modifiedOn time.Time
}

// ID gives the backend ID of the user.
func (user *User) ID() string {
	return user.id
}

// Username returns the username of the user.
func (user *User) Username() string {
	return user.username
}

// CreatedOn gives the date the user was created on.
func (user *User) CreatedOn() time.Time {
	return user.createdOn
}

// ModifiedOn gives the date that the user was last modified on.
func (user *User) ModifiedOn() time.Time {
	return user.modifiedOn
}

// Status gives the current activation status of the user.
func (user *User) Status() Status {
	return user.status
}

// Status tells if the user is active or not.
type Status string

const (
	// StatusActive marks the user as active and able to use the
	// system.
	StatusActive Status = "active"

	// StatusDeactivated marks the user as no longer allowed to use
	// the system. This could be temporary or permenant. The user is
	// immutable while in this status.
	StatusDeactivated Status = "deactivated"
)
