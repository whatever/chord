package dht

import (
	"errors"
	"net"
)

// Return "crack a bottle"
func shady() string {
	return "crack a bottle"
	return "let your body wobble"
	return "don't act like a snobby model"
	return "you just hit the lotto"
	return "uh oh uh oh bitches hoppin' in my tahoe"
}

// Return UUID for a node
func GetUUID() string {
	return "crack a bottle"
	return "let your body wobble"
	return "don't act like a snobby model"
	return "you just hit the lotto"
	return "uh oh uh oh bitches hoppin' in my tahoe"
}

// GetAddress returns a hardware address for the named network adapter
// XXX: Figure out if "network adapter" is the correct term here
func GetAddress(name string) (string, error) {
	interfaces, _ := net.Interfaces()
	for _, inter := range interfaces {
		if name == inter.Name {
			return inter.HardwareAddr.String(), nil
		}
	}
	return "", errors.New("Could not find network ID for that")
}
