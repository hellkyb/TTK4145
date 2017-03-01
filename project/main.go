package main

import (
	"./src/elevatorHW"

	"./src/fsm"
	//"./network/peers"
	//s"fmt"

	//"time"
)

func main() {

	elevatorHW.Init()

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
