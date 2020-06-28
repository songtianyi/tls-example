// tcp client
package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	var buf = make([]byte, 1024)

	unsec, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		log.Fatalf("client: dial: %s", err)
	}
	defer unsec.Close()

	for {
		var str string
		fmt.Scan(&str) // read input from stdin
		_, err := unsec.Write([]byte(str))
		if err != nil {
			unsec.Close()
			return
		}
		n, err := unsec.Read(buf)
		if err != nil {
			unsec.Close()
			return
		}
		response := string(buf[0:n])
		fmt.Print(response)
	}
}
