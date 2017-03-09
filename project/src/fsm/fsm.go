package fsm

import (
	"../elevatorHW"
	//"../olasnetwork"
	"fmt"
	//"math"
	"time"
)

type HelloMsg struct {
	Message          string
	Iter             int
	MyElevatorNumber int // This number identifies the elevator
	CurrentState     int // This number, says if the elevator is moving up (1) / down (-1) / idle (0)
	LastFloor        int // The last floor the elevator visited
	//GlobalQueue [][]int
}

var OperatingElevators int
var OperatingElevatorsPrt *int

var timeStamp int64
var currentTime int64
var timeStampPtr *int64
var latestFloorPtr *int
var LatestFloor int

/*
Order type - description
{Floor, 0} Calldown from #Floor
{Floor, 1} Callup from #Floor
{Floor, 2} InsideOrder to #Floor

*/

type Order struct {
	Floor  int
	Button elevatorHW.ButtonType //0 is callDown, 1 callUp, 2 callInside
}

type ElevatorStatus struct {
	lastFloor int
	direction elevatorHW.MotorDirection
}

func ArrivedAtFloorSetDoorOpen(floor int) {
	timeStampPtr = &timeStamp
	*timeStampPtr = time.Now().Unix()
	elevatorHW.SetFloorIndicator(floor)
	elevatorHW.SetDoorLight(true)
}

func GetButtonsPressed(buttonCh chan<- Order) {
	for {
		upOrder := elevatorHW.GetUpButton()
		downOrder := elevatorHW.GetDownButton()
		insideOrder := elevatorHW.GetInsideElevatorButton()

		var order Order
		if upOrder != 0 {
			order.Floor = upOrder
			order.Button = elevatorHW.ButtonCallUp
			buttonCh <- order
			time.Sleep(60 * time.Millisecond)
		}
		if downOrder != 0 {
			order.Floor = downOrder
			order.Button = elevatorHW.ButtonCallDown
			buttonCh <- order
			time.Sleep(60 * time.Millisecond)
		}
		if insideOrder != 0 {
			order.Floor = insideOrder
			order.Button = elevatorHW.ButtonCommand
			buttonCh <- order
			time.Sleep(60 * time.Millisecond)
		}
	}
}

func PutOrderInLocalQueue(newOrder Order) {
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
	if newOrder.Floor != 0 {
		buttontype := newOrder.Button
		switch buttontype {
		case elevatorHW.ButtonCallDown:
			AppendDownOrder(newOrder.Floor)
			elevatorHW.SetDownLight(newOrder.Floor, true)
		case elevatorHW.ButtonCallUp:
			AppendUpOrder(newOrder.Floor)
			elevatorHW.SetUpLight(newOrder.Floor, true)
		case elevatorHW.ButtonCommand:
			AppendInsideOrder(newOrder.Floor)
			elevatorHW.SetInsideLight(newOrder.Floor, true)
		}
	}
}

func PutInsideOrderInLocalQueue(newOrder Order) {
	insideOrder := elevatorHW.GetInsideElevatorButton()
	if insideOrder != 0 {
		AppendInsideOrder(insideOrder)
		elevatorHW.SetInsideLight(insideOrder, true)
	}
	AppendInsideOrder(newOrder.Floor)
}

/*func PutOrderInGlobalQueue() {
	upOrder := elevatorHW.GetUpButton()
	downOrder := elevatorHW.GetDownButton()
	if upOrder != 0 {
		AppendGlobalUpOrder(upOrder)
		elevatorHW.SetUpLight(upOrder, true)
	}
	if downOrder != 0 {
		AppendGlobalDownOrder(downOrder)
		elevatorHW.SetDownLight(downOrder, true)
	}
}*/

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
	currentDirection := elevatorHW.GetElevatorDirection() // 1 is going down, 0 is going up

	for i := 0; i < 3; i++ {

		for j := range localQueue[i] {
			if currentFloor == localQueue[i][j] {
				if (i == 0 || i == 1) && (currentDirection == 0 || currentFloor == 1) {
					elevatorHW.SetMotor(elevatorHW.DirectionStop)
					ArrivedAtFloorSetDoorOpen(currentFloor)
					//elevatorHW.SetDownLight(currentFloor, false)
					elevatorHW.SetUpLight(currentFloor, false)
					elevatorHW.SetInsideLight(currentFloor, false)
					elevatorHW.SetFloorIndicator(currentFloor)
					DeleteIndexLocalQueue(i, j)
					return
				} else if (i == 0 || i == 2) && (currentDirection == 1 || currentFloor == 4) {
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
	PrintQueues()

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

func SetLatestFloor() {
	latestFloorPtr = &LatestFloor
	if elevatorHW.GetFloorSensorSignal() == 1 || elevatorHW.GetFloorSensorSignal() == 2 || elevatorHW.GetFloorSensorSignal() == 3 || elevatorHW.GetFloorSensorSignal() == 4 {
		*latestFloorPtr = elevatorHW.GetFloorSensorSignal()
		elevatorHW.SetFloorIndicator(elevatorHW.GetFloorSensorSignal())
	}
}
func StartUpMessage() {
	fmt.Println("DO YOU EVEN LIFT BRO?")
	time.Sleep(140 * time.Millisecond)
	fmt.Print("\nFuck yeeah bro!\n\n")
	time.Sleep(100 * time.Millisecond)

}
func RunElevator() {
	for {
		SetLatestFloor()
		SetElevatorDirection()
		//PutOrderInLocalQueue(Order{}) // This must be replaced by "PutOrderInGlobalQueue"
		//PutInsideOrderInLocalQueue(Order{})
		/*PutOrderInGlobalQueue()*/
		StopAtThisFloor()
		TurnOffDoorLight()
		StopButtonPressed()
	}
}

/*func GetElevatorState() elevatorState {
	elevatorState.previousFloor
}
*/

/*func costFunc(elevatorStates []elevatorState, newOrder order) {
	distanceCost := 0
	directionCost := 0
	totalCost := 0

	costs := [3]float64{math.Inf(1), math.Inf(1), math.Inf(1)}

	for i := 0; i < len(elevatorStates); i++ {
		distanceCost = math.Abs(elevatorStates[i].previousFloor - newOrder.Floor)
		if (elevatorStates[i].previousFloor == newOrder.Floor) && (elevatorStates[i].direction == DirectionStop){
			distanceCost = 0
		}
			//NEED PLENTY MORE LOGIC HERE
	}


}*/
