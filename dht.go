package main

// Defines the things you need to be a distribute-hash-table
type Dht interface {
	Listen()
	Join(string) ([]string, bool)
	Close()
	Get(string) (string, bool)
	Put(string) (string, bool)
}
