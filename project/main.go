package main

import (
	"./src/network"
	"./src/elevatorHW"
	"./src/fsm"
	//"./src/io"
	//"./network/peers"
	"fmt"
	//"time"
)

func main() {
	//start init
	fmt.Println("Starting system")
	elevatorHW.Init()
	//finished init
	fsm.CreateQueueSlice()

	go fsm.RunElevator()
	go network.Main()
	for {
		//fsm.PrintLocalQueue()
		//time.Sleep(1*time.Second)
	}
}
