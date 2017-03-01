package fsm

import (
	"../elevatorHW"
	//"../io"
	//"../network/peers"
	"fmt"
	//"math"
	//"time"
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

func timer() {
	start := time.Now()
	elapsed := 0
	for elapsed < 3*time.Second {
		elapsed = time.Since(start)
	}
	return true
}

func ArrivedAtFloorSetDoorOpen(floor int) {
	elevatorHW.SetMotor(DirectionStop)
	elevatorHW.SetFloorIndicator(floor)
	elevatorHW.SetDoorLight(true)
	for !timer() {
		continue
	}
	elevatorHW.SetDoorLight(false)
}
