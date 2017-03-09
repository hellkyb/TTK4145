package olasnetwork

import (
	"flag"
	"fmt"
	"os"
	"time"

	"./network/bcast"
	"./network/localip"
	"./network/peers"

	"../elevatorHW"
	"../fsm"
)

// We define some custom struct to send over the network.
// Note that all members we want to transmit must be public. Any private members
//  will be received as zero-values.
type HelloMsg struct {
	Iter         int
	ElevatorID   string // This number identifies the elevator
	CurrentState int    // This number, says if the elevator is moving up (1) / down (-1) / idle (0)
	LastFloor    int    // The last floor the elevator visited
	Order        OrderMsg
}

type OrderMsg struct {
	Order                   fsm.Order
	ElevatorToTakeThisOrder string
}

var OperatingElevators int
var OperatingElevatorsPtr *int
var newElevStatus ElevatorStatus
var TimeOutQueue []int64
var Elevators []fsm.ElevatorStatus

type ElevatorStatus struct {
	ElevatorID string
	Alive      bool
}


func putNetworkOrderInLocalQueue(receivedOrder OrderMsg, myElevatorID string) {
	if receivedOrder.ElevatorToTakeThisOrder == myElevatorID {
		fsm.PutOrderInLocalQueue(receivedOrder.Order)
	}
}

func UpdateElevatorStates(newMsg HelloMsg, OperatingElevators int, operatingElevatorStates []HelloMsg) []HelloMsg {

	if len(operatingElevatorStates) == 0 {
		operatingElevatorStates = append(operatingElevatorStates, newMsg)
		return operatingElevatorStates
	}
	for i := range operatingElevatorStates {

		if newMsg.ElevatorID == operatingElevatorStates[i].ElevatorID {
			operatingElevatorStates[i] = newMsg
			return operatingElevatorStates
		} else if i == len(operatingElevatorStates)-1 {
			operatingElevatorStates = append(operatingElevatorStates, newMsg)
		}
	}
	if len(operatingElevatorStates) > OperatingElevators {
		return operatingElevatorStates[0:]
	}
	return operatingElevatorStates
}

func GetLocalID() string {
	localIP, _ := localip.LocalIP()
	return localIP[12:]
}

func NetworkMain(messageCh chan<- HelloMsg) {
	// Our id can be anything. Here we pass it on the command line, using
	//  `go run main.go -id=our_id`
	var id string
	var elevatorID string
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.Parse()

	// ... or alternatively, we can use the local IP address.
	// (But since we can run multiple programs on the same PC, we also append the
	//  process ID)
	if id == "" {
		localIP, err := localip.LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}
		id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
		elevatorID = localIP[12:]
	}

	// We make a channel for receiving updates on the id's of the peers that are
	//  alive on the network
	peerUpdateCh := make(chan peers.PeerUpdate)
	// We can disable/enable the transmitter after it has been started.
	// This could be used to signal that we are somehow "unavailable".
	peerTxEnable := make(chan bool)
	go peers.Transmitter(15647, id, peerTxEnable)
	go peers.Receiver(15647, peerUpdateCh)

	// We make channels for sending and receiving our custom data types
	helloTx := make(chan HelloMsg)
	helloRx := make(chan HelloMsg)

	// ... and start the transmitter/receiver pair on some port
	// These functions can take any number of channels! It is also possible to
	//  start multiple transmitters/receivers on the same port.
	go bcast.Transmitter(16167, helloTx)
	go bcast.Receiver(16167, helloRx)
	OperatingElevatorsPtr = &OperatingElevators

	// The example message. We just send one of these every second.
	go func() {
		order := 
		helloMsg := HelloMsg{0, elevatorID, 0, 5, order}
		for {
			helloMsg.Iter++
			helloMsg.CurrentState = elevatorHW.GetElevatorState()
			helloMsg.Order = order
			helloMsg.LastFloor = fsm.LatestFloor
			helloTx <- helloMsg

			time.Sleep(100 * time.Millisecond)
		}
	}()

	//fmt.Println("Started")

	for {
		select {
		case p := <-peerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)
			*OperatingElevatorsPtr = len(p.Peers)

		case a := <-helloRx:
			//fmt.Printf("Received: %#v\n", a)
			fmt.Print("\n\n")
			fmt.Println("Message recieved")
			fmt.Println("---------------------")
			fmt.Print("From              : ")
			fmt.Print(a.ElevatorID)
			fmt.Println(" ")
			fmt.Print("Direction         : ")
			switch a.CurrentState {
			case 0:
				fmt.Println("Not moving")
			case 1:
				fmt.Println("Going up")
			case -1:
				fmt.Println("Going down")
			}
			fmt.Print("My last floor     : ")
			fmt.Println(a.LastFloor)
			fmt.Print("Message number    : ")
			fmt.Println(a.Iter)
			fmt.Print("#Elevators online : ")
			fmt.Println(OperatingElevators)
			fmt.Print("New Order         :")
			fmt.Println(" ")
			fmt.Println("---------------------")
			fmt.Print("\n\n\n")
			messageCh <- a
		}
	}
}
