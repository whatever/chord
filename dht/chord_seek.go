package dht

// PingRequest
type SeekRequest struct {
	Type   string    `json:"type"`
	Source ChordNode `json:"source,omitempty"`
	Token  string    `json:"token,omitempty"`
}

// OkayResponse
type OkayResponse struct {
	Type  string `json:"type"`
	Token string `json:"token"`
}

// PingResponse
type SeekResponse struct {
	Type   string     `json:"type"`
	Source *ChordNode `json:"source,omitempty"`
	Prev   *ChordNode `json:"next,omitempty"`
	Next   *ChordNode `json:"prev,omitempty"`
	Token  string     `json:"token,omitempty"`
}

// SendSeek sends a request to find out where a node belongs, and
// returns a token to verify the authenticity of the future request
func SendSeek(from, to ChordNode, token string) (string, error) {
	return "--nah--", nil
}

// HandleSeek
func (self *ChordTable) HandleSeek(req SeekRequest) string {
	return "--nah--"
}
