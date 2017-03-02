package fsm

import (
	"../elevatorHW"
	//"../io"
	//"../network/peers"
	"fmt"
	//"math"
	"time"
)

type elevatorState struct {
	previousFloor int
	direction     elevatorHW.MotorDirection
}
type order struct {
	direction elevatorHW.MotorDirection
	fromFloor int
}

/*func costFunc(elevatorStates []elevatorState, newOrder order) {
	distanceCost := 0
	directionCost := 0
	totalCost := 0

	costs := [3]float64{math.Inf(1), math.Inf(1), math.Inf(1)}

	for i := 0; i < len(elevatorStates); i++ {
		distanceCost = math.Abs(elevatorStates[i].previousFloor - newOrder.fromFloor)
		if (elevatorStates[i].previousFloor == newOrder.fromFLoor) && (elevatorStates[i].direction == DirectionStop){
			distanceCost = 0
		}
			//NEED PLENTY MORE LOGIC HERE
	}


}*/

func timer() bool {
	start := time.Now()
	elapsed := time.Since(start)
	for elapsed < 3 {
		elapsed = time.Since(start)
	}
	return true
}

func ArrivedAtFloorSetDoorOpen(floor int) {
	elevatorHW.SetFloorIndicator(floor)
	elevatorHW.SetDoorLight(true)
	time.Sleep(3000 * time.Millisecond)
	elevatorHW.SetDoorLight(false)
}

func PutOrderInLocalQueue() {
	upOrder := elevatorHW.GetUpButton()
	downOrder := elevatorHW.GetDownButton()
	insideOrder := elevatorHW.GetInsideElevatorButton()
	if upOrder != 0 {
		AppendUpOrder(upOrder)
		elevatorHW.SetUpLight(upOrder, true)
	}
	if downOrder != 0 {
		AppendDownOrder(downOrder)
		elevatorHW.SetDownLight(downOrder, true)
	}
	if insideOrder != 0 {
		AppendInsideOrder(insideOrder)
		elevatorHW.SetInsideLight(insideOrder, true)
	}
}

func SetElevatorDirection() {
	currentDirection := elevatorHW.GetElevatorDirection()
	currentFloor := elevatorHW.GetFloorSensorSignal()
	if currentDirection == 1 || currentDirection == 0 {
		if currentFloor != 0 {
			if len(localQueue[0]) > 0 {
				if localQueue[0][0] < currentFloor {
					elevatorHW.SetMotor(elevatorHW.DirectionDown)
				} else if localQueue[0][0] > currentFloor {
					elevatorHW.SetMotor(elevatorHW.DirectionUp)
				}
			} else if len(localQueue[1]) > 0 {
				if localQueue[1][0] < currentFloor {
					elevatorHW.SetMotor(elevatorHW.DirectionDown)
				} else if localQueue[1][0] > currentFloor {
					elevatorHW.SetMotor(elevatorHW.DirectionUp)
				}
			} else if len(localQueue[2]) > 0 {
				if localQueue[2][0] < currentFloor {
					elevatorHW.SetMotor(elevatorHW.DirectionDown)
				} else if localQueue[2][0] > currentFloor {
					elevatorHW.SetMotor(elevatorHW.DirectionUp)
				}
			}
		}
	}

}

func StopAtThisFloor() {
	currentFloor := elevatorHW.GetFloorSensorSignal()

	for i := 0; i < 3; i++ {
		if len(localQueue[i]) < 0 {
			continue
		}
		for j := range localQueue[i] {
			if currentFloor == localQueue[i][j] {
				elevatorHW.SetMotor(elevatorHW.DirectionStop)
				ArrivedAtFloorSetDoorOpen(currentFloor)
				elevatorHW.SetDownLight(currentFloor, false)
				elevatorHW.SetUpLight(currentFloor, false)
				elevatorHW.SetInsideLight(currentFloor, false)
				elevatorHW.SetFloorIndicator(currentFloor)
				switch i {
				case 0:
					DeleteOldestOrderInside()
				case 1:
					DeleteOldestOrderUp()
				case 2:
					DeleteOldestOrderDown()
				}
				return
			}
		}
	}
}

func PrintElevatorStatus() {
	currentFloor := elevatorHW.GetFloorSensorSignal()
	currentDirection := elevatorHW.GetElevatorDirection()
	for i := 0; i < 5; i++ {
		fmt.Println(" ")
	}
	fmt.Println("------------------------------------")
	fmt.Println("Elevator status")
	fmt.Println("floor: ", currentFloor)
	fmt.Println("Direction: ", currentDirection)
	fmt.Println(" ")
	fmt.Println("Queue: ")
	PrintLocalQueue()

}

func RunElevator() {
	printCounter := 0
	CreateQueueSlice()
	for {
		PutOrderInLocalQueue()
		SetElevatorDirection()
		StopAtThisFloor()
		printCounter++
		if printCounter > 100000 {
			printCounter = 0
		}
		if printCounter%10000 == 0 {
			PrintElevatorStatus()
		}
	}
}
