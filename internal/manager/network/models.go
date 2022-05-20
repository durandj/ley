package network

import (
	"time"

	"inet.af/netaddr"
)

// Network represents a virtual network powered by WireGuard.
type Network struct {
	id        string
	name      string
	ipv4CIDR  *netaddr.IPPrefix
	ipv6CIDR  *netaddr.IPPrefix
	createdOn time.Time
	// TODO: createdBy
	modifiedOn time.Time
	// TODO: modifiedBy
	// TODO: add list of nodes that belong to network
	// TODO: ACL/permissions policy
	// TODO: add ingress settings
	// TODO: add egress settings
}

// ID is the database ID of the network.
func (network *Network) ID() string {
	return network.id
}

// Name is the name of the network.
func (network *Network) Name() string {
	return network.name
}

// IPv4CIDR is the set of IPv4 addresses that this network can use.
func (network *Network) IPv4CIDR() *netaddr.IPPrefix {
	return network.ipv4CIDR
}

// IPv6CIDR is the set of IPv6 addresses that this network can use.
func (network *Network) IPv6CIDR() *netaddr.IPPrefix {
	return network.ipv6CIDR
}

// CreatedOn is the date and time that the network was created on.
func (network *Network) CreatedOn() time.Time {
	return network.createdOn
}

// ModifiedOn is the date and time that the network was last modified on.
func (network *Network) ModifiedOn() time.Time {
	return network.modifiedOn
}
