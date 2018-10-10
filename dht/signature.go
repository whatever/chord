package dht

/**
 * SIGNATUR
 * SIGNATUR
 * SIGNATUR
 * SIGNATURE
 *
 * SIGNATUR
 */

// Signature
type Signature struct {
	public  string
	private string
	Node    ChordNode
}

// GetPublic returns a public signature
func (self *Signature) GetPublic() PublicSignature {
	return PublicSignature{
		Public: self.public,
		Node:   self.Node,
	}
}

// Matches returns whether a public signature matches its private one
func (self *Signature) Matches(rhs PublicSignature) bool {
	return self.private == rhs.Public
}

// ...
type PublicSignature struct {
	Public string
	Node   ChordNode
}
