package network

import (
	"errors"
	"net/http"

	"github.com/durandj/ley/internal/manager/errortypes"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// Controller handles all the HTTP requests for network related API's.
type Controller struct {
	NetworkService *Service
}

// RegisterRoutes registers HTTP request handlers for all network API's.
func (controller *Controller) RegisterRoutes(router chi.Router) {
	router.Get("/", controller.ListNetworks)
	router.Post("/", controller.CreateNetwork)
}

// CreateNetworkRequest is the expected request body for creating a new
// network.
type CreateNetworkRequest struct {
	Name string `json:"name"`
}

// Bind is used to determine how to map from a request body to a
// network creation request.
func (createNetworkRequest *CreateNetworkRequest) Bind(request *http.Request) error {
	return nil
}

var _ render.Binder = (*CreateNetworkRequest)(nil)

// CreateNetworkResponse is the response body for a successful network
// creation request.
type CreateNetworkResponse struct {
	*RenderableNetwork
}

var _ render.Renderer = (*CreateNetworkResponse)(nil)

// CreateNetwork handles requests to create a new network.
func (controller *Controller) CreateNetwork(
	response http.ResponseWriter,
	request *http.Request,
) {
	ctx := request.Context()

	defer func() {
		_ = request.Body.Close()
	}()

	var createNetworkRequest CreateNetworkRequest
	if err := render.Bind(request, &createNetworkRequest); err != nil {
		// TODO: write error response
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	network, err := controller.NetworkService.CreateNetwork(
		ctx,
		CreateNetworkOpts(createNetworkRequest),
	)
	if err != nil {
		handleError(response, request, err)
		return
	}

	createNetworkResponse := CreateNetworkResponse{
		RenderableNetwork: &RenderableNetwork{
			Name: network.Name(),
		},
	}

	response.WriteHeader(http.StatusCreated)
	_ = render.Render(response, request, createNetworkResponse)
}

// ListNetworksResponse is the response for requesting all networks.
type ListNetworksResponse struct {
	Networks []RenderableNetwork `json:"networks"`
}

// NewListNetworksResponse creates a network list response.
func NewListNetworksResponse(networks []Network) ListNetworksResponse {
	renderableNetworks := make([]RenderableNetwork, len(networks))
	for index := range networks {
		renderableNetworks[index] = NewRenderableNetwork(&networks[index])
	}

	return ListNetworksResponse{
		Networks: renderableNetworks,
	}
}

// Render customizes the rendering process for a response object.
func (listNetworksResponse *ListNetworksResponse) Render(
	response http.ResponseWriter,
	request *http.Request,
) error {
	return nil
}

// ListNetworks handles requests to list the available networks.
func (controller *Controller) ListNetworks(
	response http.ResponseWriter,
	request *http.Request,
) {
	ctx := request.Context()

	networks, err := controller.NetworkService.ListNetworks(ctx)
	if err != nil {
		handleError(response, request, err)
		return
	}

	listNetworksResponse := NewListNetworksResponse(networks)

	response.WriteHeader(http.StatusOK)
	_ = render.Render(response, request, &listNetworksResponse)
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
		_ = render.Render(response, request, &ErrorResponse{
			Message: validationError.SafeMessage,
		})

		return

	case errors.As(err, &userError):
		response.WriteHeader(http.StatusBadRequest)
		_ = render.Render(response, request, &ErrorResponse{
			Message: userError.SafeMessage,
		})

		return

	case errors.As(err, &systemError):
		response.WriteHeader(http.StatusInternalServerError)
		_ = render.Render(response, request, &ErrorResponse{
			Message: systemError.SafeMessage,
		})

		return

	case err != nil:
		response.WriteHeader(http.StatusInternalServerError)
		_ = render.Render(response, request, &ErrorResponse{
			Message: "Internal server error, please try again later",
		})

		return
	}
}

// RenderableNetwork defines what should returned to a user for a
// network.
type RenderableNetwork struct {
	Name string `json:"name"`
}

// NewRenderableNetwork creates a new renderable network from a backend
// network instance.
func NewRenderableNetwork(network *Network) RenderableNetwork {
	return RenderableNetwork{
		Name: network.Name(),
	}
}

// Render provides a hook to customize the render process.
func (renderableNetwork *RenderableNetwork) Render(
	response http.ResponseWriter,
	request *http.Request,
) error {
	return nil
}

var _ render.Renderer = (*RenderableNetwork)(nil)

// ErrorResponse defines what is response to the client if there's an
// error.
type ErrorResponse struct {
	Message string `json:"message"`
}

// Render provides a hook to customize the render process.
func (errorResponse *ErrorResponse) Render(
	response http.ResponseWriter,
	request *http.Request,
) error {
	return nil
}

var _ render.Renderer = (*ErrorResponse)(nil)
