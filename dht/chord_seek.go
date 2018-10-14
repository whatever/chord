package dht

import (
	"errors"
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
func SendSeek(from, to ChordNode, token string) (StatusResponse, error) {

	encoded := EncodeStruct(SeekRequest{"seek", from, "ayyyy"})

	r, err := sendMessage(to.GetAddress(), encoded)

	log.Println(r, err)

	res := StatusResponse{
		Type:  "okay",
		Token: "Busted",
	}
	return res, nil
}

// RequestSeek sends a request to find where it belongs in the network
func (self *ChordTable) RequestSeek(addr DhtAddress) error {
	return errors.New("FAILED")
}

// HandleSeek
func (self *ChordTable) HandleSeek(req SeekRequest) (StatusResponse, error) {
	res := StatusResponse{
		Type:    "okay",
		Message: "what",
		Token:   "---",
	}

	var node ChordNode

	if self.IsSuccessor(req.Source) {
		r := SeekResponse{
			Type:   "seek!",
			Source: nil,
			Prev:   &node,
			Next:   self.Next,
			Token:  req.Token,
		}
		msg, _ := sendMessage(req.Source.GetAddress(), EncodeStruct(r))
		log.Println("MSG:", msg)
	} else {
	}

	return res, nil
}
