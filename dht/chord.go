package dht

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strings"
)

type ChordTable struct {
	server net.Listener
	seeds  DhtAddresses
	Id     string
	Port   int
	Alive  bool
}

type ChordTables []ChordTable

// Reads and closes 1024 bytes of TCP connection
// - It routes traffic to appropriate subhandler
func (self *ChordTable) handle(conn net.Conn) {
	buf := bytes.Trim(make([]byte, 1024), "\n\r")
	length, _ := conn.Read(buf)

	// Note that we need to truncate the byte array
	msg := strings.TrimRight(string(buf[0:length]), "\n\r")

	switch msg {
	case "who":
		n, err := conn.Write(([]byte)(self.Id))
		log.Println(">>", n, err)
	default:
		log.Printf("I don't recognize message \"%s\"\n", msg)
	}

	conn.Close()
}

// Listens for messages from valid messages from other nodes in the hash table
// - It passes all raw connections to {{handle}}.
func (self *ChordTable) Listen() {

	server, err := net.Listen("tcp", fmt.Sprintf(":%d", self.Port))

	if err != nil {
		log.Println("Failed listening on port:", err)
		return
	}

	self.server = server
	self.Alive = true

	// Turn on TCP listen-loop in the background
	go func() {
		// Continually accept
		for {
			conn, err := self.server.Accept()
			if err == nil {
				go self.handle(conn)
			}
		}
	}()
}

func (self ChordTable) String() string {
	return fmt.Sprintf("[%s // %d]", self.Id, self.Port)
}

// Join a chord table
// - It returns a string of node ids
func (self *ChordTable) Join() ([]string, bool) {
	fmt.Println(*self, "::", self.seeds)
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
func NewChordServer(port int, bootstrap DhtAddresses) (*ChordTable, error) {

	// ...
	if port > 0 {
		seeds := make(DhtAddresses, 0)

		// Only try talking to people who aren't you
		for _, v := range bootstrap {
			if v.Port != port {
				seeds = append(seeds, v)
			}
		}

		// For now... initiate with random integer as node id
		// In the future this will be deterministic
		table := ChordTable{
			Id:    fmt.Sprintf("%d", rand.Int()),
			Port:  port,
			seeds: seeds,
		}

		// Return with everyything good
		return &table, nil
	}

	// Return with everythhing bad
	return nil, errors.New("Invalid port")
}
