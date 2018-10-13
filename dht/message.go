package dht

import (
	"encoding/json"
	"log"
)

// Generic Requests
// Generic Requests
// Generic Requests
// Generic Requests
// Generic Requests

// EncodeWireMessage takes a wire message and returns a byte array
func EncodeWireMessage(decoded interface{}) (encoded []byte) {
	encoded, err := json.Marshal(decoded)
	if err != nil {
		return encoded
	}
	return encoded
}

// DecodeWireMessage takes a byte array and returns a DHT WireMessage
func DecodeWireMessage(encoded []byte) (decoded ChordWireMessage) {
	err := json.Unmarshal(encoded, &decoded)
	if err != nil {
		return decoded
	}
	return decoded
}

// Wire Protocol REQUEST Formats
// Wire Protocol REQUEST Formats
// Wire Protocol REQUEST Formats
// Wire Protocol REQUEST Formats
// Wire Protocol REQUEST Formats

// WireMessageRequestFormat
type ChordWireMessage struct {
	// Destination ChordNode `json:"destination"`
	Type   string    `json:"type"`
	Source ChordNode `json:"source"`
	Hops   int       `json:"hops"`
}

// Equals returns whether two messages are equivalent
func Equals(lhs, rhs ChordWireMessage) bool {
	return lhs.Type == rhs.Type && true
}

// Wire Protocol RESPONSE Formats
// Wire Protocol RESPONSE Formats
// Wire Protocol RESPONSE Formats
// Wire Protocol RESPONSE Formats
// Wire Protocol RESPONSE Formats

type TopologyRequest struct {
	Type   string `json:"type"`
	Source uint   `json:"source"`
	Path   []uint `json:"path"`
}

// Response Interface
type WireResponse interface {
	Bytes() []byte
}

// InfoResponse
type InfoResponse struct {
	Id   uint `json:"id"`
	Prev uint `json:"prev"`
	Next uint `json:"next"`
}

// Return
func EncodeStruct(decoded interface{}) (encoded []byte) {
	encoded, err := json.Marshal(decoded)
	if err != nil {
		log.Println("!!!!")
		return encoded
	}
	return encoded
}

// DecodeWireMessage takes a byte array and returns a DHT WireMessage
func DecodeStruct(encoded []byte, decoded interface{}) {
	err := json.Unmarshal(encoded, decoded)
	if err != nil {
		log.Println("ERROR", err)
	}
}

// Bytes returns a byte-array representing
func (self *InfoResponse) Bytes() []byte {
	bytes := EncodeStruct(*self)
	return bytes
}

// JoinedResponse
type JoinedResponse struct {
	Prev ChordNode `json:"prev"`
	Self ChordNode `json:"self"`
	Next ChordNode `json:"next"`
}
