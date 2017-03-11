package fsm

import (
	"../elevatorHW"
	"fmt"
	"time"
)

var OperatingElevators int
var OperatingElevatorsPrt *int
var latestFloorPtr *int
var LatestFloor int

/*
Order type - description
{Floor, 0} Calldown from IFloor
{Floor, 1} Callup from IFloor
{Floor, 2} InsideOrder to IFloor
*/
type Order struct {
	Floor  int
	Button elevatorHW.ButtonType //0 is callDown, 1 callUp, 2 callInside
}

type ElevatorStatus struct {
	alive      bool
	elevatorID string
}

func ArrivedAtFloorSetDoorOpen(floor int, timeOut chan<- bool) {
	elevatorHW.SetFloorIndicator(floor)
	elevatorHW.SetDoorLight(true)
	time.Sleep(3 * time.Second)
	timeOut <- true
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
	if newOrder.Floor != -1 {
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

func HandleTimeOutOrder(hallButtonsMap map[Order]int64){
	for{
		if len(hallButtonsMap) > 0{
			for key,value := range hallButtonsMap{
				if value < time.Now().Unix() + 6 {
					PutOrderInLocalQueue(key)
				}
			}
		}
	}
}

func PutInsideOrderInLocalQueue() {
	insideOrder := elevatorHW.GetInsideElevatorButton()
	if insideOrder != 0 {
		AppendInsideOrder(insideOrder)
		elevatorHW.SetInsideLight(insideOrder, true)
	}
}

func SetElevatorDirection() {
	if elevatorHW.GetDoorLight() != 0 {
		return
	}
	currentDirection := elevatorHW.GetElevatorDirection()
	currentState := elevatorHW.GetElevatorState()
	currentFloor := elevatorHW.GetFloorSensorSignal()

	if currentDirection == 1 || currentDirection == 0 {
		if currentFloor != 0 {
			if len(localQueue[0]) > 0 && currentState == 0 {
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

func StopAtThisFloor(timeOut chan<- bool, orderCompletedCh chan<- Order) {
	currentState := elevatorHW.GetElevatorState()
	currentFloor := elevatorHW.GetFloorSensorSignal()
	currentDirection := elevatorHW.GetElevatorDirection() // 1 is going down, 0 is going up
	downOrders := len(localQueue[2])
	upOrders := len(localQueue[1])
	localOrders := len(localQueue[0])

	for i := 0; i < 3; i++ {
		for j := range localQueue[i] {
			if currentFloor == localQueue[i][j] {
				if len(localQueue[0]) == 0 {
					if localQueue[i][j] == currentFloor {
						if currentDirection == 1 { // Going up
							elevatorHW.SetMotor(elevatorHW.DirectionStop)
							ArrivedAtFloorSetDoorOpen(currentFloor, timeOut)
							elevatorHW.SetUpLight(currentFloor, false)

						} else {
							elevatorHW.SetMotor(elevatorHW.DirectionStop)
							ArrivedAtFloorSetDoorOpen(currentFloor, timeOut)
							elevatorHW.SetDownLight(currentFloor, false)
						}
					}
				}
				if (i == 0 || i == 1) && (currentDirection == 0 || currentFloor == 1) {
					elevatorHW.SetMotor(elevatorHW.DirectionStop)
					ArrivedAtFloorSetDoorOpen(currentFloor, timeOut)
					elevatorHW.SetUpLight(currentFloor, false)
					elevatorHW.SetInsideLight(currentFloor, false)
					elevatorHW.SetFloorIndicator(currentFloor)
					DeleteIndexLocalQueue(i, j)
					if i != 0 && i == 1 {

						orderCompletedCh <- Order{currentFloor, elevatorHW.ButtonCallUp}
					} else if i != 0 {
						orderCompletedCh <- Order{currentFloor, elevatorHW.ButtonCallDown}
					}
					return

				} else if (i == 0 || i == 2) && (currentDirection == 1 || currentFloor == 4) {
					elevatorHW.SetMotor(elevatorHW.DirectionStop)
					ArrivedAtFloorSetDoorOpen(currentFloor, timeOut)
					elevatorHW.SetDownLight(currentFloor, false)
					elevatorHW.SetInsideLight(currentFloor, false)
					elevatorHW.SetFloorIndicator(currentFloor)
					DeleteIndexLocalQueue(i, j)
					if i != 0 && i == 1 {

						orderCompletedCh <- Order{currentFloor, elevatorHW.ButtonCallUp}
					} else if i != 0 {
						orderCompletedCh <- Order{currentFloor, elevatorHW.ButtonCallDown}
					}
					return
				}
			}
		}
	}
	if (downOrders > 0) && (localOrders > 0) && currentState == 1 {
		for i := range localQueue[0] {
			for j := range localQueue[2] {
				if (localQueue[0][i] < localQueue[2][j]) && (localQueue[2][j] == currentFloor) {
					elevatorHW.SetMotor(elevatorHW.DirectionStop)
					ArrivedAtFloorSetDoorOpen(currentFloor, timeOut)
					elevatorHW.SetDownLight(currentFloor, false)
					elevatorHW.SetInsideLight(currentFloor, false)
					elevatorHW.SetFloorIndicator(currentFloor)
					DeleteIndexLocalQueue(2, j)
					orderCompletedCh <- Order{currentFloor, elevatorHW.ButtonCallDown}
				}
			}
		}
	}
	if (upOrders > 0) && (localOrders > 0) && currentState == -1 {
		for i := range localQueue[0] {
			for j := range localQueue[1] {
				if (localQueue[0][i] > localQueue[1][j]) && (localQueue[1][j] == currentFloor) {
					elevatorHW.SetMotor(elevatorHW.DirectionStop)
					ArrivedAtFloorSetDoorOpen(currentFloor, timeOut)
					elevatorHW.SetUpLight(currentFloor, false)
					elevatorHW.SetInsideLight(currentFloor, false)
					elevatorHW.SetFloorIndicator(currentFloor)
					DeleteIndexLocalQueue(1, j)
					orderCompletedCh <- Order{currentFloor, elevatorHW.ButtonCallUp}
				}
			}
		}
	}
}

func StopButtonPressed(timeOut chan<- bool) {
	if !elevatorHW.GetStopButtonPressed() {
		return
	}
	fmt.Println("Initiating emergency procedure")
	elevatorHW.SetStopButton(true)
	elevatorHW.SetMotor(elevatorHW.DirectionStop)
	DeleteLocalQueue()
	time.Sleep(2000 * time.Millisecond)
	elevatorHW.SecondInit()
	ArrivedAtFloorSetDoorOpen(elevatorHW.GetFloorSensorSignal(), timeOut)
	fmt.Println("Operaing Normally")
}

func SetLatestFloor() {
	latestFloorPtr = &LatestFloor
	if elevatorHW.GetFloorSensorSignal() == 1 || elevatorHW.GetFloorSensorSignal() == 2 || elevatorHW.GetFloorSensorSignal() == 3 || elevatorHW.GetFloorSensorSignal() == 4 {
		*latestFloorPtr = elevatorHW.GetFloorSensorSignal()
		elevatorHW.SetFloorIndicator(elevatorHW.GetFloorSensorSignal())
	}
}


func RunElevator(timeOut chan<- bool, orderCompletedCh chan<- Order) {
	for {
		SetLatestFloor()
		StopAtThisFloor(timeOut, orderCompletedCh)
		SetElevatorDirection()
		StopButtonPressed(timeOut)
	}
}


func StartUpMessage() {
	fmt.Println("\n\nDO YOU EVEN LIFT BRO?")
	time.Sleep(800 * time.Millisecond)
	fmt.Print("\nFuck yeeah bro!\n\n")
	time.Sleep(800 * time.Millisecond)
	fmt.Println("Im fucking Arnold bruvh!!")
	time.Sleep(800 * time.Millisecond)
	fmt.Println("Here's a picture of me:")
	time.Sleep(1000 * time.Millisecond)
	fmt.Println(`##(*#########################(####(##################################################(((((####((#(((((#(((((((((((((###(#########
###%/*(########(#######((####((((################################################(((((((((((#(((((((((((#(((((((##(((###(########
#####&(**##############(((###(((((((#############################################((((((((((##(((((((###(#(((((((######((((((#####
#####(&&(***((######(((((((((((((((#######(#####################################((((((((((((####((((((((((((((((((#(((((((((#####
#######%&&&%///###((((#((((#((((((((((((((((##########################(((((((((((((((((((((((((((((((((((((((((((((((((((((((((
#########%&&&/**#&&(###((((((((((((((((((###((((####################(######(((((#((((((((((((((((((((((((((((((((((((((((((((((((
#######(###&&%#%%%%&%#((((((((((((((((((((((((((####################%&&@@@@@@@@@@@@%&((((((((((((((((((((((((((((((((((((((#(((((
##########(###&&&&&((#(((((((((((((((((((((#((((((((#############@@@@&@@@@@&@&@&@@@@@#(((((((((((((((((((((((((((((((((((((((((
#######(((((#((#&@@@@##%(((((((((((((((((((((((((((((((###((####(@&@@&@@@&&@&&%&@&&(((((((((((((((((((((#((((((((((((((((((
#######((((((((#(((&@@#####&(((((((((((((((((((((((((((((##((((#(##//#&&&&&&&&&&@@@@@@@%((((((((((((((((((((#((((((((((((((((((
#######(((((((((((#(&&%(#((%#%(((((((((((((((((((((((((#((((((((##((//*,//%&@@@&@@@@@&@@@@(((((((((((((((((((((((((((((((((((((((
#######((((((((((((((&%((#(####%((((((((((((((((((((((((((((((((#(#&%%#/**/@@@@@&@@@@@@@@@#((((((((((((((((((((((((((((((((((((((
########((((((((((((#(&((((####((##(((((((((((((((((((((((((((((#(,,((*,*,*%@&&&&@&@@@@@@@@((((((((((((((((((((((((((((((((((((((
########(#(((((((((((((&(((((((((((##%((((((((((((((((#(((((((##(,,//,,*//((#/%@@@@@@@@@@@@&(((((((((((((((((((((((((((((((((((((
((##(###(#((((((((((((((%#(((/(/((###%%#&(((((((((((((((((((((((%&%*(/((//(((((%@@@@@@@@@@@@(((((((((((((((((((((((((((((((((((((
((##(((((((((((((((((((((%#(##(//(##%%%%%%%((((((((((((((((((((((#%#%//(//(((((((@@@@@@@@@@#(((((((((((((((((((((((((((((((((((((
##((#(((((((((((((((((((((#%####(#%%%%%%%%%%#/((((((((((((((((((((/*/((((##%%#///(((@@@@@@%((((((((((((((((((((((((((((((((((((((
((((((((((((((((((((((((((((&%#(((###((#/&%#(/**/(((((((((((((((((/////(%&&///((((#%&&@((((((((((((((((((((((((((((((((((((((
(#((((((((((((((((((((((((((((&&%%##(/(((,(..,,,*/((((((((((((((((%#%&&&&&&%((/(((((/////(#/(((((((((((((((((((((((((((((((((((((
##((((((((((((((((((((((((((((((&&&&&%%##%*/*,,,,**//((#%(((((((((((((((((###(((//(*,/**///(((#((((((((((((((((((((((((((((((((((
((((((((((#(((((((((((((((((((((((#@@@&@@#(((((((/(/((#&(((((((((((((((((@&((((/*/*//////(///(%((((((((((((((((((((((((((((((((
((((((((((((((((((((((((((((((((((((#@%(**/#%%%%%##(((((#%(//((#(##(##(((%/(((((/(((//////(((((#%((((((((((((((((((((((((((((((
(((((((((((((((((((((((((((((((((((((((&%((#(%&%&&&%%##(#####/(((((((((((((((((((((//////***/////(((%%%/(((((((((((((((((((((((((
(((((((((((((((((((((((((((((((((((((((((&@&%@@&&&@@@@@&&&&%##((((((((//(((((((((((////*,*.....,*/((#%#%%((((((((((((((((((((((((
(((((((((((((((((((((((((((((((((((((((((((@@&%##(/(##&&@@&%#((((/////,/((((((/((/*,..***/***,..,//(%##%&%(((((((((((((((((((((((
(((((((((((((((((((((((((((((((((((((((((((((&%##(((((((((#%(((((((/*.,((((////*,.,/(((((/((/////((#%##%%&/((((((((((((((((((((((
((((((((((((((((((((((((((((((((((((((((((((((%@&%%%%##%%%%%((((((/,,,((((//,..,*/((#%%#(((((((((((##(((#%%((((((((((((((((((((((
(((((((((((((((((((((((((((((((((((((((((((((((((&@@@@@@@@@&((((((/.*(((/////((((##%&@@@&%((((//(((///(((##&/((((((((((((((((((((
(((((((((((((((((((((((((((((((((((((((((((((((((((/#@@@@@@@#(((///*(((((((###%%&&&&&@@@@%*///********/(((#%&/(((((((((((((((((((
(((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((#((//((#%%%%&&&@@@@@@&&@@&@@/(////((/((((/((((#%&/((((((((((((((((((
(((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((%%###%(&&&@&@@&&@@@@&&&//((((/*,**/,,((**(((#%&%((((((((((((((((((
(((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((&(((((##%%(/,,,*/(((//***(*((***///((**/((#%&@/(((((((((((((((((
(((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((/((((((///***,*,../(((//*,*###//(*/*/((/**(##%&&((((((((((((((((((
((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((#((##(((//,.,,.,/((/*/***/(#%##%#(/*/((//(###%&((((((((((((((((((
(((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((/%%%%%%%%(//,*/(((((((#((/((#%((((*,,/(((#(###%&((((((((((((((((((
((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((%##%%#%%(//((//(((((((####(((//(((((##%#(##%&((((((((((((((((((
(((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((/%%#%%####%%#(((######%%%%(((((/(##%%%####*/##%&((((((((((((((((((
(((((((((((((((((((((((((((((((((((((((((((/((((((((((((((((((((&%%%%%##%%#%&&%#(##%%%#((((//(,/(((((//***(&%&@#(/(((((((((((((((
(((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((/(((#@&%%#%%#((//(((/,*//(/*/**/(((((((###(#&@&&@&((/(((((((((((((((
((((((((((((((((((((((((((((((////(/(((((((((((((((((((((((/((((((&%#(%(((/////(*(////(/((((#####%&&&&@@@@@@%/(/(((((((((((((((((
(((((((((((((((((((((((((((((((((/((((((((/(((((((/((((/((///(/(/(#(((((((////,,*////(#####%&@@@@@@@@@@@(((/((/(((/((((((((((((
((((((((((((((((((((((((((((//(((//(((((((((((((((/((((((((/(/(((/((((((/*///**//((#%%&&&&@@@@@@@@@@@@#/(((((/((/((((((((((((((((
((((((((((((((((((((((((((((((/(((((((((/(((/(((((((//(//((((/(//#(((((/*////%%&@@@@@@@@@%//(/////(//((/((((((//((/(///(/(/////((
((((((((((((((((((((((((((((/(/(((((///((/(///(//(/((//((/((/(//(#(((((((((*,,/#%&@@@@/////(((/////(//(///////////////////(((((
(((((((((((((((/((((((/((((((/(/(/((//////////////(/////////////(#((/((*(//***(%@@%%&@%//////////////////////////////(/////(/((/(
(((((((((((/(((/(/(//(/((((((((//(/////////////////////////////(##((%&%&(****(##,&@@@&%//////////////////////////////////////(/((
((((((((((((((((((/(/(/(///////////////////////////////////////(#%&&&&&&%%&&&&%%@&@@@@@*/////////////////////////////////////(///
((((((((/(((/////(//////////((///////////////////////////////%%%&&&&&&&&&@@@@@@@@@@@@@@@(/////////////////////////////(/////(((//
((((((((/(///////////////////////////////////////////////////&&&&&&&&&&&&&&@@@@@@@@@@@@@@//////////////////(/////////////////////
(//((((////////////////////////////////////////////////////&&&&&&&&&/*//,*(/(&@@@@@@@@@@@@///////////////////////////////////////
((/(/(/////////////////////////////////////////////////////&&&&@&&(//,*/**/((#%&@@@@@@@@@@///////////////////////////////////////
(((((//((((////////////////////////////////////////////////&&&@@@&(/**,/*,//((#%##&@@@@@@@#//////////////////////////////////((/(
((/((//(((/////////////////////////////////////////////////&@@@@@(//*,,*/**/((###((((@@@@@#//////////////////////////////////////
((/((((///(//////////////////////////////////////////////(/*&@@@#(/*,,,*,*/((((((##((#%@@@///////////////////////////////////////
///((/(///(/////////////////////////////////////////////(////&@@((/*,.,,*,*/((((((###%&@@&///////////////////////////////////////
(((((//////////////////////////////////////////////////((///***%(//*,..****////(((#%%&&@&////////////////////////////////////////
//(((((((((////////////////////////////////////////////(/***/**,(//*,..*/**////(((#%&@@&/////////////////////////////////////////
///(//////////////////////////////////////////////////(//***//**((//,,.*/*/////((((%&@@//////////////////////////////////////////
((((/(//((///////////////////////////////////////////((/****/(//((//**,,////(((((((%&&%*/////////////////////////////////////////
(//(///(((///////////////////////////////////////////(//***/////(((//*,,/((((((((((#%&&//////////////////////////////////////////
////((((/////////////////////////////////////////////(///***//(/((((//**/(((((((((((%&@&/////////////////////////////////////////
////////////////////////////////////////////////////(((//****///((((//////((((((((((%&@@/////////////////////////////////////////
////(////////////////////////////////////////////////((//*,***//(#((((///(((((((((((%&&@&////////////////////////////////////////
/////((//////////////////////////////////////////////((//*.***///#(((((((/(((((((((#%&@@@////////////////////////////////////////`)

}
