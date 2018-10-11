package dht

import (
	"fmt"
	"testing"
)

// XXX: Try writing some tests here, so we can have explicit assertions on how we're organized
// Test that things start correctly
func TestMultipleJoins(t *testing.T) {
	seed, _ := newChordServer(uint(1), 9000, nil)
	seed.Listen()
	defer seed.Close()

	node := seed.GetNode()
	bootstrap := []DhtAddress{}
	bootstrap = append(bootstrap, node.GetAddress())

	alice, _ := newChordServer(uint(2), 9001, bootstrap)
	alice.Listen()
	defer alice.Close()
	alice.RequestJoin()

	if seed.Next.Id != alice.Id {
		t.Fail()
	}

	/*
		if alice.Next.Id != seed.Id {
			t.Fail()
		}

			bob, _ := newChordServer("003", 9002, bootstrap)
			bob.Listen()
			defer bob.Close()
			bob.RequestJoin()

			// XXX: If this works, then we have some basic joining
			if alice.Next.Id != bob.Id {
				t.Fail()
			}
	*/
}

func xxx() {
	fmt.Println("x_x")
}
