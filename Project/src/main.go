package main

import (
	"driver"
	"elevatorHW"
	"fmt"
	"time"
)

func main(){
	const motorSpeed = 2800
	status := driver.Init()
	if(status){
		fmt.Println("Ja")
	} else {
		fmt.Println("Nei")
	}
	elevatorHW.SetDoorOpen()
	elevatorHW.SetMotor(1, motorSpeed)
	time.Sleep(1*time.Second)
	elevatorHW.SetMotor(1,0)
}