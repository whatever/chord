package dht

import (
	"encoding/json"
	"log"
)

type ChordWireMessage struct {
	// Source      ChordNode `json:"source"`
	// Destination ChordNode `json:"destination"`
	Type string `json:"type"`
}

// Equals returns whether two messages are equivalent
func Equals(lhs, rhs ChordWireMessage) bool {
	return lhs.Type == rhs.Type
}

// MESSAGE TYPES:
// JON

// EncodeWireMessage takes a wire message and returns a byte array
func EncodeWireMessage(decoded ChordWireMessage) (encoded []byte) {
	encoded, err := json.Marshal(decoded)
	if err != nil {
		return encoded
	}
	return encoded
}

// DecodeWireMessage takes a byte array and returns a DHT WireMessage
func DecodeWireMessage(encoded []byte) (decoded ChordWireMessage) {
	err := json.Unmarshal(encoded, &decoded)
	log.Println("!!!", string(encoded), decoded, err)
	return decoded
}

// Response

// JoinedResponse
type JoinedResponse struct {
	Prev  string `json:"prev"`
	Next  string `json:"next"`
	Id    string `json:"id"`
	Chain string `json:"chain"`
}
