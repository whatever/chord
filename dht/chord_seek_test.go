package dht

import (
	"fmt"
	"testing"
)

func TestSeek(t *testing.T) {
	s, _ := NewChordServer(8000, nil)
	s.Listen()
	defer s.Close()

	resp, _ := SendSeek(
		ChordNode{23, "127.0.0.1", 9000},
		ChordNode{24, "127.0.0.1", 9001},
		"Busted",
	)

	if resp != "Busted" {
		t.Error("response should contain type everything is okay")
	}
}

func zzz() {
	fmt.Println("zzz")
}
