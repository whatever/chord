package dht

import (
	"fmt"
	"testing"
)

func x_o() {
	fmt.Println("x_x")
}

// Ensure that ping messages are properly structured
func TestPingMessage(t *testing.T) {

	message := ChordWireMessage{
		"join",
		ChordNode{1, "127.0.0.1", 8080},
		3,
	}

	encoded := EncodeWireMessage(message)
	decoded := DecodeWireMessage(encoded)

	if !Equals(message, decoded) {
		t.Fail()
	}
}

// Ensure that wire messages are bijective
func TestWireMessage(t *testing.T) {
}

// Ensure that wire messages are decodable
func TestDecodeWireMessage(t *testing.T) {
}

// Ensure that wire messages are encodable
func TestEncodeWireMessage(t *testing.T) {
}

func TestGetNodeID(t *testing.T) {
	fmt.Println(GetNodeID("en0", 8081))
}
