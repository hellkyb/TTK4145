package main

import (
	. "fmt"
	. "net"
	"time"
)

func main() {
	go sender()

	addr, err := ResolveUDPAddr("udp4", ":20001")
	Println(err)
	conn, err := ListenUDP("udp4", addr)
	Println(err)

	for {
		b := make([]byte, 1024)
		n, err := conn.Read(b)
		Println(n, err, string(b))
	}
}

func sender() {
	addr, err := ResolveUDPAddr("udp4", "129.241.187.149:20001")
	Println(err)
	conn, err := DialUDP("udp4", nil, addr)
	Println(err)

	for {

		n, err := conn.Write([]byte("hello ernst"))
		Println(n, err)
		time.Sleep(1 * time.Second)
	}
}
