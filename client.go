// tls client
package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
)

var (
	config tls.Config
)

func init() {
	certPool := x509.NewCertPool()
	serverCert, err := ioutil.ReadFile("./server.crt")
	if err != nil {
		log.Fatal("Could not load server certificate!")
	}
	certPool.AppendCertsFromPEM(serverCert)

	// write per-session secrets
	w, err := os.OpenFile("sslkeylog", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		log.Fatal(err)
	}

	config = tls.Config{
		KeyLogWriter:       w, // writer
		RootCAs:            certPool,
		ServerName:         "127.0.0.1",
		InsecureSkipVerify: true,  // ignore self signed cert
		MaxVersion:         0x303, // tls version number 1.2
		MinVersion:         0x303,
		CipherSuites: []uint16{
			tls.TLS_RSA_WITH_RC4_128_SHA,
			tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
			tls.TLS_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_128_CBC_SHA256,
			tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
		},
	}
}

func main() {
	var buf = make([]byte, 1024)

	unsec, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		log.Fatalf("client: dial: %s", err)
	}
	defer unsec.Close()

	conn := tls.Client(unsec, &config)
	err = conn.Handshake()
	if err != nil {
		log.Fatalf("tls: handshake: %s", err)
	}
	for {
		var str string
		fmt.Scan(&str) // read input from stdin
		_, err := conn.Write([]byte(str))
		if err != nil {
			conn.Close()
			return
		}
		n, err := conn.Read(buf)
		if err != nil {
			conn.Close()
			return
		}
		response := string(buf[0:n])
		fmt.Print(response)
	}
}
