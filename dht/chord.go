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

func (self *ChordNode) GetAddress() DhtAddress {
	return DhtAddress{self.Ip, self.Port}
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

	// This tells me who I belong to
	belonger   chan JoinedResponse
	randomizer Randomizer
}

// Plural of ChordTable
type ChordTables []ChordTable

// GetNode
func (self *ChordTable) GetNode() ChordNode {
	return ChordNode{
		Id:   self.Id,
		Ip:   self.Ip,
		Port: self.Port,
	}
}

// Return info about myself
func (self *ChordTable) Info() string {
	return self.Id
}

// handleJoin tries to find a place in the request chain for a node
func (self *ChordTable) handleJoin(joining *ChordNode) (JoinedResponse, error) {
	if self.Prev == nil && self.Next == nil {
		self.Prev = joining
		self.Next = joining
		node := self.GetNode()
		return JoinedResponse{
			Prev: node,
			Self: node,
			Next: node,
		}, nil
	}

	if self.Id == joining.Id {
		log.Println("WEIRD... I don't even know myself")
		return JoinedResponse{}, nil
	}

	// XXX: Need to figure out how to differentiate between
	// - "I want to join the network" and
	// - "Hey find a place for this node, then organize the update"

	// We found the right node to update the successor
	if self.Id < joining.Id && joining.Id < self.Next.Id {
		log.Println(self.Id, "<", joining.Id, "<", self.Next.Id)
		return JoinedResponse{}, nil
	}

	log.Println("so sure me")

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

// handleTopology closes the connection and passes on the request until it comes back to itself
func (self *ChordTable) handleTopology(req TopologyRequest) {

	req.Type = "topology"

	if req.Path == nil {
		req.Path = []string{}
	}

	// This is a complicated scenario over here
	// Basically:
	// - Empty requests pretend that this is the start, and pass the request on
	// - If they match, we've down a full circuit and should end early
	// - If the are not contained in the {{path}}, then we should append and continue
	// - If they are caontained in the {{path}}, then we should end, because we met a weird loop
	if req.Source == "" {
		req.Source = self.Id
		req.Path = []string{req.Source}
	} else if req.Source == self.Id {
		req.Path = append(req.Path, self.Id)
		return
	} else if !Contains(self.Id, req.Path) {
		req.Path = append(req.Path, self.Id)
	} else {
		log.Println("[[[WEIRD]]]")
		return
	}

	// XXX: Figure out why we're not getting the kind of result that we expect here

	b := EncodeStruct(req)
	log.Println(string(b))

	resp, err := sendMessage(
		DhtAddress{self.Next.Ip, self.Next.Port},
		b,
	)

	_ = resp
	_ = err
	// fmt.Println(resp, err)
}

// Reads and closes 1024 bytes of TCP connection
// - It routes traffic to appropriate subhandler
func (self *ChordTable) handle(conn net.Conn) {
	// What else
	buf := bytes.Trim(make([]byte, 1024), "\n\r")
	buf_length, _ := conn.Read(buf)
	buf = buf[0:buf_length]

	var n int
	var err error

	// Note that we need to truncate the byte array
	// log.Println("+++", string(buf))
	msg := DecodeWireMessage(buf)
	// log.Println("---", string(buf))

	log.Println("Received message:", msg.Type)

	switch msg.Type {
	case "who":
		n, err = conn.Write([]byte(self.Info()))
	case "info":
		r := InfoResponse{
			Id:   self.Id,
			Prev: self.Prev.Id,
			Next: self.Next.Id,
		}
		n, err = conn.Write(r.Bytes())
	case "topology":
		req := TopologyRequest{}
		DecodeStruct(buf, &req)
		self.handleTopology(req)
		n, err = conn.Write([]byte(self.Id))
	case "join":
		resp, _ := self.handleJoin(&msg.Source)
		n, err = conn.Write(EncodeStruct(resp))
	default:
		n, err = conn.Write([]byte("IGNORED"))
	}

	_ = n
	_ = err

	// {}}{

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

	// Listen to requests
	go func() {
		for joinRequest := range self.belonger {
			fmt.Println(joinRequest)
		}
	}()

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
	resp, _ := sendMessage(addr, msg)

	r := JoinedResponse{}

	DecodeStruct([]byte(resp), &r)

	self.Prev = &r.Prev
	self.Next = &r.Prev

	// log.Println(self.Id, ": JOIN RESPONSE :", string(resp), err)
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

		if err != nil {
			log.Panic("things went wrong")
		}

		table := ChordTable{
			Id:         GetNodeID("en0", port),
			Ip:         ip[0],
			Port:       port,
			seeds:      seeds,
			randomizer: NewRandomizer("see you in hell"),
		}

		// Return with everyything good
		return &table, nil
	}

	// Return with everythhing bad
	return nil, errors.New("Invalid port")
}

// GetSignature returns a message signature that needs to be caught
func (self *ChordTable) GetSignature() Signature {
	node := self.GetNode()
	nonce := self.randomizer.GetToken()
	return Signature{
		public:  nonce,
		private: nonce,
		Node:    node,
	}
}
