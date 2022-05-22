package user

import (
	"errors"
	"net/http"

	"github.com/durandj/ley/internal/manager/errortypes"
	"github.com/durandj/ley/internal/manager/renderable"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// Controller handles user related HTTP requests.
type Controller struct {
	UserService *Service
}

// RegisterRoutes adds HTTP routes to the parent router.
func (controller *Controller) RegisterRoutes(router chi.Router) {
	router.Post("/", controller.CreateUser)
}

// CreateUserRequest holds the request body for creating a new user.
type CreateUserRequest struct {
	Name string `json:"name"`
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

// CreateUser handles requests to create a new user.
func (controller *Controller) CreateUser(
	response http.ResponseWriter,
	request *http.Request,
) {
	ctx := request.Context()

	defer func() {
		_ = request.Body.Close()
	}()

	var createUserRequest CreateUserRequest
	if err := render.Bind(request, &createUserRequest); err != nil {
		response.WriteHeader(http.StatusBadRequest)
		_ = render.Render(response, request, &renderable.ErrorResponse{
			Message: err.Error(),
		})

		return
	}

	user, err := controller.UserService.CreateUser(
		ctx,
		CreateUserOpts(createUserRequest),
	)
	if err != nil {
		handleError(response, request, err)
		return
	}

	createUserResponse := CreateUserResponse{
		RenderableUser: NewRenderableUser(user),
	}

	response.WriteHeader(http.StatusCreated)
	// TODO: set location header
	_ = render.Render(response, request, &createUserResponse)
}

func handleError(
	response http.ResponseWriter,
	request *http.Request,
	err error,
) {
	var validationError errortypes.ValidationError
	var userError errortypes.UserError
	var systemError errortypes.SystemError
	switch {
	case errors.As(err, &validationError):
		response.WriteHeader(http.StatusBadRequest)
		_ = render.Render(response, request, &renderable.ErrorResponse{
			Message: validationError.SafeMessage,
		})

		return

	case errors.As(err, &userError):
		response.WriteHeader(http.StatusBadRequest)
		_ = render.Render(response, request, &renderable.ErrorResponse{
			Message: userError.SafeMessage,
		})

		return

	case errors.As(err, &systemError):
		response.WriteHeader(http.StatusInternalServerError)
		_ = render.Render(response, request, &renderable.ErrorResponse{
			Message: systemError.SafeMessage,
		})

		return

	case err != nil:
		response.WriteHeader(http.StatusInternalServerError)
		_ = render.Render(response, request, &renderable.ErrorResponse{
			Message: "Internal server error, please try again later",
		})

		return
	}
}

// RenderableUser gives the safe version of a user that can be returned
// over HTTP.
type RenderableUser struct {
	Name       string          `json:"name"`
	Status     Status          `json:"status"`
	CreatedOn  renderable.Time `json:"createdOn"`
	ModifiedOn renderable.Time `json:"modifiedOn"`
}

// NewRenderableUser creates a renderable user from a backend user
// instance.
func NewRenderableUser(user *User) RenderableUser {
	return RenderableUser{}
}

// Render provides a hook into the rendering process.
func (renderableUser *RenderableUser) Render(
	response http.ResponseWriter,
	request *http.Request,
) error {
	return nil
}

var _ render.Renderer = (*RenderableUser)(nil)
