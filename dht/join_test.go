package dht

import (
	"fmt"
	"testing"
)

// XXX: Try writing some tests here, so we can have explicit assertions on how we're organized
// Test that things start correctly
func TestStart(t *testing.T) {
	seed, _ := newChordServer("001", 9000, nil)
	seed.Listen()
	defer seed.Close()

	node := seed.GetNode()
	bootstrap := []DhtAddress{}
	bootstrap = append(bootstrap, node.GetAddress())

	alice, _ := newChordServer("002", 9001, bootstrap)
	alice.Listen()
	defer alice.Close()
	alice.Join()

	if seed.Next.Id != alice.Id {
		t.Fail()
	}

	if alice.Next.Id != seed.Id {
		t.Fail()
	}
}

func xxx() {
	fmt.Println("x_x")
}
