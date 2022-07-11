package user

import (
	"net/http"
	"time"

	"github.com/go-chi/render"
)

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

// CreateUserRequest holds the request body for creating a new user.
type CreateUserRequest struct {
	Name string `json:"name"`
}

// SetName sets the name of the user in the creation request and returns
// an instance of the request object. This can be used as a factory
// builder style pattern.
func (createUserRequest *CreateUserRequest) SetName(name string) *CreateUserRequest {
	createUserRequest.Name = name

	return createUserRequest
}

// Bind is a hook into the process for converting an HTTP request body
// into a request object.
func (createUserRequest *CreateUserRequest) Bind(request *http.Request) error {
	return nil
}

var _ render.Binder = (*CreateUserRequest)(nil)

// CreateUserResponse holds the response object for creating a user user.
type CreateUserResponse struct {
	RenderableUser
}

var _ render.Renderer = (*CreateUserResponse)(nil)

// GetUserByUsernameResponse is the response returned when requesting
// a user by their username.
type GetUserByUsernameResponse struct {
	RenderableUser
}

var _ render.Renderer = (*GetUserByUsernameResponse)(nil)
