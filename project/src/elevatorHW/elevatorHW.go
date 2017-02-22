package elevatorHW

import (
	"../io"
	"fmt"
)

const motorSpeed = 2800
const NFloors = 4
const NButtons = 3
const NLights = 4

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

func Init() {
	status := io.Init()
	if status {
		fmt.Println("Initialization executed successfully")
	} else {
		fmt.Println("Initialization error")
	}
	SetMotor(DirectionDown)
	for {
		floor := GetFloorSensorSignal()
		switch floor {
		case 1:
			SetMotor(DirectionStop)
			SetFloorIndicator(1)
			break
		case 2:
			SetMotor(DirectionStop)
			SetFloorIndicator(2)
			break
		case 3:
			SetMotor(DirectionStop)
			SetFloorIndicator(3)
			break
		case 4:
			SetMotor(DirectionStop)
			SetFloorIndicator(4)
			break
		}
	}
}

func SetDoorLight(open bool) {
	if open {
		io.SetBit(lightDoorOpen)
	} else {
		io.ClearBit(lightDoorOpen)
	}
}

func SetMotor(dir int) {
	if dir == DirectionStop {
		io.WriteAnalog(motor, 0)
	} else if dir == DirectionDown {
		io.SetBit(motordir)
		io.WriteAnalog(motor, motorSpeed)
	} else if dir == DirectionUp {
		io.ClearBit(motordir)
		io.WriteAnalog(motor, motorSpeed)
	}
}

func SetFloorIndicator(floor int) {
	switch floor {
	case 1:
		io.ClearBit(lightFloorInd1)
		io.ClearBit(lightFloorInd2)
	case 2:
		io.SetBit(lightFloorInd2)
		io.ClearBit(lightFloorInd1)
	case 3:
		io.ClearBit(lightFloorInd2)
		io.SetBit(lightFloorInd1)
	case 4:
		io.SetBit(lightFloorInd1)
		io.SetBit(lightFloorInd2)
	}
}

func GetFloorSensorSignal() int {
	if io.ReadBit(sensorFloor1) == 1 {
		return 1
	} else if io.ReadBit(sensorFloor2) == 1 {
		return 2
	} else if io.ReadBit(sensorFloor3) == 1 {
		return 3
	} else if io.ReadBit(sensorFloor4) == 1 {
		return 4
	} else {
		return 0
	}
}
