package dht

import (
	"log"
)

// StatusResponse
type StatusResponse struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Token   string `json:"token"`
}

// SeekRequest
type SeekRequest struct {
	Type   string    `json:"type"`
	Source ChordNode `json:"source,omitempty"`
	Token  string    `json:"token,omitempty"`
}

// SeekResponse
type SeekResponse struct {
	Type   string     `json:"type"`
	Source *ChordNode `json:"source,omitempty"`
	Prev   *ChordNode `json:"next,omitempty"`
	Next   *ChordNode `json:"prev,omitempty"`
	Token  string     `json:"token,omitempty"`
}

// SendSeek sends a request to find out where a node belongs, and
// returns a token to verify the authenticity of the future request
func SendSeek(from ChordNode, to DhtAddress, token string) (StatusResponse, error) {

	encoded := EncodeStruct(SeekRequest{"seek?", from, "ayyyy"})

	r, err := sendMessage(to, encoded)

	log.Println(r, err)

	res := StatusResponse{
		Type:  "okay",
		Token: "Busted",
	}
	return res, nil
}

// RequestSeek sends a request to find where it belongs in the network
func (self *ChordTable) RequestSeek(addr DhtAddress) error {
	SendSeek(
		self.GetNode(),
		addr,
		"what",
	)
	return nil
}

// HandleSeek
func (self *ChordTable) HandleSeek(req SeekRequest) (StatusResponse, error) {
	res := StatusResponse{
		Type:    "okay",
		Message: "what",
		Token:   "---",
	}

	node := self.GetNode()

	if self.IsSuccessor(req.Source) {
		var next *ChordNode

		if self.Next == nil {
			next = &node
		} else {
			next = self.Next
		}

		prev := &node

		r := SeekResponse{
			Type:   "seek!",
			Source: &req.Source,
			Prev:   prev,
			Next:   next,
			Token:  req.Token,
		}
		sendMessage(
			req.Source.GetAddress(),
			EncodeStruct(r),
		)
	} else {
	}

	return res, nil
}
