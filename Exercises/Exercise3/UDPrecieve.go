package main
// This is a comment

import (
	"fmt"
	"net"
)



func main(){
	addr, err := net.ResolveUDPAddr("udp4", ":30000	")
	fmt.Println(err)
	conn, err := net.ListenUDP("udp4", addr)
	fmt.Println(err)

	for {
		b := make([]byte, 1024)
		n, err := conn.Read(b)
		fmt.Println(n, err, string(b))
	}
}
