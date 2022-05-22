package user

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/durandj/ley/internal/manager/errortypes"
	"github.com/google/uuid"
)

var (
	userNameRegex = regexp.MustCompile(`^\w[-\w_ ']+$`)
)

// Service is a service for working with user objects.
type Service struct {
	// TODO: replace this with a DB for persistence
	users []User
}

// NewService creates a new user service.
func NewService() *Service {
	return &Service{
		users: nil,
	}
}

// CreateUserOpts gives the options for creating a new user.
type CreateUserOpts struct {
	Name string
}

// Validate checks that the user creation options are valid.
func (opts *CreateUserOpts) Validate() error {
	if !userNameRegex.MatchString(opts.Name) {
		return fmt.Errorf("Invalid user name '%s'", opts.Name)
	}

	return nil
}

// CreateUser creates a new user object.
func (service *Service) CreateUser(
	ctx context.Context,
	opts CreateUserOpts,
) (*User, error) {
	if err := opts.Validate(); err != nil {
		return nil, errortypes.NewWrappedValidationError(err, "Unable to create user: %v", err)
	}

	// TODO: validate name is unique
	for _, user := range service.users {
		if user.Username() == opts.Name {
			return nil, errortypes.NewValidationError("User name is already taken")
		}
	}

	creationTime := time.Now().UTC()

	user := User{
		id:         uuid.NewString(),
		username:   opts.Name,
		status:     StatusActive,
		createdOn:  creationTime,
		modifiedOn: creationTime,
	}

	service.users = append(service.users, user)

	return &user, nil
}

// GetUserByUsername fetches a user by their username.
func (service *Service) GetUserByUsername(
	ctx context.Context,
	name string,
) (*User, error) {
	return nil, fmt.Errorf("Not implemented")
}
