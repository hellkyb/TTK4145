package main

import (
	"driver"
	"elevatorHW"
	"fmt"
)

func main(){
	status := driver.Init()
	if(status){
		fmt.Println("Ja")
	} else {
		fmt.Println("Nei")
	}
	elevatorHW.SetDoorOpen()
}