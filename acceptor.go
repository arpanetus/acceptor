package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
)

func handleConnection(c net.Conn) {
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())

	if _, err := c.Write([]byte(string("ok"))); err != nil {
		fmt.Printf("couldn't write into the conn: %+v", err)
	}

	if err := c.Close(); err != nil {
		fmt.Printf("couldn't close the conn: %+v", err)
		return
	}
}

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide a port range.")
		return
	}

	fromPortStr := arguments[1]
	tilPortStr := arguments[2]

	fromPort, err := strconv.ParseInt(fromPortStr, 10, 17)
	if err != nil {
		log.Panicf("can't parse the from port: %+v", err)
	}

	tilPort, err := strconv.ParseInt(tilPortStr, 10, 17)
	if err != nil {
		log.Panicf("can't parse the til port: %+v", err)
	}
	log.Printf("spawning tcp servers in range from %d til %d", fromPort, tilPort)

	toListen := make([]net.Listener, tilPort-fromPort+1)

	for i := 0; fromPort <= tilPort; i, fromPort = i+1, fromPort+1 {
		addr := ":" + strconv.FormatInt(fromPort, 10)
		log.Printf("trying to bind %s", addr)
		l, err := net.Listen("tcp", addr)
		if err != nil {
			log.Printf("can't create the {%d}th of port listener: %+v", i, err)
			goto toDie
		}

		toListen[i] = l
	}

	for {
		for _, l := range toListen {
			c, err := l.Accept()
			if err != nil {
				fmt.Printf("can't accept the conn:%+v", err)
				return
			}
			go handleConnection(c)
		}
	}

toDie:
	for _, l := range toListen {
		if l != nil {
			if err := l.Close(); err != nil {
				fmt.Printf("can't close the listener:%+v", err)
				return
			}
		}
	}
}
