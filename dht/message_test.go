package dht

import (
	"testing"
)

// Ensure that ping messages are properly structured
func TestPingMessage(t *testing.T) {

	message := ChordWireMessage{
		ChordNode{},
		ChordNode{},
		"join",
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
