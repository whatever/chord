package dht

import (
	"testing"
)

// Ensure that ping messages are properly structured
func TestPing(t *testing.T) {
	s, _ := NewChordServer(8000, nil)
	s.Listen()
	defer s.Close()

	resp := SendPing(
		ChordNode{0, "127.0.0.1", 7000},
		ChordNode{1, "127.0.0.1", 8000},
		"zap",
	)

	if resp.Source == nil || resp.Source.Port != 8000 {
		t.Error("Port should be 8000")
	}

	if resp.Token != "zap" {
		t.Error("Token is not equal to zap")
	}
}
