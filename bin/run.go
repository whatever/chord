package main

import (
	"flag"
	"fmt"
	dht "github.com/internet-research-labs/chord/dht"
	_ "log"
	"os"
	"os/signal"
	"strings"
)

// Setup interrupt handler quietly in the background
func setupInterruptHandler(f func()) {
	interrupts := make(chan os.Signal, 1)
	signal.Notify(interrupts, os.Interrupt)
	go func() {
		<-interrupts
		f()
	}()
}

// Main
func main() {

	// Parse command-line args
	seeds := flag.String("seeds", "0.0.0.0:8080", "comma-delimeted list of bootstrap locations")
	addr := flag.String("addr", "0.0.0.0", "ip-address")
	port := flag.Int("l", -1, "port")
	flag.Parse()

	addresses := strings.Split(*seeds, ",")

	// Setup blocking channel
	done := make(chan bool, 1)

	// Setup interrupt handler quietly in the background
	setupInterruptHandler(func() {
		done <- true
	})

	for _, v := range addresses {
		fmt.Println(dht.Address(v))
	}

	server, err := dht.NewChordServer(*port)

	// Activate server conditional on what args are passed in
	if err == nil {
		server.Listen()
		defer server.Close()
	} else {
		done <- true
	}

	// Shutup, golang
	_ = addr
	_ = port

	// Block
	<-done
	fmt.Println("<EO> Things are properly closing")

	/*
		// x
		x := dht.ChordTable{
			Port: *listen,
		}

		// Listen for connections
		x.Listen()

		// Join and close
		x.Join()
		defer x.Close()

		time.Sleep(5 * time.Second)

		// Ayy
		fmt.Println(*addr, *port, x)
	*/
}
