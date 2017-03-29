package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/url"
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
	// fmt.Println(string(buffer[:]))
	/*
	   http
	   GET http://www.localhost:8080 HTTP/1.1
	   Host www.localhost
	   ..
	   ..
	*/
	_, line, err := bufio.ScanLines(buffer[:], true)
	if err != nil {
		log.Println(err)
		return
	}
	var method string
	var absPath string
	_, err = fmt.Sscanf(string(line), "%s%s", &method, &absPath)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(method, absPath)
	destUrl, err := url.Parse(absPath)
	fmt.Printf("%+v", *destUrl)
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
