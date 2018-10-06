package dht

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
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

// x_x
func x_x() {
	fmt.Println("...")
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

// GetNodeID returns a UUID for a node on a network with that port number.
func GetNodeID(name string, port int) string {
	hasher := sha256.New()
	addr, err := GetAddress(name)
	var id string
	if err == nil {
		hasher.Write([]byte(fmt.Sprintf("%s////%d", addr, port)))
		id = base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	} else {
		id = "..."
	}
	return id
}
