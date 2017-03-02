package fsm

import (
	"../elevatorHW"
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

func Timer() {
	time.Sleep(1000 * time.Millisecond)

}

func ArrivedAtFloorSetDoorOpen(floor int) {
	elevatorHW.SetFloorIndicator(floor)
	elevatorHW.SetDoorLight(true)
	Timer()
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

		for j := range localQueue[i] {
			if currentFloor == localQueue[i][j] {
				elevatorHW.SetMotor(elevatorHW.DirectionStop)
				ArrivedAtFloorSetDoorOpen(currentFloor)
				elevatorHW.SetDownLight(currentFloor, false)
				elevatorHW.SetUpLight(currentFloor, false)
				elevatorHW.SetInsideLight(currentFloor, false)
				elevatorHW.SetFloorIndicator(currentFloor)
				DeleteIndexLocalQueue(i, j)
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
	CreateQueueSlice()
	t := 0
	for {
		t++
		SetElevatorDirection()
		PutOrderInLocalQueue()
		StopAtThisFloor()
		if t%5000 == 0 {
			PrintLocalQueue()
		}
		if t > 100000 {
			t = 0
		}
	}
}
