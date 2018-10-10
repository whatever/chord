package dht

import (
	"fmt"
	"testing"
)

// Test that we make things
func TestRandomizer(t *testing.T) {
	randomizer := NewRandomizer("salt")
	fmt.Println((randomizer.GetToken()))
}
