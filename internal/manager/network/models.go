package network

// Network represents a virtual network powered by WireGuard.
type Network struct {
	id   string
	name string
}

// ID is the database ID of the network.
func (network *Network) ID() string {
	return network.id
}

// Name is the name of the network.
func (network *Network) Name() string {
	return network.name
}
