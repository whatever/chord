package dht

import (
	"fmt"
	"log"
	"net"
)

type ChordTable struct {
	server net.Listener
	Port   int
}

type ChordTables []ChordTable

// Reads and closes 1024 bytes of TCP connection
// - It routes traffic to appropriate subhandler
func handle(conn net.Conn) {
	buf := make([]byte, 1024)
	conn.Read(buf)
	log.Println("Received message: ", string(buf))
	conn.Close()
}

// Listens for messages from valid messages from other nodes in the hash table
// - It passes all raw connections to {{handle}}.
func (self *ChordTable) Listen() {

	server, err := net.Listen("tcp", fmt.Sprintf(":%d", self.Port))

	if err != nil {
		log.Println("Failed listening on port")
		return
	}

	self.server = server

	// Turn on TCP listen-loop in the background
	go func() {

		// Continually accept
		for {
			conn, err := self.server.Accept()
			if err == nil {
				go handle(conn)
			}
		}
	}()
}

// Join a chord table
// - It returns a string of node ids
func (self *ChordTable) Join() ([]string, bool) {
	return []string{"red"}, true
}

// ...
func (self *ChordTable) Close() {
	log.Println("Closing")
	go self.server.Close()
}

// Get
func (self *ChordTable) Get(key string) (string, bool) {
	return "whatever", true
}

// Put records
func (self *ChordTable) Put(key string) (string, bool) {
	return "whatever", true
}

// Create Node in a Chord Table
func NewChordTable() *ChordTable {
	return nil
}
