package dht

import (
	"crypto/sha256"
	_ "encoding/base64"
	"encoding/binary"
	"fmt"
	"net"
)

// Get networks
func ExternalIp() ([]string, error) {
	ips := make([]string, 0)
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			ips = append(ips, ip.String())
			// return ip.String(), nil
		}
	}
	return ips, nil
}

// GetNodeID returns a UUID for a node on a network with that port number.
func GetNodeID(name string, port int) uint {
	hasher := sha256.New()
	addr, err := GetAddress(name)
	var id uint
	if err == nil {
		hasher.Write([]byte(fmt.Sprintf("%s x_x %d", addr, port)))
		b := hasher.Sum(nil)[0:8]
		id := uint(binary.LittleEndian.Uint64(b))
		return id
	} else {
		id = 0
	}
	return id
}
