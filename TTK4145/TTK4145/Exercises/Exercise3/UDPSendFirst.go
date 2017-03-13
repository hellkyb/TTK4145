package main

import (
	."fmt"
	."net"
	"time"
)

func main(){
	conn, err := Dial("udp4", "129.241.187.43:20001")
	Println(err)

	for{
	Fprint(conn, "Hello Server, this is Badger")
	time.Sleep(1*time.Second)	
	}
	conn.Close()
}
