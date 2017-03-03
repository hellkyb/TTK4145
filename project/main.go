package main

import (
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
	fsm.CreateQueueSlice()
	//finished init

	fsm.RunElevator()
}
