package fsm

import (
	"../elevatorHW"
	"../io"
	//"../network/peers"
	"fmt"
	//"math"
	"time"
)

var timeStamp int64
var currentTime int64
var timeStampPtr *int64

type Order struct{
	FromFloor int
	Button elevatorHW.ButtonType
}

func ArrivedAtFloorSetDoorOpen(floor int) {
	timeStampPtr = &timeStamp
	*timeStampPtr = time.Now().Unix()
	elevatorHW.SetFloorIndicator(floor)
	elevatorHW.SetDoorLight(true)
}

func PutOrderInLocalQueue(newRemoteOrder Order) {
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
	if newRemoteOrder.FromFloor != -1 {
		buttontype := newRemoteOrder.Button
		switch buttontype {
		case elevatorHW.ButtonCallDown:
			AppendDownOrder(newRemoteOrder.FromFloor)
			elevatorHW.SetDownLight(downOrder, true)
		case elevatorHW.ButtonCallUp:
			AppendUpOrder(newRemoteOrder.FromFloor)
			elevatorHW.SetDownLight(downOrder, true)
		case elevatorHW.ButtonCommand:
			AppendInsideOrder(insideOrder)
			elevatorHW.SetInsideLight(insideOrder, true)
		}
	}
}

func SetElevatorDirection() {
	currentTime := int64(time.Now().Unix())
	if currentTime-timeStamp < 4 {
		return
	}

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
	currentDirection := io.ReadBit(elevatorHW.Motordir) // 1 is going down, 0 is going up

	for i := 0; i < 3; i++ {

		for j := range localQueue[i] {		
			if currentFloor == localQueue[i][j] {
				if ((i==0 || i==1) && (currentDirection == 0 || currentFloor == 1)){ 
					elevatorHW.SetMotor(elevatorHW.DirectionStop)
					ArrivedAtFloorSetDoorOpen(currentFloor)
					//elevatorHW.SetDownLight(currentFloor, false)
					elevatorHW.SetUpLight(currentFloor, false)
					elevatorHW.SetInsideLight(currentFloor, false)
					elevatorHW.SetFloorIndicator(currentFloor)
					DeleteIndexLocalQueue(i, j)
					return
				} else if((i==0 || i==2) && (currentDirection == 1 || currentFloor == 4)){
					elevatorHW.SetMotor(elevatorHW.DirectionStop)
					ArrivedAtFloorSetDoorOpen(currentFloor)
					elevatorHW.SetDownLight(currentFloor, false)
					//elevatorHW.SetUpLight(currentFloor, false)
					elevatorHW.SetInsideLight(currentFloor, false)
					elevatorHW.SetFloorIndicator(currentFloor)
					DeleteIndexLocalQueue(i, j)
					return
				}
			} 
		}
	}
}

func PrintElevatorStatus() { //For debbung Purposes
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

func TurnOffDoorLight() {
	currentTime := int64(time.Now().Unix())
	if currentTime-timeStamp > 3 {
		elevatorHW.SetDoorLight(false)
	}
}

func StopButtonPressed() {
	if !elevatorHW.GetStopButtonPressed() {
		return
	}
	fmt.Println("Initiating emergency procedure")
	elevatorHW.SetStopButton(true)

	elevatorHW.SetMotor(elevatorHW.DirectionStop)
	DeleteLocalQueue()
	time.Sleep(2000 * time.Millisecond)
	elevatorHW.SecondInit()
	ArrivedAtFloorSetDoorOpen(elevatorHW.GetFloorSensorSignal())
	fmt.Println("Operaing Normally")
}

func RunElevator() {
	//CreateQueueSlice()
	t := 0
	for {
		t++
		SetElevatorDirection()
		PutOrderInLocalQueue(Order {})
		StopAtThisFloor()
		TurnOffDoorLight()
		StopButtonPressed()
		if t%20000 == 0 {
			PrintLocalQueue()
		}
		if t > 100000 {
			t = 0
		}
	}
}


// Un-used code below
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
