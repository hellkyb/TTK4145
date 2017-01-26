package main

import (
	"net"
	"fmt"
	"time"
)

const (
	serverIP = "129.241.187.43"
	udpPort = string(20000 + 24)
)

func udpSend(done chan bool, port, saddr *udpAddr){
	conn, err := net.DialUDP("udp", nil, saddr)
	if err != nil{
		fmt.Println("Error connecting to " + saddr)
	}
	for{
		time.Sleep(1000*time.Millisecond)
		conn.Write([]byte("Hello there Sanntidssalen"))
		fmt.Println("Message sent on udp")
	}
	done <- true
}

func udpReceive(done chan bool, port, saddr *udpAddr){
	buff := make([]byte, 1024)

	l, err := net.ListenUDP("udp4", saddr)
	if err != nil{
		fmt.Println("Error listening to "+saddr)
	}

	_,_, err = l.ReadFromUDP(buff)

	if err != nil{
		fmt.Println(err)
	}
	fmt.Println(string(buff[:]))
}

func main(){
	done := make (chan bool)

	saddr,err := net.ResolveUDPAddr("udp4", net.JoinHostPort(host, udpPort))

	if err != nil{
		fmt.Println("Failed to resolve address for: " + udpPort)

	}

	go udpSend(done, udpPort, saddr)
	go udpReceive(done, udpPort, saddr)

	<-done
	<-done
}