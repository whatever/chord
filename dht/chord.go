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
	// XXX: port over to integer to make circle topology easier to think about
	Id   uint   `json:"id,omitempty"`
	Ip   string `json:"ip,omitempty"`
	Port int    `json:"port,omitempty"`
}

func (self *ChordNode) GetAddress() DhtAddress {
	return DhtAddress{self.Ip, self.Port}
}

// Return a string representation of a node
func (self ChordNode) String() string {
	return fmt.Sprintf("(%d, %s, %d)", self.Id, self.Ip, self.Port)
}

// An actual table that one can use to get data
type ChordTable struct {
	server net.Listener
	seeds  DhtAddresses
	Id     uint
	Ip     string
	Port   int

	// /shrug
	Alive bool

	// No finger table yet
	Prev *ChordNode
	Next *ChordNode

	// list of nodes that I've met in the past
	seen      []ChordNode
	Responses chan bool

	// This tells me who I belong to
	belonger   chan JoinedResponse
	randomizer Randomizer
	waiting    map[string]ChordWireMessage
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
func (self *ChordTable) Info() uint {
	return self.Id
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
		req.Path = []uint{}
	}

	// This is a complicated scenario over here
	// Basically:
	// - Empty requests pretend that this is the start, and pass the request on
	// - If they match, we've down a full circuit and should end early
	// - If the are not contained in the {{path}}, then we should append and continue
	// - If they are caontained in the {{path}}, then we should end, because we met a weird loop
	if req.Source == 0 && req.Path == nil {
		req.Source = self.Id
		req.Path = []uint{req.Source}
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

	switch msg.Type {

	// Just saying "sup... are you alive?"
	case "ping":
		req := PingRequest{}
		DecodeStruct(buf, &req)
		res := self.HandlePing(req)
		n, err = conn.Write([]byte(EncodeWireMessage(res)))

	// SEEK: Finding
	case "seek":
		req := SeekRequest{}
		DecodeStruct(buf, &req)
		res, _ := self.HandleSeek(req)
		log.Println("???", req)
		n, err = conn.Write([]byte(EncodeWireMessage(res)))

	// JOIN: Find and join in a single operation (might have incosistencies)
	case "join":
		resp, _ := self.handleJoin(msg)
		n, err = conn.Write(EncodeStruct(resp))

	// INFO: Get some details about the node
	case "info":
		r := InfoResponse{
			Id:   self.Id,
			Prev: self.Prev.Id,
			Next: self.Next.Id,
		}
		n, err = conn.Write(r.Bytes())

	// TOPOLOGY: Kick off a topology scan
	case "topology":
		req := TopologyRequest{}
		DecodeStruct(buf, &req)
		self.handleTopology(req)
		n, err = conn.Write([]byte(string(self.Id)))

	// DEFAULT: x_x
	default:
		n, err = conn.Write([]byte("IGNORED"))
		fmt.Println("!!!<<<")
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
	return fmt.Sprintf("[%d // %d]", self.Id, self.Port)
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

// Printlnf prints a log line
func (self *ChordTable) Println(line string, vars ...interface{}) {
	id := fmt.Sprintf("[%02d]", self.Id)
	line = fmt.Sprintf(line, vars...)
	log.Printf(id + " " + line)
}

// Returns whether this node needs to be bootstrapped (has no one in the world)
func (self *ChordTable) NeedsBootstrap() bool {
	return self.Prev == nil && self.Next == nil
}

// Returns whether this node discovered its predecessor
func (self *ChordTable) IsPredecessor(cand ChordNode) bool {
	switch {

	// If we're boostrapping
	case self.Prev == nil:
		return true

	// If it's the actual predecessor
	case self.Id > cand.Id && cand.Id > self.Prev.Id:
		return true

	// If we're the first node
	case self.Id < self.Prev.Id:
		return self.Id > cand.Id

	default:
		return false
	}
}

// Returns whether this node discovered a its successor
func (self *ChordTable) IsSuccessor(cand ChordNode) bool {
	switch {

	// If we're bootstrapping a successor
	case self.Next == nil:
		return true

	// If this node actually fits
	case self.Id < cand.Id && cand.Id < self.Next.Id:
		return true

	// If we're the last node
	case self.Id > self.Next.Id:
		return self.Id < cand.Id

	default:
		return false
	}
}

// SendUpdate gives some details
func RequestTopology(addr DhtAddress) string {
	return "x_x"
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
func newChordServer(id uint, port int, bootstrap DhtAddresses) (*ChordTable, error) {
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
			Id:         id,
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

func NewChordServer(port int, bootstrap DhtAddresses) (*ChordTable, error) {
	id := GetNodeID("en0", port)
	return newChordServer(id, port, bootstrap)
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
