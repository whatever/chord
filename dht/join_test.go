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

	if seed.Next == nil || seed.Next.Id != alice.Id {
		t.Error("seed does not point to alice")
	}

	if alice.Prev == nil || alice.Prev.Id != seed.Id {
		t.Error("alice does not point back to seed", alice, seed)
	}

	if alice.Next == nil || alice.Next.Id != seed.Id {
		t.Error("alice does point forward to seed")
	}

	bob, _ := newChordServer(3, 9002, bootstrap)
	bob.Listen()
	defer bob.Close()
	bob.RequestJoin()

	// XXX: If this works, then we have some basic joining
	if bob.Next == nil {
		t.Error("bob did not properly join the circle:", bob.Next)
	}

	if alice.Next == nil || alice.Next.Id != bob.Id {
		t.Error("alice does not point to bob:", alice, "!=", bob)
	}
}

func xxx() {
	fmt.Println("x_x")
}
