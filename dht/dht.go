package dht

import (
	"fmt"
	"strconv"
	"strings"
)

// Defines the things you need to be a distribute-hash-table
type Dht interface {
	Listen()
	Join([]DhtAddress) ([]string, bool)
	Close()
	Get(string) (string, bool)
	Put(string) (string, bool)
}

// Specifies a node destination
// - This could probably be more generic
type DhtAddress struct {
	Ip   string
	Port int
}

// Return an DHT Address from a string
// - It should just handle the bullshit associated with mapping a string to a valid address
func Address(s string) DhtAddress {
	pieces := strings.SplitN(s, ":", 2)
	address := "127.0.01"
	port := 8080

	address = pieces[0]
	port, err := strconv.Atoi(pieces[1])

	if address == "" {
		address = "127.0.0.1"
	}

	if err != nil {
		port = 8081
	}

	return DhtAddress{address, port}
}

func (self *DhtAddress) String() string {
	return fmt.Sprintf("%s:%d", self.Ip, self.Port)
}

// Define type plural of DhtAddress
type DhtAddresses []DhtAddress
