package main

import (
	"./src/elevatorHW"
	"./src/fsm"
	"./src/network"
	"./src/backup"
	"fmt"
	"time"
	"sync"
)

// This function returns how suitet the elevator is to handle a global call
func costFunction(dir int, lastFloor int, order fsm.Order) int {
	var distanceToTarget int
	var dirCost int
	if lastFloor > order.Floor {
		distanceToTarget = lastFloor - order.Floor
	} else {
		distanceToTarget = order.Floor - lastFloor
	}

	distCost := distanceToTarget

	if dir == 0 && order.Floor == lastFloor { // Elevator is Idle at floor being called
		return 0
	}
	if dir == 0 {
		return distanceToTarget * 2
	}
	if order.Button == 1 { //UpType Order
		if dir == -1 { // Moving in opposite direction
			if lastFloor < order.Floor {
				distCost = 4 * distCost
				dirCost = 8
			} else {
				dirCost = 10
			}
		} else if dir == 0 { //Elevator is idle
			dirCost = 0
		} else if dir == 1 { // Elevator is moving up
			if lastFloor > order.Floor {
				dirCost = 8
			} else {
				dirCost = 10
			}
		}

	} else if order.Button == 0 { //DownType order
		if dir == 1 { // Oposite directioin
			if lastFloor > order.Floor {
				distCost = 4 * distCost
				dirCost = 8
			} else {
				dirCost = 10
			}
		} else if dir == 0 {
			dirCost = 0
		} else if dir == -1 {
			if lastFloor < order.Floor {
				dirCost = 8
			} else {
				dirCost = 10
			}
		}
	}
	return (2 * dirCost) + distCost
}

func decitionmaker(onlineElevatorStates map[string]network.HelloMsg, newOrder fsm.Order) (string, int) {
	numberOfElevatorsInNetwork := len(onlineElevatorStates)
	if numberOfElevatorsInNetwork == 0 || numberOfElevatorsInNetwork == 1 {
		return network.GetLocalID(), 0
	}
	var elevatorWithLowestCost string
	lowestCost := 1337
	if len(onlineElevatorStates) < 2 {
		return network.GetLocalID(), 0
	}
	for key, value := range onlineElevatorStates {
		thisCost := costFunction(value.CurrentState, value.LastFloor, newOrder)
		// fmt.Print("\n\nStates of ")
		// fmt.Print(key)
		// fmt.Println(": ")
		// fmt.Print("Current state value: ")
		// fmt.Print(value.CurrentState)
		// fmt.Print("\nIts Last Floor: ")
		// fmt.Println(value.LastFloor)
		// fmt.Print("The Order is : ")
		// fmt.Println(newOrder)
		// fmt.Println(" ")
		// fmt.Print(key)
		// fmt.Print(" has a cost of ")
		// fmt.Println(thisCost)
		// fmt.Print("\n\n")
		if thisCost < lowestCost {
			lowestCost = thisCost
			elevatorWithLowestCost = key
		}
	}
	return elevatorWithLowestCost, lowestCost
}

func PrintLocalData(operatingElevatorStates map[string]network.HelloMsg){
	for{
		up := ColG + "↑" + ColN
		down := ColG + "↓" + ColN
		idle := ColR + "↺" + ColN
		localCurrentState := elevatorHW.GetElevatorState()
		lastFloor := fsm.LatestFloor
		otherElevatorsOnline := len(operatingElevatorStates)
		ordersBeingHandled := fsm.GetNumberOfOrders()
		myID := network.GetLocalID()
		flatLine := ColY +"_______________________________"+ ColN
		roof := ColY + "¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯" + ColN
		wall := ColY + "|" + ColN
		localElev := wall + ColY +" I am elevator " + string(myID) + "           "+wall+"\n"+ ColN
		
		fmt.Print(flatLine + "\n")
		fmt.Print(localElev)
		fmt.Print(roof + "\n")
		fmt.Print(ColY + wall + ColM + "Direction:           ")
		switch localCurrentState {
		case 1:
			fmt.Print(up+ "       " + wall)
		case -1:
			fmt.Print(down+ "       " + wall)
		case 0:
			fmt.Print(idle+ "       " + wall)

		}
		fmt.Print("\n")
		fmt.Print(wall + ColM + "Last floor:          ")
		fmt.Print(lastFloor)
		fmt.Print("       " + wall)
		fmt.Print("\n")
		fmt.Print(wall + ColM + "Number of orders:    ")
		fmt.Print(ordersBeingHandled)
		fmt.Print("       ")
		fmt.Print(wall+"\n")
		fmt.Print(wall + ColM + "Elevators online:    ")
		fmt.Print(otherElevatorsOnline)
		fmt.Print("       " + wall + "\n")
		fmt.Print(roof)
		fmt.Print("\n")
		fmt.Print(wall + ColY + "Online elevators             "+ ColY + wall + "\n")
		fmt.Print(roof + "\n")

		if len(operatingElevatorStates) > 0 {
			for key, value := range operatingElevatorStates{
				fmt.Print(wall + ColM + "Elevator ID:       " + key + "       " + wall + "\n")
				switch value.CurrentState{
				case 0:
					fmt.Print(wall + ColM + "Direction:          " +idle + "        " + wall + "\n")
				case 1:
					fmt.Print(wall + ColM + "Direction:          " +up + "        " + wall + "\n")
				case -1:
					fmt.Print(wall + ColM + "Direction:          "  + down + "        " + wall + "\n")
				}
				fmt.Print(wall + ColM + "LastFloor:           ")
				fmt.Print(value.LastFloor)
				fmt.Print("       " + wall + "\n")
				fmt.Print(roof + "\n")
			}
		}
		fmt.Print("\n\n")
		time.Sleep(400 * time.Millisecond)
		for k := 0; k < 15; k++{
			fmt.Print("\n")
		}
	}
}


