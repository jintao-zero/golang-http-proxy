package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"strings"
)

func handleConn(conn net.Conn) {
	defer conn.Close()
	log.Println("client conn form ", conn.RemoteAddr())
	var buffer [1024]byte
	n, err := conn.Read(buffer[:])
	if err != nil {
		return
	}
	log.Println(string(buffer[:]))
	_, line, err := bufio.ScanLines(buffer[:], true)
	if err != nil {
		log.Println(err)
		return
	}
	var method string
	var absUri string
	_, err = fmt.Sscanf(string(line), "%s%s", &method, &absUri)
	if err != nil {
		log.Println(err)
		return
	}
	absUrl, err := url.Parse(absUri)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("%+v", *absUrl)
	var hostPort string
	if absUrl.Scheme == "http" {
		if strings.Index(absUrl.Host, ":") == -1 {
			hostPort = absUrl.Host + ":80"
		} else {
			hostPort = absUrl.Host
		}
	} else {
		hostPort = absUrl.Scheme + ":" + absUrl.Opaque
	}
	log.Println("hostPort", hostPort)
	serverConn, err := net.Dial("tcp", hostPort)
	if err != nil {
		log.Println(err)
		return
	}
	defer serverConn.Close()

	if method == "CONNECT" {
		n, err = conn.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("connection estabilished", n)
	} else {
		serverConn.Write(buffer[:n])
	}
	//
	go io.Copy(conn, serverConn)
	io.Copy(serverConn, conn)
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
