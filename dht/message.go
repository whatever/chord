package dht

type ChordWireMessage struct {
	Source ChordNode
	Type   string // JOIN
	Other  []string
}

// MESSAGE TYPES:
// JON

// EncodeWireMessage takes a wire message and returns a byte array
func EncodeWireMessage(decoded ChordWireMessage) (encoded []byte) {
	return []byte{}
}

// DecodeWireMessage takes a byte array and returns a DHT WireMessage
func DecodeWireMessage(encoded []byte) (decoded ChordWireMessage) {
	return ChordWireMessage{
		Source: ChordNode{},
		Type:   "ping",
		Other:  nil,
	}
}
