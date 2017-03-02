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
	fmt.Println("Starting system")
	elevatorHW.Init()
	fsm.RunElevator()
}
