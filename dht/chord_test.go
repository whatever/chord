package dht

import (
	"fmt"
	"net"
	"testing"
	"time"
)

// MockConn
// MockConn
// MockConn

type MockConn struct {
	Message []byte
}

func (self MockConn) Read(b []byte) (n int, err error) {
	for i, c := range self.Message {
		b[i] = c
	}
	return len(self.Message), nil
}

func (self MockConn) Write(b []byte) (n int, err error) {
	return 0, nil
}

func (self MockConn) Close() error {
	return nil
}

func (self MockConn) LocalAddr() net.Addr {
	return MockAddr{}
}

func (self MockConn) RemoteAddr() net.Addr {
	return MockAddr{}
}

func (self MockConn) SetDeadline(t time.Time) error {
	return nil
}

func (self MockConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (self MockConn) SetWriteDeadline(t time.Time) error {
	return nil
}

// Mock Address
// Mock Address
// Mock Address

type MockAddr struct {
}

func (self MockAddr) Network() string {
	return "matt"
}

func (self MockAddr) String() string {
	return ""
}

func TestHandle(t *testing.T) {
	table := ChordTable{
		Id:    1,
		Port:  -1,
		Alive: true,
	}

	_ = table

	// conn := MockConn{Message: []byte("you know")}
	//// table.handle(conn)
}

func x_x() {
	fmt.Println("x_x")
}
