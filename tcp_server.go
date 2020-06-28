// tls server
package main

import (
	"log"
	"net"
)

func main() {
	listener, _ := net.Listen("tcp", "127.0.0.1:8000")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("server: accept: %s", err)
			break
		}
		go handleTLSConnection(conn)
	}
}

func handleTLSConnection(conn net.Conn) {
	var buf = make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			conn.Close()
			return
		}
		response := string(buf[0:n])
		log.Printf("Get [%s], echo back", response)
		// echo
		conn.Write([]byte(response + "\n"))
	}
}
