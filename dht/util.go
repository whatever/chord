package dht

// FigureOutTopology contains a method for figuring out the topology of the network
// XXX: Figure out the best way to do something like this
func FigureOutTopology() {
}

func Contains(needle string, haystack []string) bool {
	for _, maybe := range haystack {
		if needle == maybe {
			return true
		}
	}
	return false
}
