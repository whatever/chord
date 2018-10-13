package dht

// PingRequest
type PingRequest struct {
	Type   string    `json:"type,omitempty"`
	Source ChordNode `json:"source,omitempty"`
	Token  string    `json:"token,omitempty"`
}

// PingResponse
type PingResponse struct {
	Source *ChordNode `json:"source,omitempty"`
	Token  string     `json:"token,omitempty"`
}

// SendPing sends a ping message to a node address
func SendPing(from, to ChordNode, token string) PingResponse {
	req := PingRequest{
		"ping",
		from,
		token,
	}
	resp, _ := sendMessage(to.GetAddress(), EncodeWireMessage(req))
	res := PingResponse{}
	DecodeStruct([]byte(resp), &res)
	return res
}

// HandlePing
func (self *ChordTable) HandlePing(req PingRequest) PingResponse {
	node := self.GetNode()
	return PingResponse{
		Source: &node,
		Token:  req.Token,
	}
}
