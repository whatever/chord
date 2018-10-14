package dht

import (
	"fmt"
	"testing"
)

func TestSeekBootstrap(t *testing.T) {
	alice, _ := NewChordServer(8000, nil)
	alice.Listen()
	defer alice.Close()

	bob, _ := NewChordServer(8001, nil)
	bob.Listen()
	defer bob.Close()

	resp := bob.RequestSeek(DhtAddress{alice.Ip, alice.Port})

	change := <-bob.ChangeLog

	if resp != nil {
		t.Error("response should contain type everything is okay", resp)
	}
}

func zzz() {
	fmt.Println("zzz")
}
