package main

import (
	"time"

	"./src/elevatorHW"
	"./src/fsm"
	"./src/olasnetwork"
	//"./src/io"
	"fmt"
)

// This function returns how suitet the elevator is to handle a global call
func costFunction(dir int, lastFloor int, order fsm.Order) int {
	var distanceToTarget int
	var dirCost int
	if lastFloor > order.Floor {
		distanceToTarget = lastFloor - order.Floor
	} else {
		distanceToTarget = order.Floor - lastFloor
	}

	distCost := 2 * distanceToTarget

	if dir == 0 && order.Floor == lastFloor { // Elevator is Idle at floor being called
		return 0
	}
	if dir == 0{
		return distanceToTarget
	}
	if order.Button == 1 { //UpType Order
		if dir == -1 { // Moving in opposite direction
			if lastFloor < order.Floor {
				dirCost = 8
			} else {
				dirCost = 10
			}
		} else if dir == 0 { //Elevator is idle
			dirCost = 0
		} else if dir == 1 { // Elevator is moving up
			if lastFloor > order.Floor {
				dirCost = 8
			} else {
				dirCost = 10
			}
		}

	} else if order.Button == 0 { //DownType order
		if dir == 1 { // Oposite directioin
			if lastFloor > order.Floor {
				dirCost = 8
			} else {
				dirCost = 10
			}
		} else if dir == 0 {
			dirCost = 0
		} else if dir == -1 {
			if lastFloor < order.Floor {
				dirCost = 8
			} else {
				dirCost = 10
			}
		}
	}
	return dirCost + distCost
}

func decitionmaker(onlineElevatorStates map[string]olasnetwork.HelloMsg) (string, int) {
	numberOfElevatorsInNetwork := olasnetwork.OperatingElevators
	if numberOfElevatorsInNetwork == 0 || numberOfElevatorsInNetwork == 1 {
		return olasnetwork.GetLocalID(), 0
	}	
	var elevatorWithLowestCost string
	lowestCost := 1337
	if len(onlineElevatorStates) < 2 {
		return olasnetwork.GetLocalID(), 0
	}
	for key, value := range onlineElevatorStates {
		thisCost := costFunction(value.CurrentState, value.LastFloor, value.Order.Order)
		if thisCost < lowestCost{
			lowestCost = thisCost
			elevatorWithLowestCost = key
		}		
	}
	return elevatorWithLowestCost, lowestCost	
	
}

func main() {
	//start init
	fmt.Println("Starting system")
	fmt.Print("\n\n")
	elevatorHW.Init()
	//finished init
	fsm.CreateQueueSlice()
	time.Sleep(1 * time.Millisecond)

	//var operatingElevatorStates []olasnetwork.HelloMsg
	operatingElevatorStates := make(map[string]olasnetwork.HelloMsg) 
	//var operatingTimeStampList []int64  Use this to delete objects??

	buttonCh := make(chan fsm.Order)
	messageCh := make(chan olasnetwork.HelloMsg)
	networkOrderCh := make(chan olasnetwork.HelloMsg)
	networkSendOrderCh := make(chan olasnetwork.OrderMsg)

	go fsm.RunElevator()
	go fsm.GetButtonsPressed(buttonCh)
	go olasnetwork.NetworkMain(messageCh, networkOrderCh, networkSendOrderCh)

	for {
		select {

		case newMsg := <-messageCh:
			olasnetwork.UpdateElevatorStates(newMsg, operatingElevatorStates)
			olasnetwork.DeleteDeadElevator(operatingElevatorStates)
			if newMsg.Order.ElevatorToTakeThisOrder == olasnetwork.GetLocalID() {				
				fsm.PutOrderInLocalQueue(newMsg.Order.Order)				
			}

		case newOrder := <-buttonCh:			
			if newOrder.Button == elevatorHW.ButtonCommand {
				fsm.PutInsideOrderInLocalQueue()
			}else if len(operatingElevatorStates) == 0 || len(operatingElevatorStates) == 1{
				fsm.PutOrderInLocalQueue(newOrder)
			}else{
				elevatorToHandleThisOrder, _ := decitionmaker(operatingElevatorStates)
				fmt.Print("I want ")
				fmt.Print(elevatorToHandleThisOrder)
				fmt.Print(" to handle this order\n\n")	
								
				networkSendOrderCh <- olasnetwork.OrderMsg{newOrder, elevatorToHandleThisOrder}
			}
			//fsm.PrintQueues()
		}
	}
}
