package fsm

import (
	"../elevatorHW"
	"../io"
	"fmt"
	"time"
)

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
