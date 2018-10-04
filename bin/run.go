package main

import (
	"flag"
	"fmt"
	dht "github.com/internet-research-labs/chord/dht"
	_ "log"
	"os"
	"os/signal"
	"strings"
	"time"
)

// Setup interrupt handler quietly in the background
func setupInterruptHandler() {
	interrupts := make(chan os.Signal, 1)
	signal.Notify(interrupts, os.Interrupt)
	go func() {
		<-interrupts
		fmt.Println("[Interrupt caught... let's clean things up quietly]")
		os.Exit(0)
	}()
}

// Main
func main() {

	// Parse command-line args
	seeds := flag.String("seeds", "0.0.0.0:8080", "comma-delimeted list of bootstrap locations")
	addr := flag.String("addr", "0.0.0.0", "ip-address")
	port := flag.Int("port", -1, "port")
	listen := flag.Int("listen", 8081, "port")
	flag.Parse()

	addresses := strings.Split(*seeds, ",")

	// Setup interrupt handler quietly in the background
	setupInterruptHandler()

	for _, v := range addresses {
		fmt.Println(dht.Address(v))
	}

	server, err := dht.NewChordServer(*port)

	if err == nil {
		server.Listen()
		defer server.Close()
	} else {
		fmt.Println(err)
	}

	_ = addr
	_ = port
	x := dht.ChordTable{Port: *listen}
	_ = x

	time.Sleep(1 * time.Second)

	if server != nil {
		for i := 0; i < 10; i++ {

			time.Sleep(1 * time.Second)
			fmt.Println("TICK", server.Alive)
		}
	}

	fmt.Println("<EO>")
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
