package dht

// JoinWaitResponse
type WaitResponse struct {
	Type  string    `json:"type"`
	Token string    `json:"token"`
	Node  ChordNode `jsoin:"node"`
}
