package dht

// FigureOutTopology contains a method for figuring out the topology of the network
// XXX: Figure out the best way to do something like this
func FigureOutTopology() {
}

// Contains returns whether a slice contains a given value
func Contains(needle uint, haystack []uint) bool {
	for _, maybe := range haystack {
		if needle == maybe {
			return true
		}
	}
	return false
}
