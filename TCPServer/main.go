package main

import (
	"log"
	"net"
	"io"
	"fmt"
)

func handleConnection(conn net.Conn) {
	log.Println("Handling connection from", conn.RemoteAddr())
	for {
		cmd, err := readCommand(conn)
		if err != nil {
			conn.Close()
			log.Println("Client disconnected", conn.RemoteAddr())
			if err == io.EOF {
				break
			}
		}
		if err = respond(cmd, conn); err != nil {
			log.Println("Error responding to client", conn.RemoteAddr())
		}
	}
	
}

func readCommand(conn net.Conn) (string, error) {
	var buf []byte = make([]byte, 512)
	n, err := conn.Read(buf)
	if err != nil {
		return "", err
	}
	return string(buf[:n]), nil
}

func respond(cmd string, conn net.Conn) error {
	cmd = fmt.Sprintf("+%s\r\n", cmd)
	if _, err := conn.Write([]byte(cmd)); err != nil {
		return err
	}
	return nil
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
