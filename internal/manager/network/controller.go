package network

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/durandj/ley/internal/manager/errortypes"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"inet.af/netaddr"
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
	Name     string            `json:"name"`
	IPv4CIDR *netaddr.IPPrefix `json:"ipv4CIDR,omitempty"`
	IPv6CIDR *netaddr.IPPrefix `json:"ipv6CIDR,omitempty"`
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
	RenderableNetwork
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
		response.WriteHeader(http.StatusBadRequest)
		_ = render.Render(response, request, &ErrorResponse{
			Message: err.Error(),
		})

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
		RenderableNetwork: NewRenderableNetwork(network),
	}

	response.WriteHeader(http.StatusCreated)
	_ = render.Render(response, request, &createNetworkResponse)
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
	Name       string            `json:"name"`
	IPv4CIDR   *netaddr.IPPrefix `json:"ipv4CIDR,omitempty"`
	IPv6CIDR   *netaddr.IPPrefix `json:"ipv6CIDR,omitempty"`
	CreatedOn  time.Time         `json:"createdOn"`
	ModifiedOn time.Time         `json:"modifiedOn"`
}

// NewRenderableNetwork creates a new renderable network from a backend
// network instance.
func NewRenderableNetwork(network *Network) RenderableNetwork {
	return RenderableNetwork{
		Name:       network.Name(),
		IPv4CIDR:   network.IPv4CIDR(),
		IPv6CIDR:   network.IPv6CIDR(),
		CreatedOn:  network.CreatedOn(),
		ModifiedOn: network.ModifiedOn(),
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

// RenderableTime makes time values renderable in API responses.
type RenderableTime time.Time

// MarshalJSON converts a time value into an ISO-8601 string.
func (renderableTime *RenderableTime) MarshalJSON() ([]byte, error) {
	stringTime := time.Time(*renderableTime).Format(time.RFC3339)

	return []byte(stringTime), nil
}

// UnmarshalJSON converts an ISO-8601 time string into a time instance.
func (renderableTime *RenderableTime) UnmarshalJSON(rawBytes []byte) error {
	timeValue, err := time.Parse(time.RFC3339, string(rawBytes))
	if err != nil {
		return fmt.Errorf("Unable to parse as ISO-8601 time string: %w", err)
	}

	*renderableTime = RenderableTime(timeValue)

	return nil
}

var _ json.Marshaler = (*RenderableTime)(nil)
var _ json.Unmarshaler = (*RenderableTime)(nil)

// RenderableIPPrefix makes an IP CIDR renderable in an API response.
type RenderableIPPrefix struct {
	netaddr.IPPrefix
}

// MarshalJSON renders an IP CIDR as a string.
func (renderableIPPrefix *RenderableIPPrefix) MarshalJSON() ([]byte, error) {
	return []byte(renderableIPPrefix.String()), nil
}

// UnmarshalJSON converts a string into an IP CIDR.
func (renderableIPPrefix *RenderableIPPrefix) UnmarshalJSON(rawBytes []byte) error {
	prefix, err := netaddr.ParseIPPrefix(string(rawBytes))
	if err != nil {
		return fmt.Errorf("Unable to parse JSON string as IP CIDR: %w", err)
	}

	renderableIPPrefix.IPPrefix = prefix

	return nil
}

var _ json.Marshaler = (*RenderableIPPrefix)(nil)
var _ json.Unmarshaler = (*RenderableIPPrefix)(nil)

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
