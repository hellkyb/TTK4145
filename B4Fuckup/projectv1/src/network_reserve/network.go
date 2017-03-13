package network

import (
	"net"
)

func SendAlive() {
	addr, err := ResolveUDPAddr("udp4", "129.241.187.43:20001")
	Println(err)
	conn, err := DialUDP("udp4", nil, addr)
	if err != nil {
		Println(err)
	}

	for {

		n, err := conn.Write([]byte("I is alive"))
		Println(n, err)
		time.Sleep(1 * time.Second)
	}
}
