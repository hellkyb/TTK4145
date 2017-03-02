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
	
	go HentKnapper
	go Kjør
	go nett
	go RegnUtKostFunksjon

	for{
		t := 0
		if t%5000 == 0 {
			PrintLocalQueue()
		}
		if t > 100000 {
			t = 0
		}
	}
	}

	fsm.RunElevator()
}
