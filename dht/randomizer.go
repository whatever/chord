package dht

import (
	"fmt"
	"math/rand"
	"time"
)

// Randomizer handles all token generation
type Randomizer struct {
	salt string
	r    *rand.Rand
}

// NewRandomizer returns an object to generate new tokens
func NewRandomizer(salt string) Randomizer {
	seed := time.Now().UTC().UnixNano()
	fmt.Println("seed:", seed)
	source := rand.NewSource(seed)
	return Randomizer{
		salt: salt,
		r:    rand.New(source),
	}
}

// GetToken returns a random string of length N
// - Used to validate the source of a
func (self *Randomizer) GetToken() string {
	return fmt.Sprintf("%032d", self.r.Int())
}
