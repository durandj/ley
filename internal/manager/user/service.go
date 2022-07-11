package user

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"regexp"
	"time"

	"github.com/durandj/ley/internal/manager/errortypes"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

var (
	userNameRegex = regexp.MustCompile(`^\w[-\w_ ']+$`)

	//go:embed create_user.sql
	createUserSQL string

	//go:embed get_user_by_username.sql
	getUserByUsernameSQL string
)

// Service is a service for working with user objects.
type Service struct {
	db *sql.DB
}

// NewService creates a new user service.
func NewService(db *sql.DB) *Service {
	return &Service{
		db: db,
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

	creationTime := time.Now().UTC()

	var user User
	err := service.db.QueryRowContext(
		ctx,
		createUserSQL,
		uuid.NewString(),
		opts.Name,
		StatusActive,
		creationTime,
	).Scan(
		&user.id,
		&user.username,
		&user.status,
		&user.createdOn,
		&user.modifiedOn,
	)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			errorName := pqErr.Code.Name()
			constraint := pqErr.Constraint
			if errorName == "unique_violation" && constraint == "users_username_key" {
				return nil, errortypes.NewValidationError("Username is already taken")
			}

			return nil, errortypes.SystemError{
				SafeMessage:   "Unable to create new user due to a system error",
				UnsafeMessage: "Unable to create new user due to a system error",
				WrappedError:  err,
			}
		}

		return nil, errortypes.SystemError{
			SafeMessage:   "Unable to create new user due to a system error",
			UnsafeMessage: "Unable to create new user due to a system error",
			WrappedError:  err,
		}
	}

	return &user, nil
}

// GetUserByUsername fetches a user by their username.
func (service *Service) GetUserByUsername(
	ctx context.Context,
	username string,
) (*User, error) {
	var user User
	err := service.db.QueryRowContext(
		ctx,
		getUserByUsernameSQL,
		username,
	).Scan(
		&user.id,
		&user.username,
		&user.status,
		&user.createdOn,
		&user.modifiedOn,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errortypes.NotFoundError{
				UserError: errortypes.UserError{
					SafeMessage:  "Could not find a user with that name",
					WrappedError: err,
				},
			}
		}

		return nil, errortypes.SystemError{
			SafeMessage:   "Unable to get user by username due to a system error",
			UnsafeMessage: "Unable to get user by username due to a system error",
			WrappedError:  err,
		}
	}

	return &user, nil
}
