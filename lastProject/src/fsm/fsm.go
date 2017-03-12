package fsm

import (
	"../elevatorHW"
	"fmt"
	"time"
	"sync"
)

var OperatingElevators int
var OperatingElevatorsPrt *int
var TimeStamp int64
var TimeStampPtr *int64
var latestFloorPtr *int
var LatestFloor int
//var mu sync.Mutex
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

/*func floorInQueue(queue []int, floor int)bool{
	for queueElement := range queue {
		if queueElement == floor{
			return true
		}
	}
	return false
}*/
	//	mu.Unlock()


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

func HandleTimeOutOrder(hallButtonsMap map[Order]int64, mutex sync.Mutex){
	for{
		if len(hallButtonsMap) > 0{
			for key,value := range hallButtonsMap{
				if (time.Now().Unix() - value) >  15 {
					mutex.Lock()
					PutOrderInLocalQueue(key)
					mutex.Unlock()
				}
			}
		}
		time.Sleep(500*time.Millisecond)
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
	lengthInsideQueue := len(localQueue[0])
	// // lengthUpQueue := len(localQueue[1])
	// // lengthDownQueue := len(localQueue[2])
	if currentDirection == 1 || currentDirection == 0 {
		if currentFloor != 0 {
			if len(localQueue[0]) > 0 && currentState == 0 {
				if (localQueue[0][0] < currentFloor) {
					elevatorHW.SetMotor(elevatorHW.DirectionDown)
					return
				} else if localQueue[0][0] > currentFloor {
					elevatorHW.SetMotor(elevatorHW.DirectionUp)
					return
				}
			} else if len(localQueue[1]) > 0 {
				if (localQueue[1][0] < currentFloor){
					if lengthInsideQueue > 0{
						SortLocalBiggestFirst()
						if localQueue[0][0] > currentFloor && currentDirection == 0{
							elevatorHW.SetMotor(elevatorHW.DirectionUp)
							return
						}
					}
					elevatorHW.SetMotor(elevatorHW.DirectionDown)
					return
				} else if localQueue[1][0] > currentFloor {
					elevatorHW.SetMotor(elevatorHW.DirectionUp)

					return
				}
			} else if len(localQueue[2]) > 0 {
				if localQueue[2][0] < currentFloor {
					elevatorHW.SetMotor(elevatorHW.DirectionDown)
					return
				} else if localQueue[2][0] > currentFloor {
					if lengthInsideQueue > 0{
						SortLocalQueue()
						if localQueue[0][0] < currentFloor && currentDirection == 1{
							elevatorHW.SetMotor(elevatorHW.DirectionDown)
							return
						}
					}
					elevatorHW.SetMotor(elevatorHW.DirectionUp)
					return
				}
			}
		}
	}
}
/*func NewStopAtThisFloor(orderCompletedCh chan<- Order){
	currentState := elevatorHW.GetElevatorState()
	currentFloor := elevatorHW.GetFloorSensorSignal()
	currentDirection := elevatorHW.GetElevatorDirection() // 1 is going down, 0 is going up
	numberOfDownOrders:= len(localQueue[2])
	numberOfUpOrders := len(localQueue[1])
	numberOfLocalOrders := len(localQueue[0])

	for i := 0; i < 3; i++ {
		if i==0{
			switch currentDirection{
				case
			}
		}
	}
}*/

func StopAtThisFloor(orderCompletedCh chan<- Order) {
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
						if currentState == 1 { // Going up
							elevatorHW.SetMotor(elevatorHW.DirectionStop)
							elevatorHW.SetInsideLight(currentFloor, false)
							ArrivedAtFloorSetDoorOpen(currentFloor)
							if i == 1{
								elevatorHW.SetUpLight(currentFloor, false)
								DeleteIndexLocalQueue(i,j)
								orderCompletedCh <- Order{currentFloor, elevatorHW.ButtonCallUp}
								return
							}else if i == 2{
								elevatorHW.SetDownLight(currentFloor, false)
								DeleteIndexLocalQueue(i,j)
								orderCompletedCh <- Order{currentFloor, elevatorHW.ButtonCallDown}
								return
							}
						} else {
							elevatorHW.SetMotor(elevatorHW.DirectionStop)
							ArrivedAtFloorSetDoorOpen(currentFloor)

							if i == 1{
								DeleteIndexLocalQueue(i,j)
								orderCompletedCh <- Order{currentFloor, elevatorHW.ButtonCallUp}
								elevatorHW.SetUpLight(currentFloor, false)
								return
							}else if i == 2{
								DeleteIndexLocalQueue(i,j)
								elevatorHW.SetDownLight(currentFloor, false)
								orderCompletedCh <- Order{currentFloor, elevatorHW.ButtonCallDown}
								return
							}
						}
					}
				}
				if (i == 0 || i == 1) && currentState == 1 {
					elevatorHW.SetMotor(elevatorHW.DirectionStop)
					ArrivedAtFloorSetDoorOpen(currentFloor)
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
					ArrivedAtFloorSetDoorOpen(currentFloor)
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
					ArrivedAtFloorSetDoorOpen(currentFloor)
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
					ArrivedAtFloorSetDoorOpen(currentFloor)
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

func ArrivedAtFloorSetDoorOpen(floor int) {
	TimeStampPtr = &TimeStamp
	elevatorHW.SetFloorIndicator(floor)
	elevatorHW.SetDoorLight(true)
	*TimeStampPtr = time.Now().Unix()
}

func DoorLightTimeOut(){
	currentTime := time.Now().Unix()
	if (currentTime - TimeStamp) > 2{
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

func RunElevator(orderCompletedCh chan<- Order) {
	for {
		DoorLightTimeOut()
		//HandleTimeOutOrder(hallButtonsMap, mutex)
		SetLatestFloor()
		StopAtThisFloor(orderCompletedCh)
		DoorLightTimeOut()
		SetElevatorDirection()
		StopButtonPressed()
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
