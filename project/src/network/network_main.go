package network

import (
	"./bcast"
	"./localip"
	"./peers"
	"../elevatorHW"
	"flag"
	"fmt"
	"os"
	"time"
	"../fsm"
)

// We define some custom struct to send over the network.
// Note that all members we want to transmit must be public. Any private members
//  will be received as zero-values.
type OrderMsg struct {
	Message string
	NewOrder fsm.Order
}

func Main() {
	// Our id can be anything. Here we pass it on the command line, using
	//  `go run main.go -id=our_id`
	var id string
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
	/*helloTx := make(chan HelloMsg)
	helloRx := make(chan HelloMsg)*/
	orderTx := make(chan OrderMsg)
	orderRx := make(chan OrderMsg)
	// ... and start the transmitter/receiver pair on some port
	// These functions can take any number of channels! It is also possible to
	//  start multiple transmitters/receivers on the same port.
	go bcast.Transmitter(16569, orderTx)
	go bcast.Receiver(16569, orderRx)

	// The example message. We just send one of these every second.
	/*go func() {
		helloMsg := HelloMsg{"Hello from " + id, order{4, elevatorHW.ButtonCallUp}}
		for {
			helloTx <- helloMsg
			time.Sleep(1 * time.Second)
		}
	}()*/

	go func() {
		orderMsg := OrderMsg{"Hello from " + id, fsm.Order{4, elevatorHW.ButtonCallDown}}
		for{
			orderTx <- orderMsg
			time.Sleep(1000 * time.Millisecond)
		}
	}()



	fmt.Println("Started networking")

	for {
		select {
		case p := <-peerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)

		case a := <-orderRx:
			/*if helloMsg.order.Floor != -1 {
				fsm.PutOrderInLocalQueue(helloRx.NewOrder)
			}*/
			fmt.Printf("Received: %#v\n", a) 
			var ReceivedOrder OrderMsg
			ReceivedOrder = <- orderRx
			if ReceivedOrder.NewOrder.Floor != -1 {
				fsm.PutOrderInLocalQueue(ReceivedOrder.NewOrder)
			}
		}
	}
}
