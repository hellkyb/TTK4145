package main

import (
	"./src/elevatorHW"
	/*"./src/fsm"
	"./src/io"
	"fmt"
	"time"
	*/)

func main() {

	elevatorHW.Init()

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
