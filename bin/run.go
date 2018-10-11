package main

import (
	"flag"
	"fmt"
	dht "github.com/internet-research-labs/chord/dht"
	"log"
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

// Return a list of addresses
func getAddresses(adds []string) dht.DhtAddresses {
	addresses := []dht.DhtAddress{}
	for _, v := range adds {
		addresses = append(addresses, dht.Address(v))
	}
	return addresses
}

// Main
func main() {
	// Parse command-line args
	seeds := flag.String("seeds", "", "comma-delimeted list of bootstrap locations")
	port := flag.Int("l", -1, "port")
	flag.Parse()

	// Setup blocking channel
	done := make(chan bool, 1)

	// Setup interrupt handler quietly in the background
	setupInterruptHandler(func() {
		done <- true
	})

	// Point to self
	if *seeds == "" {
		tmp := fmt.Sprintf(":%d", *port)
		seeds = &tmp
	}

	// Initialize server
	server, err := dht.NewChordServer(
		*port,
		getAddresses(strings.Split(*seeds, ",")),
	)

	// Activate server conditional on what args are passed in
	if err == nil {
		log.Printf("Listening on %d\n", *port)
		server.Listen()
		server.RequestJoin()
		// XXX: Get this script up-and-running
		defer server.Close()
	} else {
		done <- true
	}

	// Block
	<-done

	// Exit message
	log.Println("<EO> Things are properly closing")
}