func main() {
	//start init
	fmt.Println("Starting system")
	fmt.Print("\n\n")
	elevatorHW.Init()
	//fsm.StartUpMessage()
	//finished init
	queueFromBackup := backup.ReadBackupFromFile()
	for i := range(queueFromBackup){
		if queueFromBackup[i] != -1{
			elevatorHW.SetInsideLight(i+1,true)
		}
	}
	fsm.CreateQueueSlice(queueFromBackup)
	time.Sleep(10 * time.Millisecond)

	operatingElevatorStates := make(map[string]network.HelloMsg)
	hallButtonsMap := make(map[fsm.Order]int64)

	buttonCh := make(chan fsm.Order)
	messageCh := make(chan network.HelloMsg)
	receivedNetworkOrderCh := make(chan network.HelloMsg)
	sendOrderToPeerCh := make(chan network.OrderMsg)
	orderCompletedCh := make(chan fsm.Order)
	sendDeletedOrderCh := make(chan fsm.Order)
	backupCh := make(chan int)
	backupRemoveOrderCh := make(chan int)

	var mutex sync.Mutex
	go fsm.RunElevator(orderCompletedCh, backupRemoveOrderCh)
	go fsm.GetButtonsPressed(buttonCh)
	go network.NetworkMain(messageCh, receivedNetworkOrderCh, sendOrderToPeerCh, orderCompletedCh, sendDeletedOrderCh)
	go fsm.HandleTimeOutOrder(hallButtonsMap, mutex)
	go backup.FileBackup(backupCh,backupRemoveOrderCh)
	go PrintLocalData(operatingElevatorStates)


	for {
		select {
		case orderIsHandled := <-orderCompletedCh:
			mutex.Lock()
			delete(hallButtonsMap, orderIsHandled)
			mutex.Unlock()
			if orderIsHandled.Button == elevatorHW.ButtonCommand{
				backupRemoveOrderCh <- orderIsHandled.Floor
			}
			switch orderIsHandled.Floor{
			case 1:
				elevatorHW.SetUpLight(orderIsHandled.Floor, false)
			case 2:
				if orderIsHandled.Button == elevatorHW.ButtonCallDown{
					elevatorHW.SetDownLight(orderIsHandled.Floor, false)
				}else if orderIsHandled.Button == elevatorHW.ButtonCallUp{
					elevatorHW.SetUpLight(orderIsHandled.Floor, false)
				}
			case 3:
				if orderIsHandled.Button == elevatorHW.ButtonCallDown{
					elevatorHW.SetDownLight(orderIsHandled.Floor, false)
				}else if orderIsHandled.Button == elevatorHW.ButtonCallUp{
					elevatorHW.SetUpLight(orderIsHandled.Floor, false)
				}
			case 4:
				elevatorHW.SetDownLight(orderIsHandled.Floor, false)
			}
			sendDeletedOrderCh <- orderIsHandled

		case newMsg := <-messageCh:
			if newMsg.Order.Order.Floor != -1 {
				mutex.Lock()
				hallButtonsMap[newMsg.Order.Order] = time.Now().Unix()
				mutex.Unlock()
				switch newMsg.Order.Order.Floor{
				case 1:
					if newMsg.Order.Order.Button == elevatorHW.ButtonCallUp{
							elevatorHW.SetUpLight(1, true)
					}
				case 2:
					if newMsg.Order.Order.Button == elevatorHW.ButtonCallUp{
							elevatorHW.SetUpLight(2, true)
					}else if newMsg.Order.Order.Button == elevatorHW.ButtonCallDown{
						elevatorHW.SetDownLight(2, true)
					}
				case 3:
					if newMsg.Order.Order.Button == elevatorHW.ButtonCallUp{
							elevatorHW.SetUpLight(3, true)
					}else if newMsg.Order.Order.Button == elevatorHW.ButtonCallDown{
						elevatorHW.SetDownLight(3, true)
					}
				case 4:
					if newMsg.Order.Order.Button == elevatorHW.ButtonCallDown{
							elevatorHW.SetDownLight(1, true)
					}
				}
			}

			network.UpdateElevatorStates(newMsg, operatingElevatorStates)
			network.DeleteDeadElevator(operatingElevatorStates)
			if newMsg.Order.ElevatorToTakeThisOrder == network.GetLocalID() {
				fsm.PutOrderInLocalQueue(newMsg.Order.Order)
			}
			if newMsg.OrderExecuted.Floor != -1 {

				mutex.Lock()
				delete(hallButtonsMap, newMsg.OrderExecuted)
				mutex.Unlock()

				switch newMsg.OrderExecuted.Floor{
				case 1:
					elevatorHW.SetUpLight(newMsg.OrderExecuted.Floor, false)
				case 2:
					if newMsg.OrderExecuted.Button == elevatorHW.ButtonCallDown{
						elevatorHW.SetDownLight(newMsg.OrderExecuted.Floor, false)
					}else if newMsg.OrderExecuted.Button == elevatorHW.ButtonCallUp{
						elevatorHW.SetUpLight(newMsg.OrderExecuted.Floor, false)
					}
				case 3:
					if newMsg.OrderExecuted.Button == elevatorHW.ButtonCallDown{
						elevatorHW.SetDownLight(newMsg.OrderExecuted.Floor, false)
					}else if newMsg.OrderExecuted.Button == elevatorHW.ButtonCallUp{
						elevatorHW.SetUpLight(newMsg.OrderExecuted.Floor, false)
					}
				case 4:
					elevatorHW.SetDownLight(newMsg.OrderExecuted.Floor, false)
				}
			}

		case newOrder := <-buttonCh:

			if newOrder.Button == elevatorHW.ButtonCommand {
				fsm.PutInsideOrderInLocalQueue()
				backupCh <- newOrder.Floor
			} else if len(operatingElevatorStates) == 0 || len(operatingElevatorStates) == 1 {
				fsm.PutOrderInLocalQueue(newOrder)
			} else {
				elevatorToHandleThisOrder, _ := decitionmaker(operatingElevatorStates, newOrder)


				mutex.Lock()
				hallButtonsMap[newOrder] = time.Now().Unix()
				mutex.Unlock()

				switch newOrder.Floor{
				case 1:
					elevatorHW.SetUpLight(newOrder.Floor, true)
				case 2:
					if newOrder.Button == elevatorHW.ButtonCallDown{
						elevatorHW.SetDownLight(newOrder.Floor, true)
					}else if newOrder.Button == elevatorHW.ButtonCallUp{
						elevatorHW.SetUpLight(newOrder.Floor, true)
					}
				case 3:
					if newOrder.Button == elevatorHW.ButtonCallDown{
						elevatorHW.SetDownLight(newOrder.Floor, true)
					}else if newOrder.Button == elevatorHW.ButtonCallUp{
						elevatorHW.SetUpLight(newOrder.Floor, true)
					}
				case 4:
					elevatorHW.SetDownLight(newOrder.Floor, true)
				}
				sendOrderToPeerCh <- network.OrderMsg{newOrder, elevatorToHandleThisOrder}
			}
		//default:
			//fsm.HandleTimeOutOrder(hallButtonsMap)
		}
	}
}


const Col0 = "\x1b[30;1m" // Dark grey
const ColR = "\x1b[31;1m" // Red
const ColG = "\x1b[32;1m" // Green
const ColY = "\x1b[33;1m" // Yellow
const ColB = "\x1b[34;1m" // Blue
const ColM = "\x1b[35;1m" // Magenta
const ColC = "\x1b[36;1m" // Cyan
const ColW = "\x1b[37;1m" // White
const ColN = "\x1b[0m" // Grey (neutral)
