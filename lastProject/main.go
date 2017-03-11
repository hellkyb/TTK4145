package main

import (
	"./src/elevatorHW"
	"./src/fsm"
	"./src/olasnetwork"
	"fmt"
	"time"
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
	if dir == 0 {
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

func decitionmaker(onlineElevatorStates map[string]olasnetwork.HelloMsg, newOrder fsm.Order) (string, int) {
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
		thisCost := costFunction(value.CurrentState, value.LastFloor, newOrder)
		fmt.Print("\n\nStates of ")
		fmt.Print(key)
		fmt.Println(": ")
		fmt.Print("Current state value: ")
		fmt.Print(value.CurrentState)
		fmt.Print("\nIts Last Floor: ")
		fmt.Println(value.LastFloor)
		fmt.Print("The Order is : ")
		fmt.Println(newOrder)
		fmt.Println(" ")
		fmt.Print(key)
		fmt.Print(" has a cost of ")
		fmt.Println(thisCost)
		fmt.Print("\n\n")
		if thisCost < lowestCost {
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
	//fsm.StartUpMessage()
	//finished init
	fsm.CreateQueueSlice()
	time.Sleep(1 * time.Millisecond)

	operatingElevatorStates := make(map[string]olasnetwork.HelloMsg)
	hallButtonsMap := make(map[fsm.Order]int64)

	buttonCh := make(chan fsm.Order)
	messageCh := make(chan olasnetwork.HelloMsg)
	networkOrderCh := make(chan olasnetwork.HelloMsg)
	networkSendOrderCh := make(chan olasnetwork.OrderMsg)
	timeOutCh := make(chan bool)
	orderCompletedCh := make(chan fsm.Order)
	sendDeletedOrderCh := make(chan fsm.Order)

	go fsm.RunElevator(timeOutCh, orderCompletedCh, hallButtonsMap)
	go fsm.GetButtonsPressed(buttonCh)
	go olasnetwork.NetworkMain(messageCh, networkOrderCh, networkSendOrderCh, orderCompletedCh, sendDeletedOrderCh)

	for {
		fsm.DoorLightTimeOut()
		select {

		case orderIsHandled := <-orderCompletedCh:
			fmt.Println("It has deleted an order!")
			delete(hallButtonsMap, orderIsHandled)
			sendDeletedOrderCh <- orderIsHandled

		case doorClose := <-timeOutCh:
			elevatorHW.SetDoorLight(!doorClose)

		case newMsg := <-messageCh:
			if newMsg.Order.Order.Floor != -1 {
				hallButtonsMap[newMsg.Order.Order] = time.Now().Unix()
			}
			fsm.DoorLightTimeOut()
			fmt.Print("This is OrderMap:  ")
			fmt.Println(hallButtonsMap)
			fmt.Print("Length of elevatorState Map:  ")
			fmt.Println(len(operatingElevatorStates))
			olasnetwork.UpdateElevatorStates(newMsg, operatingElevatorStates)
			olasnetwork.DeleteDeadElevator(operatingElevatorStates)
			if newMsg.Order.ElevatorToTakeThisOrder == olasnetwork.GetLocalID() {
				fsm.PutOrderInLocalQueue(newMsg.Order.Order)
				fmt.Println("I recieved an order! Local Queue:  ")
				fsm.PrintQueues()
			}
			if newMsg.OrderExecuted.Floor != -1 {
				fmt.Println(hallButtonsMap)
				delete(hallButtonsMap, newMsg.OrderExecuted)
				fmt.Println("Some elevator has actually done its job!!!! Hurray")
				fmt.Println(hallButtonsMap)
			}

		case newOrder := <-buttonCh:

			if newOrder.Button == elevatorHW.ButtonCommand {
				fsm.PutInsideOrderInLocalQueue()
			} else if len(operatingElevatorStates) == 0 || len(operatingElevatorStates) == 1 {
				fsm.PutOrderInLocalQueue(newOrder)
			} else {
				elevatorToHandleThisOrder, cost := decitionmaker(operatingElevatorStates, newOrder)

				fmt.Print("I want ")
				fmt.Print(elevatorToHandleThisOrder)
				fmt.Print(" to handle this order with cost of ")
				fmt.Print(cost)
				fmt.Println(" ")
				hallButtonsMap[newOrder] = time.Now().Unix()
				fmt.Println(hallButtonsMap)
				networkSendOrderCh <- olasnetwork.OrderMsg{newOrder, elevatorToHandleThisOrder}
			}
		}
	}
}
