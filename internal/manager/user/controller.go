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
	router.Get("/", controller.GetUserByUsername)
}

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
	_ = render.Render(response, request, &createUserResponse)
}

// GetUserByUsername fetches user information given a username.
func (controller *Controller) GetUserByUsername(
	response http.ResponseWriter,
	request *http.Request,
) {
	queryParams := request.URL.Query()
	username := queryParams.Get("username")

	if username == "" {
		response.WriteHeader(http.StatusBadRequest)
		_ = render.Render(response, request, &renderable.ErrorResponse{
			Message: "Missing query parameter 'username'",
		})

		return
	}

	ctx := request.Context()

	user, err := controller.UserService.GetUserByUsername(ctx, username)
	if err != nil {
		handleError(response, request, err)
		return
	}

	getUserByUsernameResponse := GetUserByUsernameResponse{
		RenderableUser: NewRenderableUser(user),
	}

	response.WriteHeader(http.StatusOK)
	_ = render.Render(response, request, &getUserByUsernameResponse)
}

func handleError(
	response http.ResponseWriter,
	request *http.Request,
	err error,
) {
	var validationError errortypes.ValidationError
	var notFoundError errortypes.NotFoundError
	var userError errortypes.UserError
	var systemError errortypes.SystemError
	switch {
	case errors.As(err, &validationError):
		response.WriteHeader(http.StatusBadRequest)
		_ = render.Render(response, request, &renderable.ErrorResponse{
			Message: validationError.SafeMessage,
		})

	case errors.As(err, &notFoundError):
		response.WriteHeader(http.StatusNotFound)
		_ = render.Render(response, request, &renderable.ErrorResponse{
			Message: notFoundError.SafeMessage,
		})

	case errors.As(err, &userError):
		response.WriteHeader(http.StatusBadRequest)
		_ = render.Render(response, request, &renderable.ErrorResponse{
			Message: userError.SafeMessage,
		})

	case errors.As(err, &systemError):
		response.WriteHeader(http.StatusInternalServerError)
		_ = render.Render(response, request, &renderable.ErrorResponse{
			Message: systemError.SafeMessage,
		})

	case err != nil:
		response.WriteHeader(http.StatusInternalServerError)
		_ = render.Render(response, request, &renderable.ErrorResponse{
			Message: "Internal server error, please try again later",
		})
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
	return RenderableUser{
		Name:       user.Username(),
		Status:     user.Status(),
		CreatedOn:  renderable.Time(user.CreatedOn()),
		ModifiedOn: renderable.Time(user.ModifiedOn()),
	}
}

// Render provides a hook into the rendering process.
func (renderableUser *RenderableUser) Render(
	response http.ResponseWriter,
	request *http.Request,
) error {
	return nil
}

var _ render.Renderer = (*RenderableUser)(nil)
