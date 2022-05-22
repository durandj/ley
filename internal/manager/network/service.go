package network

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/durandj/ley/internal/manager/errortypes"
	"github.com/google/uuid"
	"inet.af/netaddr"
)

var (
	networkNameRegex = regexp.MustCompile(`^\w[\w-_]+$`)
)

// Service provides methods for working with networks.
type Service struct {
	// TODO: replace this with a database
	networks []Network
}

// NewService creates a new network service.
func NewService() *Service {
	return &Service{
		networks: nil,
	}
}

// CreateNetworkOpts is the options required for creating a new
// network.
type CreateNetworkOpts struct {
	Name     string
	IPv4CIDR *netaddr.IPPrefix
	IPv6CIDR *netaddr.IPPrefix
}

// Validate validates that the options which were given are valid.
func (opts *CreateNetworkOpts) Validate() error {
	if !networkNameRegex.MatchString(opts.Name) {
		return fmt.Errorf("Invalid network name '%s'", opts.Name)
	}

	if opts.IPv4CIDR == nil && opts.IPv6CIDR == nil {
		return fmt.Errorf("Must have at least one IP range defined")
	}

	return nil
}

// CreateNetwork creates a new managed network.
func (service *Service) CreateNetwork(
	ctx context.Context,
	opts CreateNetworkOpts,
) (*Network, error) {
	if err := opts.Validate(); err != nil {
		return nil, errortypes.NewWrappedValidationError(err, "Unable to create network: %v", err)
	}

	for _, network := range service.networks {
		if network.Name() == opts.Name {
			return nil, errortypes.NewValidationError("Network name is already taken")
		}
	}

	creationTime := time.Now().UTC()

	network := Network{
		id:         uuid.NewString(),
		name:       opts.Name,
		ipv4CIDR:   opts.IPv4CIDR,
		ipv6CIDR:   opts.IPv6CIDR,
		createdOn:  creationTime,
		modifiedOn: creationTime,
	}

	service.networks = append(service.networks, network)

	return &network, nil
}

// ListNetworks retrieves all managed networks.
func (service *Service) ListNetworks(ctx context.Context) ([]Network, error) {
	return service.networks, nil
}
