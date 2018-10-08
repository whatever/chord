package dht

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net"
)

// A representation of a node - contains and ID and an address
type ChordNode struct {
	Id   string
	Ip   string
	Port int
}

// Return a string representation of a node
func (self ChordNode) String() string {
	return fmt.Sprintf("(%s, %s, %d)", self.Id, self.Ip, self.Port)
}

// An actual table that one can use to get data
type ChordTable struct {
	server net.Listener
	seeds  DhtAddresses
	Id     string
	Port   int

	// /shrug
	Alive bool

	// No finger table yet
	Prev *ChordNode
	Next *ChordNode
}

// Plural of ChordTable
type ChordTables []ChordTable

// Return info about myself
func (self *ChordTable) Info() string {
	return self.Id
}

// handleJoin tries to connect
// - It ...
func (self *ChordTable) handleJoin(joining *ChordNode) (bool, error) {
	if self.Prev == nil && self.Next == nil {
		fmt.Println("!!!!", *joining)
	}
	return true, nil
}

// Reads and closes 1024 bytes of TCP connection
// - It routes traffic to appropriate subhandler
func (self *ChordTable) handle(conn net.Conn) {
	buf := bytes.Trim(make([]byte, 1024), "\n\r")
	_, _ = conn.Read(buf)

	// Note that we need to truncate the byte array
	// msg := strings.TrimRight(string(buf[0:length]), "\n\r")
	msg := DecodeWireMessage(buf)

	var n int
	var err error

	switch msg.Type {
	case "who":
		n, err = conn.Write(([]byte)(self.Info()))
	case "join":
		fmt.Println("::", DecodeWireMessage(buf))
		// TODO: Define wire protocol
		// self.handleJoin(nil)
		n, err = conn.Write(([]byte)(self.Info()))
	case "ping":
		n, err = conn.Write(([]byte)("-"))
	case "put":
		n, err = conn.Write(([]byte)("-"))
	case "get":
		n, err = conn.Write(([]byte)("-"))
	default:
		n, err = conn.Write([]byte("IGNORED"))
	}

	// {}}{
	log.Println(self.Id, ":Received message:", msg, n, err)

	// -_-
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

// String outputs a human readable representation of the ChordTable
func (self ChordTable) String() string {
	return fmt.Sprintf("[%s // %d]", self.Id, self.Port)
}

// readConn returns a message <= 1024 bytes from a TCP connection
func readConn(conn net.Conn) string {
	buff := make([]byte, 1024)
	n, err := conn.Read(buff)
	if err != nil {
		return ""
	}
	return string(buff)[0:n]
}

// Send message to DHT
func sendMessage(addr DhtAddress, message string) (string, error) {
	conn, err := net.Dial("tcp", addr.String())

	if err != nil {
		return "", errors.New("Could not connect to address")
	}

	conn.Write([]byte(message))

	buff := make([]byte, 1024)
	n, err := conn.Read(buff)

	if err != nil {
		return "", errors.New("Could not read from address")
	}

	return string(buff)[0:n], nil
}

// hello pings a server
func (self *ChordTable) hello(addr DhtAddress) ChordNode {
	resp, err := sendMessage(addr, "who")
	log.Println("Got response: ", resp, err)
	return ChordNode{
		Id:   self.Id,
		Ip:   addr.Ip,
		Port: addr.Port,
	}
}

// join sends a request to join
func (self *ChordTable) join(addr DhtAddress) ChordNode {
	result, err := sendMessage(addr, "join")
	_ = result
	_ = err
	return ChordNode{}
}

// Join a chord table
// - It returns a string of node ids
func (self *ChordTable) Join() ([]string, bool) {
	if len(self.seeds) > 0 {
		addr := self.seeds[0]
		log.Println("HELLO:", self.hello(addr))
		log.Println("JOIN:", self.join(addr))
	}
	return []string{"red"}, true
}

// Close closes the TCP server
func (self *ChordTable) Close() {
	log.Println("Closing")
	go self.server.Close()
}

// Get returns the value of an entry
func (self *ChordTable) Get(key string) (string, bool) {
	return "whatever", true
}

// Put sets the value of a record on the DHT
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
			Id:    GetNodeID("en0", port),
			Port:  port,
			seeds: seeds,
		}

		// Return with everyything good
		return &table, nil
	}

	// Return with everythhing bad
	return nil, errors.New("Invalid port")
}
