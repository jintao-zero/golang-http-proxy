package main

import (
	"log"
	"net"
)

func handleConn(conn net.Conn) {
	defer conn.Close()
	log.Println("client conn form ", conn.RemoteAddr())
	var buffer [1024]byte
	_, err := conn.Read(buffer[:])
	if err != nil {
		return
	}
	log.Println(string(buffer[:]))
}

func main() {
	// listen on tcp port
	l, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Println(err)
		return
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
			return
		}
		go handleConn(conn)
	}
}
