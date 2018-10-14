package dht

import (
	"fmt"
	"testing"
)

func TestSeek(t *testing.T) {
	alice, _ := NewChordServer(8000, nil)
	alice.Listen()
	defer alice.Close()

	bob, _ := NewChordServer(8001, nil)
	bob.Listen()
	defer bob.Close()

	resp := bob.RequestSeek(DhtAddress{alice.Ip, alice.Port})

	// XXX: Make this channel live to help synchronize events
	<-bob.Responses

	if resp != nil {
		t.Error("response should contain type everything is okay")
	}
}

func zzz() {
	fmt.Println("zzz")
}
