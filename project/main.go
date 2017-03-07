package main

import (
	"./src/olasnetwork"
	"./src/elevatorHW"
	"./src/fsm"
	//"./src/io"
	//"./network/peers"
	"fmt"
	"time"
)




// This function returns how suitet the elevator is to handle a global call
func costFunction(state int, lastFloor int, globalOrder fsm.Order) int{ 
	distanceToTarget := 0
	if lastFloor > globalOrder.Floor {
		distanceToTarget = lastFloor - globalOrder.Floor
	}else {
		distanceToTarget = globalOrder.Floor - lastFloor
	}

	if state == 0 && globalOrder.Floor == lastFloor { // Elevator is Idle at floor being called
		return 10 
	}
	if globalOrder.Button == 1 { //UpType Order
		if state == -1 { // Moving in opposite direction
			switch distanceToTarget{
			case 2: // Distance to target is max and moving in oposite direction is worstcase
				return 0
			case 1:
				return 1
			case 0:
				return 2
			}
		}else if state == 0 { //Elevator is idle
			switch distanceToTarget{
			case 0:
				return 10 //BestCase
			case 1:
				return 9
			case 2:
				return 8
			case 3:
				return 7
			}

		}else if state == 1{  // Elevator is moving up
			switch distanceToTarget{
			case 0:
				return 3 //Elevator just left
			case 1:
				return 9
			case 2:
				return 8
			case 3:
				return 7
			}
		}


	}else if globalOrder.Button == 0 { //DownType order
		if state == 1{ // Oposite directioin
			switch distanceToTarget{
			case 0:
				return 2
			case 1:
				return 8
			case 2:
				return 7
			case 3:
				return 6
			}

		}else if state == 0{
			switch distanceToTarget{
			case 0:
				return 10
			case 1:
				return 9
			case 2:
				return 8
			case 3:
				return 7
			}

		}else if state == -1{
			switch distanceToTarget{
			case 0:
				return 2
			case 1:
				return 8
			case 2:
				return 7
			case 3:
				return 6
			}

		}

	}
	return 0
}

//This function looks on the oldest global queue order
//uses the costFunction on all elevators and determine if it is best suited for the job
// If this is the case, it should append the global order into its localqueue
func decitionmaker(){
	numberOfElevatorsInNetwork := olasnetwork.OperatingElevators
	fmt.Println(numberOfElevatorsInNetwork)

	if numberOfElevatorsInNetwork == 1 || numberOfElevatorsInNetwork == 0 {
		// Alone mode, handle the globalQueue as LocalQueue
	}
	/*localState = elevatorHW.GetElevatorState()
	localCost := costFunction(localState, fsm.LastFloor, )*/
}

func main() {

	
	//start init
	fmt.Println("Starting system")
	elevatorHW.Init()
	//finished init
	fsm.CreateQueueSlice()
	//fsm.CreateGlobalQueueSlice()
	//testOrder := fsm.Order{2,0}
	//fsm.PutOrderInLocalQueue(testOrder)
	go func (){
		for{
			/*upOrder := elevatorHW.GetUpButton()
			downOrder := elevatorHW.GetDownButton()
			insideOrder := elevatorHW.GetInsideElevatorButton()*/
			state := elevatorHW.GetElevatorState()
			
			fmt.Println(state)
			time.Sleep(1000*time.Millisecond)
		}
	}()
	go fsm.RunElevator()
	go olasnetwork.NetworkMain()
	for {
		fsm.PrintQueues()
		decitionmaker()
		time.Sleep(1*time.Second)
	}
}