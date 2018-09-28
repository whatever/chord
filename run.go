package main

import (
	"flag"
	"fmt"
	_ "log"
	"time"
)

func main() {

	// Parse command-line args
	addr := flag.String("addr", "0.0.0.0", "ip-address")
	port := flag.Int("port", 8081, "port")
	listen := flag.Int("listen", 8081, "port")
	flag.Parse()

	// x
	x := ChordTable{
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
}
