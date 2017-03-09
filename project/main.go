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

func decitionmaker(onlineElevatorStates []olasnetwork.HelloMsg) (string, int) {
	numberOfElevatorsInNetwork := olasnetwork.OperatingElevators
	fmt.Print("D-NumElevs")
	fmt.Println(numberOfElevatorsInNetwork)
	if numberOfElevatorsInNetwork == 0 || numberOfElevatorsInNetwork == 1 {
		return olasnetwork.GetLocalID(), 0
	}
	var costs []int

	lowestCost := 152
	var minPos int
	for i := 0; i < numberOfElevatorsInNetwork; i++ {
		thisCost := costFunction(onlineElevatorStates[i].CurrentState, onlineElevatorStates[i].LastFloor, onlineElevatorStates[i].Order.Order)
		costs = append(costs, thisCost)
		if lowestCost > costs[i] {
			lowestCost = costs[i]
			minPos = i
		}
	}
	return onlineElevatorStates[minPos].ElevatorID, lowestCost
}

func main() {

	//start init
	fmt.Println("Starting system")
	fmt.Print("\n\n")
	fsm.StartUpMessage()
	elevatorHW.Init()
	//finished init
	fsm.CreateQueueSlice()
	var operatingElevatorStates []olasnetwork.HelloMsg

	//stateRx := make(chan fsm.ElevatorStatus)

	buttonCh := make(chan fsm.Order)
	messageCh := make(chan olasnetwork.HelloMsg)

	// go func() {
	// 	for {
	// 		/*upOrder := elevatorHW.GetUpButton()
	// 		downOrder := elevatorHW.GetDownButton()
	// 		insideOrder := elevatorHW.GetInsideElevatorButton()*/
	// 		/*state := elevatorHW.GetElevatorState()
	//
	// 		fmt.Println(state)*/
	//
	// 	}
	// }()
	//go RunElevator()
	go fsm.RunElevator()
	go fsm.GetButtonsPressed(buttonCh)
	go olasnetwork.NetworkMain(messageCh)
	time.Sleep(1000 * time.Millisecond)

	mylist := make([]olasnetwork.HelloMsg, 2)
	mylist[0] = olasnetwork.HelloMsg{0, "Winner", 0, 4, olasnetwork.OrderMsg{fsm.Order{0, 2}, " "}}
	mylist[1] = olasnetwork.HelloMsg{0, "Looser", 0, 3, olasnetwork.OrderMsg{fsm.Order{0, 2}, " "}}
	theChosenOne, cost := decitionmaker(mylist)
	fmt.Println(theChosenOne+" ", cost)

	for {
		select {
		case newMsg := <-messageCh:
			operatingElevatorStates = olasnetwork.UpdateElevatorStates(newMsg, olasnetwork.OperatingElevators, operatingElevatorStates)
			fmt.Println(operatingElevatorStates)
		case newOrder := <-buttonCh:
			fmt.Print("You made an order: ")
			fmt.Println(newOrder)
			if newOrder.Button == elevatorHW.ButtonCommand {
				fsm.PutInsideOrderInLocalQueue(newOrder)
			} else {

				decitionmaker(olasnetwork.OperatingElevatorStates)
			}
			fsm.PrintQueues()
		}
		//time.Sleep(200 * time.Millisecond)
	}
}
