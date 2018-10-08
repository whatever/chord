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
	Ip     string
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

// handleJoin tries to find a place in the request chain for a node
func (self *ChordTable) handleJoin(joining *ChordNode) (JoinedResponse, error) {
	if self.Prev == nil && self.Next == nil {
		self.Prev = joining
		self.Next = joining
		return JoinedResponse{
			self.Id,
			self.Id,
			self.Id,
			"",
		}, nil
	} else {
		fmt.Println("Let the wild rumpus start")
	}
	return JoinedResponse{}, nil
}

// GetInfo returns meta information about Prev -> Self -> Next nodes in the chain
func (self *ChordTable) GetInfo() InfoResponse {
	return InfoResponse{
		Id:   self.Id,
		Prev: self.Prev.Id,
		Next: self.Next.Id,
	}
}

// Reads and closes 1024 bytes of TCP connection
// - It routes traffic to appropriate subhandler
func (self *ChordTable) handle(conn net.Conn) {
	// What else
	buf := bytes.Trim(make([]byte, 1024), "\n\r")
	buf_length, _ := conn.Read(buf)
	buf = buf[0:buf_length]

	/*
		fmt.Println("<<<<")
		fmt.Println(conn.LocalAddr())
		fmt.Println(conn.RemoteAddr())
		fmt.Println(">>>>")
	*/

	// Note that we need to truncate the byte array
	msg := DecodeWireMessage(buf)

	var n int
	var err error

	switch msg.Type {
	case "info":
		r := InfoResponse{
			Id:   self.Id,
			Prev: self.Prev.Id,
			Next: self.Next.Id,
		}
		n, err = conn.Write(r.Bytes())
	case "who":
		n, err = conn.Write(([]byte)(self.Info()))
	case "join":
		resp, _ := self.handleJoin(&msg.Source)
		n, err = conn.Write(EncodeStruct(resp))
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
func sendMessage(addr DhtAddress, message []byte) (string, error) {
	conn, err := net.Dial("tcp", addr.String())

	if err != nil {
		return "", errors.New("Could not connect to address")
	}

	conn.Write(message)

	buff := make([]byte, 1024)
	n, err := conn.Read(buff)

	fmt.Println("~~~~", string(buff))

	if err != nil {
		return "", errors.New("Could not read from address")
	}

	return string(buff)[0:n], nil
}

// getNode returns
func (self *ChordTable) getNode() ChordNode {
	return ChordNode{
		self.Id,
		"127.0.0.1",
		self.Port,
	}
}

// join sends a request to join
func (self *ChordTable) join(addr DhtAddress) ChordNode {
	msg := EncodeWireMessage(ChordWireMessage{
		"join",
		self.getNode(),
	})
	resp, err := sendMessage(addr, msg)
	log.Println(self.Id, ": JOIN RESPONSE :", string(resp), err)
	return ChordNode{}
}

// Join a chord table
// - It returns a string of node ids
func (self *ChordTable) Join() ([]string, bool) {
	if len(self.seeds) > 0 {
		addr := self.seeds[0]
		self.join(addr)
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
		ip, err := ExternalIp()
		_ = ip
		if err != nil {
			log.Panic("things went wrong")
		}
		table := ChordTable{
			Id:    GetNodeID("en0", port),
			Ip:    "",
			Port:  port,
			seeds: seeds,
		}

		// Return with everyything good
		return &table, nil
	}

	// Return with everythhing bad
	return nil, errors.New("Invalid port")
}
