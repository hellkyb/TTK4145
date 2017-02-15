package main

import (
	"elevatorProject/io"
	"elevatorProject/elevatorHW"
	"fmt"
	"time"
)

func main(){
	status := io.Init()
	if(status){
		fmt.Println("Ja")
	} else {
		fmt.Println("Nei")
	}
	elevatorHW.SetMotor(elevatorHW.DirectionUp)
	time.Sleep(1*time.Second)
	elevatorHW.SetMotor(elevatorHW.DirectionDown)
	time.Sleep(1*time.Second)
	elevatorHW.SetMotor(elevatorHW.DirectionStop)
}