package main

import (
	."fmt"
	."net"
	"time"
)

func main() {
	addr, err := ResolveUDPAddr("udp4", "129.241.187.43:20001")
	Println(err)
	conn, err := DialUDP("udp4", nil, addr)
	Println(err)

	for{
		
		n,err := conn.Write([]byte("hello etc"))
		Println(n,err)
		time.Sleep(1*time.Second)
	}
}