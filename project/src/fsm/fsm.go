package fsm

import (
	"../elevatorHW"
	"../io"
	"fmt"
	"time"
	"math"
)

type elevatorState struct{
	previousFloor int
	direction elevatorHW.MotorDirection
}
type order struct{
	direction elevatorHW.MotorDirection
	fromFloor int 
}

func costFunc(elevatorStates []elevatorState, newOrder order){
	distanceCost := 0
	directionCost := 0
	totalCost := 0

	costs := [3]float64{math.Inf(1), math.Inf(1), math.Inf(1)} 

	for i = 0; i < len(elevatorStates); i++{
		distanceCost = math.Abs(elevatorStates[i].previousFloor - newOrder.fromFloor)
		directionCost = math.Abs(elevatorStates[i].direction)

	}
	
}


func TestFunction() {

	if io.ReadAnalog(elevatorHW.ButtonUp3) == 1 {

		elevatorHW.SetMotor(elevatorHW.DirectionUp)
		time.Sleep(2680 * time.Millisecond)
		elevatorHW.SetMotor(elevatorHW.DirectionStop)

	} else if io.ReadAnalog(elevatorHW.ButtonDown3) == 1 {
		elevatorHW.SetMotor(elevatorHW.DirectionDown)
		time.Sleep(2680 * time.Millisecond)
		elevatorHW.SetMotor(elevatorHW.DirectionStop)
	}

}

func StupidFunc() {
	fmt.Println(io.ReadAnalog(elevatorHW.SensorFloor1))
}


