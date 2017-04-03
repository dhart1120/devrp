package main

import (
	"fmt"
	"flag"
	"os"
	"strings"
	"strconv"
	"time"
	"net"
	"log"
	"io"
)

var portForwards = make(map[int]int)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "devrp is a simple reverse proxy for use during development. Use ctrl+c to stop forwarding.\n")
		flag.PrintDefaults()
	}
	ports := flag.String("p", "", "Forward ports with a comma seperated list of src:dest. Example -p 8080:80,2222:22 forwards 8080 to 80 and 2222 to 22.")
	flag.Parse()

	parsePortForwards(ports)
}

func parsePortForwards(ports *string) {
	if *ports != "" {
		// Split -p argument into pairs of ports
		// Example: 3333:3000,2222:22 {3333->3000, 2222->22}
		var portPairs = strings.Split(*ports, ",")
		for _, pair := range portPairs {
			forward := strings.SplitN(pair, ":", 2)
			if len(forward) == 2 {
				src, srcErr := strconv.Atoi(forward[0])
				dest, destErr := strconv.Atoi(forward[1])
				if srcErr != nil || destErr != nil || src == dest {
					log.Fatalf("Cannot forward ports %s\n", forward)
				}
				portForwards[src] = dest
			} else {
				log.Fatalf("Cannot forward ports %s\n", forward)
			}
		}
	}
}

func main() {
	for src, dest := range portForwards {
		fmt.Printf("Forwarding port %d to %d\n", src, dest)
		go acceptConnections(src, dest)
	}

	// Wait for connection on listening ports
	for len(portForwards) > 0 {
		time.Sleep(time.Second * 60)
	}

	// Print usage if there were no ports to forward
	flag.Usage()
}

func acceptConnections(src int, dest int) {
	// Listen on src port for tcp connections
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(src))
	if err != nil {
		log.Fatalf("Could not listen on source port %s\n", src)
	}

	// when a connection is accepted, forward to dest port
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Could not accept tcp connection on port %s\n", src)
		}
		fmt.Printf("Forwarding request on %d to %d\n", src, dest)
		go handleRequest(conn, dest)
	}
}

func handleRequest(conn net.Conn, destPort int) {
	client, err := net.Dial("tcp", ":"+strconv.Itoa(destPort))
	if err != nil {
		log.Fatalf("Could not forward to destination port %s\n", destPort)
	}

	// copy data both ways and cleanup
	go copyStream(conn, client)
	go copyStream(client, conn)
}

func copyStream(src net.Conn, dest net.Conn) {
	defer src.Close()
	defer dest.Close()
	io.Copy(src, dest)
}
