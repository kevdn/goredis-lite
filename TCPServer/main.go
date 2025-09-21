package main

import (
	"log"
	"net"
	"time"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	log.Println(conn.RemoteAddr())
	var buf []byte = make([]byte, 1000)
	_, err := conn.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Second * 2)
	log.Println("write response")
	conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\nHello!\r\n"))
}

func main() {
	listener, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConnection(conn)
	}
}
