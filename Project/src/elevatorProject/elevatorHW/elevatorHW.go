package elevatorHW

import(
	"elevatorProject/io"
)

const motorSpeed = 2800
const NFloors	 = 4
const NButtons	 = 3
const NLights 	 = 4

type MotorDirection int
const (
	DirectionStop = iota
	DirectionUp
	DirectionDown
)

type ButtunType int
const (
	ButtonCallDown = iota
	ButtonCallUp
	ButtonCommand
)



func SetDoorLight(open bool) {
	if open {
		io.SetBit(lightDoorOpen)
	} else {
		io.ClearBit(lightDoorOpen)
	}
}

func SetMotor(dir int) {
	if (dir == DirectionStop) {
		io.WriteAnalog(motor, 0)
	} else if (dir == DirectionDown) {
		io.SetBit(motordir)
		io.WriteAnalog(motor,motorSpeed)
	} else if(dir == DirectionUp) {
		io.ClearBit(motordir)
		io.WriteAnalog(motor, motorSpeed)
	}
}