package dht

import (
	"fmt"
	"testing"
)

// XXX: Try writing some tests here, so we can have explicit assertions on how we're organized
// Test that things start correctly
func TestMultipleJoins(t *testing.T) {
	seed, _ := newChordServer(1, 9000, nil)
	seed.Listen()
	defer seed.Close()

	node := seed.GetNode()
	bootstrap := []DhtAddress{}
	bootstrap = append(bootstrap, node.GetAddress())

	alice, _ := newChordServer(2, 9001, bootstrap)
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

	bob, _ := newChordServer(9, 9002, bootstrap)
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

	carlos, _ := newChordServer(4, 9003, bootstrap)
	carlos.Listen()
	defer carlos.Close()
	carlos.RequestJoin()

	// alice -> carlos -
	if alice.Next.Id != carlos.Id && carlos.Prev.Id != alice.Id && carlos.Next.Id != bob.Id {
		t.Error("Carlos did not properly enter the pool")
	}
}

// Test that more complicated starts work
func TestRequestToplogy(t *testing.T) {

	// Get seed
	seed, _ := newChordServer(420, 8000, nil)
	seed.Listen()
	defer seed.Close()

	node := seed.GetNode()
	bootstrap := []DhtAddress{}
	bootstrap = append(bootstrap, node.GetAddress())

	// Make network
	network := make([]*ChordTable, 10)
	for i, _ := range network {
		network[i], _ = newChordServer(96+1000*uint(i), 9000+i, bootstrap)
	}

	// Put everyone on the network
	for _, v := range network {
		v.Listen()
		defer v.Close()
	}

	// Request some joins
	for _, v := range network {
		v.RequestJoin()
	}

	if network[2].Next.Id != network[3].Id {
		t.Error("2 should point at 3")
	}

	resp := RequestTopology(bootstrap[0])

	if resp != "what?" {
		t.Error("Need to figure out this thing")
	}

	// x_x
	network = append(network, seed)
}

// xxx does nothing
func xxx() {
	fmt.Println("x_x")
}
