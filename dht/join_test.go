package dht

import (
	"fmt"
	"testing"
)

// XXX: Try writing some tests here, so we can have explicit assertions on how we're organized
// Test that things start correctly
func TestStart(t *testing.T) {
	seed, _ := NewChordServer(9000, nil)
	seed.Listen()
	defer seed.Close()

	fmt.Println(">>", seed.GetNode())
}
