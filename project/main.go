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

	/*	for {
		if i%40000000 == 0 {
			go fsm.PrintLocalQueue()
			go fsm.PutOrderInLocalQueue()
			go fsm.SetElevatorDirection()
			go fsm.StopAtThisFloor()
			fmt.Println(elevatorHW.GetElevatorDirection())
		}

		i++
	}*/

	/*go elevatorHW.GetFloorSensorSignal()
	fmt.Println(currentfloor)
	*/
	//peers []string

	//elevatorHW.SetMotor(elevatorHW.DirectionUp)
	//time.Sleep(1 * time.Second)
	//elevatorHW.SetMotor(elevatorHW.DirectionDown)
	//time.Sleep(1 * time.Second)
	//elevatorHW.SetMotor(elevatorHW.DirectionStop)
	/*
		status := io.Init()
		if status {
			fmt.Println("Ok init")
		} else {
			fmt.Println("Error at init")
		}
		if io.ReadAnalog(elevatorHW.SensorFloor1) != 1 {
			elevatorHW.SetMotor(elevatorHW.DirectionDown)
			for {
				if io.ReadAnalog(elevatorHW.SensorFloor1) == 1 {
					elevatorHW.SetMotor(elevatorHW.DirectionStop)
					break
				}
			}
		}
	*/
}
