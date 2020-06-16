// tls server
package main

import (
	"crypto/tls"
	"log"
	"net"
)

var (
	config tls.Config
)

func init() {
	cert, _ := tls.LoadX509KeyPair("./server.crt", "./server.key")

	config = tls.Config{
		Certificates: []tls.Certificate{cert},
		MaxVersion:   0x303, // tls version number 1.2
		MinVersion:   0x303,
		//CipherSuites: []uint16{tls.TLS_RSA_WITH_AES_256_GCM_SHA384, tls.TLS_RSA_WITH_AES_128_CBC_SHA256, tls.TLS_RSA_WITH_AES_128_GCM_SHA256},
	}
}

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

func handleTLSConnection(unsec net.Conn) {
	conn := tls.Server(unsec, &config)
	var buf = make([]byte, 1024)
	if err := conn.Handshake(); err != nil {
		log.Printf("%s\n", err)
	}
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
