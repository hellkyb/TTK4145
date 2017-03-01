package fsm

import (
	"../elevatorHW"
	//"../io"
	//"../network/peers"
	//"fmt"
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
	elevatorHW.SetMotor(elevatorHW.DirectionStop)
	elevatorHW.SetFloorIndicator(floor)
	elevatorHW.SetDoorLight(true)
	for !timer() {
		continue
	}
	elevatorHW.SetDoorLight(false)
}

func PutOrderInLocalQueue() {
	upOrder := elevatorHW.GetUpButton()
	downOrder := elevatorHW.GetDownButton()
	insideOrder := elevatorHW.GetInsideElevatorButton()
	if upOrder != 0 {
		AppendUpOrder(upOrder)
	}
	if downOrder != 0 {
		AppendDownOrder(downOrder)
	}
	if insideOrder != 0 {
		AppendInsideOrder(insideOrder)
	}
}

func SetElevatorDirection() {
	/*if elevatorHW.GetElevatorDirection() == 1 && elevatorHW.GetFloorSensorSignal() == 1 {
		elevatorHW.Init() // Run Init if elevator is stationary between floors
	}*/
	if elevatorHW.GetElevatorDirection() == 1 || elevatorHW.GetElevatorDirection() == 0 {
		if elevatorHW.GetFloorSensorSignal() != 0 {
			if len(localQueue[0]) > 0 {
				if localQueue[0][0] < elevatorHW.GetFloorSensorSignal() {
					elevatorHW.SetMotor(elevatorHW.DirectionDown)
				} else if localQueue[0][0] > elevatorHW.GetFloorSensorSignal() {
					elevatorHW.SetMotor(elevatorHW.DirectionUp)
				}
			} else if len(localQueue[1]) > 0 {
				if localQueue[1][0] < elevatorHW.GetFloorSensorSignal() {
					elevatorHW.SetMotor(elevatorHW.DirectionDown)
				} else if localQueue[1][0] > elevatorHW.GetFloorSensorSignal() {
					elevatorHW.SetMotor(elevatorHW.DirectionUp)
				}
			} else if len(localQueue[2]) > 0 {
				if localQueue[2][0] < elevatorHW.GetFloorSensorSignal() {
					elevatorHW.SetMotor(elevatorHW.DirectionDown)
				} else if localQueue[2][0] > elevatorHW.GetFloorSensorSignal() {
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
				DeleteIndexLocalQueue(i, j)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}

}
